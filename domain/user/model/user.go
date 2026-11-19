package model

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(150);not null;uniqueIndex"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Role      string    `gorm:"type:varchar(20);not null;default:user"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`

	// Bookings []bookingmodel.Booking `gorm:"foreignKey:UserID"` // Satu user punya banyak booking
}
