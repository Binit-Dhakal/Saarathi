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

type TripOffer struct {
	OfferID    string     `json:"offerId"`
	TripID     string     `json:"tripId"`
	DriverID   string     `json:"driverId"`
	PickUp     [2]float64 `json:"pickUp"`
	DropOff    [2]float64 `json:"dropOff"`
	CarType    string     `json:"carType"`
	ExpiresAt  time.Time  `json:"expiresAt"`
	InstanceID string     `json:"instanceID"`
}

func (TripOffer) EventName() string {
	return "trip.offer"
}

func init() {
	RegisterEvent(EventOfferCreated, func() Event { return &TripOffer{} })
}

type TripOfferResponse struct {
	OfferID  string    `json:"offerId"`
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
