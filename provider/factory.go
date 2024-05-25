package provider

type Unmarshaler interface {
	Unmarshal(v any) error
}

type Factory interface {
	Name() string
	NewProvider(cfg Unmarshaler) (Provider, error)
}

var (
	factories = make(map[string]Factory)
)

func RegisterFactory(factory Factory) {
	if _, ok := factories[factory.Name()]; ok {
		panic("duplicate provider factory name registered")
	}
}

func GetFactory(name string) Factory {
	return factories[name]
}
