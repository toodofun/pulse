package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
