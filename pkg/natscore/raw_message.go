package natscore

import "github.com/Binit-Dhakal/Saarathi/pkg/am"

type rawMessage struct {
	id   string
	name string
	data []byte
	ack  bool
}

var _ am.IncomingMessage = (*rawMessage)(nil)

func (m *rawMessage) ID() string          { return m.id }
func (m *rawMessage) MessageName() string { return m.name }
func (m *rawMessage) Data() []byte        { return m.data }

func (m *rawMessage) Ack() error {
	if m.ack {
		return nil
	}

	m.ack = true
	return nil
}

func (m *rawMessage) NAck() error {
	if m.ack {
		return nil
	}

	m.ack = true
	return nil
}

func (m *rawMessage) Extend() error {
	return nil
}

func (m *rawMessage) Kill() error {
	if m.ack {
		return nil
	}

	m.ack = true
	return nil
}

