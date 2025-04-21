package config

import (
	"io/ioutil"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version string `yaml:"version"`
	Nodes   []Node `yaml:"nodes"`
	Secret  Secret `yaml:"secret"`
}

type Node struct {
	Name    string                 `yaml:"name"`
	Targets []string               `yaml:"targets"`
	Type    NodeType               `yaml:"type"`
	Labels  []string               `yaml:"labels"`
	Taints  []string               `yaml:"taints"`
	Params  map[string]interface{} `yaml:"params"`
}

type Secret struct {
	URL       string `env:"K3S_URL,required"`
	DataStore string `env:"K3S_DATASTORE,required"`
	Token     string `env:"K3S_TOKEN,required"`
}

type NodeType string

var (
	NodeTypeServer NodeType = "server"
	NodeTypeAgent  NodeType = "agent"
)

func Load(path string) (*Config, error) {
	c := &Config{}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(buf, c); err != nil {
		return nil, err
	}
	if c.Secret.URL == "" || c.Secret.DataStore == "" || c.Secret.Token == "" {
		s, err := loadSecretFromEnv()
		if err != nil {
			return nil, err
		}
		c.Secret = s
	}
	return c, nil
}

func loadSecretFromEnv() (Secret, error) {
	ret := Secret{}
	if err := env.Parse(&ret); err != nil {
		return ret, err
	}
	return ret, nil
}
