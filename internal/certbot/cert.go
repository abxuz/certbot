package certbot

import (
	"certbot/reciever"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

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
