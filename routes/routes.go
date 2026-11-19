package routes

import (
	"net/http"
	"roomly-api/common/handler"
	"roomly-api/domain/auth"
	"roomly-api/domain/auth/repository"
	"roomly-api/domain/auth/service"
	"roomly-api/domain/booking"
	"roomly-api/domain/facility"
	"roomly-api/domain/room"
	"roomly-api/domain/user"
	"roomly-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Mendaftarkan endpoint
func RegisterRoutes(router *gin.Engine, userModule *user.Module, roomModule *room.Module, facilityModule *facility.Module, bookingModule *booking.Module, authModule *auth.Module, sessionRepository repository.SessionRepository, jwtService *service.JWTServiceStruct, redisClient *redis.Client) {
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

	api := router.Group("/api/v1")

	// Auth Routes (Public Routes)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authModule.Handler.Register)
		auth.POST("/login", authModule.Handler.Login)
	}

	// AUTHENTICATION MIDDLEWARE
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(
		jwtService,
		sessionRepository,
	))

	// Protected Routes
	{
		// ROLE MIDDLEWARE
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.POST("/users", userModule.Handler.CreateUser)
			admin.GET("/users", userModule.Handler.GetUserList)
		}

		// Auth Routes
		protected.POST("/auth/logout", authModule.Handler.Logout)

		// User Routes
		// Hubungkan ke method User Handler.
		protected.GET("/users/:id", userModule.Handler.GetUserByID)
		protected.GET("/users/profile", userModule.Handler.GetMyProfile)

		// Booking routes
		// Hubungkan ke method Booking Handler.
		protected.POST("/bookings", bookingModule.Handler.CreateBooking)
		protected.GET("/bookings/:id", bookingModule.Handler.GetBookingByID)
		protected.GET("/bookings", bookingModule.Handler.GetBookingList)

		// Room routes
		// Hubungkan ke method Room Handler.
		protected.GET("/rooms/:id", roomModule.Handler.GetRoomByID)
		protected.GET("/rooms", roomModule.Handler.GetRoomList)

		// Facility routes
		// Hubungkan ke method Facility Handler.
		protected.GET("/facilities/:id", facilityModule.Handler.GetFacilityByID)
		protected.GET("/facilities", facilityModule.Handler.GetFacilityList)
	}
}
