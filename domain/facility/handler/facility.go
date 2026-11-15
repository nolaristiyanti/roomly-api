package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	commonhandler "roomly-api/common/handler"
	"roomly-api/domain/facility/dto"
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
