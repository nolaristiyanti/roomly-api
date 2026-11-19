package entity

import "roomly-api/domain/facility/entity"

type Room struct {
	ID          int64
	Name        string
	Description *string
	Capacity    int
	Location    string
	IsActive    bool
	Facilities  []entity.Facility
}
