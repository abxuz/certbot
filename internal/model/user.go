package model

import (
	"certbot/internal/config"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/registration"
)

type User struct {
	Email        string
	Key          crypto.PrivateKey
	Registration *registration.Resource
}

func NewUser(u *config.User) (*User, error) {

	m := &User{Email: u.Email}

	if data, err := os.ReadFile(u.Key); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		data = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
		if err := os.MkdirAll(filepath.Dir(u.Key), 0655); err != nil {
			return nil, err
		}
		if err := os.WriteFile(u.Key, data, 0655); err != nil {
			return nil, err
		}

		m.Key = key
	} else {
		var block *pem.Block
		for {
			block, data = pem.Decode(data)
			if block == nil || block.Type == "RSA PRIVATE KEY" {
				break
			}
		}
		if block == nil {
			return nil, errors.New("invalid private key file")
		}

		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		m.Key = key
	}

	return m, nil
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}
