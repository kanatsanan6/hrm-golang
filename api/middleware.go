package api

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/kanatsanan6/hrm/utils"
	"github.com/spf13/viper"
)

func AuthMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(viper.GetString("app.jwt_secret")),
		ErrorHandler: authError,
	})
}

func authError(c *fiber.Ctx, e error) error {
	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "unauthorized")
}
