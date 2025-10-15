package setup

import (
	"time"

	"github.com/nats-io/nats.go"
)

func SetupJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js, err
}

func SetupStreams(js nats.JetStreamContext, tripStream, sagaStream string) error {
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     tripStream,
		Subjects: []string{"trips.>"},
		Replicas: 3,
	})
	if err != nil {
		return err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     sagaStream,
		Subjects: []string{"trips.requested", "offers.>", "rms.>"},
		MaxAge:   24 * time.Hour,
	})
	return err
}
