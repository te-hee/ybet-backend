package service

import (
	"errors"
	"fmt"
	"log"
	"loginService/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LogInService struct{
	storage storage.Storage
	authTokenDuration time.Duration
	refreshTokenDuration time.Duration
	issuer string
	jwtKey string
}

func NewLogInService(storage storage.Storage, authTokenDuration time.Duration, refreshTokenDuration time.Duration, issuer string, jwtKey string) *LogInService{
	return &LogInService{storage: storage, authTokenDuration: authTokenDuration, refreshTokenDuration:refreshTokenDuration, issuer: issuer, jwtKey: jwtKey} 
}

type TokenData struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *LogInService) newTokenData(username string, userId uuid.UUID, expiersAt time.Time) TokenData{
	return TokenData{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiersAt),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: s.issuer,
			Subject: userId.String(),
		}}
}

func (s *LogInService) generateJWT(claims TokenData) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)	
	tokenString, err := token.SignedString([]byte(s.jwtKey))

	if err != nil{
		return "", err
	}

	return tokenString, nil;
}

func (s *LogInService) LogIn(username string, password string) (*string ,*string, error){
	user, err := s.storage.GetUserWithUsername(username)	

	if err != nil{
		return nil,nil,err 
	}

	if user.Password != password{
		return nil,nil,errors.New("bad password")
	}

	authToken, err := s.generateJWT(s.newTokenData(username, user.Id, time.Now().Add(s.authTokenDuration)))

	if err != nil{
		return nil, nil,err
	}

	refreshToken, err := s.generateJWT(s.newTokenData(username, user.Id, time.Now().Add(s.refreshTokenDuration)))

	if err != nil{
		return nil, nil,err
	}

	return &authToken,&refreshToken,nil
}

func (s *LogInService) SignIn(username string, password string) (*uuid.UUID, *string, *string,error){
	_, err := s.storage.GetUserWithUsername(username)

	if err == nil{
		return nil, nil,nil,errors.New("username is taken")
	}

	id , err := s.storage.AddUser(username, password)

	if(err != nil){
		return nil,nil,nil, err
	}

	authToken, err := s.generateJWT(s.newTokenData(username, id, time.Now().Add(s.authTokenDuration)))	

	if err != nil{
		return nil,nil,nil,err
	}

	refreshToken, err := s.generateJWT(s.newTokenData(username, id, time.Now().Add(s.refreshTokenDuration)))	

	if err != nil{
		return nil,nil,nil,err
	}

	return &id,&authToken,&refreshToken, nil 
}

func (s *LogInService) GenerateNewAuthToken(refreshToken string) (*string, error){
	claims := &TokenData{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		return []byte(s.jwtKey), nil
	})

	if err != nil || !token.Valid{
		return nil, err
	}

	claims, ok := token.Claims.(*TokenData);

	if !ok{
		return nil, errors.New("token claims error")
	}

	log.Println(claims.Subject)
	id, err := uuid.Parse(claims.Subject)

	if err != nil{
		log.Println(err)
		return nil, err
	}
	
	user, err := s.storage.GetUserWithID(id)

	if err != nil{
		return nil, err
	}

	if user.Username != claims.Username{
		return nil, errors.New("wrong username in token")
	}

	authToken, err := s.generateJWT(*claims)
	
	if err != nil{
		return nil, err
	}

	return &authToken, err
}
