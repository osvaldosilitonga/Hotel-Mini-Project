package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func LogrusConfig() middleware.RequestLoggerConfig {
	log := logrus.New()
	return middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":     values.URI,
				"Method":  c.Request().Method,
				"Status":  values.Status,
				"Latency": values.Latency,
			}).Info("request")

			return nil
		},
	}
}
