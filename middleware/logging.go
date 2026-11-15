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

		// Mengambil Request ID dari Gin context
		requestID := ctx.GetString("RequestID")

		// Mencatat informasi HTTP request
		slog.Info(
			"HTTP request",
			"request_id", requestID,
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"status", ctx.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
		)
	}
}
