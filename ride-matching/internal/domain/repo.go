package domain

// test
type RedisRideMatchingRepository interface {
	FindNearestDriver(lat, lon float64) []string
	IsDriverAvailable(driverID string) bool
}

type RedisMetaRepository interface {
	BulkSearchDriverMeta(driverIDs []string) ([]DriverVehicleMetadata, error)
	BulkInsertDriverMeta(metas []DriverVehicleMetadata) error
}

type PGMetaRepository interface {
	BulkSearchMeta(driverIDs []string) ([]DriverVehicleMetadata, error)
}
