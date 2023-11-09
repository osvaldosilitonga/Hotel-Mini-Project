package controller

import (
	"errors"
	"hotel/dto"
	"hotel/entity"
	"hotel/handler"
	"hotel/helpers"
	"hotel/repository"
	"net/http"
	"time"

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
		Data:    &user,
	}
	return c.JSON(http.StatusCreated, response)
}

// GetRooms godoc
// @Summary      Get Rooms
// @Description  get all available room
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.RoomsResponse
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /rooms [get]
func (controller User) GetRooms(c echo.Context) error {
	rooms, err := repository.GetRooms(controller.DB)
	if len(*rooms) < 1 {
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "no room available")
	}
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

	response := dto.RoomsResponse{
		Message: "OK",
		Data:    *rooms,
	}

	return c.JSON(http.StatusOK, response)
}

// CreateOrder godoc
// @Summary      Create Order
// @Description  create new order by giving order information in request body
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "JWT Token"
// @Param        request   body      dto.OrderBody  true  "Order Data"
// @Success      201  {object}  dto.CreateOrderResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /user/orders [post]
func (controller User) CreateOrder(c echo.Context) error {
	userID := c.Get("id").(uint)
	body := dto.OrderBody{}

	c.Bind(&body)
	if err := c.Validate(&body); err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, err.Error())
	}

	checkIn, err := time.Parse("2006-01-02", body.CheckIn)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "invalid check_in date format (YYYY-MM-DD)")
	}
	checkOut, err := time.Parse("2006-01-02", body.CheckOut)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "invalid check_out date format (YYYY-MM-DD)")
	}

	days := checkOut.Sub(checkIn).Hours() / 24

	tx := controller.DB.Begin()

	room, err := repository.GetRoomById(body.RoomID, tx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "room not found")
	}
	if err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

	if room.Status != "ready" {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "room not available")
	}

	_, err = repository.GetOrdersByDateRange(room.ID, body.CheckIn, body.CheckOut, tx)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "room not available")
	}

	newOrder := entity.Orders{
		UserID:   userID,
		RoomID:   room.ID,
		Adult:    body.Adult,
		Child:    body.Child,
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Status:   "booked",
		Payments: entity.Payments{
			Amount: room.Price * int(days),
			Status: "unpaid",
		},
	}

	order, err := repository.CreateOrder(newOrder, tx)
	if err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

	tx.Commit()

	data := dto.OrderData{
		ID:        order.ID,
		RoomID:    order.RoomID,
		Adult:     order.Adult,
		Child:     order.Child,
		CheckIn:   order.CheckIn,
		CheckOut:  order.CheckOut,
		Status:    order.Status,
		Amount:    order.Payments.Amount,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, dto.CreateOrderResponse{
		Message: "create order success",
		Data:    data,
	})
}

// TopupBalance godoc
// @Summary      Top Up Balance
// @Description  top up balance by giving nominal in request body
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "JWT Token"
// @Param        request   body      dto.TopUpBody  true  "Top Up Nominal"
// @Success      200  {object}  dto.TopUpResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /user/payments/topup [put]
func (controller User) UserTopUp(c echo.Context) error {
	userId := c.Get("id")

	body := dto.TopUpBody{}
	c.Bind(&body)
	if err := c.Validate(&body); err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, err.Error())
	}

	user, res := repository.TopUpSaldo(userId.(uint), body.Nominal, controller.DB)
	if res.Error != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, res.Error.Error())
	}

	response := dto.TopUpResponse{
		Message: "top up success",
		Data:    *user,
	}

	mailData := dto.MailData{
		Name:    user.Name,
		Email:   user.Email,
		Nominal: body.Nominal,
		Saldo:   user.Saldo,
		Status:  "Success",
		Date:    time.Now(),
	}

	handler.SendMail(mailData)

	return c.JSON(http.StatusOK, response)
}
