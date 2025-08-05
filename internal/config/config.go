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

package config

import (
	"time"

	"github.com/mcuadros/go-defaults"
)

const (
	ContextKeyConfig     ContextKey = "config"
	ContextKeyInstanceId ContextKey = "instanceId"
	ContextKeyUser       ContextKey = "user"
)

var config *Config

type ContextKey string

type Config struct {
	Server      Server                 `json:"server"   yaml:"server"`
	JWT         JWT                    `json:"jwt"      yaml:"jwt"`
	OAuthConfig map[string]OAuthConfig `json:"oauth"    yaml:"oauth"`
	Database    Database               `json:"database" yaml:"database"`
}

func Current() *Config {
	return config
}

func New(path string) *Config {
	config = new(Config)
	if path != "" {
		_ = NewStructFromFile(path, config)
	}
	defaults.SetDefaults(config)
	return config
}

type Server struct {
	BaseURL     string        `json:"baseURL"     yaml:"baseURL"     default:"http://pulse.toodo.fun"`
	Port        int           `json:"port"        yaml:"port"        default:"80"`
	Prefix      string        `json:"prefix"      yaml:"prefix"      default:"/api/v1"`
	Debug       bool          `json:"debug"       yaml:"debug"       default:"false"`
	GracePeriod time.Duration `json:"gracePeriod" yaml:"gracePeriod" default:"30s"`
}

type JWT struct {
	Secret string        `json:"secret" yaml:"secret" default:"aurora"`
	Issuer string        `json:"issuer" yaml:"issuer" default:"fun.toodo.aurora"`
	Expire time.Duration `json:"expire" yaml:"expire" default:"720h"`
}

type OAuthConfig struct {
	AuthType     string `json:"authType"     yaml:"authType"`
	AuthURL      string `json:"authURL"      yaml:"authURL"`
	TokenURL     string `json:"tokenURL"     yaml:"tokenURL"`
	ClientId     string `json:"clientId"     yaml:"clientId"`
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
}

type Database struct {
	Driver string `json:"driver" yaml:"driver" default:"sqlite"`
	DSN    string `json:"dsn"    yaml:"dsn"    default:"db.sqlite"`

	MaxIdleConn int           `json:"maxIdleConn" yaml:"maxIdleConn" default:"10"`
	MaxOpenConn int           `json:"maxOpenConn" yaml:"maxOpenConn" default:"40"`
	ConnMaxLift time.Duration `json:"connMaxLift" yaml:"connMaxLift" default:"0s"`
	ConnMaxIdle time.Duration `json:"connMaxIdle" yaml:"connMaxIdle" default:"0s"`
}
