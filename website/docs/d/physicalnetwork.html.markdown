---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_physicalnetwork"
sidebar_current: "docs-cloudstack-datasource-physicalnetwork"
description: |-
  Gets information about a physical network.
---

# cloudstack_physicalnetwork

Use this data source to get information about a physical network.

## Example Usage

```hcl
data "cloudstack_physicalnetwork" "default" {
  filter {
    name = "name"
    value = "test-physical-network"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to filter off of.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the physical network.
* `name` - The name of the physical network.
* `zone` - The name of the zone where the physical network belongs to.
* `broadcast_domain_range` - The broadcast domain range for the physical network.
* `isolation_methods` - The isolation method for the physical network.
* `network_speed` - The speed for the physical network.
* `vlan` - The VLAN for the physical network.