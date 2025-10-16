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

type OfferStatus string

const (
	OfferPending   OfferStatus = "pending"
	OfferDelivered OfferStatus = "delivered"
	OfferAccepted  OfferStatus = "accepted"
	OfferRejected  OfferStatus = "rejected"
	OfferTimedOut  OfferStatus = "timed_out"
	OfferCancelled OfferStatus = "cancelled"
)

type Offer struct {
	ddd.Aggregate
	TripID   string
	SagaID   string
	DriverID string
	Price    int32
	Distance float64

	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    OfferStatus
}

func NewOffer(tripID string, sagaID string, driverID string, price int32, distance float64) Offer {
	ttl := 15 * time.Second
	return Offer{
		Aggregate: ddd.NewAggregate(uuid.NewString(), "drivers.Offer"),
		TripID:    tripID,
		DriverID:  driverID,
		SagaID:    sagaID,
		Price:     price,
		Distance:  distance,

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(ttl),
		Status:    OfferPending,
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
