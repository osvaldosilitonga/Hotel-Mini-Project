package controller

import (
	"errors"
	"hotel/dto"
	"hotel/entity"
	"hotel/helpers"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// WalletPayment godoc
// @Summary      Wallet Pay
// @Description  payment using user wallet
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "JWT Token"
// @Param        OrderID   path      int  true "Order ID"
// @Success      200  {object}  dto.PaymentResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      401  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /user/payments/process/wallet/:id [post]
func (controller User) PaymentWallet(c echo.Context) error {
	userId := c.Get("id")
	paramId := c.Param("id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "invalid param id")
	}

	tx := controller.DB.Begin()

	order := entity.Orders{}

	err = tx.Preload("Payments").Preload("Users").Where("id = ? AND user_id = ?", orderId, userId.(uint)).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "order not found")
	}
	if err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err)
	}

	totalAmount := order.Payments.Amount
	balance := order.Users.Saldo

	isValid := BalanceCheck(balance, totalAmount)
	if !isValid {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "not enough balance. please top up first")
	}

	// kurangin saldo user
	order.Users.Saldo -= order.Payments.Amount

	order.Payments.Status = "paid"
	order.Payments.Method = "wallet"
	order.Payments.UpdatedAt = time.Now()
	order.Status = "paid"
	order.UpdatedAt = time.Now()

	// update user saldo
	u := entity.Users{}
	u = order.Users
	if err := tx.Save(&u).Error; err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err)
	}

	// update orders status dan updated_at
	o := entity.Orders{}
	o = order
	if err := tx.Save(&o).Error; err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err)
	}

	// update payment status, method, updated_at
	p := entity.Payments{}
	p = order.Payments
	if err := tx.Save(&p).Error; err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err)
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

	return c.JSON(http.StatusOK, dto.PaymentResponse{
		Message: "payment success",
		Data:    data,
	})
}
