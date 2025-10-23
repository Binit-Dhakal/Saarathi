package rmspb

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/registry"
	serdes "github.com/Binit-Dhakal/Saarathi/pkg/registry/serde"
)

const (
	RMSAggregateChannel       = "saarathi.rms.events"
	RMSCandidatesMatchedEvent = "rms.candidates"
)

func Registration(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	if err = serde.Register(&CandidatesMatched{}); err != nil {
		return err
	}
	return nil
}

func (*CandidatesMatched) Key() string { return RMSCandidatesMatchedEvent }
