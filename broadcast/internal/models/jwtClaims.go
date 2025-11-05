package models

import "github.com/google/uuid"

type UserClaims struct {
	Uuid           uuid.UUID
	Username       string
	ExpirationDate int64
}
