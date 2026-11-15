package usecase

import (
	"context"
	"errors"
	"net/http"

	"gorm.io/gorm"

	apperrors "roomly-api/common/errors"
	"roomly-api/domain/facility/entity"
	"roomly-api/domain/facility/repository"
)

type FacilityUsecase struct {
	facilityRepository repository.FacilityRepository
}

func NewFacilityUsecase(facilityRepository repository.FacilityRepository) *FacilityUsecase {
	return &FacilityUsecase{
		facilityRepository: facilityRepository,
	}
}

// GetFacilityByID mengambil detail facility berdasarkan ID
func (u *FacilityUsecase) GetFacilityByID(ctx context.Context, id int64) (*entity.Facility, error) {
	facility, err := u.facilityRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(
				http.StatusNotFound,
				"FACILITY_NOT_FOUND",
				"Facility not found",
				nil,
			)
		}

		return nil, err
	}

	return facility, nil
}

// GetFacilityList mengambil seluruh daftar facility
func (u *FacilityUsecase) GetFacilityList(ctx context.Context) ([]entity.Facility, error) {
	return u.facilityRepository.FindList(ctx)
}
