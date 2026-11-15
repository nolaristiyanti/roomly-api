package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"roomly-api/common/handler"
	"roomly-api/domain/user/dto"
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

// func (h *UserHandlerStruct) GetUserByEmail(ctx *gin.Context) {
// 	// 1. Ambil email dari URL
// 	email := ctx.Param("email")

// 	// 2. Panggil Usecase
// 	user, err := h.userUsecase.GetUserByEmail(ctx.Request.Context(), email)
// 	if err != nil {
// 		handler.HandleError(ctx, err)
// 		return
// 	}

// 	// 3. Mapping Entity ke DTO
// 	response := dto.UserResponse{
// 		ID:    user.ID,
// 		Name:  user.Name,
// 		Email: user.Email,
// 		Role:  user.Role,
// 	}

// 	// 4. Kirim response
// 	handler.SuccessResponse(
// 		ctx,
// 		http.StatusOK,
// 		"User retrieved successfully",
// 		response,
// 	)
// }
