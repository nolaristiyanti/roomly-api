package model

type RoomFacility struct {
	RoomID     int64 `gorm:"type:bigint;primaryKey"`
	FacilityID int64 `gorm:"type:bigint;primaryKey"`
}

// Function untuk secara eksplisit memberi tahu GORM bahwa model RoomFacility menggunakan tabel room_facilities dengan format snake_case
func (RoomFacility) TableName() string {
	return "room_facilities"
}
