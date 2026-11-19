package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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

// Create menyimpan facility baru
func (r *FacilityRepositoryStruct) Create(ctx context.Context, facility entity.Facility) (*entity.Facility, error) {
	// 1. Mapping Entity ke Model
	facilityModel := model.Facility{
		Name:        facility.Name,
		Description: facility.Description,
	}

	// 2. Simpan ke database
	err := r.db.WithContext(ctx).Create(&facilityModel).Error
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, repository.ErrFacilityAlreadyExists
		}

		return nil, err
	}

	// 3. Mapping Model ke Entity
	facilityEntity := &entity.Facility{
		ID:          facilityModel.ID,
		Name:        facilityModel.Name,
		Description: facilityModel.Description,
	}

	return facilityEntity, nil
}

// Update memperbarui data facility
func (r *FacilityRepositoryStruct) Update(ctx context.Context, facility entity.Facility) (*entity.Facility, error) {
	// 1. Cari Facility berdasarkan ID
	var facilityModel model.Facility

	err := r.db.WithContext(ctx).First(&facilityModel, facility.ID).Error

	// SELECT * FROM facilities WHERE id = 1;

	if err != nil {
		return nil, err
	}

	// 2. Mapping Entity ke Model
	facilityModel.Name = facility.Name
	facilityModel.Description = facility.Description

	// 3. Update data Facility
	err = r.db.WithContext(ctx).
		Model(&facilityModel).
		Select("Name", "Description").
		Updates(&facilityModel).Error

	// UPDATE facilities
	// SET
	//     name = 'Smart TV',
	//     description = 'Smart TV untuk presentasi meeting',
	//     updated_at = NOW()
	// WHERE id = 1;

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique rule
			return nil, repository.ErrFacilityAlreadyExists
		}

		return nil, err
	}

	// 4. Mapping Model ke Entity
	facilityEntity := &entity.Facility{
		ID:          facilityModel.ID,
		Name:        facilityModel.Name,
		Description: facilityModel.Description,
	}

	return facilityEntity, nil
}

func (r *FacilityRepositoryStruct) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&model.Facility{}, id).
		Error
}
