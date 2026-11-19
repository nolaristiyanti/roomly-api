package usecase

import (
	"context"
	"errors"
	"net/http"

	apperrors "roomly-api/common/errors"
	"roomly-api/domain/user/entity"
	"roomly-api/domain/user/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserUsecase membutuhkan repository yang didapatkan dari luar melalui constructor.
// Usecase dependency dengan interface Repository
type UserUsecase struct {
	userRepository repository.UserRepository
}

// Constructor
func NewUserUsecase(userRepository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	// Usecase memanggil Repository melalui interface
	user, err := u.userRepository.FindByID(ctx, id)

	if err != nil {
		// mengenali error tertentu meskipun error tersebut dibungkus (wrapped error) oleh error lain
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(
				http.StatusNotFound,
				"USER_NOT_FOUND",
				"User not found",
				nil,
			)
		}

		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// Validasi role
	if user.Role != "user" && user.Role != "admin" {
		return nil, apperrors.New(
			http.StatusBadRequest,
			"INVALID_ROLE",
			"Invalid user role",
			nil,
		)
	}

	// Check duplicate email
	_, err := u.userRepository.FindByEmail(ctx, user.Email)

	if err == nil {
		return nil, apperrors.New(
			http.StatusConflict,
			"EMAIL_ALREADY_EXISTS",
			"Email is already registered",
			nil,
		)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	// Simpan user melalui repository
	createdUser, err := u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *UserUsecase) GetUserList(ctx context.Context, page int, limit int, search string) ([]entity.User, int64, error) {
	// 1. Validasi page
	if page < 1 {
		return nil, 0, apperrors.New(
			http.StatusBadRequest,
			"INVALID_PAGE",
			"Page must be greater than 0",
			nil,
		)
	}

	// 2. Validasi limit
	if limit < 1 || limit > 100 {
		return nil, 0, apperrors.New(
			http.StatusBadRequest,
			"INVALID_LIMIT",
			"Limit must be between 1 and 100",
			nil,
		)
	}

	// 3. Panggil Repository
	users, total, err := u.userRepository.FindList(ctx, page, limit, search)

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	// 1. Pastikan User tersedia
	_, err := u.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// 2. Validasi role
	if user.Role != "user" && user.Role != "admin" {
		return nil, apperrors.New(
			http.StatusBadRequest,
			"INVALID_ROLE",
			"Invalid user role",
			nil,
		)
	}

	// 3. Periksa apakah email sudah digunakan
	existingUser, err := u.userRepository.FindByEmail(ctx, user.Email)

	if err == nil && existingUser.ID != user.ID {
		return nil, apperrors.New(
			http.StatusConflict,
			"EMAIL_ALREADY_EXISTS",
			"Email is already registered",
			nil,
		)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 4. Update User melalui Repository
	updatedUser, err := u.userRepository.Update(ctx, user)

	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, apperrors.New(
				http.StatusConflict,
				"EMAIL_ALREADY_EXISTS",
				"Email is already registered",
				nil,
			)
		}

		return nil, err
	}

	return updatedUser, nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int64) error {
	// 1. Pastikan User tersedia
	_, err := u.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// 2. Hapus User melalui Repository
	err = u.userRepository.Delete(ctx, id)
	if err != nil {
		// 3. Periksa apakah User masih memiliki booking
		if errors.Is(err, repository.ErrUserHasBookings) {
			return apperrors.New(
				http.StatusConflict,
				"USER_HAS_BOOKINGS",
				"User cannot be deleted because they have existing bookings",
				nil,
			)
		}

		return err
	}

	return nil
}

func (u *UserUsecase) UpdateMyProfile(ctx context.Context, userID int64, name string, email string) (*entity.User, error) {
	// 1. Ambil data User yang sedang login
	existingUser, err := u.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Periksa apakah email sudah digunakan User lain
	emailUser, err := u.userRepository.FindByEmail(ctx, email)

	if err == nil && emailUser.ID != userID {
		return nil, apperrors.New(
			http.StatusConflict,
			"EMAIL_ALREADY_EXISTS",
			"Email is already registered",
			nil,
		)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 3. Buat Entity dengan mempertahankan Role lama
	user := entity.User{
		ID:    userID,
		Name:  name,
		Email: email,
		Role:  existingUser.Role,
	}

	// 4. Update melalui Repository
	updatedUser, err := u.userRepository.Update(ctx, user)

	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, apperrors.New(
				http.StatusConflict,
				"EMAIL_ALREADY_EXISTS",
				"Email is already registered",
				nil,
			)
		}

		return nil, err
	}

	return updatedUser, nil
}
