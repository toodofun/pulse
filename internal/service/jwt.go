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

package service

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"github.com/toodofun/pulse/internal/config"
	"github.com/toodofun/pulse/internal/model"
)

var (
	onceJWT    sync.Once
	jwtService *JWTService
)

type JWTService struct {
	jwt config.JWT
}

func GetJWTService() *JWTService {
	onceJWT.Do(func() {
		jwtService = &JWTService{
			jwt: config.Current().JWT,
		}
	})
	return jwtService
}

func (s *JWTService) CreateToken(user *model.User) (string, error) {
	logrus.Debugf("create token for user: %+v", user)
	now := time.Now()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			User: user,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    s.jwt.Issuer,
				ExpiresAt: jwt.NewNumericDate(now.Add(s.jwt.Expire)),
				NotBefore: jwt.NewNumericDate(now),
				ID:        user.ID,
			},
		},
	)
	return token.SignedString([]byte(s.jwt.Secret))
}

func (s *JWTService) ParseToken(tokenString string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.User, nil
	}
	return nil, errors.New("invalid token")
}

type CustomClaims struct {
	User *model.User
	jwt.RegisteredClaims
}
