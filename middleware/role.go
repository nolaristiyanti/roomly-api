package middleware

import (
	"net/http"

	responsehandler "roomly-api/common/handler"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Mengambil role dari Context (Authentication Middleware sebelumnya)
		role := ctx.GetString("role")

		if role == "" {
			responsehandler.ErrorResponse(
				ctx,
				http.StatusForbidden,
				"Access denied",
				"FORBIDDEN",
				nil,
			)
			ctx.Abort()
			return
		}

		allowed := false

		// Memeriksa role yang diperbolehkan
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			responsehandler.ErrorResponse(
				ctx,
				http.StatusForbidden,
				"You do not have permission to access this resource",
				"FORBIDDEN",
				nil,
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
