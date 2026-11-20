// @title           Roomly API
// @version         1.0
// @description     Room & Facility Booking Management API
// @termsOfService  http://swagger.io/terms/

// @contact.name    Roomly API Support

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token with the Bearer prefix. Example: Bearer eyJhbGciOiJIUzI1NiIs...
package main

import (
	"log"
	"roomly-api/common/logger"
	"roomly-api/config"
	authdomain "roomly-api/domain/auth"
	authredis "roomly-api/domain/auth/repository/redis"
	authservice "roomly-api/domain/auth/service"
	bookingdomain "roomly-api/domain/booking"
	facilitydomain "roomly-api/domain/facility"
	roomdomain "roomly-api/domain/room"
	userdomain "roomly-api/domain/user"
	"roomly-api/middleware"
	"roomly-api/routes"
	"time"

	"github.com/gin-gonic/gin"

	_ "roomly-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Load application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize application logger
	logger.InitLogger()

	// Connect PostgreSQL
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// menutup koneksi database
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	// Connect Redis
	redisClient, err := config.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	// Koneksi Redis harus ditutup ketika fungsi main() selesai / aplikasi selesai berjalan
	defer redisClient.Close()

	// Initialize Module (Dependency Injection)
	// memberikan koneksi GORM (db) kepada constructor New() ke tiap domain
	userModule := userdomain.New(db)
	roomModule := roomdomain.New(db)
	facilityModule := facilitydomain.New(db)
	bookingModule := bookingdomain.New(db)

	jwtService := authservice.NewJWTService(cfg.JWTSecret, time.Duration(cfg.JWTExpireHours)*time.Hour)
	sessionRepository := authredis.NewSessionRepository(redisClient)
	authModule := authdomain.New(db, jwtService, sessionRepository)

	// Create Gin router
	router := gin.New()

	// Register middleware
	router.Use(middleware.RequestID()) // membuat identitas request
	router.Use(middleware.Logging())   // mencatat informasi request setelah proses selesai
	router.Use(gin.Recovery())         // menangani panic agar server tidak berhenti akibat panic pada satu request

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	// Meneruskan Modules ke Routes dan register routes dan pasang auth middleware ke routes
	routes.RegisterRoutes(router, userModule, roomModule, facilityModule, bookingModule, authModule, sessionRepository, jwtService, redisClient)

	// Menjalankan HTTP server
	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
