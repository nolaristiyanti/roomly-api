package database

import (
	"os"
	userModel "roomly-api/domain/user/model"
	"testing"

	"roomly-api/config"
	"roomly-api/domain/room/model"
)

func TestPreloadRoomFacilities(t *testing.T) {
	// Simpan working directory saat ini
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal("Failed to get working directory:", err)
	}

	// Pindah dari folder database ke root project
	err = os.Chdir("..")
	if err != nil {
		t.Fatal("Failed to change working directory:", err)
	}

	// Kembalikan working directory setelah test selesai
	defer os.Chdir(originalDir)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal("Failed to load config:", err)
	}

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		t.Fatal("Failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal("Failed to get SQL DB:", err)
	}
	defer sqlDB.Close()

	// Test Room -> Facilities
	t.Logf("======== Room -> Facilities ==========")
	var rooms []model.Room

	err = db.Preload("Facilities").Find(&rooms).Error
	if err != nil {
		t.Fatal("Failed to preload facilities:", err)
	}

	// SELECT *
	// FROM rooms r
	// JOIN room_facilities rf
	// 	ON r.id = rf.room_id
	// JOIN facilities f
	// 	ON f.id = rf.facility_id

	if len(rooms) == 0 {
		t.Fatal("No rooms found")
	}

	for _, r := range rooms {
		t.Logf("Room: %s", r.Name)
		t.Logf("Capacity: %d", r.Capacity)

		for _, f := range r.Facilities {
			t.Logf("Facility: %s", f.Name)
		}
	}

	// Test Room -> Bookings
	t.Logf("======== Room -> Bookings ==========")
	var roomsWithBookings []model.Room

	err = db.Preload("Bookings").Find(&roomsWithBookings).Error
	if err != nil {
		t.Fatal("Failed to preload bookings from rooms:", err)
	}

	for _, r := range roomsWithBookings {
		t.Logf("Room: %s", r.Name)

		for _, b := range r.Bookings {
			t.Logf(
				"Booking ID: %d | Start: %s | End: %s | Status: %s",
				b.ID,
				b.StartTime.Format("2006-01-02 15:04"),
				b.EndTime.Format("2006-01-02 15:04"),
				b.Status,
			)
		}
	}

	// Test User -> Bookings
	t.Logf("======== User -> Bookings ==========")
	var usersWithBookings []userModel.User

	err = db.Preload("Bookings").Find(&usersWithBookings).Error
	if err != nil {
		t.Fatal("Failed to preload bookings from users:", err)
	}

	for _, u := range usersWithBookings {
		t.Logf("User: %s | Email: %s", u.Name, u.Email)

		for _, b := range u.Bookings {
			t.Logf(
				"Booking ID: %d | Room ID: %d | Status: %s",
				b.ID,
				b.RoomID,
				b.Status,
			)
		}
	}
}
