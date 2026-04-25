package server

import (
	userv1 "backend/proto/user/v1"
	"context"
	"crypto/sha512"
	"userService/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)



type UserServer struct{
	service service.UserService
	passwordSalt string
	userv1.UnimplementedUserServiceServer
}

func NewUserInServer(service service.UserService, passwordSalt string) *UserServer{
	return &UserServer{service: service, passwordSalt: passwordSalt}
}
	
func (s *UserServer)hashPassword(password string) string{
	hash := sha512.New()
	hash.Write([]byte(password+s.passwordSalt))
	hashed := hash.Sum(nil)
	return string(hashed)
}

func (s* UserServer) SignIn(ctx context.Context,request *userv1.SignInRequest) (*userv1.SignInResponse, error) {
	password := s.hashPassword(request.Password)	
	username := request.Username
	id, authToken, refreshToken, err := s.service.SignIn(username, password)
	
	if err != nil{
		return nil, err
	}

	response := userv1.SignInResponse{UserId: id.String(), RefreshToken: *refreshToken, AuthToken: *authToken}

	return &response, nil

}
func (s* UserServer) LogIn(ctx context.Context,request *userv1.LogInRequest) (*userv1.LogInResponse, error) {
	password := s.hashPassword(request.Password)
	username := request.Username
	
	id, authToken, refreshToken, err := s.service.LogIn(username, password)

	if err != nil {
		return nil, err
	}

	response := userv1.LogInResponse{RefreshToken: *refreshToken, AuthToken: *authToken, UserId: id.String()}

	return &response, nil
}

func (s* UserServer) GetNewAuthToken(ctx context.Context,request *userv1.GetNewAuthTokenRequest) (*userv1.GetNewAuthTokenResponse, error) {
	authToken, err := s.service.GenerateNewAuthToken(request.RefreshToken)

	if err != nil{
		return nil, err
	}

	response := userv1.GetNewAuthTokenResponse{AuthToken: *authToken}	

	return &response, nil 
}

func (s* UserServer) DoesUserExistsAndHasValidName(c context.Context, r *userv1.DoesUserExistsAndHasValidNameRequest) (*userv1.DoesUserExistsAndHasValidNameResponse, error) {
	id, err := uuid.Parse(r.UserId)

	if err != nil{
		return nil, status.Error(codes.InvalidArgument,  "id is not in uuid format")
	}

	valid := s.service.DoesUserExistsAndHasValidName(id, r.Username)

	return &userv1.DoesUserExistsAndHasValidNameResponse{Valid: valid}, nil
}
