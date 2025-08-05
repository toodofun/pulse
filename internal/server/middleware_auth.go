package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/glob"
	"pulse/internal/config"
	"pulse/internal/controller"
	"pulse/internal/service"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func skipAuthCheck(path string) bool {
	skipPaths := []string{
		"/api/v1/login/*",
		"/api/v1/monitor/*/daily/public",
	}

	for _, p := range skipPaths {
		g := glob.MustCompile(p)
		if g.Match(path) {
			return true
		}
	}
	return false
}

func AuthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !skipAuthCheck(ctx.Request.URL.Path) {
			authHeader := ctx.GetHeader(authorizationHeader)
			token := authHeader
			if len(authHeader) == 0 {
				token, _ = ctx.Cookie("token")
			}

			if len(token) == 0 {
				controller.Reply(ctx, controller.CodeNotAuthorized, nil)
				ctx.Abort()
				return
			}

			u, err := service.GetJWTService().ParseToken(strings.TrimPrefix(token, "Bearer "))
			if err != nil {
				controller.Reply(ctx, controller.CodeNotAuthorized, nil)
				ctx.Abort()
				return
			}
			ctx.Request = ctx.Request.WithContext(
				context.WithValue(ctx.Request.Context(), config.ContextKeyUser, u.Username),
			)
		}
		ctx.Next()
	}
}
