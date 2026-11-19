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

func (r *UserRepositoryStruct) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	// Mapping Entity ke Database Model
	userModel := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}

	// Simpan user ke PostgreSQL
	err := r.db.WithContext(ctx).Create(&userModel).Error
	if err != nil {
		return nil, err
	}

	// Mapping Database Model ke Entity User
	userEntity := &entity.User{
		ID:       userModel.ID,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.Password,
		Role:     userModel.Role,
	}

	return userEntity, nil
}

func (r *UserRepositoryStruct) FindList(ctx context.Context, page int, limit int, search string) ([]entity.User, int64, error) {
	var userModels []model.User
	var total int64

	// Menghitung total seluruh user
	countQuery := r.db.WithContext(ctx).Model(&model.User{})

	if search != "" {
		keyword := "%" + search + "%"
		countQuery = countQuery.Where(
			"name ILIKE ? OR email ILIKE ?",
			keyword,
			keyword,
		)
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// SELECT COUNT(*)
	// FROM users
	// WHERE name ILIKE '%nola%'
	//    OR email ILIKE '%nola%';

	// Mengambil data user
	dataQuery := r.db.WithContext(ctx).Model(&model.User{})
	if search != "" {
		keyword := "%" + search + "%"
		dataQuery = dataQuery.Where(
			"name ILIKE ? OR email ILIKE ?",
			keyword,
			keyword,
		)
	}

	// offset pagination
	offset := (page - 1) * limit

	err := dataQuery.
		Select("id", "name", "email", "role").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&userModels).Error

	// SELECT id, name, email, role
	// FROM users
	// WHERE name ILIKE '%nola%'
	//    OR email ILIKE '%nola%'
	// ORDER BY id ASC
	// LIMIT 10
	// OFFSET 0;

	if err != nil {
		return nil, 0, err
	}

	// 4. Mapping Model ke Entity
	users := make([]entity.User, 0, len(userModels))

	for _, userModel := range userModels {
		user := entity.User{
			ID:    userModel.ID,
			Name:  userModel.Name,
			Email: userModel.Email,
			Role:  userModel.Role,
		}

		users = append(users, user)
	}

	return users, total, nil
}
