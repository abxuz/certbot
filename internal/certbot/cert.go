package certbot

import (
	"certbot/reciever"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
)

type CertRequest struct {
	Name      string
	File      string
	Client    *lego.Client
	Domains   []string
	Recievers []reciever.Reciever
}

func (r *CertRequest) LoadCert() (*x509.Certificate, error) {
	data, err := os.ReadFile(r.File)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, err
	}

	var block *pem.Block
	for {
		block, data = pem.Decode(data)
		if block == nil || block.Type == "CERTIFICATE" {
			break
		}
	}
	if block == nil {
		return nil, errors.New("cert not found")
	}
	return x509.ParseCertificate(block.Bytes)
}

func (r *CertRequest) ObtainCert() ([]byte, error) {
	obtain := certificate.ObtainRequest{
		Domains: r.Domains,
		Bundle:  true,
	}

	certificates, err := r.Client.Certificate.Obtain(obtain)
	if err != nil {
		return nil, err
	}

	certStr := string(certificates.Certificate) + string(certificates.IssuerCertificate) + string(certificates.PrivateKey)
	certData := []byte(certStr)
	return certData, nil
}

func (r *CertRequest) WriteCert(data []byte) error {
	err := os.MkdirAll(filepath.Dir(r.File), 0655)
	if err != nil {
		return err
	}
	return os.WriteFile(r.File, data, 0655)
}

func (r *CertRequest) PushCert(data []byte) error {
	errs := make([]error, 0)
	for _, recv := range r.Recievers {
		errs = append(errs, recv.PushCert(r.Name, data))
	}
	return errors.Join(errs...)
}
