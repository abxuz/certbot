package certbot

import (
	"certbot/provider"
	"fmt"
)

type delegatedProvider struct {
	Host     string
	Domain   string
	Provider provider.Provider
}

type Delegater struct {
	providers map[string]*delegatedProvider
}

func NewDelegater() *Delegater {
	return &Delegater{providers: make(map[string]*delegatedProvider)}
}

func (d *Delegater) fqdn(host, domain string) string {
	return provider.ChallengeHost(host) + "." + domain + "."
}

func (d *Delegater) SetProvider(host, domain string, p provider.Provider) error {
	fqdn := d.fqdn(host, domain)
	if _, ok := d.providers[fqdn]; ok {
		return fmt.Errorf("duplicate host [%v] and domain [%v]", host, domain)
	}
	d.providers[fqdn] = &delegatedProvider{
		Host:     host,
		Domain:   domain,
		Provider: p,
	}
	return nil
}

func (d *Delegater) Present(domain, token, keyAuth string) error {
	fqdn := provider.ChallengeHost(domain) + "."
	provider, ok := d.providers[fqdn]
	if !ok {
		return fmt.Errorf("provider not found for domain %v", domain)
	}
	return provider.Provider.Present(provider.Host, provider.Domain, token, keyAuth)
}

func (d *Delegater) CleanUp(domain, token, keyAuth string) error {
	fqdn := provider.ChallengeHost(domain) + "."
	provider, ok := d.providers[fqdn]
	if !ok {
		return fmt.Errorf("provider not found for domain %v", domain)
	}
	return provider.Provider.CleanUp(provider.Host, provider.Domain, token, keyAuth)
}
