package rmspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	// Events
	RMSAggregateChannel = "saarathi.rms.events"
	RMSTripMatched      = "rms.TripMatched"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&TripMatched{}); err != nil {
		return err
	}

	return nil
}

func (r *TripMatched) Key() string { return RMSTripMatched }
