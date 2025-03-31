locals {
  # Combine VM details from all nodes
  all_vms = merge(
    { for name, vm in proxmox_vm_qemu.node1_vm : name => {
      id          = vm.id
      name        = vm.name
      target_node = vm.target_node
      ip_address  = split("/", split("=", vm.ipconfig0)[1])[0]
    }},
    { for name, vm in proxmox_vm_qemu.node2_vm : name => {
      id          = vm.id
      name        = vm.name
      target_node = vm.target_node
      ip_address  = split("/", split("=", vm.ipconfig0)[1])[0]
    }},
    { for name, vm in proxmox_vm_qemu.node3_vm : name => {
      id          = vm.id
      name        = vm.name
      target_node = vm.target_node
      ip_address  = split("/", split("=", vm.ipconfig0)[1])[0]
    }}
  )
}

output "vm_details" {
  description = "Details of created VMs"
  value       = local.all_vms
}

output "vm_ips" {
  description = "IP addresses of created VMs"
  value = {
    for name, vm in local.all_vms : name => vm.ip_address
  }
}