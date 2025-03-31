provider "proxmox" {
  pm_api_url      = var.proxmox_api_url
  pm_user         = var.proxmox_user
  pm_password     = var.proxmox_password
  pm_tls_insecure = var.proxmox_tls_insecure
  
  # Debug options
  pm_debug = var.proxmox_debug
  pm_log_enable = var.proxmox_log_enable
  pm_log_file   = var.proxmox_log_file
}