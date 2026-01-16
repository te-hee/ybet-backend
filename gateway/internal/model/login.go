package model

type LoginRequest struct {
	Username string `json:"username,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}
