package provider

import (
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
)

type CloudflareProviderConfig struct {
	AuthEmail string `yaml:"auth_email"`

	AuthToken string `yaml:"auth_token"`
	ZoneToken string `yaml:"zone_token"`

	AuthKey string `yaml:"auth_key"`
}

func NewCloudflareProvider(cfg ConfigDecoder) (challenge.Provider, error) {
	c := new(CloudflareProviderConfig)
	if err := cfg.Decode(c); err != nil {
		return nil, err
	}

	config := cloudflare.NewDefaultConfig()
	config.AuthEmail = c.AuthEmail
	config.AuthKey = c.AuthKey
	config.AuthToken = c.AuthToken
	config.ZoneToken = c.ZoneToken

	return cloudflare.NewDNSProviderConfig(config)
}
