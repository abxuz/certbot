package cloudflare

import (
	"certbot/provider"

	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

func Init() {
	provider.RegisterFactory(&factory{})
}

type factory struct{}

func (f *factory) Name() string {
	return "cloudflare"
}

func (f *factory) NewProvider(cfg provider.Unmarshaler) (provider.Provider, error) {
	var pcfg *Config
	err := cfg.Unmarshal(&pcfg)
	if err != nil {
		return nil, err
	}

	config := cloudflare.NewDefaultConfig()
	config.AuthEmail = pcfg.AuthEmail
	config.AuthKey = pcfg.AuthKey
	config.AuthToken = pcfg.AuthToken
	config.ZoneToken = pcfg.ZoneToken

	cfProvider, err := cloudflare.NewDNSProviderConfig(config)
	if err != nil {
		return nil, err
	}

	p := new(Provider)
	p.p = cfProvider
	return p, nil
}
