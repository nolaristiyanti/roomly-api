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

// CreateFacility membuat facility baru
func (u *FacilityUsecase) CreateFacility(ctx context.Context, facility entity.Facility) (*entity.Facility, error) {
	// Simpan Facility
	result, err := u.facilityRepository.Create(ctx, facility)
	if err != nil {
		if errors.Is(err, repository.ErrFacilityAlreadyExists) {
			return nil, apperrors.New(
				http.StatusConflict,
				"FACILITY_ALREADY_EXISTS",
				"Facility already exists",
				nil,
			)
		}

		return nil, err
	}

	return result, nil
}

func (u *FacilityUsecase) UpdateFacility(ctx context.Context, facility entity.Facility) (*entity.Facility, error) {
	// 1. Pastikan Facility tersedia
	_, err := u.GetFacilityByID(ctx, facility.ID)
	if err != nil {
		return nil, err
	}

	// 2. Update Facility melalui Repository
	result, err := u.facilityRepository.Update(ctx, facility)
	if err != nil {
		if errors.Is(err, repository.ErrFacilityAlreadyExists) {
			return nil, apperrors.New(
				http.StatusConflict,
				"FACILITY_ALREADY_EXISTS",
				"Facility already exists",
				nil,
			)
		}

		return nil, err
	}

	return result, nil
}

func (u *FacilityUsecase) DeleteFacility(ctx context.Context, id int64) error {
	// 1. Pastikan Facility tersedia
	_, err := u.GetFacilityByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Hapus Facility melalui Repository
	err = u.facilityRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
