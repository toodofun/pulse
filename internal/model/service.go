package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	ID        string      `json:"id" gorm:"type:varchar(64);primary_key"`
	Title     string      `json:"title" gorm:"type:varchar(255);not null"`
	IsSuccess bool        `json:"isSuccess" gorm:"-"`
	Interval  int         `json:"interval" gorm:"not null;default:300"`
	Type      CheckerType `json:"type" gorm:"index;type:varchar(8);not null"`
	Enabled   bool        `json:"enabled" gorm:"not null;default:true"`
	Private   bool        `json:"private" gorm:"not null;default:true"`
	Fields    string      `json:"fields"`
	Records   []Record    `json:"records" gorm:"-"`

	CreatedBy string         `json:"createdBy" gorm:"type:varchar(64);not null"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type CheckerType string

func (m *Service) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return nil
}
