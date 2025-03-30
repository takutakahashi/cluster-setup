package proxmox

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/cluster-setup/pkg/config"
)

// Manager handles Proxmox VM lifecycle operations
type Manager struct {
	config *config.Config
}

// NewManager creates a new Proxmox Manager
func NewManager(cfg *config.Config) (*Manager, error) {
	if cfg == nil {
		return nil, errors.New("configuration is nil")
	}
	
	if cfg.Proxmox == nil {
		return nil, errors.New("proxmox configuration is missing")
	}

	return &Manager{
		config: cfg,
	}, nil
}

// CreateNodes creates all VMs specified in the config
func (m *Manager) CreateNodes() error {
	if len(m.config.ProxmoxNodes) == 0 {
		logrus.Info("No Proxmox nodes to create")
		return nil
	}

	logrus.Infof("Creating %d Proxmox nodes", len(m.config.ProxmoxNodes))
	
	// Track successful and failed creations
	var createdNodes []string
	var failedNodes []string
	
	for _, nodeSpec := range m.config.ProxmoxNodes {
		logrus.Infof("Creating node: %s", nodeSpec.Name)
		
		client, err := NewClient(m.config.Proxmox, nodeSpec.Name)
		if err != nil {
			logrus.Errorf("Failed to create Proxmox client for node %s: %v", nodeSpec.Name, err)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		vmid, err := client.CreateVM(nodeSpec.VMConfig)
		if err != nil {
			logrus.Errorf("Failed to create VM for node %s: %v", nodeSpec.Name, err)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		logrus.Infof("Successfully created node %s with VMID %d", nodeSpec.Name, vmid)
		
		// If the VM has an IP address, update targets in Node configuration
		if nodeSpec.VMConfig.IPAddress != "" {
			m.addNodeTarget(nodeSpec.Name, nodeSpec.VMConfig.IPAddress)
			logrus.Infof("Added target %s to node %s", nodeSpec.VMConfig.IPAddress, nodeSpec.Name)
		}
		
		createdNodes = append(createdNodes, nodeSpec.Name)
	}
	
	// Summary of operations
	if len(createdNodes) > 0 {
		logrus.Infof("Successfully created nodes: %v", createdNodes)
	}
	
	if len(failedNodes) > 0 {
		return fmt.Errorf("failed to create nodes: %v", failedNodes)
	}
	
	return nil
}

// DeleteNodes deletes all VMs specified in the config
func (m *Manager) DeleteNodes() error {
	if len(m.config.ProxmoxNodes) == 0 {
		logrus.Info("No Proxmox nodes to delete")
		return nil
	}

	logrus.Infof("Deleting %d Proxmox nodes", len(m.config.ProxmoxNodes))
	
	// Track successful and failed deletions
	var deletedNodes []string
	var failedNodes []string
	
	for _, nodeSpec := range m.config.ProxmoxNodes {
		logrus.Infof("Deleting node: %s", nodeSpec.Name)
		
		client, err := NewClient(m.config.Proxmox, nodeSpec.Name)
		if err != nil {
			logrus.Errorf("Failed to create Proxmox client for node %s: %v", nodeSpec.Name, err)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		// Skip if VMID is not specified and no current VM exists with the name
		if nodeSpec.VMConfig.VMID == 0 {
			logrus.Warnf("VMID not specified for node %s, skipping deletion", nodeSpec.Name)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		if err := client.DeleteVM(nodeSpec.VMConfig.VMID); err != nil {
			logrus.Errorf("Failed to delete VM for node %s: %v", nodeSpec.Name, err)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		logrus.Infof("Successfully deleted node %s", nodeSpec.Name)
		deletedNodes = append(deletedNodes, nodeSpec.Name)
	}
	
	// Summary of operations
	if len(deletedNodes) > 0 {
		logrus.Infof("Successfully deleted nodes: %v", deletedNodes)
	}
	
	if len(failedNodes) > 0 {
		return fmt.Errorf("failed to delete nodes: %v", failedNodes)
	}
	
	return nil
}

// UpdateNodes updates existing VMs based on configuration
func (m *Manager) UpdateNodes() error {
	if len(m.config.ProxmoxNodes) == 0 {
		logrus.Info("No Proxmox nodes to update")
		return nil
	}

	logrus.Infof("Updating %d Proxmox nodes", len(m.config.ProxmoxNodes))
	
	// Track successful and failed updates
	var updatedNodes []string
	var failedNodes []string
	
	for _, nodeSpec := range m.config.ProxmoxNodes {
		logrus.Infof("Updating node: %s", nodeSpec.Name)
		
		client, err := NewClient(m.config.Proxmox, nodeSpec.Name)
		if err != nil {
			logrus.Errorf("Failed to create Proxmox client for node %s: %v", nodeSpec.Name, err)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		// Skip if VMID is not specified
		if nodeSpec.VMConfig.VMID == 0 {
			logrus.Warnf("VMID not specified for node %s, skipping update", nodeSpec.Name)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		// Apply configuration updates to the VM
		if err := client.configureVM(nodeSpec.VMConfig.VMID, nodeSpec.VMConfig); err != nil {
			logrus.Errorf("Failed to update VM for node %s: %v", nodeSpec.Name, err)
			failedNodes = append(failedNodes, nodeSpec.Name)
			continue
		}
		
		logrus.Infof("Successfully updated node %s", nodeSpec.Name)
		updatedNodes = append(updatedNodes, nodeSpec.Name)
	}
	
	// Summary of operations
	if len(updatedNodes) > 0 {
		logrus.Infof("Successfully updated nodes: %v", updatedNodes)
	}
	
	if len(failedNodes) > 0 {
		return fmt.Errorf("failed to update nodes: %v", failedNodes)
	}
	
	return nil
}

// addNodeTarget adds a target to the appropriate Node configuration
func (m *Manager) addNodeTarget(nodeName, ipAddress string) {
	// Find the matching node in the config
	for i, node := range m.config.Nodes {
		if node.Name == nodeName {
			// Format the target as username@ip
			// We'll use 'ubuntu' as a default username but this could be made configurable
			target := fmt.Sprintf("ubuntu@%s", ipAddress)
			
			// Add to targets if not already present
			found := false
			for _, existingTarget := range node.Targets {
				if existingTarget == target {
					found = true
					break
				}
			}
			
			if !found {
				m.config.Nodes[i].Targets = append(m.config.Nodes[i].Targets, target)
			}
			return
		}
	}
	
	// If no matching node was found, create a new one
	newNode := config.Node{
		Name:    nodeName,
		Targets: []string{fmt.Sprintf("ubuntu@%s", ipAddress)},
	}
	
	// Determine if this should be server or agent
	// For simplicity, the first node is a server, rest are agents
	if len(m.config.Nodes) == 0 {
		newNode.Type = config.NodeTypeServer
	} else {
		newNode.Type = config.NodeTypeAgent
	}
	
	m.config.Nodes = append(m.config.Nodes, newNode)
}