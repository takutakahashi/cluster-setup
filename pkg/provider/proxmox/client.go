package proxmox

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/cluster-setup/pkg/config"
)

type Client struct {
	client     *proxmox.Client
	config     *config.ProxmoxConfig
	nodeName   string
	targetNode string
}

// NewClient creates a new Proxmox client instance
func NewClient(cfg *config.ProxmoxConfig, nodeName string) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("proxmox configuration is nil")
	}

	tlsconfig := &proxmox.TLSConfiguration{
		InsecureSkipVerify: cfg.SkipVerify,
	}

	c, err := proxmox.NewClient(
		cfg.ApiUrl,
		nil, // Pass nil for full URL configuration
		tlsconfig,
		"",
		1*time.Minute,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Proxmox client")
	}

	// Authenticate with either APIToken or username/password
	if cfg.APIToken != "" && cfg.TokenID != "" {
		if err := c.SetAPIToken(cfg.TokenID, cfg.APIToken); err != nil {
			return nil, errors.Wrap(err, "failed to set API token")
		}
	} else if cfg.Username != "" && cfg.Password != "" {
		if err := c.Login(cfg.Username, cfg.Password, ""); err != nil {
			return nil, errors.Wrap(err, "failed to login with username and password")
		}
	} else {
		return nil, errors.New("authentication failed: either API token or username/password must be provided")
	}

	// Ensure we have a target node specified in config
	if cfg.TargetNode == "" {
		return nil, errors.New("target Proxmox node must be specified")
	}

	return &Client{
		client:     c,
		config:     cfg,
		nodeName:   nodeName,
		targetNode: cfg.TargetNode,
	}, nil
}

// CreateVM creates a new virtual machine in Proxmox
func (c *Client) CreateVM(vmConfig *config.VMConfig) (int, error) {
	logrus.Infof("Creating VM: %s", vmConfig.Name)
	
	// Validate required fields
	if vmConfig.Template == "" {
		return 0, errors.New("template is required")
	}

	// Find a free VMID if not specified
	vmid := vmConfig.VMID
	if vmid == 0 {
		id, err := c.client.GetNextID(0)
		if err != nil {
			return 0, errors.Wrap(err, "failed to get next VMID")
		}
		vmid = id
		logrus.Infof("Assigned VMID: %d", vmid)
	}

	// Determine if we're cloning from a template or creating from an ISO
	if vmConfig.Template != "" {
		return c.cloneFromTemplate(vmid, vmConfig)
	}

	// If we got here without returning, something went wrong
	return 0, errors.New("unsupported VM creation method")
}

// cloneFromTemplate clones a VM from an existing template
func (c *Client) cloneFromTemplate(vmid int, vmConfig *config.VMConfig) (int, error) {
	logrus.Infof("Cloning from template: %s", vmConfig.Template)

	// Parse template name to get source VMID
	sourceVMID, err := strconv.Atoi(vmConfig.Template)
	if err != nil {
		return 0, errors.Wrap(err, "template must be a valid VMID")
	}

	// Clone parameters
	cloneParams := map[string]interface{}{
		"newid":  vmid,
		"name":   vmConfig.Name,
		"target": c.targetNode,
		"full":   1, // Full clone
	}

	// If storage is specified
	if vmConfig.Storage != "" {
		cloneParams["storage"] = vmConfig.Storage
	}

	// Clone the VM
	taskID, err := c.client.CloneQemuVm(sourceVMID, cloneParams)
	if err != nil {
		return 0, errors.Wrap(err, "failed to clone VM")
	}
	logrus.Infof("Clone task started: %s", taskID)

	// Wait for the clone task to complete
	err = c.waitForTaskCompletion(taskID)
	if err != nil {
		return 0, errors.Wrap(err, "VM clone task failed")
	}

	// Set VM configuration
	err = c.configureVM(vmid, vmConfig)
	if err != nil {
		return 0, errors.Wrap(err, "failed to configure VM")
	}

	// Start VM if required
	if vmConfig.AutoStart {
		logrus.Infof("Starting VM: %d", vmid)
		_, err = c.client.StartVm(vmid)
		if err != nil {
			return vmid, errors.Wrap(err, "failed to start VM")
		}

		// Wait for VM to be ready (get an IP)
		if vmConfig.WaitForIP {
			logrus.Info("Waiting for VM to get an IP address")
			ipAddress, err := c.waitForIP(vmid, 300) // 5 minute timeout
			if err != nil {
				return vmid, errors.Wrap(err, "timeout waiting for VM IP")
			}
			logrus.Infof("VM IP address: %s", ipAddress)
			
			// Add the IP to VM metadata
			vmConfig.IPAddress = ipAddress
		}
	}

	return vmid, nil
}

// configureVM applies configuration to a VM after creation
func (c *Client) configureVM(vmid int, vmConfig *config.VMConfig) error {
	vmParams := map[string]interface{}{}

	// Set CPU and Memory
	if vmConfig.Cores > 0 {
		vmParams["cores"] = vmConfig.Cores
	}
	if vmConfig.Memory > 0 {
		vmParams["memory"] = vmConfig.Memory
	}

	// Network configuration
	if vmConfig.Network != nil {
		for i, net := range vmConfig.Network {
			nicKey := fmt.Sprintf("net%d", i)
			nicValue := fmt.Sprintf("model=%s,bridge=%s", net.Model, net.Bridge)
			if net.VLAN > 0 {
				nicValue += fmt.Sprintf(",tag=%d", net.VLAN)
			}
			vmParams[nicKey] = nicValue
		}
	}

	// CloudInit configurations if specified
	if vmConfig.CloudInit != nil {
		if vmConfig.CloudInit.User != "" {
			vmParams["ciuser"] = vmConfig.CloudInit.User
		}
		if vmConfig.CloudInit.Password != "" {
			vmParams["cipassword"] = vmConfig.CloudInit.Password
		}
		if len(vmConfig.CloudInit.SSHKeys) > 0 {
			vmParams["sshkeys"] = url.QueryEscape(vmConfig.CloudInit.SSHKeys)
		}
		if vmConfig.CloudInit.IPConfig != "" {
			vmParams["ipconfig0"] = vmConfig.CloudInit.IPConfig
		}
	}

	// Description
	if vmConfig.Description != "" {
		vmParams["description"] = vmConfig.Description
	}

	// Apply configuration
	if len(vmParams) > 0 {
		err := c.client.SetVmConfig(vmid, vmParams)
		if err != nil {
			return errors.Wrap(err, "failed to set VM configuration")
		}
		logrus.Infof("VM configured successfully: %d", vmid)
	}

	return nil
}

// waitForTaskCompletion waits for a Proxmox task to complete
func (c *Client) waitForTaskCompletion(taskID string) error {
	timeout := time.After(10 * time.Minute)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeout:
			return errors.New("task timed out")
		case <-tick:
			taskStatus, err := c.client.GetTaskStatus(taskID)
			if err != nil {
				return errors.Wrap(err, "failed to get task status")
			}

			if taskStatus["status"] == "stopped" {
				if taskStatus["exitstatus"] == "OK" {
					return nil // Task completed successfully
				}
				return fmt.Errorf("task failed: %v", taskStatus)
			}
			logrus.Debugf("Task %s status: %s", taskID, taskStatus["status"])
		}
	}
}

// waitForIP waits for the VM to get an IP address
func (c *Client) waitForIP(vmid int, timeoutSeconds int) (string, error) {
	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <-timeout:
			return "", errors.New("timeout waiting for IP address")
		case <-tick:
			// Get VM configuration from Proxmox
			config, err := c.client.GetVmConfig(vmid)
			if err != nil {
				logrus.Warnf("Failed to get VM config: %v", err)
				continue
			}

			// Get the agent network interfaces
			if ipconfig, ok := config["ipconfig0"].(string); ok && ipconfig != "" {
				if ipconfig != "ip=dhcp" {
					// For static IP from config
					return parseIPFromIPConfig(ipconfig), nil
				}
			}

			// Try to get IP from agent
			interfaces, err := c.client.GetVmAgentNetworkIfaces(vmid)
			if err != nil {
				logrus.Debugf("Failed to get network interfaces from agent: %v", err)
				continue
			}

			for _, iface := range interfaces {
				for _, addr := range iface.IPAddresses {
					// Skip IPv6 and loopback
					if addr.IsIPv4() && !addr.IsLoopback() {
						return addr.String(), nil
					}
				}
			}
		}
	}
}

// parseIPFromIPConfig extracts IP from proxmox ipconfig format
func parseIPFromIPConfig(ipconfig string) string {
	// Format is typically "ip=192.168.1.100/24,gw=192.168.1.1"
	// Just a simple parsing for now
	for _, part := range strings.Split(ipconfig, ",") {
		if strings.HasPrefix(part, "ip=") {
			ip := strings.TrimPrefix(part, "ip=")
			if idx := strings.Index(ip, "/"); idx > 0 {
				return ip[:idx]
			}
			return ip
		}
	}
	return ""
}

// DeleteVM deletes a virtual machine
func (c *Client) DeleteVM(vmid int) error {
	logrus.Infof("Deleting VM: %d", vmid)

	// First, stop the VM if it's running
	status, err := c.client.GetVmStatus(vmid)
	if err != nil {
		return errors.Wrap(err, "failed to get VM status")
	}

	if status["status"] == "running" {
		logrus.Infof("Stopping VM %d before deletion", vmid)
		stopTask, err := c.client.StopVm(vmid)
		if err != nil {
			return errors.Wrap(err, "failed to stop VM")
		}

		// Wait for the VM to stop
		if err := c.waitForTaskCompletion(stopTask); err != nil {
			return errors.Wrap(err, "failed to stop VM before deletion")
		}
	}

	// Now delete the VM
	deleteTask, err := c.client.DeleteVm(vmid)
	if err != nil {
		return errors.Wrap(err, "failed to delete VM")
	}

	// Wait for deletion to complete
	return c.waitForTaskCompletion(deleteTask)
}

// StartVM starts a virtual machine
func (c *Client) StartVM(vmid int) error {
	logrus.Infof("Starting VM: %d", vmid)
	
	task, err := c.client.StartVm(vmid)
	if err != nil {
		return errors.Wrap(err, "failed to start VM")
	}
	
	return c.waitForTaskCompletion(task)
}

// StopVM stops a virtual machine
func (c *Client) StopVM(vmid int) error {
	logrus.Infof("Stopping VM: %d", vmid)
	
	task, err := c.client.StopVm(vmid)
	if err != nil {
		return errors.Wrap(err, "failed to stop VM")
	}
	
	return c.waitForTaskCompletion(task)
}