package repository

import (
	"context"
	"errors"

	"roomly-api/domain/auth/entity"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository interface {
	// Menyimpan session ke Redis
	Create(ctx context.Context, session entity.Session) error

	// Mengecek session berdasarkan JTI
	FindByID(ctx context.Context, jti string) (*entity.Session, error)

	// Menghapus session saat logout
	Delete(ctx context.Context, jti string) error
}
