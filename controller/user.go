package controller

import (
	"errors"
	"hotel/dto"
	"hotel/helpers"
	"hotel/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

func (controller User) LoginUser(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "OK",
	})
}

func (controller User) RegisterUser(c echo.Context) error {
	body := dto.RegisterBody{}

	c.Bind(&body)
	if err := c.Validate(&body); err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

	body.Password = string(hash)

	// insert to db
	user, err := repository.RegisterUser(&body, controller.DB)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "email is already registered")
	}
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

	response := dto.RegisterResponse{
		Message: "register success",
		Data:    user,
	}
	return c.JSON(http.StatusOK, response)
}
