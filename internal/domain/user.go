package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Name              string     `json:"name" gorm:"not null"`
	Email             string     `json:"email" gorm:"uniqueIndex;not null"`
	Password          string     `json:"-" gorm:"not null"` // "-" excludes from JSON
	Verified          bool       `json:"verified" gorm:"default:false"`
	VerificationToken *string    `json:"-"`
	ResetToken        *string    `json:"-"`
	ResetTokenExpiry  *time.Time `json:"-"`
	Avatar            *string    `json:"avatar"`
	Bio               *string    `json:"bio"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type JWTClaims struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty"`
	Avatar *string `json:"avatar,omitempty"`
	Bio    *string `json:"bio,omitempty"`
}

// UserResponse - response struct for user (excludes sensitive data)
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse - converts User to UserResponse (hides password)
func (u *User) ToResponse() *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Verified:  u.Verified,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UserStats struct {
	Posts     int64 `json:"posts"`
	Comments  int64 `json:"comments"`
	Likes     int64 `json:"likes"`
	Followers int64 `json:"followers"`
	Following int64 `json:"following"`
}

type UserProfileResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"createdAt"`
	Stats     UserStats `json:"stats"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=6"`
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
}
