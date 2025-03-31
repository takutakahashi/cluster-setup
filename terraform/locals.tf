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
}