package model

import (
	"certbot/internal/config"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/registration"
)

type User struct {
	Email        string
	Key          crypto.PrivateKey
	KeyType      certcrypto.KeyType
	Registration *registration.Resource
}

func NewUser(u *config.User) (*User, error) {

	m := &User{Email: u.Email}

	if data, err := os.ReadFile(u.Key); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		der, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return nil, err
		}
		x509.MarshalPKCS8PrivateKey(key)
		data = pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: der,
		})
		if err := os.MkdirAll(filepath.Dir(u.Key), 0655); err != nil {
			return nil, err
		}
		if err := os.WriteFile(u.Key, data, 0655); err != nil {
			return nil, err
		}

		m.Key = key
	} else {
		loop := true
		for loop {
			var block *pem.Block
			block, data = pem.Decode(data)
			if block == nil {
				break
			}
			switch block.Type {
			case "PRIVATE KEY":
				pk, err := x509.ParsePKCS8PrivateKey(block.Bytes)
				if err != nil {
					return nil, err
				}
				m.Key = pk
				loop = false
			case "RSA PRIVATE KEY":
				pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
				if err != nil {
					return nil, err
				}
				m.Key = pk
				loop = false
			case "EC PRIVATE KEY":
				pk, err := x509.ParseECPrivateKey(block.Bytes)
				if err != nil {
					return nil, err
				}
				m.Key = pk
				loop = false
			}
		}
		if m.Key == nil {
			return nil, errors.New("no private key found")
		}
	}

	switch k := m.Key.(type) {
	case *rsa.PrivateKey:
		switch k.N.BitLen() {
		case 2048:
			m.KeyType = certcrypto.RSA2048
		case 3072:
			m.KeyType = certcrypto.RSA3072
		case 4096:
			m.KeyType = certcrypto.RSA4096
		case 8192:
			m.KeyType = certcrypto.RSA8192
		default:
			return nil, errors.New("unsupported rsa key len")
		}
	case *ecdsa.PrivateKey:
		switch k.Curve {
		case elliptic.P256():
			m.KeyType = certcrypto.EC256
		case elliptic.P384():
			m.KeyType = certcrypto.EC384
		default:
			return nil, errors.New("unsupported ecc curve")
		}
	default:
		return nil, errors.New("unsupported private key type")
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
