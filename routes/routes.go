package routes

import (
	"net/http"
	"roomly-api/common/handler"
	"roomly-api/domain/booking"
	"roomly-api/domain/facility"
	"roomly-api/domain/room"
	"roomly-api/domain/user"

	"github.com/gin-gonic/gin"
)

// Mendaftarkan endpoint
func RegisterRoutes(router *gin.Engine, userModule *user.Module, roomModule *room.Module, facilityModule *facility.Module, bookingModule *booking.Module) {
	// Endpoint untuk mengecek kesehatan API
	router.GET("/health", func(ctx *gin.Context) {
		// Mengirim JSON response
		handler.SuccessResponse(
			ctx,
			http.StatusOK,
			"Roomly API is running",
			nil,
		)
	})

	// User routes
	api := router.Group("/api/v1")

	// Hubungkan ke method User Handler.
	api.GET("/users/:id", userModule.Handler.GetUserByID)
	// api.GET("/users/:email", userModule.Handler.GetUserByEmail)

	// Room routes
	// Hubungkan ke method Room Handler.
	api.GET("/rooms/:id", roomModule.Handler.GetRoomByID)
	api.GET("/rooms", roomModule.Handler.GetRoomList)

	// Facility routes
	// Hubungkan ke method Facility Handler.
	api.GET("/facilities/:id", facilityModule.Handler.GetFacilityByID)
	api.GET("/facilities", facilityModule.Handler.GetFacilityList)

	// Booking routes
	// Hubungkan ke method Booking Handler.
	api.POST("/bookings", bookingModule.Handler.CreateBooking)
	api.GET("/bookings/:id", bookingModule.Handler.GetBookingByID)
	api.GET("/bookings", bookingModule.Handler.GetBookingList)
}
