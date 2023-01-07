package config

import (
	"io/ioutil"

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
	URL       string `yaml:"url"`
	DataStore string `yaml:"datastore"`
	Token     string `yaml:"token"`
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
	return c, nil
}
