package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// konfigurasi yang dibutuhkan untuk membuat JWT
type JWTServiceStruct struct {
	secretKey []byte        // digunakan untuk menandatangani JWT
	expiresIn time.Duration // durasi masa berlaku JWT
}

// informasi atau claims yang akan dimasukkan ke dalam JWT
type JWTClaims struct {
	// Custom Claims
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`

	// Registered Claims : claim yang sudah disediakan library
	jwt.RegisteredClaims
}

// constructor untuk membuat instance JWTServiceStruct
func NewJWTService(secretKey string, expiresIn time.Duration) *JWTServiceStruct {
	return &JWTServiceStruct{
		secretKey: []byte(secretKey),
		expiresIn: expiresIn,
	}
}

// Membuat dan menandatangani JWT
func (s *JWTServiceStruct) GenerateToken(userID int64, role string) (string, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.expiresIn) // waktu kedaluwarsa
	jti := uuid.NewString()           // identitas session

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // buat token JWT

	signedToken, err := token.SignedString(s.secretKey) // JWT yang sudah ditandatangani
	if err != nil {
		return "", "", time.Time{}, err
	}

	return signedToken, jti, expiresAt, nil
}

// Memverifikasi JWT dan membaca claims
func (s *JWTServiceStruct) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims( // Membaca dan memvalidasi JWT
		tokenString,
		&JWTClaims{}, // Menentukan struktur claims yang digunakan
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 { // Memeriksa algoritma signature (HS256)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return s.secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Pastikan token valid dan claims sesuai tipe JWTClaims.
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
