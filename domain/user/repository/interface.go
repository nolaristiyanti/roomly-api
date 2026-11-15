package repository

import (
	"context"

	"roomly-api/domain/user/entity"
)

// Interface ini adalah kontrak yang menyatakan operasi apa saja yang harus tersedia.
// Artinya hanya mendefinisikan kemampuan yang harus disediakan oleh implementasi Repository.
type UserRepository interface {
	// Context dapat meneruskan informasi pembatalan dan batas waktu operasi dari request sampai database.
	// Contohnya, ketika client membatalkan request, operasi database yang menggunakan context tersebut dapat ikut dibatalkan.

	// Repository menggunakan Database Model untuk berkomunikasi dengan PostgreSQL, tetapi hasil yang diberikan kepada Usecase adalah Entity.
	// Usecase tidak perlu mengetahui keberadaan GORM maupun struktur tabel database.
	FindByID(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}
