package usecase

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	apperrors "roomly-api/common/errors"
	facilityrepository "roomly-api/domain/facility/repository"
	"roomly-api/domain/room/entity"
	"roomly-api/domain/room/repository"
)

type RoomUsecase struct {
	roomRepository     repository.RoomRepository
	facilityRepository facilityrepository.FacilityRepository
}

func NewRoomUsecase(roomRepository repository.RoomRepository, facilityRepository facilityrepository.FacilityRepository) *RoomUsecase {
	return &RoomUsecase{
		roomRepository:     roomRepository,
		facilityRepository: facilityRepository,
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

func (u *RoomUsecase) CreateRoom(ctx context.Context, room entity.Room) (*entity.Room, error) {
	// Validasi setiap Facility yang dipilih
	for _, facility := range room.Facilities {

		_, err := u.facilityRepository.FindByID(ctx, facility.ID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.New(
					http.StatusNotFound,
					"FACILITY_NOT_FOUND",
					"Facility not found",
					nil,
				)
			}

			return nil, err
		}
	}

	// Simpan Room beserta relasi Facility
	return u.roomRepository.Create(ctx, room)
}

func (u *RoomUsecase) UpdateRoom(ctx context.Context, room entity.Room) (*entity.Room, error) {
	// 1. Validasi Room berdasarkan ID
	_, err := u.GetRoomByID(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	// 2. Validasi setiap Facility yang dipilih
	for _, facility := range room.Facilities {

		_, err := u.facilityRepository.FindByID(ctx, facility.ID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperrors.New(
					http.StatusNotFound,
					"FACILITY_NOT_FOUND",
					"Facility not found",
					nil,
				)
			}

			return nil, err
		}
	}

	// 3. Update Room beserta relasi Facility
	return u.roomRepository.Update(ctx, room)
}

func (u *RoomUsecase) DeleteRoom(ctx context.Context, id int64) error {
	// 1. Pastikan Room tersedia
	_, err := u.GetRoomByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Hapus Room
	err = u.roomRepository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRoomHasBookings) {
			return apperrors.New(
				http.StatusConflict,
				"ROOM_HAS_BOOKINGS",
				"Room cannot be deleted because it has bookings",
				nil,
			)
		}

		return err
	}

	return nil
}
