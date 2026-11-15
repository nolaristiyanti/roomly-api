package repository

import (
	"context"

	"roomly-api/domain/facility/entity"
)

type FacilityRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.Facility, error)
	FindList(ctx context.Context) ([]entity.Facility, error)
}
