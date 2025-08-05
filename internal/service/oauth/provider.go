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

package oauth

import (
	"errors"
	"sync"

	"github.com/toodofun/pulse/internal/config"
)

var (
	once    sync.Once
	manager *AuthManager
)

const (
	AuthTypeGithub = "github"
	AuthTypeGitlab = "gitlab"
)

var (
	ErrAuthTypeNotSupport = errors.New("auth type not support")
	ErrAuthFailed         = errors.New("auth failed, please try again")
)

type Provider interface {
	GetInfo(code string) (*UserInfo, error)
	GetAuthURL(redirectURL string) string
}

type AuthManager struct {
	config map[string]config.OAuthConfig
}

func GetOAuthManager() *AuthManager {
	once.Do(func() {
		manager = &AuthManager{
			config: config.Current().OAuthConfig,
		}
	})
	return manager
}

type AvailableOAuth struct {
	OAuth string `json:"oauth" yaml:"oauth"`
	Type  string `json:"type"  yaml:"type"`
}

func (m *AuthManager) GetAvailableOAuth() []AvailableOAuth {
	res := make([]AvailableOAuth, 0)
	for k, v := range m.config {
		res = append(res, AvailableOAuth{OAuth: k, Type: v.AuthType})
	}
	return res
}

func (m *AuthManager) GetAuthProvider(authName string) (Provider, error) {
	if conf, ok := m.config[authName]; !ok {
		return nil, ErrAuthTypeNotSupport
	} else {
		switch conf.AuthType {
		case AuthTypeGithub:
			return NewGithubProvider(conf), nil
		default:
			return nil, ErrAuthTypeNotSupport
		}
	}
}

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	UserType string `json:"userType"`
}
