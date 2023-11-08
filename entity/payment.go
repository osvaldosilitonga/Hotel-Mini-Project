package entity

import "time"

type Payments struct {
	ID        uint
	OrderID   int
	Method    string
	Amount    int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
