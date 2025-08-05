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

package checker

import (
	"fmt"

	"github.com/toodofun/pulse/internal/checker/http"
	"github.com/toodofun/pulse/internal/model"
)

type Checker interface {
	Check(fields string) *model.Record
	Validate(fields string) error
}

func GetChecker(t model.CheckerType) (Checker, error) {
	switch t {
	case http.CheckerTypeHTTP:
		return &http.Checker{}, nil
	default:
		return nil, fmt.Errorf("unknown checker: %s", t)
	}
}
