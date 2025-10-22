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
	Aggregate ddd.Aggregate `json:"-"`
	TripID    string
	SagaID    string
	DriverID  string
	PickUp    [2]float64
	DropOff   [2]float64
	Price     int32
	Distance  float64

	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    OfferStatus
}

func NewOffer(tripID string, sagaID string, driverID string, pickUp [2]float64, dropOff [2]float64, price int32, distance float64) Offer {
	ttl := 15 * time.Second
	return Offer{
		Aggregate: ddd.NewAggregate(uuid.NewString(), "drivers.Offer"),
		TripID:    tripID,
		DriverID:  driverID,
		PickUp:    pickUp,
		DropOff:   dropOff,
		SagaID:    sagaID,
		Price:     price,
		Distance:  distance,

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(ttl),
		Status:    OfferPending,
	}
}

func (o *Offer) AddEvent(name string, payload ddd.EventPayload, options ...ddd.EventOption) {
	if o.Aggregate != nil {
		o.Aggregate.AddEvent(name, payload, options...)
	}
}

func (o *Offer) Events() []ddd.AggregateEvent {
	if o.Aggregate == nil {
		return nil
	}
	return o.Aggregate.Events()
}

func (o *Offer) ClearEvents() {
	if o.Aggregate != nil {
		o.Aggregate.ClearEvents()
	}
}

func (o *Offer) AggregateName() string {
	if o.Aggregate == nil {
		return "drivers.Offer"
	}
	return o.Aggregate.AggregateName()
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
