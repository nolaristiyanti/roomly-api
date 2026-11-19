package entity

import "time"

type Session struct {
	JTI       string // identitas unik JWT (JWT ID)
	UserID    int64
	ExpiresAt time.Time
}
