package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	commonhandler "roomly-api/common/handler"
	"roomly-api/domain/booking/dto"
	"roomly-api/domain/booking/entity"
	"roomly-api/domain/booking/usecase"
)

type BookingHandlerStruct struct {
	bookingUsecase *usecase.BookingUsecase
}

func NewBookingHandler(bookingUsecase *usecase.BookingUsecase) *BookingHandlerStruct {
	return &BookingHandlerStruct{
		bookingUsecase: bookingUsecase,
	}
}

// CreateBooking membuat booking baru
func (h *BookingHandlerStruct) CreateBooking(ctx *gin.Context) {
	var request dto.CreateBookingRequest

	// 1. Binding request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		commonhandler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid booking request",
			"INVALID_REQUEST",
			nil,
		)
		return
	}

	userID := ctx.GetInt64("user_id")

	// 2. Mapping DTO ke Entity
	booking := entity.Booking{
		// UserID:    request.UserID,
		UserID:    userID, // Ambil dari Gin Context
		RoomID:    request.RoomID,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
	}

	// 3. Panggil Usecase
	result, err := h.bookingUsecase.CreateBooking(ctx.Request.Context(), booking)
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	// 4. Mapping Entity ke DTO
	response := dto.BookingResponse{
		ID:        result.ID,
		UserID:    result.UserID,
		RoomID:    result.RoomID,
		StartTime: result.StartTime,
		EndTime:   result.EndTime,
		Status:    result.Status,
	}

	// 5. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusCreated,
		"Booking created successfully",
		response,
	)
}

// GetBookingByID mengambil detail booking berdasarkan ID
func (h *BookingHandlerStruct) GetBookingByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		commonhandler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid booking ID",
			"INVALID_BOOKING_ID",
			nil,
		)
		return
	}

	userID := ctx.GetInt64("user_id")
	role := ctx.GetString("role")

	booking, err := h.bookingUsecase.GetBookingByID(ctx.Request.Context(), id, userID, role)
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	response := dto.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		RoomID:    booking.RoomID,
		StartTime: booking.StartTime,
		EndTime:   booking.EndTime,
		Status:    booking.Status,
	}

	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Booking retrieved successfully",
		response,
	)
}

// GetBookingList mengambil seluruh daftar booking
func (h *BookingHandlerStruct) GetBookingList(ctx *gin.Context) {
	userID := ctx.GetInt64("user_id")
	role := ctx.GetString("role")

	bookings, err := h.bookingUsecase.GetBookingList(ctx.Request.Context(), userID, role)
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	responses := make([]dto.BookingResponse, 0, len(bookings))

	for _, booking := range bookings {
		responses = append(responses, dto.BookingResponse{
			ID:        booking.ID,
			UserID:    booking.UserID,
			RoomID:    booking.RoomID,
			StartTime: booking.StartTime,
			EndTime:   booking.EndTime,
			Status:    booking.Status,
		})
	}

	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Bookings retrieved successfully",
		responses,
	)
}
