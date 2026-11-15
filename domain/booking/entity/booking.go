package entity

import "time"

type Booking struct {
	ID        int64
	UserID    int64
	RoomID    int64
	StartTime time.Time
	EndTime   time.Time
	Status    string
}
