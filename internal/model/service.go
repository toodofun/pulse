// Copyright 2025 The Toodofun Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID        string      `json:"id"        gorm:"type:varchar(64);primary_key"`
	Title     string      `json:"title"     gorm:"type:varchar(255);not null"`
	IsSuccess bool        `json:"isSuccess" gorm:"-"`
	Interval  int         `json:"interval"  gorm:"not null;default:300"`
	Type      CheckerType `json:"type"      gorm:"index;type:varchar(8);not null"`
	Enabled   bool        `json:"enabled"   gorm:"not null;default:true"`
	Private   bool        `json:"private"   gorm:"not null;default:true"`
	Fields    string      `json:"fields"`
	Records   []Record    `json:"records"   gorm:"-"`

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
