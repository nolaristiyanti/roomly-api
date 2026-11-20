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
// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a room booking for the authenticated user.
// @Tags Booking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateBookingRequest true "Booking data"
// @Success 201 {object} map[string]interface{} "Booking created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid booking request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Room not found"
// @Failure 409 {object} map[string]interface{} "Booking schedule conflict"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /bookings [post]
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
	// response := dto.BookingResponse{
	// 	ID:        result.ID,
	// 	UserID:    result.UserID,
	// 	RoomID:    result.RoomID,
	// 	StartTime: result.StartTime,
	// 	EndTime:   result.EndTime,
	// 	Status:    result.Status,
	// }
	response := toBookingResponse(result)

	// 5. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusCreated,
		"Booking created successfully",
		response,
	)
}

// GetBookingByID mengambil detail booking berdasarkan ID
// GetBookingByID godoc
// @Summary Get booking detail
// @Description Retrieve booking details by ID. Access depends on the authenticated user's role and ownership.
// @Tags Booking
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]interface{} "Booking retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid booking ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Booking not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /bookings/{id} [get]
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

	// response := dto.BookingResponse{
	// 	ID:        booking.ID,
	// 	UserID:    booking.UserID,
	// 	RoomID:    booking.RoomID,
	// 	StartTime: booking.StartTime,
	// 	EndTime:   booking.EndTime,
	// 	Status:    booking.Status,
	// }
	response := toBookingResponse(booking)

	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Booking retrieved successfully",
		response,
	)
}

// GetBookingList mengambil seluruh daftar booking
// GetBookingList godoc
// @Summary Get booking list
// @Description Retrieve bookings belonging to the authenticated user. Admins can retrieve all bookings.
// @Tags Booking
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Bookings retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /bookings [get]
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
		// responses = append(responses, dto.BookingResponse{
		// 	ID:        booking.ID,
		// 	UserID:    booking.UserID,
		// 	RoomID:    booking.RoomID,
		// 	StartTime: booking.StartTime,
		// 	EndTime:   booking.EndTime,
		// 	Status:    booking.Status,
		// })
		responses = append(responses, toBookingResponse(&booking))
	}

	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Bookings retrieved successfully",
		responses,
	)
}

// CancelBooking membatalkan booking
// CancelBooking godoc
// @Summary Cancel booking
// @Description Cancel an existing booking. Users can cancel their own bookings, while admins can cancel any booking.
// @Tags Booking
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]interface{} "Booking cancelled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid booking ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Booking not found"
// @Failure 409 {object} map[string]interface{} "Booking cannot be cancelled"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /bookings/{id} [delete]
func (h *BookingHandlerStruct) CancelBooking(ctx *gin.Context) {
	// 1. Ambil dan validasi booking ID
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

	// 2. Ambil identitas user dari JWT context
	userID := ctx.GetInt64("user_id")
	role := ctx.GetString("role")

	// 3. Panggil Usecase
	err = h.bookingUsecase.CancelBooking(
		ctx.Request.Context(),
		id,
		userID,
		role,
	)

	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	// 4. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Booking cancelled successfully",
		nil,
	)
}

func toBookingResponse(item *entity.Booking) dto.BookingResponse {
	return dto.BookingResponse{
		ID: item.ID,
		// UserID:    item.UserID,
		// RoomID:    item.RoomID,
		StartTime: item.StartTime,
		EndTime:   item.EndTime,
		Status:    item.Status,

		User: dto.BookingUserResponse{
			ID:    item.User.ID,
			Name:  item.User.Name,
			Email: item.User.Email,
		},

		Room: dto.BookingRoomResponse{
			ID:       item.Room.ID,
			Name:     item.Room.Name,
			Capacity: item.Room.Capacity,
			Location: item.Room.Location,
		},
	}
}
