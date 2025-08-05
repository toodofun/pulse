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

package oauth

import (
	"context"
	"fmt"

	"github.com/google/go-github/v61/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/toodofun/pulse/internal/config"
)

type GithubProvider struct {
	config *oauth2.Config
}

func NewGithubProvider(conf config.OAuthConfig) *GithubProvider {
	return &GithubProvider{
		config: &oauth2.Config{
			Scopes: []string{"user:email", "read:user"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			ClientID:     conf.ClientId,
			ClientSecret: conf.ClientSecret,
			RedirectURL:  config.Current().Server.BaseURL + "/login/callback?oauth=github",
		},
	}
}

func (p *GithubProvider) GetAuthURL(redirectURL string) string {
	return p.config.AuthCodeURL(redirectURL, oauth2.AccessTypeOffline)
}

func (p *GithubProvider) GetInfo(code string) (*UserInfo, error) {
	logrus.Debugf("GetInfo.code: %s", code)

	token, err := p.config.Exchange(context.TODO(), code)
	if err != nil {
		logrus.Errorf("get token failed, error: %v", err)
		return nil, err
	}

	logrus.Debugf("GetInfo.token: %v", token)

	client := github.NewClient(nil).WithAuthToken(token.AccessToken)

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		logrus.Errorf("get user failed, error: %v", err)
		return nil, err
	}

	logrus.Debugf("GetInfo.user: %+v", user)
	userInfo := new(UserInfo)
	userInfo.ID = fmt.Sprintf("%d.github", user.GetID())
	userInfo.Nickname = user.GetName()
	userInfo.Avatar = user.GetAvatarURL()
	userInfo.Email = user.GetEmail()
	userInfo.UserType = AuthTypeGithub
	userInfo.Username = fmt.Sprintf("%s.github", user.GetLogin())

	logrus.Infof("GetInfo.userInfo: %+v", userInfo)

	return userInfo, nil
}
