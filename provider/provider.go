package provider

type Provider interface {
	PresentTXTRecord(host string, domain string, txt string) error
	CleanUpTXTRecord(host string, domain string, txt string) error
}
type ProviderConfig interface {
	Unmarshal(v any) error
}
type ProviderFactory interface {
	Name() string
	NewProvider(cfg ProviderConfig) (Provider, error)
}

var (
	gProviders = make(map[string]ProviderFactory)
)

func RegisterProvider(factory ProviderFactory) {
	name := factory.Name()
	if _, ok := gProviders[name]; ok {
		panic("duplicate provider name:" + name)
	}
	gProviders[name] = factory
}

func LookupProviderFactory(name string) (ProviderFactory, bool) {
	factory, ok := gProviders[name]
	return factory, ok
}
