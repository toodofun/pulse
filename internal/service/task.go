package service

import (
	"github.com/sirupsen/logrus"
	"pulse/internal/checker"
	"pulse/internal/infra"
	"pulse/internal/model"
	"time"
)

type CheckTask struct {
	service *model.Service
	db      *infra.Database
}

func NewCheckTask(service *model.Service, db *infra.Database) *CheckTask {
	return &CheckTask{
		service: service,
		db:      db,
	}
}

func (t *CheckTask) Run() {
	logrus.Debugf("checking service start: %s", t.service.Title)
	c, err := checker.GetChecker(t.service.Type)
	if err != nil {
		logrus.Errorf("get checker error: %s", err.Error())
		t.db.Create(&model.Record{
			ServiceID:    t.service.ID,
			IsSuccess:    false,
			ResponseTime: 0,
			Message:      err.Error(),
			MonitorAt:    time.Now(),
		})
		return
	}

	r := c.Check(t.service.Fields)
	r.ServiceID = t.service.ID
	t.db.Create(r)
	logrus.Debugf("checking service end: %s", t.service.Title)
}
