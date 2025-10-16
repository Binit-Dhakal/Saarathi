package domain

import "context"

type TripReadModelRepository interface {
	SaveTrip(ctx context.Context, payload TripReadModelDTO) error
	GetTripDetails(ctx context.Context, tripID string) (TripReadModelDTO, error)
}
