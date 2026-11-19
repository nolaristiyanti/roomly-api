package room

import (
	"gorm.io/gorm"

	facilitypostgres "roomly-api/domain/facility/repository/postgres"
	roomhandler "roomly-api/domain/room/handler"
	roompostgres "roomly-api/domain/room/repository/postgres"
	roomusecase "roomly-api/domain/room/usecase"
)

type Module struct {
	Handler *roomhandler.RoomHandlerStruct
}

func New(db *gorm.DB) *Module {
	// 1. Buat repository
	roomRepository := roompostgres.NewRoomRepository(db)

	// Buat Facility Repository
	facilityRepository := facilitypostgres.NewFacilityRepository(db)

	// 2. Buat usecase dengan dependency repository
	roomUsecase := roomusecase.NewRoomUsecase(roomRepository, facilityRepository)

	// 3. Buat handler dengan dependency usecase
	roomHandler := roomhandler.NewRoomHandler(roomUsecase)

	// 4. Kembalikan module
	return &Module{
		Handler: roomHandler,
	}
}
