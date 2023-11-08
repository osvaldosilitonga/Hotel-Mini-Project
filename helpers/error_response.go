package helpers

import "github.com/labstack/echo/v4"

func ErrorMessage(c echo.Context, apiErr *APIError, detail any) error {
	apiErr.Detail = detail

	return c.JSON(apiErr.Code, apiErr)
}
