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

	ReplyChannel       = "saarahi.drivers.reply"
	OfferAcceptedReply = "driversapi.offer.accepted"
	OfferRejectedReply = "driversapi.offer.rejected"
	OfferTimedoutReply = "driversapi.offer.timedout"
	OfferAckReply      = "driversapi.offer.ack"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripOfferRequest{}); err != nil {
		return err
	}

	if err = serde.Register(&OfferAccepted{}); err != nil {
		return err
	}

	if err = serde.Register(&OfferRejected{}); err != nil {
		return err
	}

	if err = serde.Register(&OfferTimedout{}); err != nil {
		return err
	}

	if err = serde.Register(&OfferAck{}); err != nil {
		return err
	}

	return nil
}

func (*TripOfferRequest) Key() string { return TripOfferCommand }
func (*OfferAccepted) Key() string    { return OfferAcceptedReply }
func (*OfferRejected) Key() string    { return OfferRejectedReply }
func (*OfferTimedout) Key() string    { return OfferTimedoutReply }
func (*OfferAck) Key() string         { return OfferAckReply }
