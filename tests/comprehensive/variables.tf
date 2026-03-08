# Variables for test infrastructure

variable "test_prefix" {
  description = "Prefix for all test resources"
  type        = string
  default     = "tf-test"
}

variable "api_url" {
  description = "CloudStack API URL"
  type        = string
}

variable "api_key" {
  description = "CloudStack API key"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "CloudStack API secret key"
  type        = string
  sensitive   = true
}

variable "zone" {
  description = "CloudStack zone name"
  type        = string
}

variable "network_offering" {
  description = "Network offering name"
  type        = string
  default     = "DefaultIsolatedNetworkSourceNatService"
}

variable "service_offering" {
  description = "VM service offering name"
  type        = string
  default     = "Small Instance"
}

variable "disk_offering" {
  description = "Disk offering for data volumes"
  type        = string
  default     = "Small"
}

variable "template" {
  description = "Template name for VMs"
  type        = string
  default     = "ubuntu-24.04-lts"
}

variable "disk_offering_id" {
  description = "Disk offering ID for data volumes (use data source or CloudStack API to find)"
  type        = string
}

variable "zone_id" {
  description = "Zone ID (use data source or CloudStack API to find)"
  type        = string
}
