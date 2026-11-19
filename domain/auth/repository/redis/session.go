package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"roomly-api/domain/auth/entity"
	authrepository "roomly-api/domain/auth/repository"
)

type SessionRepositoryStruct struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) *SessionRepositoryStruct {
	return &SessionRepositoryStruct{
		client: client,
	}
}

// Membuat session baru di Redis
func (r *SessionRepositoryStruct) Create(ctx context.Context, session entity.Session) error {

	// Menghitung sisa masa berlaku session
	ttl := time.Until(session.ExpiresAt)

	if ttl <= 0 {
		return errors.New("session expiration must be in the future")
	}

	// Mengubah entity menjadi JSON
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Membuat Redis key berdasarkan JTI
	key := fmt.Sprintf("session:%s", session.JTI)

	// Menyimpan session dengan TTL
	return r.client.Set(ctx, key, data, ttl).Err()
}

// Mencari session berdasarkan JTI
func (r *SessionRepositoryStruct) FindByID(ctx context.Context, jti string) (*entity.Session, error) {

	key := fmt.Sprintf("session:%s", jti)

	// Mengambil data dari Redis
	data, err := r.client.Get(ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil, authrepository.ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	// Mengubah JSON menjadi entity
	var session entity.Session

	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Menghapus session berdasarkan JTI
func (r *SessionRepositoryStruct) Delete(ctx context.Context, jti string) error {

	key := fmt.Sprintf("session:%s", jti)

	return r.client.Del(ctx, key).Err()
}
