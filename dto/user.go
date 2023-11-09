package dto

import (
	"hotel/entity"
	"time"
)

// --------- REQUEST BODY -----------
type RegisterBody struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type OrderBody struct {
	RoomID   uint   `json:"room_id" validate:"required"`
	Adult    int    `json:"adult" validate:"required,min=1"`
	Child    int    `json:"child"`
	CheckIn  string `json:"check_in" validate:"required"`
	CheckOut string `json:"check_out" validate:"required"`
}

// --------- RESPONSE BODY -----------
type RegisterResponse struct {
	Message string         `json:"message"`
	Data    **entity.Users `json:"data"`
}
type LoginResponse struct {
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}
type RoomsResponse struct {
	Message string
	Data    []entity.Rooms `json:"data"`
}
type CreateOrderResponse struct {
	Message string    `json:"message"`
	Data    OrderData `json:"data"`
}
type UserOrderByIdResponse struct {
	Message string    `json:"message"`
	Data    OrderData `json:"data"`
}

// New Data Format
type OrderData struct {
	ID        uint      `json:"order_id"`
	RoomID    uint      `json:"room_id"`
	Adult     int       `json:"adult"`
	Child     int       `json:"child"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
	Status    string    `json:"status"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
