---
layout: default
page_title: "CloudStack: cloudstack_network_offering"
sidebar_current: "docs-cloudstack-resource-network_offering"
description: |-
    Creates a Network Offering
---

# CloudStack: cloudstack_network_offering

A `cloudstack_network_offering` resource manages a network offering within CloudStack.

## Example Usage

```hcl
resource "cloudstack_network_offering" "example" {
    name = "example-network-offering"
    display_text = "Example Network Offering"
    guest_ip_type = "Shared"
    traffic_type = "Guest"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network offering.
* `display_text` - (Required) The display text of the network offering.
* `guest_ip_type` - (Required) The type of IP address allocation for the network offering. Possible values are "Shared" or "Isolated".
* `traffic_type` - (Required) The type of traffic for the network offering. Possible values are "Guest" or "Management".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network offering.
* `name` - The name of the network offering.
* `display_text` - The display text of the network offering.
* `guest_ip_type` - The type of IP address allocation for the network offering.
* `traffic_type` - The type of traffic for the network offering.

## Import

Network offerings can be imported; use `<NETWORKOFFERINGID>` as the import ID. For example:

```shell
$ terraform import cloudstack_network_offering.example <NETWORKOFFERINGID>
```
