package offerspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	RideMatchingRequestedEvent = "offers.rms.request_matching"
	TripOfferRequestedEvent    = "offers.drivers.requested"
	TripOfferAcceptedEvent     = "offers.trips.accepted"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&RideMatchingRequested{}); err != nil {
		return err
	}

	if err = serde.Register(&TripOfferRequested{}); err != nil {
		return err
	}

	if err = serde.Register(&TripOfferAccepted{}); err != nil {
		return err
	}

	return nil
}

func (*RideMatchingRequested) Key() string { return RideMatchingRequestedEvent }
func (*TripOfferRequested) Key() string    { return TripOfferRequestedEvent }
func (*TripOfferAccepted) Key() string     { return TripOfferAcceptedEvent }
