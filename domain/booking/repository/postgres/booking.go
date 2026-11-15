package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	bookingentity "roomly-api/domain/booking/entity"
	bookingmodel "roomly-api/domain/booking/model"
	"roomly-api/domain/booking/repository"
	roommodel "roomly-api/domain/room/model"
)

type BookingRepositoryStruct struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) repository.BookingRepository {
	return &BookingRepositoryStruct{
		db: db,
	}
}

// FindByID mengambil booking berdasarkan ID
func (r *BookingRepositoryStruct) FindByID(ctx context.Context, id int64) (*bookingentity.Booking, error) {
	var bookingModel bookingmodel.Booking

	err := r.db.WithContext(ctx).First(&bookingModel, id).Error
	if err != nil {
		return nil, err
	}

	return toEntity(bookingModel), nil
}

// FindList mengambil seluruh booking
func (r *BookingRepositoryStruct) FindList(ctx context.Context) ([]bookingentity.Booking, error) {
	var bookingModels []bookingmodel.Booking

	err := r.db.WithContext(ctx).Order("start_time ASC").Find(&bookingModels).Error

	if err != nil {
		return nil, err
	}

	bookingEntities := make([]bookingentity.Booking, 0, len(bookingModels))

	for _, item := range bookingModels {
		bookingEntities = append(bookingEntities, *toEntity(item))
	}

	return bookingEntities, nil
}

// Create membuat booking dengan transaction dan row locking
func (r *BookingRepositoryStruct) Create(ctx context.Context, booking bookingentity.Booking) (*bookingentity.Booking, error) {
	var createdBooking bookingmodel.Booking

	// Semua operasi di dalam callback berada dalam satu transaction
	// 		BEGIN;
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Lock room yang akan dipesan
		var room roommodel.Room

		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).First(&room, booking.RoomID).Error

		// SELECT *
		// FROM rooms
		// WHERE id = 1
		// ORDER BY id
		// LIMIT 1
		// FOR UPDATE; // mengunci row data yang di select selama transaksi berlangsung

		// Jika dua user mencoba booking room disaat yang sama dimana :
		// Request User A sampai ke PostgreSQL sedikit lebih dahulu
		// PostgreSQL akan membuat transaksi B menunggu sampai transaksi A melepaskan lock-nya.

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return repository.ErrRoomNotFound
			}
			return err
		}

		// 2. Pastikan room aktif
		if !room.IsActive {
			return repository.ErrRoomInactive
		}

		// -- Cek jadwal
		// 3. Periksa booking yang bertabrakan
		var count int64

		// existing_start < requested_end AND existing_end > requested_start
		err = tx.Model(&bookingmodel.Booking{}).
			Where("room_id = ?", booking.RoomID).
			Where("status = ?", "confirmed").
			Where("start_time < ?", booking.EndTime).
			Where("end_time > ?", booking.StartTime).
			Count(&count).Error

		// SELECT COUNT(*)
		// FROM bookings
		// WHERE room_id = 1
		//   AND status = 'confirmed'
		//   AND start_time < '2026-12-21 10:00:00'
		//   AND end_time > '2026-12-21 09:00:00';

		if err != nil {
			return err
		}

		if count > 0 {
			return repository.ErrBookingConflict
		}

		// -- Hasil: tidak ada jadwal bentrok

		// 4. Simpan booking
		createdBooking = bookingmodel.Booking{
			UserID:    booking.UserID,
			RoomID:    booking.RoomID,
			StartTime: booking.StartTime,
			EndTime:   booking.EndTime,
			Status:    "confirmed",
		}

		// INSERT INTO bookings (...);
		return tx.Create(&createdBooking).Error
	}) // GORM melakukan COMMIT atau ROLLBACK

	if err != nil {
		return nil, err
	}

	return toEntity(createdBooking), nil
}

// Mapping Model ke Entity
func toEntity(item bookingmodel.Booking) *bookingentity.Booking {
	return &bookingentity.Booking{
		ID:        item.ID,
		UserID:    item.UserID,
		RoomID:    item.RoomID,
		StartTime: item.StartTime,
		EndTime:   item.EndTime,
		Status:    item.Status,
	}
}
