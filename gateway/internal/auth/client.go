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

type MinimalClient struct{
	issuer string
	authTokenTime time.Duration
}

func NewMinimalClient(issuer string, authTokenTime time.Duration) *MinimalClient {
	return &MinimalClient{issuer: issuer, authTokenTime: authTokenTime}
}

func (c MinimalClient) GenerateToken(username string) (string, error) {

	claims := model.NewUserClaims(
		username,
		uuid.New(),
		"Auth",
		time.Now().Add(c.authTokenTime),
		c.issuer,
	)	

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtkey := []byte(config.Cfg.Auth.JwtSecret)

	signed, err := token.SignedString(jwtkey)

	return signed, err
}
