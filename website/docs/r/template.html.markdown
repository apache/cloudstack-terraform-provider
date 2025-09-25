---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_template"
sidebar_current: "docs-cloudstack-resource-template"
description: |-
  Registers a template into the CloudStack cloud, including support for CloudStack Kubernetes Service (CKS) templates.
---

# cloudstack_template

Registers a template into the CloudStack cloud. This resource supports both regular VM templates and specialized templates for CloudStack Kubernetes Service (CKS) clusters.

## Example Usage

### Basic Template

```hcl
resource "cloudstack_template" "centos64" {
  name       = "CentOS 6.4 x64"
  format     = "VHD"
  hypervisor = "XenServer"
  os_type    = "CentOS 6.4 (64bit)"
  url        = "http://example.com/template.vhd"
  zone       = "zone-1"
}
```

### CKS Template for Kubernetes

```hcl
resource "cloudstack_template" "cks_ubuntu_template" {
  name         = "cks-ubuntu-2204-template"
  display_text = "CKS Ubuntu 22.04 Template for Kubernetes"
  url          = "http://example.com/cks-ubuntu-2204-kvm.qcow2.bz2"
  format       = "QCOW2"
  hypervisor   = "KVM"
  os_type      = "Ubuntu 22.04 LTS"
  zone         = "zone1"

  # CKS specific flag
  for_cks = true

  # Template properties
  is_extractable          = false
  is_featured             = false
  is_public               = false
  password_enabled        = true
  is_dynamically_scalable = true

  # Wait for template to be ready
  is_ready_timeout = 1800

  tags = {
    Environment = "CKS"
    Purpose     = "Kubernetes"
    OS          = "Ubuntu-22.04"
  }
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the template.
* `format` - (Required) The format of the template. Valid values are `QCOW2`, `RAW`, `VHD`, `OVA`, and `ISO`.
* `hypervisor` - (Required) The target hypervisor for the template. Valid values include `KVM`, `XenServer`, `VMware`, `Hyperv`, and `LXC`. Changing this forces a new resource to be created.
* `os_type` - (Required) The OS Type that best represents the OS of this template.
* `url` - (Required) The URL of where the template is hosted. Changing this forces a new resource to be created.

### Optional Arguments

* `display_text` - (Optional) The display name of the template. If not specified, defaults to the `name`.
* `zone` - (Optional) The name or ID of the zone where this template will be created. Changing this forces a new resource to be created.
* `project` - (Optional) The name or ID of the project to create this template for. Changing this forces a new resource to be created.
* `account` - (Optional) The account name for the template.
* `domain_id` - (Optional) The domain ID for the template.

### CKS-Specific Arguments

* `for_cks` - (Optional) Set to `true` to indicate this template is for CloudStack Kubernetes Service (CKS). CKS templates have special requirements and capabilities. Defaults to `false`.

### Template Properties

* `is_dynamically_scalable` - (Optional) Set to indicate if the template contains tools to support dynamic scaling of VM cpu/memory. Defaults to `false`.
* `is_extractable` - (Optional) Set to indicate if the template is extractable. Defaults to `false`.
* `is_featured` - (Optional) Set to indicate if the template is featured. Defaults to `false`.
* `is_public` - (Optional) Set to indicate if the template is available for all accounts. Defaults to `true`.
* `password_enabled` - (Optional) Set to indicate if the template should be password enabled. Defaults to `false`.
* `sshkey_enabled` - (Optional) Set to indicate if the template supports SSH key injection. Defaults to `false`.
* `is_ready_timeout` - (Optional) The maximum time in seconds to wait until the template is ready for use. Defaults to `300` seconds.

### Metadata and Tagging

* `tags` - (Optional) A mapping of tags to assign to the template.

* `is_dynamically_scalable` - (Optional) Set to indicate if the template contains
    tools to support dynamic scaling of VM cpu/memory (defaults false)

* `is_extractable` - (Optional) Set to indicate if the template is extractable
    (defaults false)

* `is_featured` - (Optional) Set to indicate if the template is featured
    (defaults false)

* `is_public` - (Optional) Set to indicate if the template is available for
    all accounts (defaults true)

* `password_enabled` - (Optional) Set to indicate if the template should be
    password enabled (defaults false)

* `is_ready_timeout` - (Optional) The maximum time in seconds to wait until the
    template is ready for use (defaults 300 seconds)

## Attributes Reference

The following attributes are exported:

* `id` - The template ID.
* `display_text` - The display text of the template.
* `is_dynamically_scalable` - Set to "true" if the template is dynamically scalable.
* `is_extractable` - Set to "true" if the template is extractable.
* `is_featured` - Set to "true" if the template is featured.
* `is_public` - Set to "true" if the template is public.
* `password_enabled` - Set to "true" if the template is password enabled.
* `is_ready` - Set to `true` once the template is ready for use.
* `created` - The timestamp when the template was created.
* `size` - The size of the template in bytes.
* `checksum` - The checksum of the template.
* `status` - The current status of the template.
* `zone_id` - The zone ID where the template is registered.
* `zone_name` - The zone name where the template is registered.
* `account` - The account name owning the template.
* `domain` - The domain name where the template belongs.
* `project` - The project name if the template is assigned to a project.

### Example CKS Template Usage

```hcl
# Data source to use existing CKS template
data "cloudstack_template" "cks_template" {
  template_filter = "executable"

  filter {
    name  = "name" 
    value = "cks-ubuntu-2204-template"
  }
}

# Use in Kubernetes cluster
resource "cloudstack_kubernetes_cluster" "example" {
  name               = "example-cluster"
  zone               = "zone1"
  kubernetes_version = "1.25.0"
  service_offering   = "Medium Instance"
  
  node_templates = {
    "control" = data.cloudstack_template.cks_template.name
    "worker"  = data.cloudstack_template.cks_template.name
  }
}
```
