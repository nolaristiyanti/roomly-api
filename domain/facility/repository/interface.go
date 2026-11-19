package repository

import (
	"context"
	"errors"

	"roomly-api/domain/facility/entity"
)

var ErrFacilityAlreadyExists = errors.New("facility already exists")

type FacilityRepository interface {
	FindByID(ctx context.Context, id int64) (*entity.Facility, error)
	FindList(ctx context.Context) ([]entity.Facility, error)
	Create(ctx context.Context, facility entity.Facility) (*entity.Facility, error)
	Update(ctx context.Context, facility entity.Facility) (*entity.Facility, error)
	Delete(ctx context.Context, id int64) error
}
