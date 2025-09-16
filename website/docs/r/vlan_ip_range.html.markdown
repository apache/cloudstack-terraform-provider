---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_vlan_ip_range"
sidebar_current: "docs-cloudstack-resource-vlan_ip_range"
description: |-
  Creates a VLAN IP range.
---

# cloudstack_zone

Creates a VLAN IP range.

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
resource "cloudstack_vlan_ip_range" "test" {
    physical_network_id = cloudstack_physical_network.test.id
    for_virtual_network = true
    zone_id = cloudstack_zone.test.id
    gateway = "10.0.0.1"
    netmask = "255.255.255.0"
    start_ip = "10.0.0.2"
    end_ip   = "10.0.0.10"
    vlan     = "vlan://123"
}
```

## Argument Reference

The following arguments are supported:

* `account` - (Optional) account who will own the VLAN. If VLAN is Zone wide, this parameter should be omitted
* `domain_id` - (Optional) domain ID of the account owning a VLAN
* `end_ip` - (Optional) the ending IP address in the VLAN IP range
* `end_ipv6` - (Optional) the ending IPv6 address in the IPv6 network range
* `for_system_vms` - (Optional) true if IP range is set to system vms, false if not
* `for_virtual_network` - (Optional) true if VLAN is of Virtual type, false if Direct
* `gateway` - (Optional) the gateway of the VLAN IP range
* `ip6_cidr` - (Optional) the CIDR of IPv6 network, must be at least /64
* `ip6_gateway` - (Optional) the gateway of the IPv6 network. Required for Shared networks and Isolated networks when it belongs to VPC
* `netmask` - (Optional) the netmask of the VLAN IP range
* `network_id` - (Optional) the network id
* `physical_network_id` - (Optional) the physical network id
* `pod_id` - (Optional) Have to be specified for Direct Untagged vlan only.
* `project_id` - (Optional) project who will own the VLAN. If VLAN is Zone wide, this parameter should be omitted
* `start_ip` - (Optional) the beginning IP address in the VLAN IP range
* `start_ipv6` - (Optional) the beginning IPv6 address in the IPv6 network range
* `vlan` - (Optional) true if network is security group enabled, false otherwise
* `start_ipv6` - (Optional) the ID or VID of the VLAN. If not specified, will be defaulted to the vlan of the network or if vlan of the network is null - to Untagged
* `zone_id` - (Optional) the Zone ID of the VLAN IP range

## Attributes Reference

The following attributes are exported:

* `id` - the ID of the VLAN IP range
* `network_id` - the network id of vlan range


## Import

Vlan ip range can be imported; use `<ID>` as the import ID. For
example:

```shell
terraform import cloudstack_vlan_ip_range.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
