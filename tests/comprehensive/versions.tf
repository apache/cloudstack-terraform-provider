terraform {
  required_providers {
    cloudstack = {
      source  = "local/cloudstack/cloudstack"
      version = "0.5.0"
    }
  }
  required_version = ">= 1.0.0"
}
