package tripspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	TripCreatedEvent = "tripsapi.TripCreated"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripCreated{}); err != nil {
		return err
	}

	return nil
}

func (*TripCreated) Key() string { return TripCreatedEvent }
