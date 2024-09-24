package reciever

type ConfigDecoder interface {
	Decode(v any) error
}

type Factory interface {
	Name() string
	NewReciever(cfg ConfigDecoder) (Reciever, error)
}

var (
	factories = make(map[string]Factory)
)

func RegisterFactory(factory Factory) {
	name := factory.Name()
	if _, ok := factories[name]; ok {
		panic("duplicate reciever factory name registered")
	}
	factories[name] = factory
}

func GetFactory(name string) Factory {
	return factories[name]
}
