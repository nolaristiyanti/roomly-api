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
// Register godoc
// @Summary      	Register new user
// @Description  	Register a new user account.
// @Tags         	Authentication
// @Accept       	json
// @Produce      	json
// @Param 			request body dto.RegisterRequest true "Registration data"
// @Success 		201 {object} map[string]interface{} "Registration successful"
// @Failure 		422 {object} map[string]interface{} "Validation failed"
// @Failure 		409 {object} map[string]interface{} "Email already registered"
// @Failure 		500 {object} map[string]interface{} "Internal server error"
// @Router       	/auth/register [post]
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

// Login godoc
// @Summary      	Login user
// @Description  	Authenticate user and generate JWT access token.
// @Tags         	Authentication
// @Accept       	json
// @Produce      	json
// @Param 			request body dto.LoginRequest true "Login credentials"
// @Success 		200 {object} map[string]interface{} "Login successful"
// @Failure 		401 {object} map[string]interface{} "Invalid credentials"
// @Failure 		422 {object} map[string]interface{} "Validation failed"
// @Failure 		500 {object} map[string]interface{} "Internal server error"
// @Router       	/auth/login [post]
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

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate the current JWT session.
// @Tags         Authentication
// @Produce      json
// @Security     BearerAuth
// @Success 	200 {object} map[string]interface{} "Logout successful"
// @Failure 	401 {object} map[string]interface{} "Invalid or expired session"
// @Failure 	500 {object} map[string]interface{} "Internal server error"
// @Router       /auth/logout [post]
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
