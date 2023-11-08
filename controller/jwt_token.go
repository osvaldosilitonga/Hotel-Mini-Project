package controller

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id uint, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 5).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET")))

	return tokenString, err
}
