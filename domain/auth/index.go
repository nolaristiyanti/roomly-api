package auth

import (
	"gorm.io/gorm"

	authhandler "roomly-api/domain/auth/handler"
	authrepository "roomly-api/domain/auth/repository"
	authservice "roomly-api/domain/auth/service"
	authusecase "roomly-api/domain/auth/usecase"
	userpostgres "roomly-api/domain/user/repository/postgres"
)

type Module struct {
	Handler           *authhandler.AuthHandlerStruct
	SessionRepository authrepository.SessionRepository
}

func New(db *gorm.DB, jwtService *authservice.JWTServiceStruct, sessionRepository authrepository.SessionRepository) *Module {
	// 1. Gunakan User Repository yang sudah tersedia
	userRepository := userpostgres.NewUserRepository(db)

	// 2. Buat Auth Usecase
	authUsecase := authusecase.NewAuthUsecase(userRepository, jwtService, sessionRepository)

	// 3. Buat Auth Handler
	authHandler := authhandler.NewAuthHandler(authUsecase)

	// 4. Kembalikan Module
	return &Module{
		Handler: authHandler,
	}
}
