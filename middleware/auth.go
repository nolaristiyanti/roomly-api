package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	responsehandler "roomly-api/common/handler"
	authrepository "roomly-api/domain/auth/repository"
	"roomly-api/domain/auth/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	jwtService *service.JWTServiceStruct,
	sessionRepository authrepository.SessionRepository,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 1. Ambil Authorization Header.
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			responsehandler.ErrorResponse(
				ctx,
				http.StatusUnauthorized,
				"Authorization token is required",
				"UNAUTHORIZED",
				nil,
			)
			ctx.Abort()
			return
		}

		// 2. Pastikan format menggunakan Bearer Token.
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {

			responsehandler.ErrorResponse(
				ctx,
				http.StatusUnauthorized,
				"Invalid authorization format",
				"UNAUTHORIZED",
				nil,
			)
			ctx.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		// 3. Validasi JWT.
		claims, err := jwtService.ValidateToken(tokenString)

		if err != nil {
			responsehandler.ErrorResponse(
				ctx,
				http.StatusUnauthorized,
				"Invalid or expired token",
				"UNAUTHORIZED",
				nil,
			)
			ctx.Abort()
			return
		}

		// 4. Pastikan claims memiliki informasi penting.
		if claims.ID == "" || claims.UserID <= 0 || claims.Role == "" {
			responsehandler.ErrorResponse(
				ctx,
				http.StatusUnauthorized,
				"Invalid token claims",
				"UNAUTHORIZED",
				nil,
			)
			ctx.Abort()
			return
		}

		// 5. Cari session berdasarkan JTI di Redis.
		session, err := sessionRepository.FindByID(
			ctx.Request.Context(),
			claims.ID,
		)

		if err != nil {
			if errors.Is(err, authrepository.ErrSessionNotFound) {
				responsehandler.ErrorResponse(
					ctx,
					http.StatusUnauthorized,
					"Session not found or expired",
					"SESSION_INVALID",
					nil,
				)
			} else {
				responsehandler.ErrorResponse(
					ctx,
					http.StatusServiceUnavailable,
					"Session service unavailable",
					"SERVICE_UNAVAILABLE",
					nil,
				)
			}

			ctx.Abort()
			return
		}

		// 6. Pastikan session sesuai dengan JWT.
		if session.UserID != claims.UserID ||
			session.JTI != claims.ID ||
			!session.ExpiresAt.After(time.Now()) {

			responsehandler.ErrorResponse(
				ctx,
				http.StatusUnauthorized,
				"Session is invalid or expired",
				"SESSION_INVALID",
				nil,
			)
			ctx.Abort()
			return
		}

		// 7. Simpan informasi user ke Gin Context.
		ctx.Set("user_id", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Set("jti", claims.ID)

		// 8. Lanjutkan ke handler berikutnya.
		ctx.Next()
	}
}
