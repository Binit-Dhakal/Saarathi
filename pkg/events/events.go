package events

import "time"

type Event interface {
	EventName() string
}

type TripEventCreated struct {
	RideID   string     `json:"tripId"`
	Distance float64    `json:"distance"`
	Price    int        `json:"price"`
	PickUp   [2]float64 `json:"pickUp"`
	DropOff  [2]float64 `json:"dropOff"`
	CarType  string     `json:"carType"`
}

func (TripEventCreated) EventName() string {
	return "trip.created"
}

func init() {
	RegisterEvent(EventTripCreated, func() Event { return &TripEventCreated{} })
}

type TripOfferRequest struct {
	TripID    string     `json:"tripId"`
	DriverID  string     `json:"driverId"`
	PickUp    [2]float64 `json:"pickUp"`
	DropOff   [2]float64 `json:"dropOff"`
	CarType   string     `json:"carType"`
	ExpiresAt time.Time  `json:"expiresAt"`
}

func (TripOfferRequest) EventName() string {
	return "trip.offer"
}

func init() {
	RegisterEvent(EventOfferRequest, func() Event { return &TripOfferRequest{} })
}

type TripOfferResponse struct {
	TripID   string    `json:"tripId"`
	DriverID string    `json:"driverId"`
	Result   string    `json:"result"`
	AtUnix   time.Time `json:"atUnix"`
}

func (TripOfferResponse) EventName() string {
	return "trip.offerResponse"
}

func init() {
	RegisterEvent(EventOfferResponse, func() Event { return &TripOfferResponse{} })
}
