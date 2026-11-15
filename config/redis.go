package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *Config) (*redis.Client, error) {
	// Membuat Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0, // menggunakan logical database Redis nomor 0, secara default Redis menyediakan beberapa logical database.
	})

	// Memeriksa koneksi Redis
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect Redis: %w", err)
	}

	return client, nil
}
