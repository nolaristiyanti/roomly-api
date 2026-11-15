package facility

import (
	"gorm.io/gorm"

	facilityhandler "roomly-api/domain/facility/handler"
	facilitypostgres "roomly-api/domain/facility/repository/postgres"
	facilityusecase "roomly-api/domain/facility/usecase"
)

type Module struct {
	Handler *facilityhandler.FacilityHandlerStruct
}

func New(db *gorm.DB) *Module {
	// 1. Buat repository
	facilityRepository := facilitypostgres.NewFacilityRepository(db)

	// 2. Buat usecase
	facilityUsecase := facilityusecase.NewFacilityUsecase(facilityRepository)

	// 3. Buat handler
	facilityHandler := facilityhandler.NewFacilityHandler(facilityUsecase)

	// 4. Kembalikan module
	return &Module{
		Handler: facilityHandler,
	}
}
