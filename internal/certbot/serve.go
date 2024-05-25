package certbot

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"golang.org/x/sync/errgroup"
)

func (cb *CertBot) Serve() {
	for {
		cb.ServeOnce()
		time.Sleep(time.Minute)
	}
}

func (cb *CertBot) ServeOnce() {
	eg := &errgroup.Group{}
	eg.SetLimit(10)

	for _, request := range cb.requests {
		cert, err := request.LoadCert()
		if err != nil {
			continue
		}
		if time.Until(cert.NotAfter) >= 2*30*24*time.Hour {
			continue
		}
		eg.Go(func() error {
			return cb.DoRequest(request)
		})
	}

	eg.Wait()
}

func (cb *CertBot) DoRequest(request *CertRequest) error {
	obtain := certificate.ObtainRequest{
		Domains: request.Domains,
		Bundle:  true,
	}

	certificates, err := request.Client.Certificate.Obtain(obtain)
	if err != nil {
		return err
	}

	certStr := string(certificates.Certificate) + string(certificates.IssuerCertificate) + string(certificates.PrivateKey)
	certData := []byte(certStr)

	for _, reciever := range request.Recievers {
		err = reciever.PushCert(request.Name, certData)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(filepath.Dir(request.File), 0655)
	if err != nil {
		return err
	}

	return os.WriteFile(request.File, certData, 0655)
}
