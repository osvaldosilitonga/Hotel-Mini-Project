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

// LoginUser godoc
// @Summary      Login
// @Description  login by giving user credential in request body
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request   body      dto.LoginBody  true  "User Credential"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /login [post]
func (controller User) LoginUser(c echo.Context) error {
	body := dto.LoginBody{}
	c.Bind(&body)
	if err := c.Validate(&body); err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, err.Error())
	}

	user, err := repository.GetUser(body.Email, controller.DB)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "email not found")
	}
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "wrong password")
	}

	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

	response := dto.LoginResponse{
		Message:     "login success",
		AccessToken: token,
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterUser godoc
// @Summary      Register
// @Description  register by giving user information in request body
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request   body      dto.RegisterBody  true  "User Data"
// @Success      201  {object}  dto.RegisterResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /register [post]
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
	return c.JSON(http.StatusCreated, response)
}
