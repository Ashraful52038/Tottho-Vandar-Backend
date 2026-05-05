package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/email"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt/hash"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type AuthUsecase interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	VerifyEmail(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error
	ResendVerificationEmail(ctx context.Context, email string) error
}

type authUsecase struct {
	userRepo     repository.UserRepository
	jwtService   *jwt.JWTService
	emailService *email.EmailService
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	jwtService *jwt.JWTService,
	emailService *email.EmailService,
) AuthUsecase {
	return &authUsecase{
		userRepo:     userRepo,
		jwtService:   jwtService,
		emailService: emailService,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user exists
	existingUser, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Generate verification token
	verificationToken := generateToken()

	// Create user
	user := &domain.User{
		Name:              req.Name,
		Email:             req.Email,
		Password:          hashedPassword,
		Verified:          false,
		VerificationToken: &verificationToken,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Send verification email
	go func() {
		if err := u.emailService.SendVerificationEmail(user.Email, verificationToken); err != nil {
			log.Printf("Register: failed to send verification email to %s: %v", user.Email, err)
		}
	}()

	// Generate JWT
	token, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user by email
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !hash.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT
	token, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (u *authUsecase) VerifyEmail(ctx context.Context, token string) error {
	user, err := u.userRepo.FindByVerificationToken(ctx, token)
	if err != nil || user == nil {
		return errors.New("invalid or expired verification token")
	}

	err = u.userRepo.VerifyEmail(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// ForgotPassword
func (u *authUsecase) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil
	}

	resetToken := generateToken()

	expiryTime := time.Now().Add(1 * time.Hour)
	err = u.userRepo.SetResetToken(ctx, user.ID, resetToken, expiryTime)
	if err != nil {
		return err
	}

	go u.emailService.SendResetPasswordEmail(user.Email, resetToken)

	return nil
}

func (u *authUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := u.userRepo.FindByResetToken(ctx, token)
	if err != nil || user == nil {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.ResetToken = nil
	user.ResetTokenExpiry = nil

	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *authUsecase) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	// Get user
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Verify current password
	if !hash.CheckPasswordHash(currentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password
	user.Password = string(hashedPassword)
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (u *authUsecase) ResendVerificationEmail(ctx context.Context, email string) error {
	// ইউজার খুঁজুন
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	// যদি ইতিমধ্যে ভেরিফাই করা হয়ে থাকে
	if user.Verified {
		return errors.New("email already verified")
	}
	// নতুন টোকেন জেনারেট করুন
	newToken := generateToken()
	user.VerificationToken = &newToken
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}
	// ইমেইল পাঠান (এরর লগ করবেন কিন্তু ফিরিয়ে দেবেন না)
	go func() {
		if err := u.emailService.SendVerificationEmail(user.Email, newToken); err != nil {
			// লগ করার জন্য log প্যাকেজ ইমপোর্ট করতে হবে
			log.Printf("ResendVerificationEmail: failed to send to %s: %v", user.Email, err)
		}
	}()
	return nil
}
