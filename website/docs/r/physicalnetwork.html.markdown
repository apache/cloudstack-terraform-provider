---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_physicalnetwork"
sidebar_current: "docs-cloudstack-resource-physicalnetwork"
description: |-
  Creates a physical network.
---

# cloudstack_physicalnetwork

Creates a physical network.

## Example Usage

```hcl
resource "cloudstack_physicalnetwork" "default" {
  name = "test-physical-network"
  zone = "zone-name"
  
  broadcast_domain_range = "ZONE"
  isolation_methods      = ["VLAN"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the physical network.
* `zone` - (Required) The name or ID of the zone where the physical network belongs to.
* `broadcast_domain_range` - (Optional) The broadcast domain range for the physical network. Defaults to `ZONE`.
* `isolation_methods` - (Optional) The isolation method for the physical network.
* `network_speed` - (Optional) The speed for the physical network.
* `vlan` - (Optional) The VLAN for the physical network.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the physical network.