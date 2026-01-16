package auth

import (
	"encoding/json"
	"gateway/internal/model"
	"log"
	"net/http"
)

type AuthHandler struct {
	service Service
}

func NewAuthHandler(service Service) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	token, err := h.service.GenerateToken(loginData.Username)
	if err != nil {
		http.Error(w, "error generating JWT token", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	resp := model.LoginResponse{
		Token: token,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response QwQ: %v", err)
	}
}
