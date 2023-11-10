package controller

import (
	"errors"
	"hotel/api"
	"hotel/dto"
	"hotel/entity"
	"hotel/helpers"
	"hotel/repository"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// XenditPayment godoc
// @Summary      Xendit Payment
// @Description  payment using xendit
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "JWT Token"
// @Param        OrderID   path      int  true "Order ID"
// @Success      200  {object}  dto.XenditPaymentResponse
// @Failure      400  {object}  helpers.APIError
// @Failure      401  {object}  helpers.APIError
// @Failure      404  {object}  helpers.APIError
// @Failure      500  {object}  helpers.APIError
// @Router       /user/payments/api/xendit/:id [post]
func (controller User) XenditPayment(c echo.Context) error {
	userId := c.Get("id")
	userEmail := c.Get("email")

	paramId := c.Param("id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrBadRequest, "invalid param id")
	}

	order, err := repository.GetUserOrderById(userId.(uint), orderId, controller.DB)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helpers.ErrorMessage(c, &helpers.ErrNotFound, "order not found")
	}
	if err != nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, err)
	}

	data := dto.CreateInvoiceData{
		Description: "Hotel - Mini Project (FTGO-P2)",
		Email:       userEmail.(string),
		ExternalID:  paramId,
		Amount:      float32(order.Payments.Amount),
	}

	resp := api.CreteInvoice(&data)
	if resp == nil {
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, nil)
	}

	response := dto.XenditPaymentResponse{
		Message: "ok",
		Data:    resp,
	}

	return c.JSON(http.StatusOK, response)
}

func (controller User) XenditProcessWebHook(c echo.Context) error {
	// userId := c.Get("id")

	cbToken := c.Request().Header.Get("x-callback-token")
	if cbToken != os.Getenv("XENDIT_WEBHOOK_VERIFICATION_TOKEN") {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "token validation error",
		})
	}

	body := dto.XenditCallbackBody{}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "binding body error",
		})
	}
	// if err := c.Validate(&body); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, echo.Map{})
	// }

	if strings.ToLower(body.Status) != "paid" {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "status not paid",
		})
	}

	externalId := body.ExternalID
	e, err := strconv.Atoi(externalId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "external id error",
		})
	}

	// dateString := body.PaidAt
	// paidAt, error := time.Parse("2006-01-02", dateString)
	// if error != nil {
	// 	return c.JSON(http.StatusInternalServerError, echo.Map{
	// 		"message": "convert paidAt string to time",
	// 	})
	// }

	tx := controller.DB.Begin()
	// get orders by id
	order, err := repository.GetOrderById(e, tx)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "get order by id error",
		})
	}

	order.Status = body.Status
	order.UpdatedAt = time.Now()
	order.Payments.Method = body.PaymentMethod
	order.Payments.Status = body.Status
	order.Payments.UpdatedAt = time.Now()

	// update orders - status(paid), updated_at(now)
	o := entity.Orders{}
	o = *order
	if err := tx.Save(&o).Error; err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, echo.Map{
			"message": "db transaction error",
		})
	}

	// update payment - method(PaymentMetho), status(paid), updated_at(PaidAt)
	p := entity.Payments{}
	p = order.Payments
	if err := tx.Save(&p).Error; err != nil {
		tx.Rollback()
		return helpers.ErrorMessage(c, &helpers.ErrInternalServer, echo.Map{
			"message": "db transaction error",
		})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, echo.Map{})
}
