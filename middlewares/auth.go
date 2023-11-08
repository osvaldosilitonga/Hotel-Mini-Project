package middlewares

import (
	"errors"
	"fmt"
	"hotel/helpers"
	"hotel/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

func NewAuthMiddleware(db *gorm.DB) Auth {
	return Auth{DB: db}
}

func (handler Auth) RequiredAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get token from header
		authHeader := c.Request().Header.Get("authorization")

		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_TOKEN_SECRET")), nil
		})
		if err != nil {
			return helpers.ErrorMessage(c, &helpers.ErrUnauthorized, err.Error())
		}

		// Token Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			return helpers.ErrorMessage(c, &helpers.ErrUnauthorized, "invalid token")
		}

		// Check token exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return helpers.ErrorMessage(c, &helpers.ErrUnauthorized, "token expired")
		}

		// Find user with email
		user, err := repository.GetUser(claims["email"].(string), handler.DB)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorMessage(c, &helpers.ErrUnauthorized, "invalid token")
		}
		if err != nil {
			return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
		}

		c.Set("email", user.Email)
		c.Set("id", user.ID)

		return next(c)
	}
}
