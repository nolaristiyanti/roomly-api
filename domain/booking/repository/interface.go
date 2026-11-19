package repository

import (
	"context"
	"errors"

	"roomly-api/domain/booking/entity"
)

// Repository memberi tahu Usecase tentang kondisi yang terjadi tanpa harus menentukan HTTP status code. Nantinya Usecase menerjemahkannya menjadi APIError
var (
	ErrRoomNotFound            = errors.New("room not found")
	ErrRoomInactive            = errors.New("room inactive")
	ErrBookingConflict         = errors.New("booking time conflict")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrBookingAlreadyCancelled = errors.New("booking already cancelled")
	ErrBookingAlreadyStarted   = errors.New("booking already started")
)

type BookingRepository interface {
	FindByID(ctx context.Context, id int64, userID *int64) (*entity.Booking, error)
	FindList(ctx context.Context, userID *int64) ([]entity.Booking, error)
	Create(ctx context.Context, booking entity.Booking) (*entity.Booking, error)
	Cancel(ctx context.Context, id int64, userID *int64) error
}
