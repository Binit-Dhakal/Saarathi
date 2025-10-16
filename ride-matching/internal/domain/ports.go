package domain

import "context"

type DriverAvailabilityChecker interface {
	GetOnlineDrivers(ctx context.Context, driverIDs []string) []string
}

type DriverMetadataFetcher interface {
	GetBulkMetada(ctx context.Context, driverIDs []string) ([]DriverVehicleMetadata, error)
}
