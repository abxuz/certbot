package model

import "gopkg.in/yaml.v3"

type YamlRaw struct {
	node *yaml.Node
}

func (c *YamlRaw) UnmarshalYAML(node *yaml.Node) error {
	c.node = node
	return nil
}

func (c *YamlRaw) Decode(v any) error {
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
	Name   string   `yaml:"name"`
	Type   string   `yaml:"type"`
	Config *YamlRaw `yaml:"config"`
}

type Cert struct {
	Name      string   `yaml:"name"`
	File      string   `yaml:"file"`
	User      string   `yaml:"user"`
	Provider  string   `yaml:"provider"`
	Domains   []string `yaml:"domains"`
	Recievers []string `yaml:"recievers"`
}
