package application

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/users/internal/domain"
)

type UserDetailService interface {
	GetRiderDetails(ctx context.Context, riderID string) (*domain.UserDetail, error)
	GetDriverDetails(ctx context.Context, driverID string) (*domain.DriverDetail, error)
	GetBulkDriverMetadata(ctx context.Context, driverIDs []string) ([]domain.DriverVehicleMetadata, error)
}

var _ UserDetailService = (*userDetailService)(nil)

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

func (u *userDetailService) GetBulkDriverMetadata(ctx context.Context, driverIDs []string) ([]domain.DriverVehicleMetadata, error) {
	return u.userRepo.BulkSearchMeta(ctx, driverIDs)
}
