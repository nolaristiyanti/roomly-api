package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Mencatat waktu ketika request mulai diproses
		start := time.Now()

		// Melanjutkan request ke middleware atau handler berikutnya
		ctx.Next()

		// Menghitung durasi pemrosesan request
		duration := time.Since(start)

		// Mengambil informasi request dari Gin context
		requestID := ctx.GetString("RequestID")
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		status := ctx.Writer.Status()
		clientIP := ctx.ClientIP() // Mengambil alamat IP client

		// Log attributes
		attrs := []any{
			"request_id", requestID,
			"method", method,
			"path", path,
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"client_ip", clientIP,
		}

		// Mengambil user_id dari Auth Middleware jika request sudah terautentikasi
		userID := ctx.GetInt64("user_id")

		if userID > 0 {
			attrs = append(attrs, "user_id", userID)
		}

		// Mencatat informasi HTTP request dan menentukan log level berdasarkan HTTP status
		switch {
		case status >= 500:
			slog.Error("HTTP request", attrs...)

		case status >= 400:
			slog.Warn("HTTP request", attrs...)

		default:
			slog.Info("HTTP request", attrs...)
		}
	}
}
