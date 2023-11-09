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
