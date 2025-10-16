package domain

import "time"

const (
	AcceptOfferEvent  = "offer.response.accept"
	RejectOfferEvent  = "offer.response.reject"
	TimeoutOfferEvent = "offer.response.timeout"
)

type AcceptOffer struct {
	OfferID  string
	DriverID string
	TripID   string
	Ts       time.Time
}

func (AcceptOffer) Key() string { return AcceptOfferEvent }

type RejectOffer struct {
	OfferID  string
	DriverID string
	TripID   string
	Ts       time.Time
}

func (RejectOffer) Key() string { return RejectOfferEvent }

type TimeoutOffer struct {
	OfferID  string
	DriverID string
	TripID   string
	Ts       time.Time
}

func (TimeoutOffer) Key() string { return TimeoutOfferEvent }
