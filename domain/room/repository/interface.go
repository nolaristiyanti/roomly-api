package repository

import (
	"context"
	"roomly-api/domain/room/entity"
)

type RoomRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.Room, error)
	FindList(ctx context.Context) ([]entity.Room, error)
}
