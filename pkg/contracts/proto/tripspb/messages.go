package tripspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	TripAggregateChannel = "saarathi.trips.events"

	TripRequestedEvent = "trips.requested"
	TripAssignedEvent  = "trips.assigned"

	// Commands
	RejectTripCommand = "tripsapi.trips.reject"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripRequested{}); err != nil {
		return err
	}

	if err = serde.Register(&TripAssigned{}); err != nil {
		return err
	}
	if err = serde.Register(&RejectTrip{}); err != nil {
		return err
	}

	return nil
}

func (*TripRequested) Key() string { return TripRequestedEvent }
func (*TripAssigned) Key() string  { return TripAssignedEvent }
func (*RejectTrip) Key() string    { return RejectTripCommand }
