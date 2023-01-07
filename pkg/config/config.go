package config

type Config struct {
	Version string `yaml:"version"`
	Nodes   []Node `yaml:"nodes"`
	Secret  Secret `yaml:"secret"`
}

type Node struct {
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
