package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type CustomClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func GenerateJWT(email string) (string, jwt.MapClaims, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("app.jwt_secret")))

	if err != nil {
		return "", nil, err
	}

	return t, claims, nil
}
