---
layout: "cloudstack"
page_title: "CloudStack: physical_network"
sidebar_current: "docs-cloudstack-resource-physical-network"
description: |-
  Creates a physical network
---

# physical_network

Creates a physical network

## Example Usage

Basic usage:

```hcl
resource "cloudstack_physical_network" "example" {
	broadcast_domain_range = "ZONE"
	isolation_methods      = "VLAN"
	name                   = "example"
	network_speed          = "10G"
	tags                   = "vlan"
	zone_id                = cloudstack_zone.example.id
}
```

## Argument Reference

The following arguments are supported:

* `broadcast_domain_range` - (Optional) changeme.
* `domain_id` - (Optional) changeme.
* `isolation_methods` - (Optional) changeme.
* `name` - (Required) changeme.
* `network_speed` - (Optional) changeme.
* `tags` - (Optional) changeme.
* `vlan` - (Optional) changeme.
* `zone_id` - (Required) changeme.


## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.


## Import

Physical networks can be imported; use `<PHYSICAL NETWORK ID>` as the import ID. For
example:

```shell
terraform import physical_network.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
