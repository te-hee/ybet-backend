package model

type LoginRequest struct {
	Username string `json:"username,omitempty" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}
