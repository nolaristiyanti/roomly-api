package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	commonhandler "roomly-api/common/handler"
	facilitydto "roomly-api/domain/facility/dto"
	facilityentity "roomly-api/domain/facility/entity"
	"roomly-api/domain/room/dto"
	"roomly-api/domain/room/usecase"
)

type RoomHandlerStruct struct {
	roomUsecase *usecase.RoomUsecase
}

func NewRoomHandler(roomUsecase *usecase.RoomUsecase) *RoomHandlerStruct {
	return &RoomHandlerStruct{
		roomUsecase: roomUsecase,
	}
}

// GetRoomByID mengambil detail room berdasarkan ID
func (h *RoomHandlerStruct) GetRoomByID(ctx *gin.Context) {
	// 1. Ambil ID dari URL
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		commonhandler.ErrorResponse(
			ctx,
			http.StatusBadRequest,
			"Invalid room ID",
			"INVALID_ROOM_ID",
			nil,
		)
		return
	}

	// 2. Panggil Usecase
	room, err := h.roomUsecase.GetRoomByID(ctx.Request.Context(), id)
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	// 3. Mapping Entity ke DTO
	response := dto.RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
		Location:    room.Location,
		IsActive:    room.IsActive,
		Facilities:  mapFacilitiesToResponse(room.Facilities),
	}

	// 4. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Room retrieved successfully",
		response,
	)
}

// GetRoomList mengambil seluruh daftar room
func (h *RoomHandlerStruct) GetRoomList(ctx *gin.Context) {
	// 1. Panggil Usecase
	rooms, err := h.roomUsecase.GetRoomList(ctx.Request.Context())
	if err != nil {
		commonhandler.HandleError(ctx, err)
		return
	}

	// 2. Mapping Entity ke DTO
	var responses []dto.RoomResponse

	for _, room := range rooms {
		response := dto.RoomResponse{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Capacity:    room.Capacity,
			Location:    room.Location,
			IsActive:    room.IsActive,
			Facilities:  mapFacilitiesToResponse(room.Facilities),
		}

		responses = append(responses, response)
	}

	// 3. Kirim response
	commonhandler.SuccessResponse(
		ctx,
		http.StatusOK,
		"Rooms retrieved successfully",
		responses,
	)
}

func mapFacilitiesToResponse(facilities []facilityentity.Facility) []facilitydto.FacilityResponse {

	responses := make([]facilitydto.FacilityResponse, 0, len(facilities))

	for _, facility := range facilities {
		responses = append(responses, facilitydto.FacilityResponse{
			ID:          facility.ID,
			Name:        facility.Name,
			Description: facility.Description,
		})
	}

	return responses
}
