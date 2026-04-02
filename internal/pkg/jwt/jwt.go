package jwt

import (
	"errors"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey   string
	tokenExpiry time.Duration
}

func NewJWTService(secretKey string, tokenExpiry time.Duration) *JWTService {
	return &JWTService{
		secretKey:   secretKey,
		tokenExpiry: tokenExpiry,
	}
}

func (s *JWTService) GenerateToken(user *domain.User) (string, error) {
	claims := &domain.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
