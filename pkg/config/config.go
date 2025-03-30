package config

import (
	"io/ioutil"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version      string            `yaml:"version"`
	Nodes        []Node            `yaml:"nodes"`
	Secret       Secret            `yaml:"secret"`
	ProxmoxNodes []ProxmoxNodeSpec `yaml:"proxmoxNodes,omitempty"`
	Proxmox      *ProxmoxConfig    `yaml:"proxmox,omitempty"`
}

type Node struct {
	Name    string                 `yaml:"name"`
	Targets []string               `yaml:"targets"`
	Type    NodeType               `yaml:"type"`
	Labels  []string               `yaml:"labels"`
	Taints  []string               `yaml:"taints"`
	Params  map[string]interface{} `yaml:"params"`
}

// ProxmoxNodeSpec defines a node to be created in Proxmox
type ProxmoxNodeSpec struct {
	Name     string    `yaml:"name"`
	VMConfig *VMConfig `yaml:"vmConfig"`
}

// ProxmoxConfig holds authentication and connection details for Proxmox
type ProxmoxConfig struct {
	// API connection details
	ApiUrl     string `yaml:"apiUrl"`
	TargetNode string `yaml:"targetNode"`
	SkipVerify bool   `yaml:"skipVerify,omitempty"`
	
	// Authentication (either token or username/password)
	APIToken string `yaml:"apiToken,omitempty"`
	TokenID  string `yaml:"tokenId,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// VMConfig defines configuration for a Proxmox VM
type VMConfig struct {
	// Basic VM properties
	Name        string `yaml:"name"`
	VMID        int    `yaml:"vmid,omitempty"`
	Description string `yaml:"description,omitempty"`
	
	// Resource allocation
	Cores  int `yaml:"cores,omitempty"`
	Memory int `yaml:"memory,omitempty"` // in MB
	
	// Template or ISO to create VM from
	Template string `yaml:"template"` // VMID of template to clone
	
	// Storage configuration
	Storage string `yaml:"storage,omitempty"`
	
	// Network configuration
	Network []*NetworkConfig `yaml:"network,omitempty"`
	
	// Cloud-Init configuration
	CloudInit *CloudInitConfig `yaml:"cloudInit,omitempty"`
	
	// VM behavior flags
	AutoStart bool `yaml:"autoStart,omitempty"`
	WaitForIP bool `yaml:"waitForIP,omitempty"`
	
	// Runtime values (not from YAML)
	IPAddress string `yaml:"-"`
}

// NetworkConfig defines a network interface for a VM
type NetworkConfig struct {
	Model  string `yaml:"model,omitempty"` // e.g., virtio, e1000, etc.
	Bridge string `yaml:"bridge,omitempty"` // e.g., vmbr0
	VLAN   int    `yaml:"vlan,omitempty"`  // VLAN tag
}

// CloudInitConfig defines cloud-init settings for VMs
type CloudInitConfig struct {
	User      string `yaml:"user,omitempty"`
	Password  string `yaml:"password,omitempty"`
	SSHKeys   string `yaml:"sshKeys,omitempty"`
	IPConfig  string `yaml:"ipConfig,omitempty"` // e.g., "ip=dhcp" or "ip=192.168.1.100/24,gw=192.168.1.1"
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
