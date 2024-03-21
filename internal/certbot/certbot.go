package certbot

import (
	"certbot/internal/config"
	"certbot/internal/model"
	"certbot/provider"
	"certbot/reciever"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"golang.org/x/sync/errgroup"
)

const CADirURL = "https://acme-v02.api.letsencrypt.org/directory"

type CertBot struct {
	Cfg *config.Config

	providers map[string]provider.Provider
	recievers map[string]reciever.Reciever
	users     map[string]*model.User
}

func (b *CertBot) Run() error {

	b.providers = make(map[string]provider.Provider)
	for _, p := range b.Cfg.Providers {
		factory, ok := provider.LookupProviderFactory(p.Provider)
		if !ok {
			return fmt.Errorf("undefined provider %v", p.Provider)
		}
		provider, err := factory.NewProvider(p.Config)
		if err != nil {
			return err
		}
		b.providers[p.Name] = provider
	}

	b.recievers = make(map[string]reciever.Reciever)
	for _, r := range b.Cfg.Recievers {
		factory, ok := reciever.LookupRecieverFactory(r.Reciever)
		if !ok {
			return fmt.Errorf("undefined reciever %v", r.Reciever)
		}
		reciever, err := factory.NewReciever(r.Config)
		if err != nil {
			return err
		}
		b.recievers[r.Name] = reciever
	}

	b.users = make(map[string]*model.User)
	for _, u := range b.Cfg.Users {
		user, err := model.NewUser(u)
		if err != nil {
			return err
		}

		config := lego.NewConfig(user)
		config.CADirURL = CADirURL
		config.Certificate.KeyType = user.KeyType
		client, err := lego.NewClient(config)
		if err != nil {
			return err
		}

		user.Registration, err = client.Registration.Register(
			registration.RegisterOptions{TermsOfServiceAgreed: true},
		)
		if err != nil {
			return err
		}

		b.users[user.Email] = user
	}

	for {
		b.process()
		time.Sleep(time.Minute)
	}
}

func (b *CertBot) process() {
	certs := make([]*config.Cert, 0)
	for _, c := range b.Cfg.Certs {
		data, err := os.ReadFile(c.File)
		if err != nil {
			if os.IsNotExist(err) {
				certs = append(certs, c)
			}
			continue
		}

		cert, err := parseCertificate(data)
		if err != nil {
			continue
		}
		if time.Until(cert.NotAfter) >= 2*30*24*time.Hour {
			continue
		}
		certs = append(certs, c)
	}

	eg := &errgroup.Group{}
	eg.SetLimit(10)

	for _, c := range certs {
		c := c
		eg.Go(func() error {
			b.processCert(c)
			return nil
		})
	}
	eg.Wait()
}

func (b *CertBot) processCert(cert *config.Cert) {
	user := b.users[cert.User]

	config := lego.NewConfig(user)
	config.CADirURL = CADirURL
	config.Certificate.KeyType = user.KeyType
	client, err := lego.NewClient(config)
	if err != nil {
		log.Println(err)
		return
	}

	delegate := NewProviderDelegate()
	domains := make([]string, 0)
	for _, d := range cert.Domains {
		delegate.SetProvider(d, b.providers[d.Provider])
		domains = append(domains, d.Hostname())
	}

	client.Challenge.SetDNS01Provider(delegate, dns01.AddRecursiveNameservers([]string{"223.5.5.5"}))
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Println(err)
		return
	}

	certStr := string(certificates.Certificate) + string(certificates.IssuerCertificate) + string(certificates.PrivateKey)
	certData := []byte(certStr)

	for _, r := range cert.Recievers {
		reciever := b.recievers[r]
		err := reciever.PushCert(cert.Name, certData)
		if err != nil {
			log.Println(err)
			return
		}
	}

	err = os.MkdirAll(filepath.Dir(cert.File), 0655)
	if err != nil {
		log.Println(err)
		return
	}
	err = os.WriteFile(cert.File, certData, 0655)
	if err != nil {
		log.Println(err)
		return
	}

}

func parseCertificate(data []byte) (*x509.Certificate, error) {
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
