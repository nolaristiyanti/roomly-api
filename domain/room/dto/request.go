package dto

type CreateRoomRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Capacity    int     `json:"capacity" binding:"required,gt=0"`
	Location    string  `json:"location" binding:"required"`
	FacilityIDs []int64 `json:"facility_ids" binding:"omitempty,dive,gt=0"` // validasi dilewati jika field tidak diisi, validasi setiap elemen di dalam slice, nilainya harus lebih besar dari nol
}

type UpdateRoomRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Capacity    int     `json:"capacity" binding:"required,gt=0"`
	Location    string  `json:"location" binding:"required"`
	IsActive    *bool   `json:"is_active" binding:"required"`
	FacilityIDs []int64 `json:"facility_ids" binding:"omitempty,dive,gt=0"`
}
