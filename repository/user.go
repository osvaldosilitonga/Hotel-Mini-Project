package repository

import (
	"hotel/dto"
	"hotel/entity"

	"gorm.io/gorm"
)

func RegisterUser(body *dto.RegisterBody, db *gorm.DB) (*entity.Users, error) {
	user := entity.Users{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password,
	}

	if err := db.Create(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil
}

func GetUser(email string, db *gorm.DB) (*entity.Users, error) {
	user := entity.Users{}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return &user, err
	}

	return &user, nil
}

func GetRooms(db *gorm.DB) (*[]entity.Rooms, error) {
	rooms := []entity.Rooms{}

	err := db.Where("status = ?", "ready").Find(&rooms).Error

	return &rooms, err
}

func GetRoomById(roomID uint, db *gorm.DB) (*entity.Rooms, error) {
	room := entity.Rooms{}

	err := db.Where("id = ?", roomID).First(&room).Error

	return &room, err
}

func GetOrdersByDateRange(roomId uint, checkIn, checkOut string, db *gorm.DB) (*entity.Orders, error) {
	order := entity.Orders{}

	err := db.Where("room_id = ? AND (check_in >= ? AND check_in < ? OR check_out > ? AND check_out <= ?)", roomId, checkIn, checkOut, checkIn, checkOut).First(&order).Error

	return &order, err
}

func CreateOrder(data entity.Orders, db *gorm.DB) (entity.Orders, error) {
	// order := entity.Orders{}

	err := db.Preload("Payments").Create(&data).Error
	return data, err
}

func GetOrderById(orderId int, db *gorm.DB) (*entity.Orders, error) {
	order := entity.Orders{}
	err := db.Preload("Payments").Where("id = ?", orderId).First(&order).Error

	return &order, err
}

func GetUserOrderById(userId uint, orderId int, db *gorm.DB) (*entity.Orders, error) {
	order := entity.Orders{}
	err := db.Preload("Payments").Where("user_id = ? AND id = ?", userId, orderId).First(&order).Error

	return &order, err
}

func CancelUserOrder(userId uint, orderId int, db *gorm.DB) (*entity.Orders, *gorm.DB) {
	order := entity.Orders{}

	result := db.Preload("Payments").Where("id = ? AND user_id = ?", orderId, userId).First(&order)
	if result.Error != nil {
		return &order, result
	}

	order.Status = "cancel"
	order.Payments.Status = "cancel"

	result = db.Save(&order)

	return &order, result
}

func GetOrderHistory(userId uint, db *gorm.DB) ([]entity.Orders, error) {
	order := []entity.Orders{}

	err := db.Where("user_id = ?", userId).Find(&order).Error

	return order, err
}

func TopUpSaldo(userId uint, nominal int, db *gorm.DB) (*entity.Users, *gorm.DB) {
	user := entity.Users{}

	result := db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return &user, result
	}

	user.Saldo = user.Saldo + nominal
	result = db.Save(&user)

	return &user, result
}
