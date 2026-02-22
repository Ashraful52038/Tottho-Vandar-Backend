package usecase

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type UserUsecase interface {
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	Update(ctx context.Context, id uint, req *domain.UpdateUserRequest) (*domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) Update(ctx context.Context, id uint, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}

	err = u.userRepo.Update(ctx, user)
	return user, err
}
