package events

import "time"

type TripEventCreated struct {
	RideID   string  `json:"tripId"`
	Distance float64 `json:"distance"`
	Price    int     `json:"price"`
	// PickupLat float64 `json:"pickupLat"`
	// PickupLon float64 `json:"pickupLon"`
	PickUp  [2]float64 `json:"pickUp"`
	DropOff [2]float64 `json:"dropOff"`
	CarType string     `json:"carType"`
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

type TripOfferResponse struct {
	OfferID  string    `json:"offerId"`
	TripID   string    `json:"tripId"`
	DriverID string    `json:"driverId"`
	Result   string    `json:"result"`
	AtUnix   time.Time `json:"atUnix"`
}
