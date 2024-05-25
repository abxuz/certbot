package bdns

import (
	"certbot/provider"
)

func Init() {
	provider.RegisterFactory(&factory{})
}

type factory struct{}

func (f *factory) Name() string {
	return "bdns"
}

func (f *factory) NewProvider(cfg provider.Unmarshaler) (provider.Provider, error) {
	p := new(Provider)
	err := cfg.Unmarshal(&p.cfg)
	if err != nil {
		return nil, err
	}
	return p, nil
}
