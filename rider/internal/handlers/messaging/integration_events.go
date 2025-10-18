package messaging

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type integrationHandlers[T ddd.Event] struct {
}

func NewIntegrationEventHandlers() ddd.EventHandler[ddd.Event] {
	return integrationHandlers[ddd.Event]{}
}

func RegisterIntegrationHandlers(subscriber am.EventSubscriber, handlers ddd.EventHandler[ddd.Event]) error {
	evtMsgHandler := am.MessageHandlerFunc[am.IncomingEventMessage](func(ctx context.Context, eventMsg am.IncomingEventMessage) error {
		return handlers.HandleEvent(ctx, eventMsg)
	})

	err := subscriber.Subscribe("trips.realtime.>", evtMsgHandler)
	if err != nil {
		return err
	}

	return nil
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	var err error
	switch event.EventName() {
	case tripspb.TripAssignedEvent:
		err = h.onTripConfirmed(ctx, event)
	}

	return err
}

func (h integrationHandlers[T]) onTripConfirmed(ctx context.Context, event T) error {

}
