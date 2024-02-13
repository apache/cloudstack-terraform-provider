---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_zone"
sidebar_current: "docs-cloudstack-resource-zone"
description: |-
  Creates a Zone.
---

# cloudstack_zone

Creates a Zone.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_zone" "example" {
	name          = "example"
	dns1          = "8.8.8.8"
	internal_dns1 = "8.8.4.4"
	network_type  = "Advanced"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_state` - (Optional) Allocation state of this Zone for allocation of new resources.
* `dns1` - (Required) the first DNS for the Zone
* `dns2` - (Optional) the second DNS for the Zone
* `domain` - (Optional) Network domain name for the networks in the zone
* `domain_id` - (Optional) the ID of the containing domain, null for public zones
* `guest_cidr_address` - (Optional) the guest CIDR address for the Zone
* `internal_dns1` - (Required) the first internal DNS for the Zone
* `internal_dns2` - (Optional) the second internal DNS for the Zone
* `ip6_dns1` - (Optional) the first DNS for IPv6 network in the Zone
* `ip6_dns2` - (Optional) the second DNS for IPv6 network in the Zone
* `local_storage_enabled` - (Optional) true if local storage offering enabled, false otherwise
* `name` - (Required) the name of the Zone
* `network_type` - (Required) network type of the zone, can be Basic or Advanced
* `security_group_enabled` - (Optional) true if network is security group enabled, false otherwise

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.
* `dhcp_provider` - the dhcp Provider for the Zone.


## Import

Zones can be imported; use `<ZONE ID>` as the import ID. For
example:

```shell
terraform import cloudstack_zone.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
