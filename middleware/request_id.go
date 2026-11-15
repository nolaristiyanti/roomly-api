package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Generate unique request ID
		requestID := uuid.NewString()

		// Simpan request ID ke Gin context
		ctx.Set("RequestID", requestID)

		// Tambahkan request ID ke response header
		ctx.Writer.Header().Set("X-Request-ID", requestID)

		// Lanjutkan ke handler berikutnya
		ctx.Next()
	}
}
