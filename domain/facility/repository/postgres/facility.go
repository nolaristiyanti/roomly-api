package postgres

import (
	"context"

	"gorm.io/gorm"

	"roomly-api/domain/facility/entity"
	"roomly-api/domain/facility/model"
	"roomly-api/domain/facility/repository"
)

type FacilityRepositoryStruct struct {
	db *gorm.DB
}

func NewFacilityRepository(db *gorm.DB) repository.FacilityRepository {
	return &FacilityRepositoryStruct{
		db: db,
	}
}

// FindByID mengambil satu facility berdasarkan ID
func (r *FacilityRepositoryStruct) FindByID(ctx context.Context, id int64) (*entity.Facility, error) {
	var facilityModel model.Facility

	err := r.db.WithContext(ctx).First(&facilityModel, id).Error
	if err != nil {
		return nil, err
	}

	facilityEntity := &entity.Facility{
		ID:          facilityModel.ID,
		Name:        facilityModel.Name,
		Description: facilityModel.Description,
	}

	return facilityEntity, nil
}

// FindList mengambil seluruh data facility
func (r *FacilityRepositoryStruct) FindList(ctx context.Context) ([]entity.Facility, error) {
	var facilityModels []model.Facility

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&facilityModels).Error

	if err != nil {
		return nil, err
	}

	facilityEntities := make([]entity.Facility, 0, len(facilityModels))

	for _, facilityModel := range facilityModels {
		facilityEntity := entity.Facility{
			ID:          facilityModel.ID,
			Name:        facilityModel.Name,
			Description: facilityModel.Description,
		}

		facilityEntities = append(facilityEntities, facilityEntity)
	}

	return facilityEntities, nil
}
