package config

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/abxuz/b-tools/bslice"
	"gopkg.in/yaml.v3"
)

type YamlRaw struct {
	node *yaml.Node
}

func (c *YamlRaw) UnmarshalYAML(node *yaml.Node) error {
	c.node = node
	return nil
}

func (c *YamlRaw) Unmarshal(v any) error {
	return c.node.Decode(v)
}

type Provider struct {
	Name     string   `yaml:"name"`
	Provider string   `yaml:"provider"`
	Config   *YamlRaw `yaml:"config"`
}

type Reciever struct {
	Name     string   `yaml:"name"`
	Reciever string   `yaml:"reciever"`
	Config   *YamlRaw `yaml:"config"`
}

type User struct {
	Email string `yaml:"email"`
	Key   string `yaml:"key"`
}

type Domain struct {
	Host     string `yaml:"host"`
	Domain   string `yaml:"domain"`
	Provider string `yaml:"provider"`
}

type Cert struct {
	Name      string    `yaml:"name"`
	File      string    `yaml:"file"`
	User      string    `yaml:"user"`
	Domains   []*Domain `yaml:"domains"`
	Recievers []string  `yaml:"recievers"`
}

type Config struct {
	Providers []*Provider `yaml:"providers"`
	Recievers []*Reciever `yaml:"recievers"`
	Users     []*User     `yaml:"users"`
	Certs     []*Cert     `yaml:"certs"`
}

func LoadConfigFromReader(r io.Reader) (*Config, error) {
	cfg := &Config{}
	err := yaml.NewDecoder(r).Decode(cfg)
	if err != nil {
		return nil, err
	}
	if err := cfg.Valid(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadConfigFromData(b []byte) (*Config, error) {
	return LoadConfigFromReader(bytes.NewReader(b))
}

func LoadConfigFromFile(p string) (*Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadConfigFromReader(f)
}

func (p *Provider) Valid() error {
	if p.Name == "" {
		return fmt.Errorf("missing provider name")
	}
	if p.Provider == "" {
		return fmt.Errorf("missing provider[%v] type", p.Name)
	}
	return nil
}

func (r *Reciever) Valid() error {
	if r.Name == "" {
		return fmt.Errorf("missing reciever name")
	}
	if r.Reciever == "" {
		return fmt.Errorf("missing reciever[%v] type", r.Name)
	}
	return nil
}

func (u *User) Valid() error {
	if u.Email == "" {
		return fmt.Errorf("missing user email")
	}
	if u.Key == "" {
		return fmt.Errorf("missing user[%v] key file", u.Email)
	}
	return nil
}

func (d *Domain) Valid() error {
	if d.Domain == "" {
		return fmt.Errorf("missing domain")
	}
	if d.Provider == "" {
		return fmt.Errorf("missing provider [%v]", d.Hostname())
	}
	return nil
}

func (d *Domain) Hostname() string {
	if d.Host == "" {
		return d.Domain
	}
	return d.Host + "." + d.Domain
}

func (c *Cert) Valid() error {
	if c.Name == "" {
		return fmt.Errorf("missing cert name")
	}
	if c.File == "" {
		return fmt.Errorf("missing cert[%v] file path", c.Name)
	}
	if c.User == "" {
		return fmt.Errorf("missing cert[%v] user", c.Name)
	}
	if len(c.Domains) == 0 {
		return fmt.Errorf("missing cert[%v] domains", c.Name)
	}
	for _, d := range c.Domains {
		if err := d.Valid(); err != nil {
			return err
		}
	}
	ok := bslice.Unique[string, *Domain](c.Domains, func(d *Domain) string { return d.Hostname() })
	if !ok {
		return fmt.Errorf("duplicate cert[%v] domain host.domain", c.Name)
	}
	return nil
}

func (c *Config) Valid() error {
	ok := bslice.Unique[string, *Provider](c.Providers, func(p *Provider) string { return p.Name })
	if !ok {
		return fmt.Errorf("duplicate provider name found")
	}
	for _, p := range c.Providers {
		if err := p.Valid(); err != nil {
			return err
		}
	}

	ok = bslice.Unique[string, *Reciever](c.Recievers, func(r *Reciever) string { return r.Name })
	if !ok {
		return fmt.Errorf("duplicate reciever name found")
	}
	for _, r := range c.Recievers {
		if err := r.Valid(); err != nil {
			return err
		}
	}

	ok = bslice.Unique[string, *User](c.Users, func(u *User) string { return u.Email })
	if !ok {
		return fmt.Errorf("duplicate user email found")
	}
	for _, u := range c.Users {
		if err := u.Valid(); err != nil {
			return err
		}
	}

	ok = bslice.Unique[string, *Cert](c.Certs, func(c *Cert) string { return c.Name })
	if !ok {
		return fmt.Errorf("duplicate cert name found")
	}
	for _, cert := range c.Certs {
		if err := cert.Valid(); err != nil {
			return err
		}
		for _, d := range cert.Domains {
			pos := bslice.FindIndex[*Provider](c.Providers, func(p *Provider) bool { return p.Name == d.Provider })
			if pos < 0 {
				return fmt.Errorf("undefined provider[%v] in cert[%v] domain[%v]", d.Provider, cert.Name, d.Hostname())
			}
		}
		ok = bslice.Unique[string, string](cert.Recievers, func(s string) string { return s })
		if !ok {
			return fmt.Errorf("duplicate reciever found in cert[%v]", cert.Name)
		}

		pos := bslice.FindIndex[*User](c.Users, func(u *User) bool { return u.Email == cert.User })
		if pos < 0 {
			return fmt.Errorf("undefined user[%v] in cert[%v]", cert.User, cert.Name)
		}
	}
	return nil
}
