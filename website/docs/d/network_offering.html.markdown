---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_network_offerings"
sidebar_current: "docs-cloudstack-cloudstack_network_offering"
description: |-
  Gets information about cloudstack network offering.
---

# cloudstack_network_offering

Use this datasource to get information about a network offering for use in other resources.

### Example Usage

```hcl
  data "cloudstack_network_offering" "net-off-data-source"{
    filter{
    name = "name"
    value="TestNetworkDisplay12"
    }
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the network offering.
* `display_text` - An alternate display text of the network offering.
* `guest_ip_type` - Guest type of the network offering, can be Shared or Isolated.
* `traffic_type` - The traffic type for the network offering, supported types are Public, Management, Control, Guest, Vlan or Storage.