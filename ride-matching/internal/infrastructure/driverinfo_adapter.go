package infrastructure

import (
	"context"
	"errors"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/application"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
)

type DriverInfoAdapter struct {
	svc application.DriverInfoService
}

var (
	_ domain.DriverAvailabilityChecker = (*DriverInfoAdapter)(nil)
	_ domain.DriverMetadataFetcher     = (*DriverInfoAdapter)(nil)
)

func NewDriverInfoAdapter(svc application.DriverInfoService) *DriverInfoAdapter {
	if svc == nil {
		panic(errors.New("DriverInfoService cannot be nil"))
	}
	return &DriverInfoAdapter{
		svc: svc,
	}
}

func (a *DriverInfoAdapter) GetOnlineDrivers(ctx context.Context, driverIDs []string) []string {
	return a.svc.GetOnlineDrivers(ctx, driverIDs)
}

func (a *DriverInfoAdapter) GetBulkMetada(ctx context.Context, driverIDs []string) ([]domain.DriverVehicleMetadata, error) {
	return a.svc.GetDriversMetadata(ctx, driverIDs)
}
