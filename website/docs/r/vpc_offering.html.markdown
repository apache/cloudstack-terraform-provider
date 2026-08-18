---
layout: default
page_title: "CloudStack: cloudstack_vpc_offering"
sidebar_current: "docs-cloudstack-resource-vpc_offering"
description: |-
    Creates a VPC Offering
---

# CloudStack: cloudstack_vpc_offering

A `cloudstack_vpc_offering` resource manages a VPC offering within CloudStack.

## Example Usage

```hcl
resource "cloudstack_vpc_offering" "example" {
    name         = "example-vpc-offering"
    display_text = "Example VPC Offering"
    enable       = true
    supported_services = ["Dhcp", "Dns", "SourceNat", "PortForwarding", "Lb", "UserData", "StaticNat", "NetworkACL"]
    service_provider_list = {
        Dhcp           = "VpcVirtualRouter"
        Dns            = "VpcVirtualRouter"
        SourceNat      = "VpcVirtualRouter"
        PortForwarding = "VpcVirtualRouter"
        Lb             = "VpcVirtualRouter"
        UserData       = "VpcVirtualRouter"
        StaticNat      = "VpcVirtualRouter"
        NetworkACL     = "VpcVirtualRouter"
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VPC offering.
* `display_text` - (Required) The display text of the VPC offering.
* `domain_id` - (Optional) The ID of the containing domain(s), null for public offerings.
* `zone_id` - (Optional) The ID of the containing zone(s), null for offerings available in all zones.
* `enable` - (Optional) Whether to enable the VPC offering. Defaults to `false`.
* `for_nsx` - (Optional) Whether this VPC offering is meant to be used for NSX. Defaults to `false`.
* `nsx_support_lb` - (Optional) Whether the NSX supports the Lb service. Defaults to `false`.
* `internet_protocol` - (Optional) The internet protocol. Possible values are "IPv4" or "dualstack". Defaults to "IPv4".
* `network_mode` - (Optional) Indicates the mode with which the VPC will operate. Possible values are "NATTED" or "ROUTED".
* `routing_mode` - (Optional) The routing mode for the VPC offering. Possible values are "Static" or "Dynamic".
* `specify_as_number` - (Optional) Whether to allow specifying an AS number. Defaults to `false`.
* `network_provider` - (Optional) The network provider for the VPC offering.
* `service_offering_id` - (Optional) The ID or name of the service offering for the VPC router appliance.
* `supported_services` - (Required) A list of supported services for this VPC offering. CloudStack always requires (and, if omitted, will silently add) `NetworkACL` support, and adds `SourceNat` support unless `network_mode` is `ROUTED`. Since this attribute forces recreation on change, omitting either of these services from your configuration will cause `terraform plan` to show a permanent diff wanting to recreate the resource - always list them explicitly.
* `service_provider_list` - (Optional) A map of service providers for the supported services.
* `service_capability_list` - (Optional) A list of desired service capabilities for this VPC offering. Each block supports:
  * `service` - (Required) The service the capability applies to, e.g. `SourceNat`.
  * `capability_type` - (Required) The capability type, e.g. `RedundantRouter`.
  * `capability_value` - (Required) The capability value, e.g. `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC offering.
* `is_default` - `true` if this is the default VPC offering.

## Import

VPC offerings can be imported; use `<VPCOFFERINGID>` as the import ID. For example:

```shell
$ terraform import cloudstack_vpc_offering.example <VPCOFFERINGID>
```
