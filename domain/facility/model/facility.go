package model

import (
	"time"
)

type Facility struct {
	ID          int64     `gorm:"type:bigint;primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}
