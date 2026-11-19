package entity

import (
	roomentity "roomly-api/domain/room/entity"
	userentity "roomly-api/domain/user/entity"
	"time"
)

type Booking struct {
	ID        int64
	UserID    int64
	RoomID    int64
	StartTime time.Time
	EndTime   time.Time
	Status    string

	User userentity.User
	Room roomentity.Room
}
