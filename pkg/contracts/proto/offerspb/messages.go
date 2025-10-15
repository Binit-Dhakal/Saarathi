package offerspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	RideMatchingRequestedEvent = "offers.request_matching"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&RideMatchingRequested{}); err != nil {
		return err
	}

	return nil
}

func (*RideMatchingRequested) Key() string { return RideMatchingRequestedEvent }
