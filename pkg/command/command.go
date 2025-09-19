package command

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type Command interface {
	ddd.IDer
	CommandName() string
	Payload() []byte
}

type Reply interface {
	CorrelationID() string
	Payload() []byte
}
