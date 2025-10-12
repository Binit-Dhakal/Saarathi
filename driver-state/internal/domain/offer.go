package domain

import (
	"errors"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/google/uuid"
)

var (
	ErrAlreadyProcessed = errors.New("Offer Already Processed")
)

type Offer struct {
	ddd.Aggregate
	TripID      string
	DriverID    string
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      OfferStatus
	Correlation string
}

func NewOffer(tripID string, driverID string, ttl time.Duration, correlation string) Offer {
	return Offer{
		Aggregate:   ddd.NewAggregate(uuid.NewString(), "drivers.Offer"),
		TripID:      tripID,
		DriverID:    driverID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		ExpiresAt:   time.Now().UTC().Add(ttl),
		Status:      OfferPending,
		Correlation: correlation,
	}
}

func (o *Offer) Accept() (ddd.Event, error) {
	if o.Status != OfferPending {
		return nil, ErrAlreadyProcessed
	}

	timeNow := time.Now().UTC()
	o.Status = OfferAccepted
	o.UpdatedAt = timeNow

	o.AddEvent(DriverOfferRespondedEvent, DriverOfferResponded{
		DriverID: o.DriverID,
		TripID:   o.TripID,
		Status:   OfferAccepted,
		Ts:       time.Now(),
	})

	return ddd.NewEvent(DriverOfferRespondedEvent, o), nil
}

func (o *Offer) Reject() (ddd.Event, error) {
	if o.Status != OfferPending {
		return nil, ErrAlreadyProcessed
	}

	timeNow := time.Now().UTC()
	o.Status = OfferRejected
	o.UpdatedAt = timeNow

	o.AddEvent(DriverOfferRespondedEvent, DriverOfferResponded{
		DriverID: o.DriverID,
		TripID:   o.TripID,
		Status:   OfferRejected,
		Ts:       timeNow,
	})

	return ddd.NewEvent(DriverOfferRespondedEvent, o), nil
}

func (o *Offer) TimeOut() (ddd.Event, error) {
	if o.Status != OfferPending {
		return nil, ErrAlreadyProcessed
	}

	o.Status = OfferTimedOut
	o.UpdatedAt = time.Now().UTC()

	o.AddEvent(DriverOfferTimedOutEvent, DriverOfferTimeout{})

	return ddd.NewEvent(DriverOfferTimedOutEvent, o), nil
}
