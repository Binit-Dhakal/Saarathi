package am

type RawMessage interface {
	Message
	Data() []byte
}

type IncomingRawMessage interface {
	IncomingMessage
	Data() []byte
}

type RawMessageStream = MessageStream[RawMessage, IncomingRawMessage]
type RawMessageHandler = MessageHandler[IncomingRawMessage]
type RawMessagePublisher = MessagePublisher[RawMessage]
type RawMessageSubscriber = MessageSubscriber[IncomingRawMessage]

type rawMessage struct {
	id   string
	name string
	data []byte
}

var _ RawMessage = (*rawMessage)(nil)

func (r rawMessage) ID() string          { return r.id }
func (r rawMessage) Data() []byte        { return r.data }
func (r rawMessage) MessageName() string { return r.name }
