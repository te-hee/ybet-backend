package utils

import (
	"gateway/internal/model"
	"log"

	"github.com/gofiber/fiber/v3"
)

func WriteJsonError(c fiber.Ctx, status int, errorObject any) error {
    return c.Status(status).JSON(errorObject)
}

func WriteJsonErrorWithLog(c fiber.Ctx, status int, errorObject any) error{
	log.Println(errorObject)
	return WriteJsonError(c, status, errorObject)
}

func WriteErrorMessageWithLog(c fiber.Ctx, status int, errorMessage string) error{
	return WriteJsonErrorWithLog(c, status, model.NewOutputError(errorMessage))
}
