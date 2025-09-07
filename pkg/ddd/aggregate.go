package ddd

const (
	AggregateNameKey    = "aggregate-name"
	AggregateIDKey      = "aggregate-id"
	AggregateVersionKey = "aggregate-version"
)

type AggregateEvent interface {
	Event
	AggregateName() string
	AggregateID() string
	AggregateVersion() string
}

type aggregate struct {
	Entity
	events []AggregateEvent
}

type Eventer interface {
	Events() []AggregateEvent
	AddEvent(string, EventPayload, ...EventOption)
	ClearEvents()
}

type Aggregate interface {
	Entity
	AggregateName() string
	Eventer
}

type aggregateEvent struct {
	event
}

var _ Aggregate = (*aggregate)(nil)

func NewAggregate(id, name string) *aggregate {
	return &aggregate{
		Entity: NewEntity(id, name),
		events: make([]AggregateEvent, 0),
	}
}

func (a *aggregate) AggregateName() string    { return a.EntityName() }
func (a *aggregate) Events() []AggregateEvent { return a.events }
func (a *aggregate) ClearEvents()             { a.events = []AggregateEvent{} }

func (a *aggregate) AddEvent(name string, payload EventPayload, options ...EventOption) {
	options = append(options,
		Metadata{
			AggregateIDKey:   a.ID(),
			AggregateNameKey: a.EntityName(),
		},
	)
	a.events = append(a.events, aggregateEvent{
		event: newEvent(name, payload, options...),
	})
}

func (e aggregateEvent) AggregateName() string { return e.Metadata().Get(AggregateNameKey).(string) }
func (e aggregateEvent) AggregateID() string   { return e.Metadata().Get(AggregateIDKey).(string) }
func (e aggregateEvent) AggregateVersion() string {
	return e.Metadata().Get(AggregateVersionKey).(string)
}
