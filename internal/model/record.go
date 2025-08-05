package model

import (
	"gorm.io/gorm"
	"time"
)

type Record struct {
	ID           uint64         `json:"id" gorm:"primary_key;"`
	ServiceID    string         `json:"serviceId" gorm:"type:varchar(64);not null;index"`
	IsSuccess    bool           `json:"isSuccess" gorm:"not null;index"`
	ResponseTime int64          `json:"responseTime" gorm:"index"`
	Message      string         `json:"message" gorm:"size:1024"`
	MonitorAt    time.Time      `json:"monitorAt" gorm:"not null;index"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}
