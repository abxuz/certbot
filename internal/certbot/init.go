package certbot

import (
	"certbot/internal/model"
	"certbot/provider"
	"certbot/reciever"
	"certbot/reciever/bvhost"
	"fmt"
	"os"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"gopkg.in/yaml.v3"
)

const CADirURL = "https://acme-v02.api.letsencrypt.org/directory"

// const CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"

func (cb *CertBot) Init(config string) error {
	cfg, err := cb.initConfig(config)
	if err != nil {
		return err
	}

	userMap, err := cb.initUsers(cfg.Users)
	if err != nil {
		return err
	}

	recieverMap, err := cb.initRecievers(cfg.Recievers)
	if err != nil {
		return err
	}

	providerMap, err := cb.initProviders(cfg.Providers)
	if err != nil {
		return err
	}

	requestMap := make(map[string]*CertRequest)
	for _, cert := range cfg.Certs {
		if _, ok := requestMap[cert.Name]; ok {
			return fmt.Errorf("duplicate cert name %v", cert.Name)
		}

		user, ok := userMap[cert.User]
		if !ok {
			return fmt.Errorf("user %v in cert not found", cert.User)
		}

		provider, ok := providerMap[cert.Provider]
		if !ok {
			return fmt.Errorf("provider %v in cert not found", cert.Provider)
		}

		legoCfg := lego.NewConfig(user)
		legoCfg.CADirURL = CADirURL
		legoCfg.Certificate.KeyType = user.KeyType
		client, err := lego.NewClient(legoCfg)
		if err != nil {
			return err
		}

		client.Challenge.SetDNS01Provider(
			provider, dns01.AddRecursiveNameservers([]string{"223.5.5.5"}),
		)

		recievers := make([]reciever.Reciever, 0)
		for _, recieverName := range cert.Recievers {
			reciever, ok := recieverMap[recieverName]
			if !ok {
				return fmt.Errorf("reciever %v in cert not found", recieverName)
			}
			recievers = append(recievers, reciever)
		}

		requestMap[cert.Name] = &CertRequest{
			Name:      cert.Name,
			File:      cert.File,
			Client:    client,
			Domains:   cert.Domains,
			Recievers: recievers,
		}
	}

	cb.requests = requestMap
	return nil
}

func (cb *CertBot) initConfig(config string) (cfg *model.Config, err error) {
	f, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&cfg)
	return
}

func (cb *CertBot) initUsers(users []*model.User) (map[string]*User, error) {
	m := make(map[string]*User)

	for _, u := range users {
		if _, ok := m[u.Email]; ok {
			return nil, fmt.Errorf("duplicate user email %v", u.Email)
		}
		user, err := NewUser(u)
		if err != nil {
			return nil, err
		}

		legoCfg := lego.NewConfig(user)
		legoCfg.CADirURL = CADirURL
		legoCfg.Certificate.KeyType = user.KeyType
		client, err := lego.NewClient(legoCfg)
		if err != nil {
			return nil, err
		}

		user.Registration, err = client.Registration.Register(
			registration.RegisterOptions{TermsOfServiceAgreed: true},
		)
		if err != nil {
			return nil, err
		}

		m[u.Email] = user
	}

	return m, nil
}

func (cb *CertBot) initRecievers(recievers []*model.Reciever) (map[string]reciever.Reciever, error) {
	bvhost.Init()

	m := make(map[string]reciever.Reciever)
	for _, cfg := range recievers {
		_, ok := m[cfg.Name]
		if ok {
			return nil, fmt.Errorf("duplicate reciever name [%v]", cfg.Name)
		}

		factory := reciever.GetFactory(cfg.Reciever)
		if factory == nil {
			return nil, fmt.Errorf("reciever type [%v] not found", cfg.Reciever)
		}

		recver, err := factory.NewReciever(cfg.Config)
		if err != nil {
			return nil, err
		}
		m[cfg.Name] = recver
	}

	return m, nil
}

func (cb *CertBot) initProviders(providers []*model.Provider) (map[string]challenge.Provider, error) {
	m := make(map[string]challenge.Provider)
	for _, cfg := range providers {
		if _, ok := m[cfg.Name]; ok {
			return nil, fmt.Errorf("duplicate provider name %v", cfg.Name)
		}

		factory := provider.GetFactory(cfg.Type)
		if factory == nil {
			return nil, fmt.Errorf("unkown provider type %v", cfg.Type)
		}

		p, err := factory.NewProvider(cfg.Config)
		if err != nil {
			return nil, err
		}
		m[cfg.Name] = p
	}

	return m, nil
}
