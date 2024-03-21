package bvhost

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type RecieverConfig struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Reciever struct {
	cfg *RecieverConfig
}

type Cert struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (r *Reciever) PushCert(name string, cert []byte) error {
	api := r.cfg.Addr + "/api/v1/cert/" + name
	data, err := r.request(http.MethodGet, api, nil)
	if err != nil {
		return err
	}
	err = r.errCheck(data)
	if err == nil {
		err = r.setCert(true, &Cert{
			Name:    name,
			Content: string(cert),
		})
	} else if err.Error() == "cert not found" {
		err = r.setCert(false, &Cert{
			Name:    name,
			Content: string(cert),
		})
	}

	if err != nil {
		return err
	}

	if err := r.reload(); err != nil {
		return err
	}
	return r.save()
}

func (r *Reciever) setCert(update bool, cert *Cert) error {
	api := r.cfg.Addr + "/api/v1/cert/"
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(cert)
	if err != nil {
		return err
	}
	var method string
	if update {
		method = http.MethodPatch
	} else {
		method = http.MethodPost
	}
	data, err := r.request(method, api, body)
	if err != nil {
		return err
	}
	return r.errCheck(data)
}

func (r *Reciever) reload() error {
	api := r.cfg.Addr + "/api/v1/reload"
	data, err := r.request(http.MethodGet, api, nil)
	if err != nil {
		return err
	}
	return r.errCheck(data)
}

func (r *Reciever) save() error {
	api := r.cfg.Addr + "/api/v1/save"
	data, err := r.request(http.MethodGet, api, nil)
	if err != nil {
		return err
	}
	return r.errCheck(data)
}

func (r *Reciever) request(method string, api string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, api, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if r.cfg.Username != "" {
		request.SetBasicAuth(r.cfg.Username, r.cfg.Password)
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (r *Reciever) errCheck(data []byte) error {
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
