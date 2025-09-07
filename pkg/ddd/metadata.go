package ddd

import "maps"

type Metadata map[string]any

func (m Metadata) Set(key string, value any) {
	m[key] = value
}

func (m Metadata) Get(key string) any {
	return m[key]
}

func (m Metadata) configureEvent(event *event) {
	maps.Copy(event.metadata, m)
}
