package auth

import (
	"gateway/config"
	"gateway/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Client interface {
	GenerateToken(username string) (string, error)
}

type MinimalClient struct{}

func NewMinimalClient() *MinimalClient {
	return &MinimalClient{}
}

func (c MinimalClient) GenerateToken(username string) (string, error) {
	claims := &model.UserClaims{
		Username: username,
		UserId:   uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtkey := []byte(config.Cfg.Auth.JwtSecret)

	signed, err := token.SignedString(jwtkey)

	return signed, err
}
