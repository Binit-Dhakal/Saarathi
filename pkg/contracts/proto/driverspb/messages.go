package driverspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	DriverAggregateChannel = "saarathi.drivers.events"

	DriverOfferEventsChannel = "drivers.events.offers.%s"
	OfferAcceptedEvent       = "driversapi.offer.accepted"
	OfferRejectedEvent       = "driversapi.offer.rejected"
	OfferTimedoutEvent       = "driversapi.offer.timedout"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&OfferAccepted{}); err != nil {
		return err
	}

	if err = serde.Register(&OfferRejected{}); err != nil {
		return err
	}

	if err = serde.Register(&OfferTimedout{}); err != nil {
		return err
	}

	return nil
}

func (*OfferAccepted) Key() string { return OfferAcceptedEvent }
func (*OfferRejected) Key() string { return OfferRejectedEvent }
func (*OfferTimedout) Key() string { return OfferTimedoutEvent }
