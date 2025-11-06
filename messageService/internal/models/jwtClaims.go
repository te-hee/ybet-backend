package models

import "github.com/google/uuid"

type UserClaims struct {
	Uuid     uuid.UUID
	Username string
	Exp      int64
}
