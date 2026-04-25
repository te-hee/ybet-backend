package handler

import (
	"errors"
	"gateway/internal/model"
	"gateway/internal/service"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct{
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler{
	return &UserHandler{service: service}	
}

func (h *UserHandler) SignIn(c fiber.Ctx) error{
	var input model.SignInRequestV2

	if err := c.Bind().JSON(&input); err != nil{
		return errors.New("Bad json")
	}

	output, err := h.service.SignIn(input.Password, input.Username)
	
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

func (h *UserHandler) LogIn(c fiber.Ctx) error{
	var input model.LogInRequestV2

	if err := c.Bind().JSON(&input); err != nil{
		return errors.New("Bad json")
	}

	output, err := h.service.LogIn(input.Password, input.Username)
	
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

func (h *UserHandler) GetNewAuthToken(c fiber.Ctx) error{
	var input model.GetNewAuthTokenRequest

	if err := c.Bind().JSON(input); err != nil{
		return errors.New("Bad json")
	}

	authToken, err := h.service.GetNewAuthToken(input.RefreshToken)
	
	if err != nil {
		return err
	}
	
	return c.Status(fiber.StatusOK).JSON(model.GetNewAuthTokenResponse{AuthToken: *authToken})
}
