package driverspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	DriverAggregateChannel = "saarathi.drivers.events"

	// command to specific instance
	CommandChannel   = "saarathi.drivers.command.%s"
	TripOfferCommand = "driversapi.trip.offer"

	ReplyChannel        = "saarahi.drivers.reply"
	DriverAcceptedReply = "driversapi.driver.accepted"
	RMSTimeoutReply     = "driversapi.rms.timeout"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripOfferRequest{}); err != nil {
		return err
	}

	if err = serde.Register(&DriverAccepted{}); err != nil {
		return err
	}

	if err = serde.Register(&RMSTimeout{}); err != nil {
		return err
	}

	return nil
}

func (*TripOfferRequest) Key() string { return TripOfferCommand }
func (*DriverAccepted) Key() string   { return DriverAcceptedReply }
func (*RMSTimeout) Key() string       { return RMSTimeoutReply }
