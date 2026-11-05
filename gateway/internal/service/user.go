package service

import (
	"gateway/internal/client"
	"gateway/internal/model"
)

type UserService struct{
	client client.UserClient
}

func NewUserService(client client.UserClient) *UserService{
	return &UserService{client: client}
}

func (s *UserService) SignIn(password string, username string) (*model.SignInResponseV2, error){
	return s.client.SignIn(password, username)
}

func (s *UserService) LogIn(password string, username string ) (*model.LogInResponseV2, error){
	return s.client.LogIn(password, username)
}

func (s *UserService) GetNewAuthToken(refreshToken string) (*string, error){
	return s.client.GetNewAuthToken(refreshToken)
}
