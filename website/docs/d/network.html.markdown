---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_network"
sidebar_current: "docs-cloudstack-datasource-network"
description: |-
  Gets information about a CloudStack network.
---

# cloudstack_network

Use this datasource to get information about a network for use in other resources.

### Example Usage

```hcl
data "cloudstack_network" "network-data-source" {
  filter {
    name  = "name"
    value = "test-network"
  }

  filter {
    name  = "cidr"
    value = "10.1.1.0/24"
  }
}
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.
* `project` - (Optional) The name or ID of the project the network belongs to.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the network.
* `display_text` - An alternate display text of the network.
* `cidr` - The CIDR block of the network.
* `gateway` - The gateway of the network.
* `network_domain` - The network domain of the network.
* `network_offering_id` - The ID of the network offering used by the network.
* `network_offering_name` - The name of the network offering used by the network.
* `project` - The project name of the network.
* `project_id` - The project ID of the network.
* `vpc_id` - The VPC ID the network belongs to.
* `acl_id` - The ACL ID attached to the network.
* `zone_id` - The ID of the zone the network belongs to.
* `zone_name` - The name of the zone the network belongs to.
* `state` - The state of the network.
* `type` - The type of the network.
* `traffic_type` - The traffic type of the network.
* `tags` - The tags assigned to the network.
