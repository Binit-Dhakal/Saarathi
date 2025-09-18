package registry

import "sync"

type Registrable interface {
	Key() string
}

type Serializer func(v any) ([]byte, error)
type Deserializer func(d []byte, v any) error

type Registry interface {
	Serialize(key string, v any) ([]byte, error)
	Deserialize(key string, data []byte, options ...BuildOption) (any, error)
	Build(key string, options ...BuildOption) (any, error)
	register(key string, fn func() any, s Serializer, d Deserializer, options []BuildOption) error
}

type registered struct {
	serializer   Serializer
	deserializer Deserializer
	factory      func() any
	options      []BuildOption
}

type registry struct {
	registered map[string]registered
	mu         sync.RWMutex
}

var _ Registry = (*registry)(nil)

func NewRegistry() *registry {
	return &registry{
		registered: map[string]registered{},
	}
}

func (r *registry) Serialize(key string, v any) ([]byte, error) {
	reg, exists := r.registered[key]
	if !exists {
		return nil, UnregisteredKey(key)
	}

	return reg.serializer(v)
}

func (r *registry) Deserialize(key string, data []byte, options ...BuildOption) (any, error) {
	v, err := r.Build(key, options...)

	if err != nil {
		return nil, err
	}

	err = r.registered[key].deserializer(data, v)

	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *registry) Build(key string, options ...BuildOption) (any, error) {
	reg, exists := r.registered[key]
	if !exists {
		return nil, UnregisteredKey(key)
	}

	v := reg.factory()
	uos := append(r.registered[key].options, options...)

	for _, option := range uos {
		err := option(v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (r *registry) register(key string, fn func() any, s Serializer, d Deserializer, o []BuildOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registered[key]; exists {
		return AlreadyRegisteredKey(key)
	}

	r.registered[key] = registered{
		factory:      fn,
		serializer:   s,
		deserializer: d,
		options:      o,
	}

	return nil
}
