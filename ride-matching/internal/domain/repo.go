package domain

type RideMatchingRepository interface {
	FindNearestDriver(lat, lon float64) []string
	BulkSearchDriverMeta(driverID []string) ([]DriverVehicleMetadata, error)
	IsDriverAvailable(driverID string) bool
}
