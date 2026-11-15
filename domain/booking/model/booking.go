package model

import "time"

type Booking struct {
	ID        int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	UserID    int64     `gorm:"type:bigint;not null"`
	RoomID    int64     `gorm:"type:bigint;not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Status    string    `gorm:"type:varchar(20);not null;default:confirmed"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}
