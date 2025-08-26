package dto

import "time"

type OfferResponse struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type OfferRequestDriver struct {
	TripID    string     `json:"tripId"`
	PickUp    [2]float64 `json:"pickUp"`
	DropOff   [2]float64 `json:"dropOff"`
	ExpiresAt time.Time  `json:"expiresAt"`
}
