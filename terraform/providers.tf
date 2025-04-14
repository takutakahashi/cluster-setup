# Default Proxmox Provider (for backward compatibility)
provider "proxmox" {
  pm_api_url      = length(var.proxmox_nodes_config) > 0 ? try(values(var.proxmox_nodes_config)[0].api_url, var.proxmox_api_url) : var.proxmox_api_url
  pm_api_token_id = length(var.proxmox_nodes_config) > 0 ? try(values(var.proxmox_nodes_config)[0].api_token_id, var.proxmox_api_token_id) : var.proxmox_api_token_id
  pm_api_token_secret = length(var.proxmox_nodes_config) > 0 ? try(values(var.proxmox_nodes_config)[0].api_token_secret, var.proxmox_api_token_secret) : var.proxmox_api_token_secret
  pm_tls_insecure = length(var.proxmox_nodes_config) > 0 ? try(values(var.proxmox_nodes_config)[0].tls_insecure, var.proxmox_tls_insecure) : var.proxmox_tls_insecure
  
  # Debug options
  pm_debug      = var.proxmox_debug
  pm_log_enable = var.proxmox_log_enable
  pm_log_file   = var.proxmox_log_file
}

# Node-specific Proxmox Providers
# Node: node-1
provider "proxmox" {
  alias           = "node-1"
  pm_api_url      = try(var.proxmox_nodes_config["node-1"].api_url, var.proxmox_api_url)
  pm_api_token_id = try(var.proxmox_nodes_config["node-1"].api_token_id, var.proxmox_api_token_id)
  pm_api_token_secret = try(var.proxmox_nodes_config["node-1"].api_token_secret, var.proxmox_api_token_secret)
  pm_tls_insecure = try(var.proxmox_nodes_config["node-1"].tls_insecure, var.proxmox_tls_insecure)
  
  # Debug options
  pm_debug      = var.proxmox_debug
  pm_log_enable = var.proxmox_log_enable
  pm_log_file   = var.proxmox_log_file
}

# Node: node-2
provider "proxmox" {
  alias           = "node-2"
  pm_api_url      = try(var.proxmox_nodes_config["node-2"].api_url, var.proxmox_api_url)
  pm_api_token_id = try(var.proxmox_nodes_config["node-2"].api_token_id, var.proxmox_api_token_id)
  pm_api_token_secret = try(var.proxmox_nodes_config["node-2"].api_token_secret, var.proxmox_api_token_secret)
  pm_tls_insecure = try(var.proxmox_nodes_config["node-2"].tls_insecure, var.proxmox_tls_insecure)
  
  # Debug options
  pm_debug      = var.proxmox_debug
  pm_log_enable = var.proxmox_log_enable
  pm_log_file   = var.proxmox_log_file
}

# Node: node-3
provider "proxmox" {
  alias           = "node-3"
  pm_api_url      = try(var.proxmox_nodes_config["node-3"].api_url, var.proxmox_api_url)
  pm_api_token_id = try(var.proxmox_nodes_config["node-3"].api_token_id, var.proxmox_api_token_id)
  pm_api_token_secret = try(var.proxmox_nodes_config["node-3"].api_token_secret, var.proxmox_api_token_secret)
  pm_tls_insecure = try(var.proxmox_nodes_config["node-3"].tls_insecure, var.proxmox_tls_insecure)
  
  # Debug options
  pm_debug      = var.proxmox_debug
  pm_log_enable = var.proxmox_log_enable
  pm_log_file   = var.proxmox_log_file
}