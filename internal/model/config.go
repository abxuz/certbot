package model

import "gopkg.in/yaml.v3"

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

type Config struct {
	Users     []*User     `yaml:"users"`
	Recievers []*Reciever `yaml:"recievers"`
	Providers []*Provider `yaml:"providers"`
	Certs     []*Cert     `yaml:"certs"`
}

type User struct {
	Email string `yaml:"email"`
	Key   string `yaml:"key"`
}

type Reciever struct {
	Name     string   `yaml:"name"`
	Reciever string   `yaml:"reciever"`
	Config   *YamlRaw `yaml:"config"`
}

type Provider struct {
	Name     string   `yaml:"name"`
	Provider string   `yaml:"provider"`
	Config   *YamlRaw `yaml:"config"`
}

type Cert struct {
	Name      string    `yaml:"name"`
	File      string    `yaml:"file"`
	User      string    `yaml:"user"`
	Domains   []*Domain `yaml:"domains"`
	Recievers []string  `yaml:"recievers"`
}

type Domain struct {
	Host     string `yaml:"host"`
	Domain   string `yaml:"domain"`
	Provider string `yaml:"provider"`
}
