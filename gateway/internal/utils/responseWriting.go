package utils

import "github.com/gofiber/fiber/v3"

func WriteJsonError(c fiber.Ctx, status int, errorMessage any) error {
    return c.Status(status).JSON(fiber.Map{
        "error": errorMessage,
    })
}
