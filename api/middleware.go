package api

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kanatsanan6/hrm/utils"
)

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: authError,
	})
}

func MeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := c.Locals("user").(*jwt.Token)
		claims := data.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		c.Locals("email", email)
		return c.Next()
	}
}

func authError(c *fiber.Ctx, e error) error {
	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
}
