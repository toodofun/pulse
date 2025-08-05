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

package server

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/glob"

	"github.com/toodofun/pulse/internal/config"
	"github.com/toodofun/pulse/internal/controller"
	"github.com/toodofun/pulse/internal/service"
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
