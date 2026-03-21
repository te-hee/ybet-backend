package utils

import (
	"gateway/internal/model"
	"log"

	"github.com/gofiber/fiber/v3"
)

func WriteJsonError(c fiber.Ctx, status int, errorObject any) error {
    return c.Status(status).JSON(errorObject)
}

func WriteJsonErrorWithLog(c fiber.Ctx, status int, errorMessage any) error{
	log.Println(errorMessage)
	return WriteJsonError(c, status, errorMessage)
}

func WriteErrorMessageWithLog(c fiber.Ctx, status int, errorMessage string) error{
	return WriteJsonErrorWithLog(c, status, model.NewOutputError(errorMessage))
}
