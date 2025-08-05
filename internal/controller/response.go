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

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess Code = 0
	CodeUnknown Code = 9999

	CodeParamError Code = 1000
	CodeNotFound   Code = 1001

	CodeNotAuthorized Code = 1002
)

var codeMap = map[Code]*CodeInfo{
	CodeSuccess: {
		HTTPCode: http.StatusOK,
		Msg:      "Success",
	},
	CodeUnknown: {
		HTTPCode: http.StatusBadRequest,
		Msg:      "Unknown error",
	},
	CodeParamError: {
		HTTPCode: http.StatusBadRequest,
		Msg:      "Parameter error",
	},
	CodeNotFound: {
		HTTPCode: http.StatusNotFound,
		Msg:      "Resource not found",
	},
	CodeNotAuthorized: {
		HTTPCode: http.StatusUnauthorized,
		Msg:      "Authorization error",
	},
}

type Response struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func Reply(ctx *gin.Context, reply ReplyMsg, data any) {
	ctx.JSON(GetCodeInfo(reply.Code()).HTTPCode, Response{
		Code: reply.Code(),
		Msg:  reply.Msg(),
		Data: data,
	})
}

type CodeWithMsg struct {
	code Code
	msg  string
}

func NewCodeWithMsg(code Code, msg string) *CodeWithMsg {
	return &CodeWithMsg{code: code, msg: msg}
}

func (cwm *CodeWithMsg) Code() Code {
	return cwm.code
}

func (cwm *CodeWithMsg) Msg() string {
	return cwm.msg
}

type ReplyMsg interface {
	Code() Code
	Msg() string
}

type Code uint16

func (c Code) Code() Code {
	return c
}

func (c Code) Msg() string {
	return GetCodeInfo(c).Msg
}

type CodeInfo struct {
	HTTPCode int    `json:"httpCode"`
	Msg      string `json:"msg"`
}

func GetCodeInfo(code Code) *CodeInfo {
	if info, ok := codeMap[code]; ok {
		return info
	} else {
		return codeMap[CodeUnknown]
	}
}
