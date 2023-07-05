package certbot

import (
	"certbot/internal/config"
	"certbot/provider"
	"fmt"

	"github.com/go-acme/lego/v4/challenge/dns01"
)

type Provider struct {
	domain   *config.Domain
	provider provider.Provider
}

type ProviderDelegate struct {
	providers map[string]*Provider
}

func NewProviderDelegate() *ProviderDelegate {
	return &ProviderDelegate{
		providers: make(map[string]*Provider),
	}
}

func (d *ProviderDelegate) SetProvider(domain *config.Domain, p provider.Provider) {
	d.providers["_acme-challenge."+domain.Hostname()+"."] = &Provider{domain: domain, provider: p}
}

func (d *ProviderDelegate) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)
	provider, ok := d.providers[info.FQDN]
	if !ok {
		return fmt.Errorf("provider not found for fqdn: %v", info.FQDN)
	}
	host := "_acme-challenge"
	if provider.domain.Host != "" {
		host += "." + provider.domain.Host
	}
	return provider.provider.PresentTXTRecord(host, provider.domain.Domain, info.Value)
}

func (d *ProviderDelegate) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)
	provider, ok := d.providers[info.FQDN]
	if !ok {
		return fmt.Errorf("provider not found for fqdn: %v", info.FQDN)
	}
	host := "_acme-challenge"
	if provider.domain.Host != "" {
		host += "." + provider.domain.Host
	}
	return provider.provider.CleanUpTXTRecord(host, provider.domain.Domain, info.Value)
}
