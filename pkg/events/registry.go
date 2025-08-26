package events

import (
	"encoding/json"
	"fmt"
)

var eventRegistry = map[string]func() Event{}

func RegisterEvent(name string, factory func() Event) {
	eventRegistry[name] = factory
}

func DecodeEvent(eventName string, data []byte) (Event, error) {
	factory, ok := eventRegistry[eventName]
	if !ok {
		return nil, fmt.Errorf("no event registered for %s", eventName)
	}

	evt := factory()
	if err := json.Unmarshal(data, evt); err != nil {
		return nil, err
	}

	return evt, nil
}
