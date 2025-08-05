package controller

import (
	"github.com/gin-gonic/gin"
	"pulse/internal/config"
)

func GetUser(ctx *gin.Context) string {
	if res, ok := ctx.Request.Context().Value(config.ContextKeyUser).(string); ok {
		return res
	}
	return ""
}
