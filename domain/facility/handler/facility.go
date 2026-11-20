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
// GetFacilityByID godoc
// @Summary Get facility detail
// @Description Retrieve facility details by ID.
// @Tags Facility
// @Produce json
// @Security BearerAuth
// @Param id path int true "Facility ID"
// @Success 200 {object} map[string]interface{} "Facility retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid facility ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Facility not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /facilities/{id} [get]
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
// GetFacilityList godoc
// @Summary Get facility list
// @Description Retrieve all available facilities.
// @Tags Facility
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Facilities retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /facilities [get]
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
// CreateFacility godoc
// @Summary Create a new facility
// @Description Create a new facility. Admin access required.
// @Tags Admin - Facilities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateFacilityRequest true "Facility data"
// @Success 201 {object} map[string]interface{} "Facility created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/facilities [post]
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

// UpdateFacility godoc
// @Summary Update facility
// @Description Update facility information. Admin access required.
// @Tags Admin - Facilities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Facility ID"
// @Param request body dto.UpdateFacilityRequest true "Updated facility data"
// @Success 200 {object} map[string]interface{} "Facility updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid facility ID or request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Facility not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/facilities/{id} [put]
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

// DeleteFacility godoc
// @Summary Delete facility
// @Description Delete a facility by ID. Admin access required.
// @Tags Admin - Facilities
// @Produce json
// @Security BearerAuth
// @Param id path int true "Facility ID"
// @Success 200 {object} map[string]interface{} "Facility deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid facility ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Facility not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/facilities/{id} [delete]
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
