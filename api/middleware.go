package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/spf13/viper"
)

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(viper.GetString("app.jwt_secret")),
		ErrorHandler: authError,
	})
}

func MeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := c.Locals("user").(*jwt.Token)
		claims := data.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		fmt.Println(email)
		c.Locals("email", email)
		return c.Next()
	}
}

func authError(c *fiber.Ctx, e error) error {
	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
}
