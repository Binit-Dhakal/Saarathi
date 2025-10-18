package domain

import "context"

type LocationRepo interface {
	SaveActiveGeoLocation(*DriverLocation) error
	RemoveActiveGeoLocation(string) error
}

type WSRepo interface {
	SaveWSDetail(driverID string) error
	DeleteWSDetail(driverID string) error
}

type OfferRepository interface {
	FindByID(ctx context.Context, id string) (*Offer, error)
	Save(ctx context.Context, offer *Offer) error
}

type TripPayloadRepository interface {
	GetTripFullPayload(ctx context.Context, tripID string) ([]byte, error)
}
