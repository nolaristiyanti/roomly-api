package model

import (
	bookingmodel "roomly-api/domain/booking/model"
	facilitymodel "roomly-api/domain/facility/model"
	"time"
)

type Room struct {
	ID          int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(150);not null"`
	Description *string   `gorm:"type:text"`
	Capacity    int       `gorm:"type:int;not null"`
	Location    string    `gorm:"type:varchar(255);not null"`
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`

	Bookings   []bookingmodel.Booking   `gorm:"foreignKey:RoomID"`          // Satu room punya banyak booking
	Facilities []facilitymodel.Facility `gorm:"many2many:room_facilities;"` // Satu room punya banyak facility
}
