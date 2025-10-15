package tripspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	// Events
	TripAggregateChannel = "saarathi.trips.events"
	TripRequestedEvent   = "trips.requested"
	TripConfirmedEvent   = "trips.TripConfirmed"

	// Commands
	CommandChannel      = "saarthi.trips.commands"
	AcceptDriverCommand = "tripsapi.driver.accept"
	RejectTripCommand   = "tripsapi.trips.reject"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripRequested{}); err != nil {
		return err
	}

	if err = serde.Register(&TripConfirmed{}); err != nil {
		return err
	}

	if err = serde.Register(&AcceptDriver{}); err != nil {
		return err
	}

	if err = serde.Register(&RejectTrip{}); err != nil {
		return err
	}

	return nil
}

func (*TripRequested) Key() string { return TripRequestedEvent }
func (*TripConfirmed) Key() string { return TripConfirmedEvent }
func (*AcceptDriver) Key() string  { return AcceptDriverCommand }
func (*RejectTrip) Key() string    { return RejectTripCommand }
