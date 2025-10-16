package domain

import "context"

type DriverAvailabilityRepo interface {
	CheckPresence(ctx context.Context, driverID string) (string, error)
	TryAcquireLock(ctx context.Context, driverID string, tripID string) (bool, error)
	ReleaseLock(ctx context.Context, driverID string, tripID string) error
}
