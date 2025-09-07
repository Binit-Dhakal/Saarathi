package ddd

type IDer interface {
	ID() string
}

type entity struct {
	id   string
	name string
}

type Entity interface {
	IDer
	EntityName() string
}

var _ Entity = (*entity)(nil)

func NewEntity(id, name string) *entity {
	return &entity{
		id:   id,
		name: name,
	}
}

func (e entity) ID() string {
	return e.id
}

func (e entity) EntityName() string {
	return e.name
}
