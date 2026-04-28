package model

import "github.com/google/uuid"

type LoginRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type LogInRequestV2 struct{
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type LogInResponseV2 struct{
	Id uuid.UUID `json:"id"`
	RefreshToken string `json:"refresh_token"`
	AuthToken string  `json:"auth_token"`
}

type SignInRequestV2 struct{
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type SignInResponseV2 struct{
	Id uuid.UUID `json:"id"`
	RefreshToken string `json:"refresh_token"`
	AuthToken string  `json:"auth_token"`
}

type GetNewAuthTokenRequest struct{
	RefreshToken string `json:"refresh_token"`
}

type GetNewAuthTokenResponse struct{
	AuthToken string `json:"auth_token"`
}
