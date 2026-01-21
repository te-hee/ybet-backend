package models

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserId   string
	Username string
	jwt.RegisteredClaims
}
