package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"roomly-api/common/handler"
	"roomly-api/common/helper"
	"roomly-api/domain/user/dto"
	"roomly-api/domain/user/entity"
	"roomly-api/domain/user/usecase"
)

// Handler membutuhkan/menerima Usecase yang didapatkan dari luar melalui constructor.
type UserHandlerStruct struct {
	userUsecase *usecase.UserUsecase
}

// Constructor
func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandlerStruct {
	return &UserHandlerStruct{
		userUsecase: userUsecase,
	}
}

func (h *UserHandlerStruct) GetUserByID(ctx *gin.Context) {
	// 1. Ambil ID dari URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid user ID",
			"INVALID_USER_ID",
			nil,
		)
		return
	}

	// Ambil user ID dan role dari Auth Middleware
	currentUserID := ctx.GetInt64("user_id")
	currentUserRole := ctx.GetString("role")

	// Authorization: user hanya boleh melihat dirinya sendiri
	if currentUserRole != "admin" && currentUserID != id {
		handler.ErrorResponse(
			ctx,
			http.StatusForbidden,
			"You are not allowed to access this user",
			"FORBIDDEN",
			nil,
		)
		return
	}

	// 2. Panggil Usecase
	user, err := h.userUsecase.GetUserByID(ctx.Request.Context(), id)
	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Mapping Entity ke DTO
	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	// 4. Kirim response
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"User retrieved successfully",
		response,
	)
}

func (h *UserHandlerStruct) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest

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

	// Mapping DTO ke Entity
	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	// 2. Panggil Register Usecase
	createdUser, err := h.userUsecase.CreateUser(ctx.Request.Context(), user)

	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Mapping Entity ke Response DTO
	response := dto.UserResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Role:  createdUser.Role,
	}

	// 4. Kirim response
	handler.SuccessResponse(
		ctx,
		http.StatusCreated,
		"User created successful",
		response,
	)
}

func (h *UserHandlerStruct) GetMyProfile(ctx *gin.Context) {
	// 1. Ambil user ID dari Authentication Middleware
	userID := ctx.GetInt64("user_id")

	// 2. Panggil Usecase
	user, err := h.userUsecase.GetUserByID(
		ctx.Request.Context(),
		userID,
	)
	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Mapping Entity ke DTO Response
	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	// 4. Kirim response
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Profile retrieved successfully",
		response,
	)
}

func (h *UserHandlerStruct) GetUserList(ctx *gin.Context) {
	// 1. Ambil pagination parameter
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid page parameter",
			"INVALID_PAGE",
			nil,
		)
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid limit parameter",
			"INVALID_LIMIT",
			nil,
		)
		return
	}

	search := ctx.Query("search")
	role := ctx.Query("role")
	sortBy := ctx.Query("sort_by")
	sortOrder := ctx.Query("sort_order")

	// 2. Panggil Usecase
	users, total, err := h.userUsecase.GetUserList(ctx.Request.Context(), page, limit, search, role, sortBy, sortOrder)

	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Mapping Entity ke DTO Response
	responses := make([]dto.UserResponse, 0, len(users))

	for _, user := range users {
		response := dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}

		responses = append(responses, response)
	}

	// 4. Menghitung total halaman
	totalPages := int(
		(total + int64(limit) - 1) / int64(limit),
	)

	// 5. Membuat response dengan pagination metadata
	data := gin.H{
		"items": responses,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}

	// 6. Kirim response
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Users retrieved successfully",
		data,
	)
}

func (h *UserHandlerStruct) UpdateUser(ctx *gin.Context) {
	// 1. Ambil ID dari URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid user ID",
			"INVALID_USER_ID",
			nil,
		)
		return
	}

	// 2. Binding dan validasi request
	var req dto.UpdateUserRequest

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

	// 3. Mapping DTO ke Entity
	user := entity.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	// 4. Panggil Usecase
	updatedUser, err := h.userUsecase.UpdateUser(
		ctx.Request.Context(),
		user,
	)

	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 5. Mapping Entity ke Response DTO
	response := dto.UserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Role:  updatedUser.Role,
	}

	// 6. Kirim response
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"User updated successfully",
		response,
	)
}

func (h *UserHandlerStruct) DeleteUser(ctx *gin.Context) {
	// 1. Ambil ID dari URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid user ID",
			"INVALID_USER_ID",
			nil,
		)
		return
	}

	// 2. Panggil Usecase
	err = h.userUsecase.DeleteUser(
		ctx.Request.Context(),
		id,
	)

	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Kirim response sukses
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"User deleted successfully",
		nil,
	)
}

func (h *UserHandlerStruct) UpdateMyProfile(ctx *gin.Context) {
	// 1. Ambil User ID dari Authentication Middleware
	userID := ctx.GetInt64("user_id")

	// 2. Binding dan validasi request
	var req dto.UpdateProfileRequest

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

	// 3. Panggil Usecase
	updatedUser, err := h.userUsecase.UpdateMyProfile(
		ctx.Request.Context(),
		userID,
		req.Name,
		req.Email,
	)

	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 4. Mapping Entity ke Response DTO
	response := dto.UserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Role:  updatedUser.Role,
	}

	// 5. Kirim response
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Profile updated successfully",
		response,
	)
}
