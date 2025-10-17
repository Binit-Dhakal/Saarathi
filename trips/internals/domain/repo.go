package domain

import (
	"context"
	"time"
)

// ephemeral save of fare data in redis
type FareRepository interface {
	CreateEphemeralFareEntry(fare *FareQuote) (string, error)
	GetEphemeralFareEntry(id string) (*FareQuote, error)
	DeleteEphemeralFareEntry(id string) error
}

type TripRepository interface {
	SaveRouteDetail(route *Route, riderID string) (string, error)
	SaveFareDetail(fareModel FareRecord) (string, error)
	SaveRideDetail(rideModel TripModel) (string, error)
	AssignDriverToTrip(tripID string, driverID string) (string, error) // riderID return
}

type TripProjectionRepository interface {
	SetTripPayload(ctx context.Context, tripID string, payload map[string]any, expiration time.Duration) error
}

type PresenceGatewayRepository interface {
	GetDriverLocation(ctx context.Context, driverID string) (*Coordinate, error)
}

type TripReadRepository interface {
	GetTripProjectionDetail(ctx context.Context, tripID string) (TripProjectionDetail, error)
}
