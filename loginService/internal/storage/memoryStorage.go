package storage 

import (
	"errors"
	"loginService/internal/model"

	"github.com/google/uuid"
)

type MemoryStorage struct{
	users map[uuid.UUID]model.User
}

func NewMemoryStorage() *MemoryStorage{
	return &MemoryStorage{users: make(map[uuid.UUID]model.User)}
}

func (r *MemoryStorage) AddUser(username string, password string) (uuid.UUID, error){
	id, err := uuid.NewUUID();
	if (err != nil){
		return id, err
	}

	user := model.NewUser(id, username, password)
	r.users[id]	= user;
	return id, nil
}

func (r *MemoryStorage) GetUserWithUsername(username string) (model.User, error){
	for _, user := range r.users{
		if user.Username == username {
			return 	user, nil
		}		
	}
	return model.User{}, errors.New("User not found")
}

func (r *MemoryStorage) GetUserWithID(id uuid.UUID) (model.User, error){
	user, ok := r.users[id];
	if !ok{
		return model.User{}, errors.New("User not found")
	}
	return user, nil
}
