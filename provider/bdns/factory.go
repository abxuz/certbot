package bdns

import (
	"certbot/provider"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const (
	ProviderName = "bdns"
)

var (
	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 5,
	}
)

type factory string

func init() {
	provider.RegisterProvider(factory(ProviderName))
}

func (f factory) Name() string {
	return string(f)
}

func (f factory) NewProvider(c provider.ProviderConfig) (provider.Provider, error) {
	cfg := &ProviderConfig{}
	if err := c.Unmarshal(cfg); err != nil {
		return nil, err
	}
	p := &Provider{
		cfg: cfg,
	}
	return p, nil
}
