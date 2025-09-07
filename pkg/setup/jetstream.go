package setup

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

func SetupJetStream(streamName string, nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{fmt.Sprintf("%s.>", streamName)},
	})

	return js, err
}
