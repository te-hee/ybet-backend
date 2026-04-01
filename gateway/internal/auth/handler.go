package auth

import (
	"gateway/internal/model"
	"gateway/internal/utils"
	"log"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service Service
}

func NewAuthHandler(service Service) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) HandleLogin(c fiber.Ctx) error{
	var input model.LoginRequest
	if err := c.Bind().Body(&input); err != nil{
		return err
	}

	token, err := h.service.GenerateToken(input.Username)

	if err != nil {
		log.Println(err)
		return utils.WriteErrorMessageWithLog(c,fiber.StatusInternalServerError,  "Error generating JWT token")
	}

	resp := model.LoginResponse{
		Token: token,
	}
	return c.JSON(resp)
}
