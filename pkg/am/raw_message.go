package am

type RawMessage interface {
	Message
	Data() []byte
}

type IncomingRawMessage interface {
	IncomingMessage
	Data() []byte
}

type RawMessageHandler = MessageHandler[IncomingRawMessage]

type rawMessage struct {
	id      string
	name    string
	data    []byte
	replyTo string // optional: for request-reply
}

var _ RawMessage = (*rawMessage)(nil)

func NewRawMessage(id string, name string, data []byte, replyTo string) rawMessage {
	return rawMessage{
		id:      id,
		name:    name,
		data:    data,
		replyTo: replyTo,
	}
}
func (r rawMessage) ID() string          { return r.id }
func (r rawMessage) Data() []byte        { return r.data }
func (r rawMessage) MessageName() string { return r.name }
