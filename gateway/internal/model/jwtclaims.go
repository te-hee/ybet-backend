package model

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	Username string
	UserId   string
	jwt.RegisteredClaims
}
