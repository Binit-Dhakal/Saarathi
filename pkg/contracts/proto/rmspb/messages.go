package rmspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	// Events
	RMSAggregateChannel = "saarathi.rms.events"
	RMSTripMatched      = "rms.TripMatched"

	// Commands
	RMSCommandChannel           = "saarathi.rms.commands"
	FindEligibleDriversCommands = "rmsapi.drivers.search"

	// Reply
	EligibleDriversListReply = "rmsapi.replies.drivers.search"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripMatched{}); err != nil {
		return err
	}

	if err = serde.Register(&FindEligibleDrivers{}); err != nil {
		return err
	}

	if err = serde.Register(&FindEligibleDriversReply{}); err != nil {
		return err
	}

	return nil
}

func (*TripMatched) Key() string              { return RMSTripMatched }
func (*FindEligibleDrivers) Key() string      { return FindEligibleDriversCommands }
func (*FindEligibleDriversReply) Key() string { return EligibleDriversListReply }
