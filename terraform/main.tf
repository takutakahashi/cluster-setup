# VM Resource for node-1
resource "proxmox_vm_qemu" "node1_vm" {
  provider = proxmox.node-1
  for_each = local.vms_by_node["node-1"]

  name        = each.key
  target_node = "node-1"
  vmid        = null  # Automatically assigned if not specified
  
  # Clone from template
  clone        = var.template_name
  full_clone   = true
  
  # VM specifications
  cores        = lookup(each.value, "cores", 2)
  sockets      = 1
  memory       = lookup(each.value, "memory", 2048)
  agent        = 1  # QEMU Guest Agent
  
  # Disk configuration
  disk {
    type     = "scsi"
    storage  = lookup(each.value, "storage_pool", "local-lvm")
    size     = lookup(each.value, "disk_size", "20G")
    iothread = 1
  }
  
  # Network configuration
  network {
    model   = "virtio"
    bridge  = lookup(each.value, "network_bridge", "vmbr0")
    tag     = lookup(each.value, "vlan", -1) != -1 ? lookup(each.value, "vlan", null) : null
  }
  
  # Cloud-init configuration
  os_type = "cloud-init"
  
  # IP configuration
  ipconfig0 = "ip=${lookup(each.value, "ip_address", "dhcp")},gw=${lookup(each.value, "gateway", "")}"
  
  # SSH keys and user
  ciuser     = var.cloud_init_user
  cipassword = var.cloud_init_password
  sshkeys    = local.ssh_keys_combined
  
  # Custom cloud-init configuration
  cicustom = local.has_custom_script ? "user=local:snippets/${each.key}-cloud-init.yml" : null
  
  # Ensure cloud-init custom file is created before VM
  depends_on = [
    local_file.cloud_init_custom_config
  ]
  
  # Wait for VM to be created and booted
  define_connection_info = true
  
  # Set reasonable pool
  pool = "terraform"
  
  # Set reasonable timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
  
  # Ensure VM always gets recreated instead of reused
  lifecycle {
    create_before_destroy = true
  }
}

# VM Resource for node-2
resource "proxmox_vm_qemu" "node2_vm" {
  provider = proxmox.node-2
  for_each = local.vms_by_node["node-2"]

  name        = each.key
  target_node = "node-2"
  vmid        = null  # Automatically assigned if not specified
  
  # Clone from template
  clone        = var.template_name
  full_clone   = true
  
  # VM specifications
  cores        = lookup(each.value, "cores", 2)
  sockets      = 1
  memory       = lookup(each.value, "memory", 2048)
  agent        = 1  # QEMU Guest Agent
  
  # Disk configuration
  disk {
    type     = "scsi"
    storage  = lookup(each.value, "storage_pool", "local-lvm")
    size     = lookup(each.value, "disk_size", "20G")
    iothread = 1
  }
  
  # Network configuration
  network {
    model   = "virtio"
    bridge  = lookup(each.value, "network_bridge", "vmbr0")
    tag     = lookup(each.value, "vlan", -1) != -1 ? lookup(each.value, "vlan", null) : null
  }
  
  # Cloud-init configuration
  os_type = "cloud-init"
  
  # IP configuration
  ipconfig0 = "ip=${lookup(each.value, "ip_address", "dhcp")},gw=${lookup(each.value, "gateway", "")}"
  
  # SSH keys and user
  ciuser     = var.cloud_init_user
  cipassword = var.cloud_init_password
  sshkeys    = local.ssh_keys_combined
  
  # Custom cloud-init configuration
  cicustom = local.has_custom_script ? "user=local:snippets/${each.key}-cloud-init.yml" : null
  
  # Ensure cloud-init custom file is created before VM
  depends_on = [
    local_file.cloud_init_custom_config
  ]
  
  # Wait for VM to be created and booted
  define_connection_info = true
  
  # Set reasonable pool
  pool = "terraform"
  
  # Set reasonable timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
  
  # Ensure VM always gets recreated instead of reused
  lifecycle {
    create_before_destroy = true
  }
}

# VM Resource for node-3
resource "proxmox_vm_qemu" "node3_vm" {
  provider = proxmox.node-3
  for_each = local.vms_by_node["node-3"]

  name        = each.key
  target_node = "node-3"
  vmid        = null  # Automatically assigned if not specified
  
  # Clone from template
  clone        = var.template_name
  full_clone   = true
  
  # VM specifications
  cores        = lookup(each.value, "cores", 2)
  sockets      = 1
  memory       = lookup(each.value, "memory", 2048)
  agent        = 1  # QEMU Guest Agent
  
  # Disk configuration
  disk {
    type     = "scsi"
    storage  = lookup(each.value, "storage_pool", "local-lvm")
    size     = lookup(each.value, "disk_size", "20G")
    iothread = 1
  }
  
  # Network configuration
  network {
    model   = "virtio"
    bridge  = lookup(each.value, "network_bridge", "vmbr0")
    tag     = lookup(each.value, "vlan", -1) != -1 ? lookup(each.value, "vlan", null) : null
  }
  
  # Cloud-init configuration
  os_type = "cloud-init"
  
  # IP configuration
  ipconfig0 = "ip=${lookup(each.value, "ip_address", "dhcp")},gw=${lookup(each.value, "gateway", "")}"
  
  # SSH keys and user
  ciuser     = var.cloud_init_user
  cipassword = var.cloud_init_password
  sshkeys    = local.ssh_keys_combined
  
  # Custom cloud-init configuration
  cicustom = local.has_custom_script ? "user=local:snippets/${each.key}-cloud-init.yml" : null
  
  # Ensure cloud-init custom file is created before VM
  depends_on = [
    local_file.cloud_init_custom_config
  ]
  
  # Wait for VM to be created and booted
  define_connection_info = true
  
  # Set reasonable pool
  pool = "terraform"
  
  # Set reasonable timeouts
  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
  
  # Ensure VM always gets recreated instead of reused
  lifecycle {
    create_before_destroy = true
  }
}