---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_network_service_provider_state"
sidebar_current: "docs-cloudstack-resource-network_service_provider_state"
description: |-
  Manage network service providers for a given physical network.
---

# cloudstack_zone

Manage network service providers for a given physical network.  If Service Provider includes an underlying `Element` (configureInternalLoadBalancerElement, configureVirtualRouterElement) it will be configured.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_zone" "test" {
	name          = "acctest"
	dns1          = "8.8.8.8"
	dns2          = "8.8.8.8"
	internal_dns1 = "8.8.4.4"
	internal_dns2 = "8.8.4.4"
	network_type  = "Advanced"
	domain        = "cloudstack.apache.org"
}
resource "cloudstack_physical_network" "test" {
	broadcast_domain_range = "ZONE"
	isolation_methods      = "VLAN"
	name                   = "test01"
	network_speed          = "1G"
	tags                   = "vlan"
	zone_id                = cloudstack_zone.test.id
}
resource "cloudstack_network_service_provider_state" "virtualrouter" {
    name                  = "VirtualRouter"
    physical_network_id   = cloudstack_physical_network.test.id
    enabled                 = true
}
resource "cloudstack_network_service_provider_state" "vpcvirtualrouter" {
    name                  = "VpcVirtualRouter"
    physical_network_id   = cloudstack_physical_network.test.id
    enabled                 = true
}
resource "cloudstack_network_service_provider_state" "internallbvm" {
    name                  = "InternalLbVm"
    physical_network_id   = cloudstack_physical_network.test.id
    enabled                 = false
}
resource "cloudstack_network_service_provider_state" "configdrive" {
    name                  = "ConfigDrive"
    physical_network_id   = cloudstack_physical_network.test.id
    enabled                 = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) account who will own the VLAN. If VLAN is Zone wide, this parameter should be omitted
* `physical_network_id` - (Required) domain ID of the account owning a VLAN
* `enabled` - (Required) the ending IP address in the VLAN IP range


## Attributes Reference

The following attributes are exported:

* `id` - uuid of the network provider


## Import

Network service providers can be imported; use `<ID>` as the import ID. For
example:

```shell
terraform import cloudstack_network_service_provider_state.example VirtualRouter 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
