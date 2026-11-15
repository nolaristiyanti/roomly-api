package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {
	// Handler untuk menentukan bagaimana log ditampilkan.
	var handler slog.Handler

	// dalam format JSON
	handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Mengatur logger sebagai default logger aplikasi
	slog.SetDefault(slog.New(handler))
}
