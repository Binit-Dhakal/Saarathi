package domain

import "context"

type RedisRideMatchingRepository interface {
	FindNearestDriver(ctx context.Context, lon, lat float64) []string
}

type RedisMetaRepository interface {
	BulkSearchDriverMeta(driverIDs []string) ([]DriverVehicleMetadata, error)
	BulkInsertDriverMeta(metas []DriverVehicleMetadata) error
}

type PGMetaRepository interface {
	BulkSearchMeta(driverIDs []string) ([]DriverVehicleMetadata, error)
}

type DriverAvailabilityRepository interface {
	IsDriverFree(driverID string) bool
	DeleteUnavailableDrivers(expiredDrivers []string)
	BulkCheckDriversOnline(driversID []string) ([]string, []string)
}
