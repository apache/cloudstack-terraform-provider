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
    guest_ip_type = "Isolated"
    traffic_type = "Guest"
    network_rate = 100
    network_mode = "NATTED"
    conserve_mode = true
    enable = true
    for_vpc = false
    specify_vlan = true
    specify_ip_ranges = true
    max_connections = 256
    supported_services = ["Dhcp", "Dns", "Firewall", "Lb", "SourceNat"]
    service_provider_list = {
        Dhcp = "VirtualRouter"
        Dns = "VirtualRouter"
        Firewall = "VirtualRouter"
        Lb = "VirtualRouter"
        SourceNat = "VirtualRouter"
    }
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network offering.
* `display_text` - (Required) The display text of the network offering.
* `guest_ip_type` - (Required) The type of IP address allocation for the network offering. Possible values are "Shared" or "Isolated".
* `traffic_type` - (Required) The type of traffic for the network offering. Possible values are "Guest" or "Management".
* `network_rate` - (Optional) The network rate in Mbps for the network offering.
* `network_mode` - (Optional) The network mode. Possible values are "DHCP" or "NATTED".
* `conserve_mode` - (Optional) Whether to enable conserve mode. Defaults to `false`.
* `enable` - (Optional) Whether to enable the network offering. Defaults to `false`.
* `for_vpc` - (Optional) Whether this network offering is for VPC. Defaults to `false`.
* `for_nsx` - (Optional) Whether this network offering is for NSX. Defaults to `false`.
* `specify_vlan` - (Optional) Whether to allow specifying VLAN ID. Defaults to `false`.
* `specify_ip_ranges` - (Optional) Whether to allow specifying IP ranges. Defaults to `false`.
* `specify_as_number` - (Optional) Whether to allow specifying AS number. Defaults to `false`.
* `internet_protocol` - (Optional) The internet protocol. Possible values are "IPv4" or "IPv6". Defaults to "IPv4".
* `routing_mode` - (Optional) The routing mode. Possible values are "Static" or "Dynamic".
* `max_connections` - (Optional) The maximum number of concurrent connections supported by the network offering.
* `supported_services` - (Optional) A list of supported services for this network offering.
* `service_provider_list` - (Optional) A map of service providers for the supported services.

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
