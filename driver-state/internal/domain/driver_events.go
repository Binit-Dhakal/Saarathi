package domain

import "time"

const (
	DriverConnectedEvent      = "driver.connected"
	DriverDisconnectedEvent   = "driver.disconnected"
	DriverOfferPublishedEvent = "driver.offerPublished"
	DriverOfferRespondedEvent = "driver.offerResponded"
	DriverOfferTimedOutEvent  = "driver.offerTimedOut"
)

type DriverConnected struct {
	DriverID     string
	ConnectionID string
	Ts           time.Time
}

func (DriverConnected) Key() string { return DriverConnectedEvent }

type DriverDisconnected struct {
	DriverID     string
	ConnectionID string
	Ts           time.Time
}

func (DriverDisconnected) Key() string { return DriverDisconnectedEvent }

type DriverOfferPublished struct {
	OfferID     string
	TripID      string
	DriverID    string
	ExpiresAt   time.Time
	Correlation string
	Ts          time.Time
}

func (DriverOfferPublished) Key() string { return DriverOfferPublishedEvent }

type DriverOfferResponded struct {
	TripID   string
	DriverID string
	Status   OfferStatus
	Ts       time.Time
}

func (DriverOfferResponded) Key() string { return DriverOfferRespondedEvent }

type DriverOfferTimeout struct {
	TripID string
	Status OfferStatus
	Ts     time.Time
}

func (DriverOfferTimeout) Key() string { return DriverOfferTimedOutEvent }
