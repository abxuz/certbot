package provider

import "github.com/go-acme/lego/v4/challenge"

type ConfigDecoder interface {
	Decode(v any) error
}

type Factory interface {
	NewProvider(cfg ConfigDecoder) (challenge.Provider, error)
}

type FactoryFunc func(cfg ConfigDecoder) (challenge.Provider, error)

func (f FactoryFunc) NewProvider(cfg ConfigDecoder) (challenge.Provider, error) {
	return f(cfg)
}

var (
	factories = make(map[string]Factory)
)

func RegisterFactory(name string, factory Factory) {
	if _, ok := factories[name]; ok {
		panic("duplicate provider factory name registered")
	}
	factories[name] = factory
}

func GetFactory(name string) Factory {
	return factories[name]
}
