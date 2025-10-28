package domain

import "context"

type DriverAvailabilityRepo interface {
	CheckPresence(ctx context.Context, driverID string) (string, error)
}

type DriverLockRepository interface {
	TryAcquireLock(ctx context.Context, driverID string, tripID string) (bool, error)
	ReleaseLock(ctx context.Context, driverID string, tripID string) error
	AcceptTrip(ctx context.Context, driverID, tripID string) error
	MarkAvailable(ctx context.Context, driverID string) error
}
