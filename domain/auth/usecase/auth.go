package usecase

import (
	"context"
	"errors"
	"net/http"
	"time"

	apierror "roomly-api/common/errors"
	authentity "roomly-api/domain/auth/entity"
	authrepository "roomly-api/domain/auth/repository"
	"roomly-api/domain/auth/service"
	"roomly-api/domain/user/entity"
	userrepository "roomly-api/domain/user/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecaseStruct struct {
	userRepository    userrepository.UserRepository
	jwtService        *service.JWTServiceStruct
	sessionRepository authrepository.SessionRepository
}

func NewAuthUsecase(
	userRepository userrepository.UserRepository,
	jwtService *service.JWTServiceStruct,
	sessionRepository authrepository.SessionRepository,
) *AuthUsecaseStruct {
	return &AuthUsecaseStruct{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		jwtService:        jwtService,
	}
}

// Register menangani proses pendaftaran user baru.
func (u *AuthUsecaseStruct) Register(ctx context.Context, name string, email string, password string) (*entity.User, error) {
	// Check duplicate email
	_, err := u.userRepository.FindByEmail(ctx, email)

	if err == nil {
		return nil, apierror.New(
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
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	// Buat entity user
	newUser := entity.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	// Simpan user melalui repository
	createdUser, err := u.userRepository.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// Register menangani proses pendaftaran user baru.
func (u *AuthUsecaseStruct) Login(ctx context.Context, email string, password string) (*authentity.LoginResult, error) {
	// Find user by email
	user, err := u.userRepository.FindByEmail(ctx, email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.New(
			401,
			"INVALID_CREDENTIALS",
			"Email or password is incorrect",
			nil,
		)
	}

	if err != nil {
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return nil, apierror.New(
			401,
			"INVALID_CREDENTIALS",
			"Email or password is incorrect",
			nil,
		)
	}

	// Generate JWT
	accessToken, jti, expiresAt, err := u.jwtService.GenerateToken(
		user.ID,
		user.Role,
	)

	if err != nil {
		return nil, err
	}

	// Create Redis Session
	session := authentity.Session{
		JTI:       jti,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	err = u.sessionRepository.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	// Return login result
	return &authentity.LoginResult{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
	}, nil
}

func (u *AuthUsecaseStruct) Logout(ctx context.Context, jti string) error {
	err := u.sessionRepository.Delete(ctx, jti)
	if err != nil {
		return apierror.New(
			http.StatusServiceUnavailable,
			"SESSION_SERVICE_UNAVAILABLE",
			"Failed to terminate session",
			nil,
		)
	}

	return nil
}
