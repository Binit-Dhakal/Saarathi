package domain

import "time"

const (
	AcceptOfferIntent  = "offer.intent.accept"
	RejectOfferIntent  = "offer.intent.reject"
	TimeoutOfferIntent = "offer.intent.timeout"
)

type AcceptOffer struct {
	OfferID  string
	DriverID string
	TripID   string
	Ts       time.Time
}

func (AcceptOffer) Key() string { return AcceptOfferIntent }

type RejectOffer struct {
	OfferID  string
	DriverID string
	TripID   string
	Ts       time.Time
}

func (RejectOffer) Key() string { return RejectOfferIntent }

type TimeoutOffer struct {
	OfferID  string
	DriverID string
	TripID   string
	Ts       time.Time
}

func (TimeoutOffer) Key() string { return TimeoutOfferIntent }

