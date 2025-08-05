package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"pulse/internal/config"
	"pulse/internal/model"
	"sync"
	"time"
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
