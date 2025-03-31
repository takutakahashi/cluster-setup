locals {
  # Create a map where keys are VM names and values are node names
  # This distributes VMs evenly across the available nodes
  vm_node_distribution = {
    for idx, vm_name in keys(var.vms) :
    vm_name => lookup(var.vms[vm_name], "node", var.proxmox_nodes[idx % length(var.proxmox_nodes)])
  }

  # Combine SSH public keys into a single string
  ssh_keys_combined = join("\n", var.ssh_public_keys)
  
  # Determine if a custom script is provided
  has_custom_script = var.custom_script != ""
  
  # Get all node names from proxmox_nodes_config or fallback to proxmox_nodes
  all_node_names = length(keys(var.proxmox_nodes_config)) > 0 ? keys(var.proxmox_nodes_config) : var.proxmox_nodes

  # Group VMs by node name
  vms_by_node = {
    for node in local.all_node_names : node => {
      for vm_name, vm in var.vms : vm_name => vm
      if lookup(vm, "node", local.vm_node_distribution[vm_name]) == node
    }
  }
}