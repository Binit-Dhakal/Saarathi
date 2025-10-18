package domain

import "context"

type TripPayloadRepository interface {
	GetTripFullPayload(ctx context.Context, tripID string) ([]byte, error)
}
