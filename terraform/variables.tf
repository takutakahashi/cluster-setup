# Proxmox Provider Variables - Default configuration
variable "proxmox_api_url" {
  description = "The default Proxmox API URL (used if node-specific URL is not provided)"
  type        = string
  default     = ""
}

variable "proxmox_user" {
  description = "Default Proxmox user with appropriate permissions (used if node-specific user is not provided)"
  type        = string
  default     = "root@pam"
}

variable "proxmox_password" {
  description = "Default password for Proxmox user (used if node-specific password is not provided)"
  type        = string
  sensitive   = true
  default     = null
}

variable "proxmox_tls_insecure" {
  description = "Skip TLS verification (true for self-signed certificates)"
  type        = bool
  default     = true
}

variable "proxmox_debug" {
  description = "Enable debug logging for Proxmox Provider"
  type        = bool
  default     = false
}

variable "proxmox_log_enable" {
  description = "Enable logging to file for Proxmox Provider"
  type        = bool
  default     = false
}

variable "proxmox_log_file" {
  description = "Log file path for Proxmox Provider"
  type        = string
  default     = "terraform-plugin-proxmox.log"
}

# Node-specific Proxmox Configuration
variable "proxmox_nodes_config" {
  description = "Configuration for Proxmox nodes with their API URLs and credentials"
  type = map(object({
    api_url      = string
    user         = optional(string)
    password     = optional(string)
    tls_insecure = optional(bool)
  }))
  default = {}
}

# Backward compatibility - List of node names
variable "proxmox_nodes" {
  description = "List of Proxmox nodes (for backward compatibility)"
  type        = list(string)
  default     = ["node-1", "node-2", "node-3"]
}

# VM Template Configuration
variable "template_name" {
  description = "Name of the VM template to clone"
  type        = string
  default     = "ubuntu-cloud-template"
}

variable "template_storage" {
  description = "Storage for VM template"
  type        = string
  default     = "local"
}

# VM Configuration
variable "vms" {
  description = "Map of VM configurations"
  type = map(object({
    node         = string
    cores        = number
    memory       = number
    disk_size    = string
    storage_pool = string
    network_bridge = string
    ip_address   = string
    gateway      = string
  }))
  default = {}
}

# Cloud-Init Configuration
variable "ssh_public_keys" {
  description = "List of SSH public keys to add to VMs"
  type        = list(string)
  default     = []
}

variable "cloud_init_user" {
  description = "Username for Cloud-Init"
  type        = string
  default     = "ubuntu"
}

variable "cloud_init_password" {
  description = "Password for Cloud-Init user (not recommended, use SSH keys instead)"
  type        = string
  default     = null
  sensitive   = true
}

# Custom Script Configuration
variable "custom_script" {
  description = "Custom script content to run during boot via cloud-init"
  type        = string
  default     = ""
}