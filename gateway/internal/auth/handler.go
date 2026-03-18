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

	var loginData model.LoginRequest

	if err := c.Bind().Body(&loginData); err != nil {
		log.Println(err)
		return utils.WriteJsonError(c, fiber.StatusBadRequest, "bad json")

	}

	token, err := h.service.GenerateToken(loginData.Username)

	if err != nil {
		log.Println(err)
		return utils.WriteJsonError(c,fiber.StatusInternalServerError,  "error generating JWT token")
	}

	resp := model.LoginResponse{
		Token: token,
	}
	return c.JSON(resp)
}
