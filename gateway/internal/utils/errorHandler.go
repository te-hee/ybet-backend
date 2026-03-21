package utils

import (
	"errors"
	"github.com/gofiber/fiber/v3"
)

func AppErrorHandler(c fiber.Ctx, err error) error{
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e){
		code = e.Code
	}

	return WriteErrorMessageWithLog(c, code, e.Message) 
}

func JwtErrorHandler(c fiber.Ctx, err error) error{
	return WriteErrorMessageWithLog(c, fiber.StatusUnauthorized, "Missing or malformed JWT")

}
