output "vm_details" {
  description = "Details of created VMs"
  value = {
    for name, vm in proxmox_vm_qemu.vm : name => {
      id          = vm.id
      name        = vm.name
      target_node = vm.target_node
      ip_address  = split("/", split("=", vm.ipconfig0)[1])[0]
    }
  }
}

output "vm_ips" {
  description = "IP addresses of created VMs"
  value = {
    for name, vm in proxmox_vm_qemu.vm : name => split("/", split("=", vm.ipconfig0)[1])[0]
  }
}