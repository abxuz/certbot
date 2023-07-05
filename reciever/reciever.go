package reciever

type Reciever interface {
	PushCert(name string, cert []byte) error
}
type RecieverConfig interface {
	Unmarshal(v any) error
}
type RecieverFactory interface {
	Name() string
	NewReciever(cfg RecieverConfig) (Reciever, error)
}

var (
	gRecievers = make(map[string]RecieverFactory)
)

func RegisterReciever(factory RecieverFactory) {
	name := factory.Name()
	if _, ok := gRecievers[name]; ok {
		panic("duplicate reciever name:" + name)
	}
	gRecievers[name] = factory
}

func LookupRecieverFactory(name string) (RecieverFactory, bool) {
	factory, ok := gRecievers[name]
	return factory, ok
}
