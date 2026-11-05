package service

import (
	"log"
	"userService/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type UserService struct {
	storage              storage.Storage
	authTokenDuration    time.Duration
	refreshTokenDuration time.Duration
	issuer               string
	jwtKey               string
}

func NewUserInService(storage storage.Storage, authTokenDuration time.Duration, refreshTokenDuration time.Duration, issuer string, jwtKey string) *UserService {
	return &UserService{storage: storage, authTokenDuration: authTokenDuration, refreshTokenDuration: refreshTokenDuration, issuer: issuer, jwtKey: jwtKey}
}

type TokenData struct {
	Username string `json:"username"`
	TokenType string `json:"token-type"`
	jwt.RegisteredClaims
}

func (s *UserService) newTokenData(username string, tokenType string, userId uuid.UUID, expiersAt time.Time) TokenData {
	return TokenData{
		Username: username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiersAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   userId.String(),
		}}
}

func (s *UserService) generateJWT(claims TokenData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(s.jwtKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) LogIn(username string, password string) (*uuid.UUID,*string, *string, error) {
	user, err := s.storage.GetUserWithUsername(username)

	if err != nil {
		return nil, nil,nil, err
	}

	if user.Password != password {
		return nil, nil,nil, status.Error(codes.InvalidArgument,"bad password")
	}

	authToken, err := s.generateJWT(s.newTokenData(username,"Auth", user.Id, time.Now().Add(s.authTokenDuration)))

	if err != nil {
		return nil, nil, nil,err
	}

	refreshToken, err := s.generateJWT(s.newTokenData(username, "Refresh",user.Id, time.Now().Add(s.refreshTokenDuration)))

	if err != nil {
		return nil, nil, nil,err
	}

	return &user.Id, &authToken, &refreshToken, nil
}

func (s *UserService) SignIn(username string, password string) (*uuid.UUID, *string, *string, error) {
	_, err := s.storage.GetUserWithUsername(username)

	if err == nil {
		return nil, nil, nil, status.Error(codes.AlreadyExists,"username is taken")
	}

	id, err := s.storage.AddUser(username, password)

	if err != nil {
		return nil, nil, nil, err
	}

	authToken, err := s.generateJWT(s.newTokenData(username,"Auth", id, time.Now().Add(s.authTokenDuration)))

	if err != nil {
		return nil, nil, nil, err
	}

	refreshToken, err := s.generateJWT(s.newTokenData(username, "Refresh",id, time.Now().Add(s.refreshTokenDuration)))

	if err != nil {
		return nil, nil, nil, err
	}

	return &id, &authToken, &refreshToken, nil
}

func (s *UserService) GenerateNewAuthToken(refreshToken string) (*string, error) {
	claims := &TokenData{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(s.jwtKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenData)

	if !ok {
		return nil, status.Error(codes.PermissionDenied,"token claims error")
	}

	if claims.TokenType != "Refresh"{
		return nil, status.Error(codes.PermissionDenied, "Not refresh token passed")
	}

	id, err := uuid.Parse(claims.Subject)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	user, err := s.storage.GetUserWithID(id)

	if err != nil {
		return nil, err
	}

	if user.Username != claims.Username {
		return nil, status.Error(codes.PermissionDenied, "wrong username in token")
	}

	authToken, err := s.generateJWT(s.newTokenData(user.Username, "Auth", user.Id, time.Now().Add(s.authTokenDuration)))

	if err != nil {
		return nil, err
	}

	return &authToken, err
}

func (s *UserService) DoesUserExistsAndHasValidName(id uuid.UUID, username string) (bool){
	user, err := s.storage.GetUserWithID(id)

	if err != nil{
		return false
	}

	return user.Username == username
}
