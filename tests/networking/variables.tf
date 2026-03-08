# Variables for CloudStack networking test configuration

variable "api_url" {
  type        = string
  description = "CloudStack API URL"
}

variable "api_key" {
  type        = string
  description = "CloudStack API key"
  sensitive   = true
}

variable "secret_key" {
  type        = string
  description = "CloudStack secret key"
  sensitive   = true
}

variable "zone" {
  type        = string
  description = "CloudStack zone name"
}

variable "network_offering" {
  type        = string
  description = "Network offering for isolated networks"
  default     = "DefaultIsolatedNetworkOfferingWithSourceNatService"
}

variable "service_offering" {
  type        = string
  description = "Service offering (instance size)"
  default     = "Small Instance"
}

variable "template" {
  type        = string
  description = "Template name for the VM"
  default     = "ubuntu-24.04-lts"
}

variable "keypair" {
  type        = string
  description = "SSH keypair name registered in CloudStack"
}

variable "root_disk_size" {
  type        = number
  description = "Root disk size in GB"
  default     = 20
}
