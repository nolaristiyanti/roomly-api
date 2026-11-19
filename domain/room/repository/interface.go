package repository

import (
	"context"
	"errors"
	"roomly-api/domain/room/entity"
)

var ErrRoomHasBookings = errors.New("room has bookings")

type RoomRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.Room, error)
	FindList(ctx context.Context) ([]entity.Room, error)
	Create(ctx context.Context, room entity.Room) (*entity.Room, error)
	Update(ctx context.Context, room entity.Room) (*entity.Room, error)
	Delete(ctx context.Context, id int64) error
}
