package usecase

import (
	"context"
	"errors"
	"net/http"

	apperrors "roomly-api/common/errors"
	"roomly-api/domain/user/entity"
	"roomly-api/domain/user/repository"

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

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	// Usecase memanggil Repository melalui interface
	user, err := u.userRepository.FindByEmail(ctx, email)

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
