package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	Username string `json:"username"`
	TokenType string `json:"token-type"`
	jwt.RegisteredClaims
}

func NewUserClaims(username string, userId uuid.UUID, tokenType string, expiersAt time.Time, issuer string) *UserClaims{
	return &UserClaims{
		Username: username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiersAt),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: issuer,
			Subject: userId.String(),
		},
	}
}
