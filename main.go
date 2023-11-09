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
// @version BETA
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

	authMiddleware := middlewares.NewAuthMiddleware(db)
	userController := controller.NewUserController(db)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/v1")
	{
		v1.POST("/login", userController.LoginUser)
		v1.POST("/register", userController.RegisterUser)
		v1.GET("/rooms", userController.GetRooms)
	}
	user := v1.Group("/user")
	user.Use(authMiddleware.RequiredAuth)
	{
		user.POST("/orders", userController.CreateOrder)
		user.GET("/orders/:id", userController.GetUserOrderById)
		user.PUT("/orders/cancel/:id", userController.CancelUserOrder)
		user.GET("/orders/history", userController.GetOrderHistory)

		user.PUT("/payments/topup", userController.UserTopUp)
		user.POST("/payments/process/wallet/:id", userController.PaymentWallet)
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
