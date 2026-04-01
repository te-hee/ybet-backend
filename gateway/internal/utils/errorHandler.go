package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func AppErrorHandler(c fiber.Ctx, err error) error{

	if validationError, ok := err.(validator.ValidationErrors); ok{
		return WriteErrorMessageWithLog(c, fiber.StatusBadRequest, formatError(validationError))
	}

	code := fiber.StatusInternalServerError
	
	var e *fiber.Error
	if errors.As(err, &e){
		code = e.Code
	}

	return WriteErrorMessageWithLog(c, code, e.Message) 
}

func JwtErrorHandler(c fiber.Ctx, err error) error{
	code := fiber.StatusUnauthorized
	message := "Missing or malformed JWT"
	
	var e *fiber.Error
	if errors.As(err, &e){
		code = e.Code
		message = e.Message
	}

	return WriteErrorMessageWithLog(c, code, message)

}
