package repository

import (
	"errors"
	"loginService/internal/model"

	"github.com/google/uuid"
)

type MemoryRepository struct{
	users map[uuid.UUID]model.User
}

func (r *MemoryRepository) AddUser(username string, password string) (uuid.UUID, error){
	id, err := uuid.NewUUID();
	if (err != nil){
		return id, err
	}

	user := model.NewUser(id, username, password)
	r.users[id]	= user;
	return id, nil
}

func (r *MemoryRepository) GetUserWithUsername(username string) (model.User, error){
	for _, user := range r.users{
		if user.Username == username {
			return 	user, nil
		}		
	}
	return model.User{}, errors.New("User not found")
}

func (r *MemoryRepository) GetUserWithID(id uuid.UUID) (model.User, error){
	user, ok := r.users[id];
	if !ok{
		return model.User{}, errors.New("User not found")
	}
	return user, nil
}
