package entity

type Rooms struct {
	ID         uint   `json:"id"`
	Category   string `json:"category"`
	RoomNumber int    `json:"room_number"`
	Status     string `json:"-"`
	Price      int    `json:"price"`
}
