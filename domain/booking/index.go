package booking

import (
	"gorm.io/gorm"

	bookinghandler "roomly-api/domain/booking/handler"
	bookingpostgres "roomly-api/domain/booking/repository/postgres"
	bookingusecase "roomly-api/domain/booking/usecase"
)

type Module struct {
	Handler *bookinghandler.BookingHandlerStruct
}

func New(db *gorm.DB) *Module {
	// 1. Buat repository
	bookingRepository := bookingpostgres.NewBookingRepository(db)

	// 2. Buat usecase
	bookingUsecase := bookingusecase.NewBookingUsecase(bookingRepository)

	// 3. Buat handler
	bookingHandler := bookinghandler.NewBookingHandler(bookingUsecase)

	// 4. Kembalikan module
	return &Module{
		Handler: bookingHandler,
	}
}
