package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"pulse/internal/checker"
	"pulse/internal/infra"
	"pulse/internal/model"
	"sort"
	"sync"
	"time"
)

const (
	ThresholdSuccessRatio = 99.9 // 成功率阈值
	ThresholdWarningRatio = 95.0 // 警告阈值

	ColorSuccess = "var(--color-green-400)"  // 成功颜色
	ColorWarning = "oklch(75% 0.183 55.934)" // 警告颜色
	ColorFail    = "var(--color-red-400)"    // 失败颜色
	ColorBlack   = "var(--color-gray-400)"   // 无数据
)

type MonitorService struct {
	cron   *cron.Cron
	jobMap sync.Map

	db *infra.Database
}

func NewMonitorService() *MonitorService {
	cronClient := cron.New(cron.WithSeconds())
	cronClient.Start()

	return &MonitorService{
		cron: cronClient,
	}
}

func (s *MonitorService) checkService(service *model.Service) error {
	if service.Interval <= 0 {
		return fmt.Errorf("service interval must be greater than 0")
	}

	if service.Type == "" {
		return fmt.Errorf("service type cannot be empty")
	}

	if service.Title == "" {
		return fmt.Errorf("service title cannot be empty")
	}

	if service.CreatedBy == "" {
		return fmt.Errorf("service createdBy cannot be empty")
	}

	c, err := checker.GetChecker(service.Type)
	if err != nil {
		return fmt.Errorf("invalid check type: %w", err)
	}
	if err = c.Validate(service.Fields); err != nil {
		return err
	}

	return nil
}

func (s *MonitorService) ListServices(operator string) ([]*model.Service, error) {
	var services []*model.Service
	if err := s.db.Where(&model.Service{CreatedBy: operator}).Find(&services).Error; err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	for _, service := range services {
		s.fitRecords(service)
	}

	return services, nil
}

func (s *MonitorService) fitRecords(service *model.Service) *model.Service {
	var records []model.Record
	if err := s.db.Where(&model.Record{ServiceID: service.ID}).
		Order("monitor_at DESC").
		Limit(25).
		Find(&records).Error; err != nil {
		logrus.Errorf("failed to list records: %v", err)
		return service
	}

	service.Records = records
	if len(records) > 0 {
		service.IsSuccess = records[0].IsSuccess
	}

	return service
}

func (s *MonitorService) GetServiceByID(id, operator string) (*model.Service, error) {
	var service model.Service
	if err := s.db.First(&service, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to find service")
	}
	if service.CreatedBy != operator {
		return nil, fmt.Errorf("no permission")
	}
	return s.fitRecords(&service), nil
}

func (s *MonitorService) UpdateService(service *model.Service, id, operator string) error {
	var res model.Service
	if err := s.db.First(&res, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to find service: %w", err)
	}

	if res.CreatedBy != operator {
		return fmt.Errorf("no permission")
	}

	if err := s.checkService(service); err != nil {
		return err
	}

	service.CreatedBy = operator
	service.ID = res.ID
	service.CreatedAt = res.CreatedAt
	service.Enabled = res.Enabled

	if err := s.db.Save(service).Error; err != nil {
		logrus.Errorf("failed to update service: %v", err)
		return fmt.Errorf("failed to update service")
	}

	return nil
}

func (s *MonitorService) AddService(service *model.Service, operator string) error {
	service.CreatedBy = operator
	if err := s.checkService(service); err != nil {
		return err
	}

	service.ID = ""

	if err := s.db.Create(service).Error; err != nil {
		return fmt.Errorf("failed to add service: %w", err)
	}

	if service.Enabled {
		if err := s.addCron(service, s.db); err != nil {
			return fmt.Errorf("failed to add cron job for service %s: %w", service.Title, err)
		}
	}

	return nil
}

func (s *MonitorService) DeleteService(serviceID string, operator string) error {
	var service model.Service
	if err := s.db.First(&service, "id = ?", serviceID).Error; err != nil {
		return fmt.Errorf("failed to find service: %w", err)
	}
	if service.CreatedBy != operator {
		return fmt.Errorf("no permission")
	}

	s.delCron(&service)

	if err := s.db.Delete(&service).Error; err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

func (s *MonitorService) SetEnabled(serviceID string, enabled bool, operator string) error {
	var service model.Service
	if err := s.db.First(&service, "id = ?", serviceID).Error; err != nil {
		return fmt.Errorf("failed to find service: %w", err)
	}

	if service.CreatedBy != operator {
		return fmt.Errorf("no permission")
	}

	service.Enabled = enabled

	if err := s.db.Save(&service).Error; err != nil {
		return fmt.Errorf("failed to pause service: %w", err)
	}

	if enabled {
		return s.addCron(&service, s.db)
	} else {
		s.delCron(&service)
	}

	return nil
}

func (s *MonitorService) GetDailySuccessRatios(serviceID string, operator string, public bool) ([]*model.Ratio, error) {
	var service model.Service
	if err := s.db.First(&service, "id = ?", serviceID).Error; err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}

	if service.Private && service.CreatedBy != operator {
		return nil, fmt.Errorf("no permission")
	}

	if public && service.Private {
		return nil, fmt.Errorf("no permission")
	}

	var results []*model.Ratio
	today := time.Now().Truncate(24 * time.Hour)

	for i := 0; i < 90; i++ {
		dayStart := today.AddDate(0, 0, -i)
		dayEnd := dayStart.Add(24 * time.Hour)

		var total int64
		var success int64

		// 查询当天总记录数
		if err := s.db.Model(&model.Record{}).
			Where("service_id = ? AND monitor_at >= ? AND monitor_at < ?", serviceID, dayStart, dayEnd).
			Count(&total).Error; err != nil {
			return nil, err
		}

		r := &model.Ratio{
			Date: dayStart.Format("2006-01-02"),
		}

		if total == 0 {
			// 没有任何数据
			r.Ratio = 0.0
			r.Color = ColorBlack
		} else {
			// 有数据，统计成功数量
			if err := s.db.Model(&model.Record{}).
				Where("service_id = ? AND is_success = ? AND monitor_at >= ? AND monitor_at < ?", serviceID, true, dayStart, dayEnd).
				Count(&success).Error; err != nil {
				return nil, err
			}

			r.Ratio = float64(success) / float64(total) * 100

			// 根据成功率决定标签
			switch {
			case r.Ratio >= ThresholdSuccessRatio:
				r.Color = ColorSuccess
			case r.Ratio >= ThresholdWarningRatio:
				r.Color = ColorWarning
			default:
				r.Color = ColorFail
			}
		}

		results = append(results, r)
	}

	// 倒序返回（从最近到最远）
	sort.Slice(results, func(i, j int) bool {
		return results[i].Date > results[j].Date
	})

	return results, nil
}

func (s *MonitorService) addCron(service *model.Service, db *infra.Database) error {
	task := NewCheckTask(service, db)
	jobID, err := s.cron.AddJob(fmt.Sprintf("@every %ds", service.Interval), task)
	if err != nil {
		return fmt.Errorf("failed to add job for service %s: %w", service.Title, err)
	}
	s.jobMap.Store(service.ID, jobID)
	logrus.Infof("Job for service %s added with ID %d", service.Title, jobID)

	// 立即执行一次任务
	go task.Run()

	return nil
}

func (s *MonitorService) delCron(service *model.Service) {
	if jobID, ok := s.jobMap.Load(service.ID); ok {
		s.cron.Remove(jobID.(cron.EntryID))
		s.jobMap.Delete(service.ID)
		logrus.Infof("Job for service %s removed", service.Title)
	}
}

func (s *MonitorService) start(db *infra.Database) error {
	services := make([]*model.Service, 0)
	if err := db.Where(&model.Service{Enabled: true}).Find(&services, &model.Service{Enabled: true}).Error; err != nil {
		return err
	}
	for _, service := range services {
		if err := s.addCron(service, db); err != nil {
			return err
		}
	}
	return nil
}

func (s *MonitorService) Initialize(db *infra.Database) error {
	if err := db.AutoMigrate(&model.Service{}, &model.Record{}); err != nil {
		return err
	}
	if err := s.start(db); err != nil {
		return err
	}
	s.db = db
	return nil
}
