package utils

import "github.com/gofiber/fiber/v2"

func ErrorResponse(c *fiber.Ctx, status int, errorMsg string) error {
	return c.Status(status).JSON(fiber.Map{
		"errors": errorMsg,
	})
}

func JsonResponse(c *fiber.Ctx, status int, response interface{}) error {
	return c.Status(status).JSON(fiber.Map{"data": response})
}
