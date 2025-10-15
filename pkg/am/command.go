package am

import "github.com/Binit-Dhakal/Saarathi/pkg/ddd"

const (
	CommandHdrPrefix       = "COMMAND_"
	CommandReplyChannelHdr = CommandHdrPrefix + "REPLY_CHANNEL"
)

type Command interface {
	ddd.Command
	Destination() string
}

type command struct {
	ddd.Command
	destination string
}

func NewCommand(name, destination string, payload ddd.CommandPayload) Command {
	return command{
		Command:     ddd.NewCommand(name, payload),
		destination: destination,
	}
}

func (c command) Destination() string {
	return c.destination
}
