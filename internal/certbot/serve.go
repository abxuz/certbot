package certbot

import (
	"log"
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
			log.Println(err)
			continue
		}
		if cert != nil && time.Until(cert.NotAfter) >= 2*30*24*time.Hour {
			continue
		}
		eg.Go(func() error {
			return cb.DoRequest(request)
		})
	}

	eg.Wait()
}

func (cb *CertBot) DoRequest(request *CertRequest) error {
	log.Printf("DoRequest for: %v", request.Name)

	obtain := certificate.ObtainRequest{
		Domains: request.Domains,
		Bundle:  true,
	}

	certificates, err := request.Client.Certificate.Obtain(obtain)
	if err != nil {
		log.Printf("DoRequest for [%v] Obtain error: %v", request.Name, err)
		return err
	}

	certStr := string(certificates.Certificate) + string(certificates.IssuerCertificate) + string(certificates.PrivateKey)
	certData := []byte(certStr)

	for _, reciever := range request.Recievers {
		err = reciever.PushCert(request.Name, certData)
		if err != nil {
			log.Printf("DoRequest for [%v] PushCert error: %v", request.Name, err)
			return err
		}
	}

	err = os.MkdirAll(filepath.Dir(request.File), 0655)
	if err != nil {
		log.Printf("DoRequest for [%v] MkdirAll error: %v", request.Name, err)
		return err
	}

	err = os.WriteFile(request.File, certData, 0655)
	if err != nil {
		log.Printf("DoRequest for [%v] WriteFile error: %v", request.Name, err)
		return err
	}

	return nil
}
