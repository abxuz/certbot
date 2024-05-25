package cloudflare

import "github.com/go-acme/lego/v4/providers/dns/cloudflare"

type Config struct {
	AuthEmail string `yaml:"auth_email"`
	AuthKey   string `yaml:"auth_key"`
	AuthToken string `yaml:"auth_token"`
	ZoneToken string `yaml:"zone_token"`
}

type Provider struct {
	p *cloudflare.DNSProvider
}

func (p *Provider) Present(host, domain, token, keyAuth string) error {
	if host == "" {
		return p.p.Present(domain, token, keyAuth)
	}
	return p.p.Present(host+"."+domain, token, keyAuth)
}

func (p *Provider) CleanUp(host, domain, token, keyAuth string) error {
	if host == "" {
		return p.p.CleanUp(domain, token, keyAuth)
	}
	return p.p.CleanUp(host+"."+domain, token, keyAuth)
}
