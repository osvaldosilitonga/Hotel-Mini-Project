package entity

import "time"

type Payments struct {
	ID        uint
	OrderID   int
	Method    string
	Amount    int
	status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
