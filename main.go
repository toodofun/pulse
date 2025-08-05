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

package main

import (
	"github.com/sirupsen/logrus"

	"github.com/toodofun/pulse/internal/config"
	"github.com/toodofun/pulse/internal/server"
)

func main() {
	svc, err := server.New(config.New("config.yaml"))
	if err != nil {
		panic(err)
	}

	logrus.Infof("running with config: %+v", config.Current())

	if err = svc.Run(); err != nil {
		panic(err)
	}
}
