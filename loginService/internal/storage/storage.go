package storage

import (
	"loginService/internal/model"

	"github.com/google/uuid"
)


type Storage interface{
	AddUser(username string, password string) (uuid.UUID, error)
	GetUserWithUsername(username string) (model.User, error)
	GetUserWithID(id uuid.UUID) (model.User, error)
}
