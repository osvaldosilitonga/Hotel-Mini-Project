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
