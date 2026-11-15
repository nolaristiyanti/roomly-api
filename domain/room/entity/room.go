package entity

type Room struct {
	ID          int64
	Name        string
	Description *string
	Capacity    int
	Location    string
	IsActive    bool
}
