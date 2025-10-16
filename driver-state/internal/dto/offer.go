package dto

import "time"

type OfferResponse struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type OfferRequestedDTO struct {
	TripID   string
	SagaID   string
	DriverID string
	Price    int32
	Distance float64
}

type EventSend struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type OfferRequestDriver struct {
	OfferID   string    `json:"offerID"`
	TripID    string    `json:"tripId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type OfferResponseDriver struct {
	TripID  string `json:"tripID"`
	OfferID string `json:"offerID"`
}
