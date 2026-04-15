---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_ipaddress"
sidebar_current: "docs-cloudstack-resource-ipaddress"
description: |-
  Acquires and associates a public IP.
---

# cloudstack_ipaddress

Acquires and associates a public IP.

## Example Usage

### Basic IP Address for Network

```hcl
resource "cloudstack_ipaddress" "default" {
  network_id = "6eb22f91-7454-4107-89f4-36afcdf33021"
}
```

### IP Address for VPC

```hcl
resource "cloudstack_vpc" "foo" {
  name         = "my-vpc"
  cidr         = "10.0.0.0/16"
  vpc_offering = "Default VPC offering"
  zone         = "zone-1"
}

resource "cloudstack_ipaddress" "vpc_ip" {
  vpc_id = cloudstack_vpc.foo.id
}
```

### IP Address with Automatic Project Inheritance

```hcl
# Create a VPC in a project
resource "cloudstack_vpc" "project_vpc" {
  name         = "project-vpc"
  cidr         = "10.0.0.0/16"
  vpc_offering = "Default VPC offering"
  project      = "my-project"
  zone         = "zone-1"
}

# IP address automatically inherits project from VPC
resource "cloudstack_ipaddress" "vpc_ip" {
  vpc_id = cloudstack_vpc.project_vpc.id
  # project is automatically inherited from the VPC
}

# Or with a network
resource "cloudstack_network" "project_network" {
  name             = "project-network"
  cidr             = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  project          = "my-project"
  zone             = "zone-1"
}

# IP address automatically inherits project from network
resource "cloudstack_ipaddress" "network_ip" {
  network_id = cloudstack_network.project_network.id
  # project is automatically inherited from the network
}
```

## Argument Reference

The following arguments are supported:

* `is_portable` - (Optional) This determines if the IP address should be transferable
    across zones (defaults false)

* `network_id` - (Optional) The ID of the network for which an IP address should
    be acquired and associated. Changing this forces a new resource to be created.

* `vpc_id` - (Optional) The ID of the VPC for which an IP address should be
   acquired and associated. Changing this forces a new resource to be created.

* `zone` - (Optional) The name or ID of the zone for which an IP address should be
   acquired and associated. Changing this forces a new resource to be created.

* `project` - (Optional) The name or ID of the project to deploy this
    IP address to. Changing this forces a new resource to be created. If not
    specified and `vpc_id` is provided, the project will be automatically
    inherited from the VPC. If not specified and `network_id` is provided,
    the project will be automatically inherited from the network.

*NOTE: `network_id` and/or `zone` should have a value when `is_portable` is `false`!*
*NOTE: Either `network_id` or `vpc_id` should have a value when `is_portable` is `true`!*

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the acquired and associated IP address.
* `ip_address` - The IP address that was acquired and associated.
