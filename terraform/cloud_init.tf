# Create cloud-init custom configuration files when a custom script is provided
resource "local_file" "cloud_init_custom_config" {
  for_each = local.has_custom_script ? var.vms : {}

  content = <<-EOT
#cloud-config
runcmd:
  - |
    cat > /tmp/custom-script.sh << 'SCRIPT'
${var.custom_script}
SCRIPT
  - chmod +x /tmp/custom-script.sh
  - /tmp/custom-script.sh
  - rm -f /tmp/custom-script.sh
EOT

  filename = "${path.module}/cloud-init/${each.key}-cloud-init.yml"
}

# If custom cloud-init files are generated, we need to
# upload them to Proxmox's storage location for snippets
resource "null_resource" "upload_cloud_init_files" {
  for_each = local.has_custom_script ? var.vms : {}

  depends_on = [local_file.cloud_init_custom_config]

  # Only execute when the content changes
  triggers = {
    cloud_init_content = local.has_custom_script ? var.custom_script : ""
    node_name = lookup(each.value, "node", local.vm_node_distribution[each.key])
  }

  # For each node, upload the cloud-init file
  provisioner "local-exec" {
    # This assumes SSH access to Proxmox nodes is configured
    # Adjust the command as needed for your environment
    command = <<-EOT
      scp ${path.module}/cloud-init/${each.key}-cloud-init.yml root@${lookup(each.value, "node", local.vm_node_distribution[each.key])}:/var/lib/vz/snippets/
    EOT
  }
}