package dto

type RoomResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    int     `json:"capacity"`
	Location    string  `json:"location"`
	IsActive    bool    `json:"is_active"`
}
