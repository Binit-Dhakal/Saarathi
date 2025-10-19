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
	PickUp   [2]float64
	DropOff  [2]float64
	Price    int32
	Distance float64
}

type EventSend struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type OfferRequestDriver struct {
	OfferID   string     `json:"offerID"`
	TripID    string     `json:"tripID"`
	PickUp    [2]float64 `json:"pickUp"`
	DropOff   [2]float64 `json:"dropOff"`
	Price     int32      `json:"price"`
	Distance  float64    `json:"distance"`
	ExpiresAt time.Time  `json:"expiresAt"`
}

type OfferResponseDriver struct {
	TripID  string `json:"tripID"`
	OfferID string `json:"offerID"`
}
