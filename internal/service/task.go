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
	"time"

	"github.com/sirupsen/logrus"

	"github.com/toodofun/pulse/internal/checker"
	"github.com/toodofun/pulse/internal/infra"
	"github.com/toodofun/pulse/internal/model"
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
