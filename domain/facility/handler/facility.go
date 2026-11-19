package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"roomly-api/common/handler"
	commonhandler "roomly-api/common/handler"
	"roomly-api/domain/facility/dto"
	"roomly-api/domain/facility/entity"
	"roomly-api/domain/facility/usecase"
)

type FacilityHandlerStruct struct {
	facilityUsecase *usecase.FacilityUsecase
}

func NewFacilityHandler(facilityUsecase *usecase.FacilityUsecase) *FacilityHandlerStruct {
	return &FacilityHandlerStruct{
		facilityUsecase: facilityUsecase,
	}
}

// GetFacilityByID mengambil detail facility berdasarkan ID
func (h *FacilityHandlerStruct) GetFacilityByID(ctx *gin.Context) {
	// 1. Ambil ID dari URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		commonhandler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid facility ID",
			"INVALID_FACILITY_ID",
			nil,
		)
		return
	}

	// 2. Panggil Usecase
	facility, err := h.facilityUsecase.GetFacilityByID(ctx.Request.Context(), id)
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	// 3. Mapping Entity ke DTO
	response := dto.FacilityResponse{
		ID:          facility.ID,
		Name:        facility.Name,
		Description: facility.Description,
	}

	// 4. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Facility retrieved successfully",
		response,
	)
}

// GetFacilityList mengambil seluruh daftar facility
func (h *FacilityHandlerStruct) GetFacilityList(ctx *gin.Context) {
	// 1. Panggil Usecase
	facilities, err := h.facilityUsecase.GetFacilityList(ctx.Request.Context())
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	// 2. Mapping Entity ke DTO
	responses := make([]dto.FacilityResponse, 0, len(facilities))

	for _, facility := range facilities {
		response := dto.FacilityResponse{
			ID:          facility.ID,
			Name:        facility.Name,
			Description: facility.Description,
		}

		responses = append(responses, response)
	}

	// 3. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Facilities retrieved successfully",
		responses,
	)
}

// CreateFacility membuat facility baru
func (h *FacilityHandlerStruct) CreateFacility(ctx *gin.Context) {
	// 1. Binding dan validasi request
	var request dto.CreateFacilityRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		commonhandler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid request",
			"INVALID_REQUEST",
			err.Error(),
		)
		return
	}

	// 2. Mapping Request DTO menjadi Entity
	facility := entity.Facility{
		Name:        request.Name,
		Description: request.Description,
	}

	// 3. Panggil Usecase
	result, err := h.facilityUsecase.CreateFacility(
		ctx.Request.Context(),
		facility,
	)
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	// 4. Mapping Entity menjadi Response DTO
	response := dto.FacilityResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}

	// 5. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusCreated,
		"Facility created successfully",
		response,
	)
}

func (h *FacilityHandlerStruct) UpdateFacility(ctx *gin.Context) {
	// 1. Ambil ID dari URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid facility ID",
			"INVALID_FACILITY_ID",
			nil,
		)
		return
	}

	// 2. Validasi request body
	var request dto.UpdateFacilityRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid request body",
			"INVALID_REQUEST",
			nil,
		)
		return
	}

	// 3. Mapping DTO ke Entity
	facility := entity.Facility{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
	}

	// 4. Panggil Usecase
	result, err := h.facilityUsecase.UpdateFacility(ctx.Request.Context(), facility)
	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 5. Mapping Entity ke Response DTO
	response := dto.FacilityResponse{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
	}

	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Facility updated successfully",
		response,
	)
}

func (h *FacilityHandlerStruct) DeleteFacility(ctx *gin.Context) {
	// 1. Ambil ID dari URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		handler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid facility ID",
			"INVALID_FACILITY_ID",
			nil,
		)
		return
	}

	// 2. Panggil Usecase
	err = h.facilityUsecase.DeleteFacility(ctx.Request.Context(), id)
	if err != nil {
		handler.HandleError(ctx, err)
		return
	}

	// 3. Kirim response sukses
	handler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Facility deleted successfully",
		nil,
	)
}
