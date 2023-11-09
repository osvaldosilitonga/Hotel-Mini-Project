package controller

import (
	"errors"
	"hotel/dto"
	"hotel/helpers"
	"hotel/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// GetUserOrderByID godoc
// @Summary      Get User Order
// @Description  get user order by giving order id in request param
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "JWT Token"
// @Param        OrderID   path      int  true "Order ID"
// @Success      200  {object}  dto.UserOrderByIdResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      401  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /user/orders/:id [get]
func (controller User) GetUserOrderById(c echo.Context) error {
	userId := c.Get("id")

	param := c.Param("id")
	orderId, err := strconv.Atoi(param)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "id must be number")
	}

	order, err := repository.GetUserOrderById(userId.(uint), orderId, controller.DB)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "data not found")
	}
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err.Error())
	}

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

	return c.JSON(http.StatusOK, dto.UserOrderByIdResponse{
		Message: "ok",
		Data:    data,
	})
}

// CancelUserOrder godoc
// @Summary      Cancel User Order
// @Description  cancel user order by giving order id in request param
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "JWT Token"
// @Param        OrderID   path      int  true "Order ID"
// @Success      200  {object}  dto.CancelOrderResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      401  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /user/orders/cancel/:id [put]
func (controller User) CancelUserOrder(c echo.Context) error {
	userId := c.Get("id")

	param := c.Param("id")
	orderId, err := strconv.Atoi(param)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "id must be number")
	}

	order, res := repository.CancelUserOrder(userId.(uint), orderId, controller.DB)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "order id not found")
	}
	if res.RowsAffected < 1 {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "no rows affected")
	}
	if res.Error != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, res.Error.Error())
	}

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

	return c.JSON(http.StatusOK, dto.CancelOrderResponse{
		Message: "cancel success",
		Data:    data,
	})
}

// GetOrderHistory godoc
// @Summary      Order History
// @Description  get all user order history
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "JWT Token"
// @Success      200  {object}  dto.OrderHistoryResponse
// @Failure      401  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /user/orders/history [get]
func (controller User) GetOrderHistory(c echo.Context) error {
	userId := c.Get("id")

	orders, err := repository.GetOrderHistory(userId.(uint), controller.DB)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "empty data")
	}
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err)
	}

	data := []dto.OrderData{}

	for _, v := range orders {
		d := dto.OrderData{
			ID:        v.ID,
			RoomID:    v.RoomID,
			Adult:     v.Adult,
			Child:     v.Child,
			CheckIn:   v.CheckIn,
			CheckOut:  v.CheckOut,
			Status:    v.Status,
			Amount:    v.Payments.Amount,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}

		data = append(data, d)
	}

	return c.JSON(http.StatusOK, dto.OrderHistoryResponse{
		Message: "ok",
		Data:    data,
	})
}
