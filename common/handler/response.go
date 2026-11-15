package handler

import (
	stderrors "errors"
	"net/http"

	apperrors "roomly-api/common/errors"

	"github.com/gin-gonic/gin"
)

// SuccessResponse mengirimkan response ketika request berhasil
func SuccessResponse(ctx *gin.Context, status int, message string, data any) {
	ctx.JSON(status, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse mengirimkan response ketika request gagal
func ErrorResponse(ctx *gin.Context, status int, message string, code string, details any) {
	ctx.JSON(status, gin.H{
		"success": false,
		"message": message,
		"error": gin.H{
			"code":    code,
			"details": details,
		},
		"data": nil,
	})
}

// HandleError menangani error dan mengirimkan response yang sesuai
func HandleError(ctx *gin.Context, err error) {
	// package custom milik aplikasi
	var apiErr *apperrors.APIError

	// cek apakah error memiliki tipe error cocok (APIError)
	// Error custom
	// err := apperrors.New(404, "ROOM_NOT_FOUND", "Room not found", nil)
	// handler.HandleError(ctx, err)
	if stderrors.As(err, &apiErr) {
		ErrorResponse(
			ctx,
			apiErr.Status,
			apiErr.Message,
			apiErr.Code,
			apiErr.Details,
		)
		return
	}

	// Menangani error yang tidak dikenali
	// Error biasa
	// err := fmt.Errorf("database connection failed")
	// handler.HandleError(ctx, err)
	ErrorResponse(
		ctx,
		http.StatusInternalServerError,
		"Internal server error",
		"INTERNAL_SERVER_ERROR",
		nil,
	)
}
