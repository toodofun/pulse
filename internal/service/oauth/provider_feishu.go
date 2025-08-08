package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/toodofun/pulse/internal/config"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

type FeishuProvider struct {
	config *oauth2.Config
}

func NewFeishuProvider(conf config.OAuthConfig) *FeishuProvider {
	return &FeishuProvider{
		config: &oauth2.Config{
			ClientID:     conf.ClientId,
			ClientSecret: conf.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  conf.AuthURL,
				TokenURL: conf.TokenURL,
			},
			RedirectURL: config.Current().Server.BaseURL + "/login/callback?oauth=feishu",
			Scopes:      []string{"user_info", "email", "contact:user.email:readonly"},
		},
	}
}

func (p *FeishuProvider) GetAuthURL(redirectURI string) string {
	return p.config.AuthCodeURL(redirectURI, oauth2.AccessTypeOffline)
}

func (p *FeishuProvider) GetInfo(code string) (*UserInfo, error) {
	token, err := p.config.Exchange(context.TODO(), code)
	if err != nil {
		logrus.Errorf("get token failed, error: %v", err)
		return nil, err
	}
	user, err := p.getFeishuUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		ID:       fmt.Sprintf("%s.feishu", user.OpenID),
		Username: fmt.Sprintf("%s.feishu", user.OpenID),
		Nickname: user.Name,
		Email:    user.Email,
		Avatar:   user.AvatarURL,
		UserType: AuthTypeFeishu,
	}

	return userInfo, nil
}

func (p *FeishuProvider) getFeishuUserInfo(accessToken string) (*FeishuUserInfo, error) {
	url := "https://open.feishu.cn/open-apis/authen/v1/user_info"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置Authorization头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求用户信息失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var userResp FeishuResponse
	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("解析用户信息响应失败: %v", err)
	}

	if userResp.Code != 0 {
		return nil, fmt.Errorf("飞书返回错误: %s (code: %d)", userResp.Msg, userResp.Code)
	}

	return &userResp.Data, nil
}

type FeishuUserInfo struct {
	OpenID    string `json:"open_id"`
	UnionID   string `json:"union_id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}

type FeishuResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data FeishuUserInfo  `json:"data"`
}
