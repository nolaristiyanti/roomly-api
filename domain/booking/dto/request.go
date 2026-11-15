package dto

import "time"

type CreateBookingRequest struct {
	UserID    int64     `json:"user_id" binding:"required,gt=0"`
	RoomID    int64     `json:"room_id" binding:"required,gt=0"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}
