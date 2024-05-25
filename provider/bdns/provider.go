package bdns

import (
	"bytes"
	"certbot/provider"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

type Config struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Provider struct {
	cfg *Config
}

type Record struct {
	Id     string `json:"id,omitempty"`
	Domain string `json:"domain"`
	Host   string `json:"host"`
	TTL    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

func (p *Provider) Present(host, domain, token, keyAuth string) error {
	return p.present(
		provider.ChallengeHost(host),
		domain,
		provider.ChallengeValue(keyAuth),
	)
}

func (p *Provider) CleanUp(host, domain, token, keyAuth string) error {
	return p.cleanup(
		provider.ChallengeHost(host),
		domain,
		provider.ChallengeValue(keyAuth),
	)
}

func (p *Provider) present(host string, domain string, txt string) error {
	api := p.cfg.Addr + "/api/v1/domain/" + domain + "/record"
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(&Record{
		Domain: domain,
		Host:   host,
		Remark: "certbot auto generated",
		TTL:    600,
		Type:   "TXT",
		Value:  txt,
	})
	if err != nil {
		return err
	}
	data, err := p.request(http.MethodPost, api, body)
	if err != nil {
		return err
	}

	return p.errCheck(data)
}

func (p *Provider) cleanup(host string, domain string, txt string) error {
	ids := make([]string, 0)
	p.WalkRecord(domain, func(r *Record) {
		if r.Host != host || r.Domain != domain {
			return
		}
		if r.Type != "TXT" || r.Value != txt {
			return
		}
		ids = append(ids, r.Id)
	})

	for _, id := range ids {
		err := p.RemoveRecord(domain, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) WalkRecord(domain string, f func(r *Record)) error {
	pageNumber := 0
	for {
		pageNumber++
		list, err := p.ListRecord(domain, pageNumber, 100)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			break
		}
		for _, r := range list {
			f(r)
		}
	}
	return nil
}

func (p *Provider) ListRecord(domain string, pageNumber, pageSize int) ([]*Record, error) {
	params := make(url.Values)
	params.Set("pageNumber", strconv.Itoa(pageNumber))
	params.Set("pageSize", strconv.Itoa(pageSize))
	api := p.cfg.Addr + "/api/v1/domain/" + domain + "/record?" + params.Encode()
	data, err := p.request(http.MethodGet, api, nil)
	if err != nil {
		return nil, err
	}

	result := &struct {
		Errno *int `json:"errno"`
		Data  *struct {
			List []*Record `json:"list"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}

	if result.Errno == nil || *result.Errno != 0 {
		return nil, errors.New(string(data))
	}
	return result.Data.List, nil
}

func (p *Provider) RemoveRecord(domain, id string) error {
	api := p.cfg.Addr + "/api/v1/domain/" + domain + "/record/" + id
	data, err := p.request(http.MethodDelete, api, nil)
	if err != nil {
		return err
	}
	return p.errCheck(data)
}

func (p *Provider) request(method string, api string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, api, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if p.cfg.Username != "" {
		request.SetBasicAuth(p.cfg.Username, p.cfg.Password)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (p *Provider) errCheck(data []byte) error {
	result := &struct {
		Errno  *int   `json:"errno"`
		Errmsg string `json:"errmsg"`
	}{}
	if err := json.Unmarshal(data, result); err != nil {
		return err
	}
	if result.Errmsg != "" {
		return errors.New(result.Errmsg)
	}
	if result.Errno == nil || *result.Errno != 0 {
		return errors.New(string(data))
	}
	return nil
}
