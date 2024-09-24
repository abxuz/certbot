package provider

import (
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
)

type AliyunProviderConfig struct {
	RAMRole string `yaml:"ram_role"`

	APIKey        string `yaml:"api_key"`
	SecretKey     string `yaml:"secret_key"`
	SecurityToken string `yaml:"security_token"`
}

func NewAliyunProvider(cfg ConfigDecoder) (challenge.Provider, error) {
	c := new(AliyunProviderConfig)
	if err := cfg.Decode(c); err != nil {
		return nil, err
	}

	config := alidns.NewDefaultConfig()
	config.APIKey = c.APIKey
	config.SecretKey = c.SecretKey
	config.SecurityToken = c.SecurityToken

	return alidns.NewDNSProviderConfig(config)
}
