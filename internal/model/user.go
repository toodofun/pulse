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

type User struct {
	ID       string `json:"id"       gorm:"primary_key;"    swaggerignore:"true"`
	Username string `json:"username" gorm:"unique;not null"                      validate:"required"`
	Nickname string `json:"nickname"                                             validate:"required"`
	Password string `json:"-"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"                                                validate:"required"`
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
