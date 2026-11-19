package dto

type CreateFacilityRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdateFacilityRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}
