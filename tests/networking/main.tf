# Test configuration for CloudStack networking with egress rules
# This demonstrates the missing piece for public SSH access: EGRESS RULES

terraform {
  required_providers {
    cloudstack = {
      source  = "cloudstack/cloudstack"
      version = ">= 0.5.0"
    }
  }
}

provider "cloudstack" {
  api_url    = var.api_url
  api_key    = var.api_key
  secret_key = var.secret_key
}

# Create a network for the VM
resource "cloudstack_network" "test_network" {
  name             = "test-egress-network"
  display_text     = "Test network for egress rules verification"
  cidr             = "10.10.10.0/24"
  network_offering = var.network_offering
  zone             = var.zone
}

# Create security group with BOTH ingress and egress rules
# CloudStack is secure-by-default and blocks egress
# Without egress rules, VMs cannot reach the internet
resource "cloudstack_security_group" "test_sg" {
  name        = "test-egress-sg"
  description = "Security group with SSH ingress and full egress access"
}

# SSH ingress rule - allows incoming SSH connections
resource "cloudstack_security_group_rule" "ssh_ingress" {
  security_group_id = cloudstack_security_group.test_sg.id

  rule {
    cidr_list   = ["0.0.0.0/0"]
    protocol    = "tcp"
    ports       = ["22"]
    traffic_type = "ingress"
  }
}

# Egress rules for internet access
# Without these, the VM cannot reach the internet (apt update, curl, etc.)
resource "cloudstack_security_group_rule" "egress_all" {
  security_group_id = cloudstack_security_group.test_sg.id

  rule {
    cidr_list   = ["0.0.0.0/0"]
    protocol    = "tcp"
    ports       = ["80", "443"]
    traffic_type = "egress"
  }
}

# Also allow DNS egress (UDP)
resource "cloudstack_security_group_rule" "dns_egress" {
  security_group_id = cloudstack_security_group.test_sg.id

  rule {
    cidr_list   = ["0.0.0.0/0"]
    protocol    = "udp"
    ports       = ["53"]
    traffic_type = "egress"
  }
}

# Acquire a public IP address
resource "cloudstack_ipaddress" "test_ip" {
  network_id = cloudstack_network.test_network.id
  zone       = var.zone
}

# Create port forwarding rule for SSH access
resource "cloudstack_port_forward" "ssh_forward" {
  ip_address_id = cloudstack_ipaddress.test_ip.id

  forward {
    protocol           = "tcp"
    private_port       = 22
    public_port        = 22
    virtual_machine_id = cloudstack_instance.test_vm.id
    vm_guest_ip        = cloudstack_instance.test_vm.ip_address
  }
}

# Create the VM instance with security group
resource "cloudstack_instance" "test_vm" {
  name             = "test-egress-vm"
  display_name     = "Test VM for egress verification"
  service_offering = var.service_offering
  template         = var.template
  zone             = var.zone
  network_id       = cloudstack_network.test_network.id

  # Attach security group with egress rules
  security_group_names = [cloudstack_security_group.test_sg.name]

  # Use a small root disk (20GB is sufficient for testing)
  root_disk_size = var.root_disk_size

  # SSH keypair for access
  keypair = var.keypair

  # Cloud-init user data for initial setup
  user_data = base64encode(<<-EOT
    #cloud-config
    package_update: true
    packages:
      - curl
    runcmd:
      - echo "Egress test VM initialized at $(date)" >> /var/log/egress-test.log
  EOT
  )

  # Ensure VM starts
  start_vm = true
}

# Outputs for verification
output "public_ip" {
  value       = cloudstack_ipaddress.test_ip.ip_address
  description = "Public IP address for SSH access"
}

output "ssh_command" {
  value       = "ssh -i ~/.ssh/${var.keypair} ubuntu@${cloudstack_ipaddress.test_ip.ip_address}"
  description = "SSH command to connect to the VM"
}

output "vm_private_ip" {
  value       = cloudstack_instance.test_vm.ip_address
  description = "VM's private IP address"
}

output "network_id" {
  value       = cloudstack_network.test_network.id
  description = "Network ID"
}

output "security_group_id" {
  value       = cloudstack_security_group.test_sg.id
  description = "Security group ID with egress rules"
}

output "egress_test_command" {
  value       = "Run: curl -s https://ifconfig.me (should return public IP if egress works)"
  description = "Command to test egress connectivity from inside VM"
}
