package main

import (
	"log"
	"roomly-api/common/logger"
	"roomly-api/config"
	bookingdomain "roomly-api/domain/booking"
	facilitydomain "roomly-api/domain/facility"
	roomdomain "roomly-api/domain/room"
	userdomain "roomly-api/domain/user"
	"roomly-api/middleware"
	"roomly-api/routes"

	"github.com/gin-gonic/gin"
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

	// Create Gin router
	router := gin.New()

	// Register middleware
	router.Use(middleware.RequestID()) // membuat identitas request
	router.Use(middleware.Logging())   // mencatat informasi request setelah proses selesai
	router.Use(gin.Recovery())         // menangani panic agar server tidak berhenti akibat panic pada satu request

	// Meneruskan Modules ke Routes dan register routes
	routes.RegisterRoutes(router, userModule, roomModule, facilityModule, bookingModule)

	// Menjalankan HTTP server
	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
