package certbot

import (
	"log"
	"time"

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

	cert, err := request.ObtainCert()
	if err != nil {
		log.Println("failed to obtain cert", err)
		return err
	}

	err = request.WriteCert(cert)
	if err != nil {
		log.Println("failed to write cert", err)
		return err
	}

	err = request.PushCert(cert)
	if err != nil {
		log.Println("failed to push cert", err)
	}
	return nil
}
