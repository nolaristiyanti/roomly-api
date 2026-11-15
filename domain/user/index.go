package user

import (
	"gorm.io/gorm"

	userhandler "roomly-api/domain/user/handler"
	userpostgres "roomly-api/domain/user/repository/postgres"
	userusecase "roomly-api/domain/user/usecase"
)

// Module membutuhkan/menerima handler yang didapatkan dari luar melalui constructor.
type Module struct {
	Handler *userhandler.UserHandlerStruct
}

// Constructor yang menerima koneksi database dari luar.
func New(db *gorm.DB) *Module {
	// 1. Buat repository
	// Repository dibuat terlebih dahulu karena dibutuhkan oleh Usecase
	userRepository := userpostgres.NewUserRepository(db)

	// 2. Buat usecase dengan dependency repository
	// Usecase dibuat karena dibutuhkan oleh Handler
	userUsecase := userusecase.NewUserUsecase(userRepository)

	// 3. Buat handler dengan dependency usecase
	// Handler disimpan dalam Module agar nantinya mudah digunakan oleh routes.
	userHandler := userhandler.NewUserHandler(userUsecase)

	// 4. Kembalikan module agar bisa digunakan routes
	return &Module{
		Handler: userHandler,
	}
}
