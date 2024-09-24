package provider

import (
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
)

type AliyunProviderConfig struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	AccessToken     string `yaml:"access_token"`

	RAMRole string `yaml:"ram_role"`
}

func NewAliyunProvider(cfg ConfigDecoder) (challenge.Provider, error) {
	c := new(AliyunProviderConfig)
	if err := cfg.Decode(c); err != nil {
		return nil, err
	}

	config := alidns.NewDefaultConfig()
	config.APIKey = c.AccessKeyId
	config.SecretKey = c.AccessKeySecret
	config.SecurityToken = c.AccessToken
	config.RAMRole = c.RAMRole

	return alidns.NewDNSProviderConfig(config)
}
