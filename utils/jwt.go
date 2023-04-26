package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type CustomClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func GenerateJWT(email string) (string, CustomClaims, error) {
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
		Email: email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(viper.GetString("app.jwt_secret")))

	if err != nil {
		return "", CustomClaims{}, err
	}

	return t, claims, nil
}
