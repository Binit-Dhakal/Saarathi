package domain

import "context"

type RedisRideMatchingRepository interface {
	FindNearestDriver(ctx context.Context, lon, lat, radius float64) []string
	GetRejectedDriver(ctx context.Context, tripID string) ([]string, error)
}

type RedisMetaRepository interface {
	BulkSearchDriverMeta(driverIDs []string) ([]DriverVehicleMetadata, error)
	BulkInsertDriverMeta(metas []DriverVehicleMetadata) error
}

type DriverAvailabilityRepository interface {
	IsDriverFree(driverID string) bool
	DeleteUnavailableDrivers(expiredDrivers []string)
	BulkCheckDriversOnline(driversID []string) ([]string, []string)
}
