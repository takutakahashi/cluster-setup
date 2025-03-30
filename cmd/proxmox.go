package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/takutakahashi/cluster-setup/pkg/config"
	"github.com/takutakahashi/cluster-setup/pkg/deploy"
	"github.com/takutakahashi/cluster-setup/pkg/provider/proxmox"
)

// proxmoxCmd represents the proxmox command
var proxmoxCmd = &cobra.Command{
	Use:   "proxmox",
	Short: "Manage Proxmox virtual machines",
	Long: `Manage Proxmox virtual machines for your Kubernetes cluster. 
	
This command allows you to create, delete, and manage VMs in your Proxmox environment
based on the configuration in your config file.`,
}

// createCmd represents the proxmox create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create VMs in Proxmox",
	Long: `Create virtual machines in Proxmox based on your configuration.
	
This will create VMs as specified in the proxmoxNodes section of your config file.
Once VMs are created, their IP addresses will be added to the targets list of the 
corresponding nodes in the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		cp, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}
		
		cfg, err := config.Load(cp)
		if err != nil {
			logrus.Fatal(err)
		}
		
		// Create Proxmox manager
		manager, err := proxmox.NewManager(cfg)
		if err != nil {
			logrus.Fatalf("Failed to create Proxmox manager: %v", err)
		}
		
		// Create nodes
		if err := manager.CreateNodes(); err != nil {
			logrus.Fatalf("Failed to create nodes: %v", err)
		}
		
		logrus.Info("VM creation completed")
	},
}

// deleteCmd represents the proxmox delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete VMs in Proxmox",
	Long: `Delete virtual machines from Proxmox based on your configuration.
	
This will delete VMs as specified in the proxmoxNodes section of your config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cp, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}
		
		cfg, err := config.Load(cp)
		if err != nil {
			logrus.Fatal(err)
		}
		
		// Create Proxmox manager
		manager, err := proxmox.NewManager(cfg)
		if err != nil {
			logrus.Fatalf("Failed to create Proxmox manager: %v", err)
		}
		
		// Delete nodes
		if err := manager.DeleteNodes(); err != nil {
			logrus.Fatalf("Failed to delete nodes: %v", err)
		}
		
		logrus.Info("VM deletion completed")
	},
}

// updateCmd represents the proxmox update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update VMs in Proxmox",
	Long: `Update virtual machines in Proxmox based on your configuration.
	
This will update VM configurations as specified in the proxmoxNodes section of your config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cp, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}
		
		cfg, err := config.Load(cp)
		if err != nil {
			logrus.Fatal(err)
		}
		
		// Create Proxmox manager
		manager, err := proxmox.NewManager(cfg)
		if err != nil {
			logrus.Fatalf("Failed to create Proxmox manager: %v", err)
		}
		
		// Update nodes
		if err := manager.UpdateNodes(); err != nil {
			logrus.Fatalf("Failed to update nodes: %v", err)
		}
		
		logrus.Info("VM update completed")
	},
}

// deployCmd represents the proxmox deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create VMs and deploy the cluster",
	Long: `Create VMs in Proxmox and deploy the Kubernetes cluster on them.
	
This command combines the 'create' and the default cluster deployment into one step.
It will:
1. Create VMs in Proxmox according to your configuration
2. Wait for VMs to be ready
3. Deploy the Kubernetes cluster to these VMs`,
	Run: func(cmd *cobra.Command, args []string) {
		cp, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}
		
		local, err := cmd.Flags().GetBool("output-local")
		if err != nil {
			logrus.Fatal(err)
		}
		
		cfg, err := config.Load(cp)
		if err != nil {
			logrus.Fatal(err)
		}
		
		// Create Proxmox manager
		manager, err := proxmox.NewManager(cfg)
		if err != nil {
			logrus.Fatalf("Failed to create Proxmox manager: %v", err)
		}
		
		// Create nodes
		logrus.Info("Creating VMs in Proxmox...")
		if err := manager.CreateNodes(); err != nil {
			logrus.Fatalf("Failed to create nodes: %v", err)
		}
		
		// Deploy the cluster
		logrus.Info("Deploying Kubernetes cluster to VMs...")
		if err := executeDeployment(cfg, local); err != nil {
			logrus.Fatalf("Failed to deploy cluster: %v", err)
		}
		
		logrus.Info("Cluster deployment to Proxmox VMs completed")
	},
}

func init() {
	rootCmd.AddCommand(proxmoxCmd)
	proxmoxCmd.AddCommand(createCmd)
	proxmoxCmd.AddCommand(deleteCmd)
	proxmoxCmd.AddCommand(updateCmd)
	proxmoxCmd.AddCommand(deployCmd)
	
	// Add the flags to all proxmox subcommands
	for _, cmd := range []*cobra.Command{createCmd, deleteCmd, updateCmd, deployCmd} {
		cmd.Flags().StringP("config", "c", "", "config path")
		cmd.MarkFlagRequired("config")
	}
	
	// Add local flag to deploy command
	deployCmd.Flags().BoolP("output-local", "l", false, "output assets dist to local")
}

// executeDeployment is a helper function to execute the deployment logic
func executeDeployment(cfg *config.Config, local bool) error {
	// Reusing the deployment logic from the root command
	return deploy.Execute(cfg, local)
}