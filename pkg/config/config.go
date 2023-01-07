package config

type Config struct {
	Version string `yaml:"version"`
	URL     string `yaml:"url"`
	Token   string `yaml:"token"`
	Nodes   []Node `yaml:"nodes"`
}

type Node struct {
	Targets []string `yaml:"targets"`
	Type    NodeType `yaml:"type"`
	Labels  []string `yaml:"labels"`
	Taints  []string `yaml:"taints"`
}

type NodeType string

var (
	NodeTypeServer NodeType = "server"
	NodeTypeAgent  NodeType = "agent"
)
