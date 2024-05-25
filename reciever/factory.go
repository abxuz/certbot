package reciever

type Unmarshaler interface {
	Unmarshal(v any) error
}

type Factory interface {
	Name() string
	NewReciever(cfg Unmarshaler) (Reciever, error)
}

var (
	factories = make(map[string]Factory)
)

func RegisterFactory(factory Factory) {
	if _, ok := factories[factory.Name()]; ok {
		panic("duplicate reciever factory name registered")
	}
	factories[factory.Name()] = factory
}

func GetFactory(name string) Factory {
	return factories[name]
}
