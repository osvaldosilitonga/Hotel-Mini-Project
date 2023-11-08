package entity

import "time"

type Users struct {
	ID       uint
	Name     string
	Email    string
	Password string `json:"-"`
	Saldo    int
}

type Orders struct {
	ID        uint
	UserID    int
	RoomID    int
	Adult     int
	Child     int
	CheckIn   time.Time
	CheckOut  time.Time
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     Users    `gorm:"foreignKey:UserID"`
	Rooms     Rooms    `gorm:"foreignKey:RoomID"`
	Payments  Payments `gorm:"foreignKey:OrderID"`
}
