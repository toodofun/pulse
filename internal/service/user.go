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

package service

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/toodofun/pulse/internal/infra"
	"github.com/toodofun/pulse/internal/model"
	"github.com/toodofun/pulse/internal/service/oauth"
)

type UserService struct {
	db *infra.Database
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUserInfo(name string) (user *model.User, err error) {
	logrus.Infof("GetUserInfo: %+v", name)
	var u *model.User
	if err = s.db.Where("username = ?", name).First(&u, &model.User{Username: name}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found") // User not found
		}
		return nil, err // Other error
	}
	return u, nil
}

func (s *UserService) GetAvailableOAuthTypes() []oauth.AvailableOAuth {
	return oauth.GetOAuthManager().GetAvailableOAuth()
}

func (s *UserService) GetOAuthURL(authType, redirectURI string) (string, error) {
	if provider, err := oauth.GetOAuthManager().GetAuthProvider(authType); err != nil {
		return "", err
	} else {
		return provider.GetAuthURL(redirectURI), nil
	}
}

func (s *UserService) GetOAuthToken(authType, code string) (string, error) {
	if provider, err := oauth.GetOAuthManager().GetAuthProvider(authType); err != nil {
		return "", err
	} else {
		userInfo, err := provider.GetInfo(code)
		if err != nil {
			return "", err
		}

		user := &model.User{
			ID:       userInfo.ID,
			Username: userInfo.Username,
			Nickname: userInfo.Nickname,
			Password: "",
			Avatar:   userInfo.Avatar,
			Email:    userInfo.Email,
			Type:     "oauth",
		}

		var u *model.User
		if err = s.db.First(&u, &model.User{ID: user.ID}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = s.db.Create(user).Error; err != nil {
					return "", err
				}
			} else {
				return "", err
			}
		} else {
			user.CreatedAt = u.CreatedAt
			if err := s.db.Save(user).Error; err != nil {
				return "", err
			}
			user = u
		}
		return GetJWTService().CreateToken(user)
	}
}

func (s *UserService) Initialize(db *infra.Database) error {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	s.db = db
	return nil
}
