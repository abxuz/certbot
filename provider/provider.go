package provider

type Provider interface {
	Present(host, domain, token, keyAuth string) error
	CleanUp(host, domain, token, keyAuth string) error
}
