package dto

import "time"

type BookingResponse struct {
	ID int64 `json:"id"`
	// UserID    int64     `json:"user_id"`
	// RoomID    int64     `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`

	User BookingUserResponse `json:"user"`
	Room BookingRoomResponse `json:"room"`
}

type BookingUserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type BookingRoomResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Location string `json:"location"`
}
