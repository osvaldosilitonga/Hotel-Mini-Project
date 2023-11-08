package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) User {
	return User{
		DB: db,
	}
}

func (controller User) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}
