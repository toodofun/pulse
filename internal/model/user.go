package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       string `json:"id" gorm:"primary_key;" swaggerignore:"true"`
	Username string `json:"username" gorm:"unique;not null" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Password string `json:"-"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" validate:"required"`
	Type     string `json:"type"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if len(u.ID) == 0 {
		u.ID = uuid.NewString()
	}
	return nil
}
