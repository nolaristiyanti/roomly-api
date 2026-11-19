package dto

import "roomly-api/domain/facility/dto"

type RoomResponse struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description *string                `json:"description"`
	Capacity    int                    `json:"capacity"`
	Location    string                 `json:"location"`
	IsActive    bool                   `json:"is_active"`
	Facilities  []dto.FacilityResponse `json:"facilities"`
}
