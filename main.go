package main

import (
	"hotel/config"
	"hotel/controller"
	"hotel/initializers"
	"hotel/middlewares"
	"os"

	_ "hotel/docs"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Hotel API (Mini Project)
// @version 1.0.0
// @description Mini project Hotel API (FTGO-P2)

// @contact.name Osvaldo Silitonga
// @contact.url https://github.com/osvaldosilitonga
// @contact.email osvaldosilitonga@gmail.com

// @securityDefinitions.basic  BasicAuth

// @host localhost:8080
// @BasePath /v1

func init() {
	initializers.LoadEnvFile()
}

func main() {
	e := echo.New()
	e.Validator = &initializers.CustomValidator{Validator: validator.New()}

	e.Use(middleware.RequestLoggerWithConfig(middlewares.LogrusConfig()))

	db := config.InitDB()

	userController := controller.NewUserController(db)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/test", userController.Test)

	e.Logger.Fatal(e.Start(":" + os.Getenv("LOCAL_PORT")))
}
