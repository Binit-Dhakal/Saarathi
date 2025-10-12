package tripspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	// Events
	TripAggregateChannel = "saarathi.trips.events"
	TripCreatedEvent     = "trips.TripCreated"
	TripConfirmedEvent   = "trips.TripConfirmed"

	// Commands
	CommandChannel      = "saarthi.trips.commands"
	AcceptDriverCommand = "tripsapi.driver.accept"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripCreated{}); err != nil {
		return err
	}

	if err = serde.Register(&TripConfirmed{}); err != nil {
		return err
	}

	if err = serde.Register(&AcceptDriver{}); err != nil {
		return err
	}

	return nil
}

func (*TripCreated) Key() string   { return TripCreatedEvent }
func (*TripConfirmed) Key() string { return TripConfirmedEvent }
func (*AcceptDriver) Key() string  { return AcceptDriverCommand }
