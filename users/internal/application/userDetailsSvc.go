package application

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/users/internal/domain"
)

type UserDetailService interface {
	GetRiderDetails(ctx context.Context, riderID string) (*domain.UserDetail, error)
	GetDriverDetails(ctx context.Context, driverID string) (*domain.DriverDetail, error)
}

type userDetailService struct {
	userRepo domain.UserRepo
}

func NewUserDetailService(userRepo domain.UserRepo) UserDetailService {
	return &userDetailService{
		userRepo: userRepo,
	}
}

func (u *userDetailService) GetRiderDetails(ctx context.Context, riderID string) (*domain.UserDetail, error) {
	return u.userRepo.GetRiderByID(ctx, riderID)
}

func (u *userDetailService) GetDriverDetails(ctx context.Context, driverID string) (*domain.DriverDetail, error) {
	return u.userRepo.GetDriverByID(ctx, driverID)
}
