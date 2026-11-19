package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"roomly-api/common/handler"
	"roomly-api/common/helper"
	"roomly-api/domain/auth/dto"
	"roomly-api/domain/auth/usecase"
)

type AuthHandlerStruct struct {
	authUsecase *usecase.AuthUsecaseStruct
}

func NewAuthHandler(authUsecase *usecase.AuthUsecaseStruct) *AuthHandlerStruct {
	return &AuthHandlerStruct{
		authUsecase: authUsecase,
	}
}

// Register menangani HTTP request pendaftaran user.
func (h *AuthHandlerStruct) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	// 1. Binding dan validasi request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.ErrorResponse(
			ctx,
			http.StatusUnprocessableEntity,
			"Validation failed",
			"VALIDATION_ERROR",
			helper.TranslateValidationError(err),
		)
		return
	}

	// 2. Panggil Register Usecase
	user, err := h.authUsecase.Register(
		ctx.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
	)

	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Mapping Entity ke Response DTO
	response := dto.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	// 4. Kirim response
	handler.SuccessResponse(
		ctx,
		http.StatusCreated,
		"Registration successful",
		response,
	)
}

func (h *AuthHandlerStruct) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	// 1. Request validation
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handler.ErrorResponse(
			ctx,
			http.StatusUnprocessableEntity,
			"Validation failed",
			"VALIDATION_ERROR",
			helper.TranslateValidationError(err),
		)
		return
	}

	// 2. Call usecase
	result, err := h.authUsecase.Login(
		ctx.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Response mapping
	response := dto.LoginResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresIn:   result.ExpiresIn,
	}

	// 4. Return response
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Login successful",
		response,
	)
}

func (h *AuthHandlerStruct) Logout(ctx *gin.Context) {
	// Mengambil JTI
	jti := ctx.GetString("jti")

	if jti == "" {
		handler.ErrorResponse(
			ctx,
			http.StatusUnauthorized,
			"Invalid session",
			"SESSION_INVALID",
			nil,
		)
		return
	}

	// Memanggil Usecase
	err := h.authUsecase.Logout(
		ctx.Request.Context(),
		jti,
	)
	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// Mengembalikan response
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Logout successful",
		nil,
	)
}
