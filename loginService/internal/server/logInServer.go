package server

import (
	loginv1 "backend/proto/login/v1"
	"context"
	"crypto/sha512"
	"loginService/internal/service"
)



type LogInServer struct{
	service service.LogInService
	passwordSalt string
	loginv1.UnimplementedLoginServiceServer
}

func NewLogInServer(service service.LogInService, passwordSalt string) *LogInServer{
	return &LogInServer{service: service, passwordSalt: passwordSalt}
}
	
func (s *LogInServer)hashPassword(password string) string{
	hash := sha512.New()
	hash.Write([]byte(password+s.passwordSalt))
	hashed := hash.Sum(nil)
	return string(hashed)
}

func (s* LogInServer) SignUp(ctx context.Context,request *loginv1.SignUpRequest) (*loginv1.SignUpResponse, error) {
	password := s.hashPassword(request.Password)	
	username := request.Username
	id, authToken, refreshToken, err := s.service.SignIn(username, password)
	
	if err != nil{
		return nil, err
	}

	response := loginv1.SignUpResponse{UserId: id.String(), RefreshToken: *refreshToken, AuthToken: *authToken}

	return &response, nil

}
func (s* LogInServer) LogIn(ctx context.Context,request *loginv1.LogInRequest) (*loginv1.LogInResponse, error) {
	password := s.hashPassword(request.Password)
	username := request.Username
	
	authToken, refreshToken, err := s.service.LogIn(username, password)

	if err != nil {
		return nil, err
	}

	response := loginv1.LogInResponse{RefreshToken: *refreshToken, AuthToken: *authToken}

	return &response, nil
}

func (s* LogInServer) GetNewAuthToken(ctx context.Context,request *loginv1.GetNewAuthTokenRequest) (*loginv1.GetNewAuthTokenResponse, error) {
	authToken, err := s.service.GenerateNewAuthToken(request.RefreshToken)

	if err != nil{
		return nil, err
	}

	response := loginv1.GetNewAuthTokenResponse{AuthToken: *authToken}	
	return &response, err
}
