package usecase

import (
	"context"
	"errors"
	"net/http"
	"time"

	apperrors "roomly-api/common/errors"
	"roomly-api/domain/booking/entity"
	"roomly-api/domain/booking/repository"

	"gorm.io/gorm"
)

type BookingUsecase struct {
	bookingRepository repository.BookingRepository
}

func NewBookingUsecase(bookingRepository repository.BookingRepository) *BookingUsecase {
	return &BookingUsecase{
		bookingRepository: bookingRepository,
	}
}

func (u *BookingUsecase) CreateBooking(ctx context.Context, booking entity.Booking) (*entity.Booking, error) {
	// 1. Validasi aturan waktu
	if err := validateBookingTime(booking.StartTime, booking.EndTime); err != nil {
		return nil, err
	}

	// 2. Simpan melalui repository
	result, err := u.bookingRepository.Create(ctx, booking)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRoomNotFound):
			return nil, apperrors.New(
				http.StatusNotFound,
				"ROOM_NOT_FOUND",
				"Room not found",
				nil,
			)

		case errors.Is(err, repository.ErrRoomInactive):
			return nil, apperrors.New(
				http.StatusUnprocessableEntity,
				"ROOM_INACTIVE",
				"Room is not active",
				nil,
			)

		case errors.Is(err, repository.ErrBookingConflict):
			return nil, apperrors.New(
				http.StatusConflict,
				"BOOKING_CONFLICT",
				"Room is already booked for the selected time",
				nil,
			)
		}

		return nil, err
	}

	return result, nil
}

func (u *BookingUsecase) GetBookingByID(ctx context.Context, id int64, userID int64, role string) (*entity.Booking, error) {
	var filterUserID *int64

	if role != "admin" {
		filterUserID = &userID
	}

	booking, err := u.bookingRepository.FindByID(ctx, id, filterUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(
				http.StatusNotFound,
				"BOOKING_NOT_FOUND",
				"Booking not found",
				nil,
			)
		}

		return nil, err
	}

	return booking, nil
}

func (u *BookingUsecase) GetBookingList(ctx context.Context, userID int64, role string) ([]entity.Booking, error) {
	var filterUserID *int64

	if role != "admin" {
		filterUserID = &userID
	}

	return u.bookingRepository.FindList(ctx, filterUserID)
}

func validateBookingTime(startTime, endTime time.Time) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	start := startTime.In(loc)
	end := endTime.In(loc)
	now := time.Now().In(loc)

	// Start harus sebelum End
	if !start.Before(end) {
		return apperrors.New(
			http.StatusBadRequest,
			"INVALID_BOOKING_TIME",
			"Start time must be before end time",
			nil,
		)
	}

	// Booking tidak boleh di masa lampau
	if start.Before(now) {
		return apperrors.New(
			http.StatusBadRequest,
			"PAST_BOOKING",
			"Booking cannot be made in the past",
			nil,
		)
	}

	// Harus pada tanggal yang sama
	if start.Year() != end.Year() ||
		start.Month() != end.Month() ||
		start.Day() != end.Day() {
		return apperrors.New(
			http.StatusBadRequest,
			"INVALID_BOOKING_DATE",
			"Booking must be on the same date",
			nil,
		)
	}

	// Hanya Senin sampai Jumat
	if start.Weekday() == time.Saturday || start.Weekday() == time.Sunday {
		return apperrors.New(
			http.StatusBadRequest,
			"INVALID_BOOKING_DAY",
			"Booking is only available on weekdays",
			nil,
		)
	}

	// Jam harus tepat pada jam penuh
	if start.Minute() != 0 || start.Second() != 0 || start.Nanosecond() != 0 ||
		end.Minute() != 0 || end.Second() != 0 || end.Nanosecond() != 0 {
		return apperrors.New(
			http.StatusBadRequest,
			"INVALID_BOOKING_INTERVAL",
			"Booking must use whole-hour intervals",
			nil,
		)
	}

	// Jam operasional 08.00 - 17.00
	if start.Hour() < 8 || end.Hour() > 17 {
		return apperrors.New(
			http.StatusBadRequest,
			"OUTSIDE_OPERATION_HOURS",
			"Booking is only available from 08:00 to 17:00",
			nil,
		)
	}

	// Durasi minimal 1 jam, maksimal 9 jam
	duration := end.Sub(start)

	if duration < time.Hour || duration > 9*time.Hour {
		return apperrors.New(
			http.StatusBadRequest,
			"INVALID_BOOKING_DURATION",
			"Booking duration must be between 1 and 9 hours",
			nil,
		)
	}

	return nil
}
