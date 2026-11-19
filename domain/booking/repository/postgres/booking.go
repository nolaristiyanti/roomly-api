package postgres

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	bookingentity "roomly-api/domain/booking/entity"
	bookingmodel "roomly-api/domain/booking/model"
	"roomly-api/domain/booking/repository"
	roomentity "roomly-api/domain/room/entity"
	roommodel "roomly-api/domain/room/model"
	userentity "roomly-api/domain/user/entity"
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
func (r *BookingRepositoryStruct) FindByID(ctx context.Context, id int64, userID *int64) (*bookingentity.Booking, error) {
	var bookingModel bookingmodel.Booking

	// Filter by ID booking
	query := r.db.WithContext(ctx).Where("id = ?", id)

	if userID != nil {
		// query ditambahkan filter pemilik booking
		query = query.Where("user_id = ?", *userID)
	}

	// err := r.db.WithContext(ctx).First(&bookingModel, id).Error
	// err := query.First(&bookingModel).Error
	err := query.
		Preload("User").
		Preload("Room").
		First(&bookingModel).Error

	// --- User
	// SELECT *
	// FROM bookings b
	// INNER JOIN users u
	// 	ON u.id = b.user_id
	// INNER JOIN rooms r
	// 	ON r.id = b.room_id
	// WHERE b.user_id = 4
	// LIMIT 1;

	// --- Admin
	// SELECT *
	// FROM bookings b
	// INNER JOIN users u
	// 	ON u.id = b.user_id
	// INNER JOIN rooms r
	// 	ON r.id = b.room_id
	// LIMIT 1;

	if err != nil {
		return nil, err
	}

	return toEntity(bookingModel), nil
}

// FindList mengambil seluruh booking
func (r *BookingRepositoryStruct) FindList(ctx context.Context, userID *int64) ([]bookingentity.Booking, error) {
	var bookingModels []bookingmodel.Booking

	query := r.db.WithContext(ctx)

	if userID != nil {
		// query ditambahkan filter pemilik booking
		query = query.Where("user_id = ?", *userID)
	}

	// err := r.db.WithContext(ctx).Order("start_time ASC").Find(&bookingModels).Error
	// err := query.Order("start_time ASC").Find(&bookingModels).Error
	err := query.
		Preload("User").
		Preload("Room").
		Order("start_time ASC").
		Find(&bookingModels).Error

	if err != nil {
		return nil, err
	}

	// --- User
	// SELECT *
	// FROM bookings b
	// INNER JOIN users u
	// 	ON u.id = b.user_id
	// INNER JOIN rooms r
	// 	ON r.id = b.room_id
	// WHERE b.user_id = 4
	// ORDER BY b.start_time ASC;

	// --- Admin
	// SELECT *
	// FROM bookings b
	// INNER JOIN users u
	// 	ON u.id = b.user_id
	// INNER JOIN rooms r
	// 	ON r.id = b.room_id
	// ORDER BY b.start_time ASC;

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
		if err := tx.Create(&createdBooking).Error; err != nil {
			return err
		}

		// 5. Ambil kembali booking beserta relasinya
		if err := tx.
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "email")
			}).
			Preload("Room", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "name", "capacity", "location")
			}).
			First(&createdBooking, createdBooking.ID).Error; err != nil {
			return err
		}

		// SELECT *
		// FROM bookings b
		// INNER JOIN users u
		// 	ON u.id = b.user_id
		// INNER JOIN rooms r
		// 	ON r.id = b.room_id
		// WHERE b.id = 10;

		return nil
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

		User: userentity.User{
			ID:    item.User.ID,
			Name:  item.User.Name,
			Email: item.User.Email,
		},

		Room: roomentity.Room{
			ID:       item.Room.ID,
			Name:     item.Room.Name,
			Capacity: item.Room.Capacity,
			Location: item.Room.Location,
		},
	}
}

// Cancel membatalkan booking
func (r *BookingRepositoryStruct) Cancel(ctx context.Context, id int64, userID *int64) error {
	// BEGIN;
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var booking bookingmodel.Booking

		// 1. Cari booking dan lock row
		query := tx.Clauses(clause.Locking{
			Strength: "UPDATE", // Mengunci row booking yang ditemukan selama transaksi berlangsung
		}).Where("id = ?", id)

		// Jika user, batasi berdasarkan pemilik booking
		if userID != nil {
			query = query.Where("user_id = ?", *userID)
		}

		err := query.First(&booking).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return repository.ErrBookingNotFound
			}
			return err
		}

		// User :
		// SELECT *
		// FROM bookings
		// WHERE id = 10
		// AND user_id = 4
		// FOR UPDATE;

		// Admin :
		// SELECT *
		// FROM bookings
		// WHERE id = 10
		// FOR UPDATE;

		// 2. Periksa status booking
		if booking.Status == "cancelled" {
			return repository.ErrBookingAlreadyCancelled
		}

		// 3. Booking hanya bisa dibatalkan sebelum waktu mulai
		if !booking.StartTime.After(time.Now()) {
			return repository.ErrBookingAlreadyStarted
		}

		// 4. Update status booking
		return tx.Model(&booking).Updates(map[string]interface{}{
			"status":     "cancelled",
			"updated_at": time.Now(),
		}).Error

		// UPDATE bookings
		// SET
		// 	status = 'cancelled',
		// 	updated_at = NOW()
		// WHERE id = 10;
	}) // COMMIT; atau ROLLBACK;
}
