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

package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/sirupsen/logrus"

	"github.com/toodofun/pulse/internal/model"
	"github.com/toodofun/pulse/internal/util"
)

const (
	CheckerTypeHTTP model.CheckerType = "http"
)

type Checker struct {
}

type fields struct {
	URL     string            `json:"url"`
	Timeout int               `json:"timeout" default:"30"`
	Code    []int             `json:"code"    default:"[200,201,202,204,301,302]"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"  default:"GET"`
	Body    string            `json:"body"    default:""`
	Cookies map[string]string `json:"cookies"`
}

func (c *Checker) Validate(fields string) error {
	fs, err := c.fromFields(fields)
	if err != nil {
		return err
	}
	_, err = url.Parse(fs.URL)
	if err != nil {
		return errors.New("invalid param: url")
	}

	if !util.IsURL(fs.URL) {
		return errors.New("invalid param: url")
	}

	return nil
}

func (c *Checker) fromFields(fieldsStr string) (*fields, error) {
	f := new(fields)
	err := json.Unmarshal([]byte(fieldsStr), &f)
	if err != nil {
		logrus.Errorf("failed to unmarshal fields: %v", err)
	}
	defaults.SetDefaults(f)
	if f.URL == "" || (!strings.HasPrefix(f.URL, "http") && !strings.HasPrefix(f.URL, "https")) {
		return nil, errors.New("url is required and must start with http:// or https://")
	}

	f.Method = strings.ToUpper(f.Method)

	if !slices.Contains([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
		http.MethodPatch,
	}, f.Method) {
		return nil, errors.New("method is not allowed")
	}

	return f, nil
}

func (c *Checker) Check(fieldStr string) *model.Record {
	fs, err := c.fromFields(fieldStr)
	if err != nil {
		return &model.Record{
			IsSuccess: false,
			Message:   err.Error(),
			MonitorAt: time.Now(),
		}
	}

	client := http.Client{
		Timeout: time.Duration(fs.Timeout) * time.Second,
	}

	urlParsed, err := url.Parse(fs.URL)
	if err != nil {
		return &model.Record{
			IsSuccess: false,
			Message:   err.Error(),
			MonitorAt: time.Now(),
		}
	}

	header := http.Header{}

	for k, v := range fs.Headers {
		header.Set(k, v)
	}

	req := &http.Request{
		Method: fs.Method,
		URL:    urlParsed,
		Header: header,
		Body:   io.NopCloser(strings.NewReader(fs.Body)),
	}

	for k, v := range fs.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	start := time.Now()
	resp, err := client.Do(req)

	record := &model.Record{
		ResponseTime: time.Since(start).Milliseconds(),
		MonitorAt:    start,
	}
	if err != nil {
		record.IsSuccess = false
		record.Message = err.Error()
		record.ResponseTime = 0
	} else {
		defer resp.Body.Close()
		record.IsSuccess = slices.Contains(fs.Code, resp.StatusCode)
		if !record.IsSuccess {
			record.Message = fmt.Sprintf("status code %d is not in expected codes %v", resp.StatusCode, fs.Code)
		} else {
			record.Message = "OK"
		}
	}

	return record
}
