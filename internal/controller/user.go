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
	"github.com/gin-gonic/gin"

	"github.com/toodofun/pulse/internal/service"
)

type UserController struct {
	svc *service.UserService
}

func NewUserController(svc *service.UserService) *UserController {
	return &UserController{
		svc: svc,
	}
}

func (c *UserController) handleGetAvailableOAuthTypes(ctx *gin.Context) {
	types := c.svc.GetAvailableOAuthTypes()
	Reply(ctx, CodeSuccess, types)
}

func (c *UserController) handleGetOAuthURL(ctx *gin.Context) {
	if oauthType, ok := ctx.GetQuery("oauth"); !ok {
		Reply(ctx, CodeParamError, nil)
	} else {
		redirectURL, ok := ctx.GetQuery("redirect")
		if !ok {
			redirectURL = "/"
		}

		if res, err := c.svc.GetOAuthURL(oauthType, redirectURL); err != nil {
			Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
		} else {
			Reply(ctx, CodeSuccess, res)
		}
	}
}

func (c *UserController) handleCallback(ctx *gin.Context) {
	if oauthName, ok := ctx.GetQuery("oauth"); !ok {
		Reply(ctx, CodeParamError, nil)
	} else {
		if code, ok := ctx.GetQuery("code"); !ok {
			Reply(ctx, CodeParamError, nil)
		} else {
			redirectURL, ok := ctx.GetQuery("state")
			if !ok {
				redirectURL = "/"
			}

			token, err := c.svc.GetOAuthToken(oauthName, code)
			if err != nil {
				Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
				return
			}

			Reply(ctx, CodeSuccess, map[string]interface{}{
				"token":    token,
				"redirect": redirectURL,
			})
		}
	}
}

func (c *UserController) handleGetUserInfo(ctx *gin.Context) {
	user, err := c.svc.GetUserInfo(GetUser(ctx))
	if err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
		return
	}
	Reply(ctx, CodeSuccess, user)
}

func (c *UserController) RegisterRoute(group *gin.RouterGroup) {
	loginApi := group.Group("/login")
	loginApi.GET("/oauth/types", c.handleGetAvailableOAuthTypes)
	loginApi.GET("/oauth", c.handleGetOAuthURL)
	loginApi.GET("/callback", c.handleCallback)

	userApi := group.Group("/user")
	userApi.GET("/info", c.handleGetUserInfo)
}
