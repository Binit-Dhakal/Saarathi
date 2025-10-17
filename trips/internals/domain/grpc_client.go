package domain

import "context"

type UsersClient interface {
	GetRiderDetails(ctx context.Context, riderID string) (*RiderDetail, error)
	GetDriverDetails(ctx context.Context, driverID string) (*DriverDetail, error)
}
