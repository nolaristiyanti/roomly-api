package usecase

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	apperrors "roomly-api/common/errors"
	"roomly-api/domain/room/entity"
	"roomly-api/domain/room/repository"
)

type RoomUsecase struct {
	roomRepository repository.RoomRepository
}

func NewRoomUsecase(roomRepository repository.RoomRepository) *RoomUsecase {
	return &RoomUsecase{
		roomRepository: roomRepository,
	}
}

// GetRoomByID mengambil detail room berdasarkan ID
func (u *RoomUsecase) GetRoomByID(ctx context.Context, id int64) (*entity.Room, error) {
	room, err := u.roomRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(
				http.StatusNotFound,
				"ROOM_NOT_FOUND",
				"Room not found",
				nil,
			)
		}

		return nil, err
	}

	return room, nil
}

// GetRoomList mengambil seluruh daftar room
func (u *RoomUsecase) GetRoomList(ctx context.Context) ([]entity.Room, error) {
	return u.roomRepository.FindList(ctx)
}
