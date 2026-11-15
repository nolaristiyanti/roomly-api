package postgres

import (
	"context"

	"gorm.io/gorm"

	"roomly-api/domain/user/entity"
	"roomly-api/domain/user/model"
	"roomly-api/domain/user/repository"
)

// Repository membutuhkan koneksi GORM untuk menjalankan query. Koneksi tersebut diberikan dari luar melalui constructor.
type UserRepositoryStruct struct {
	db *gorm.DB
}

// Constructor
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryStruct{
		db: db,
	}
}

func (r *UserRepositoryStruct) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	var userModel model.User

	// GORM mengambil data berdasarkan ID
	err := r.db.WithContext(ctx).First(&userModel, id).Error
	if err != nil {
		return nil, err
	}

	// Mapping Model ke Entity User
	userEntity := &entity.User{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.Password,
		Role:     userModel.Role,
	}

	return userEntity, nil
}

func (r *UserRepositoryStruct) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.User

	// GORM mengambil data berdasarkan email
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error

	if err != nil {
		return nil, err
	}

	// Mapping Model ke Entity User
	userEntity := &entity.User{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.Password,
		Role:     userModel.Role,
	}

	return userEntity, nil
}
