---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_network_service_provider"
sidebar_current: "docs-cloudstack-resource-network-service-provider"
description: |-
  Adds a network service provider to a physical network.
---

# cloudstack_network_service_provider

Adds or updates a network service provider on a physical network.

~> **NOTE:** Network service providers are often created automatically when a physical network is created. This resource can be used to manage those existing providers or create new ones.

~> **NOTE:** Some providers like SecurityGroupProvider don't allow updating the service list. For these providers, the service list specified in the configuration will be used only during creation.

~> **NOTE:** Network service providers are created in a "Disabled" state by default. You can set `state = "Enabled"` to enable them. Note that some providers like VirtualRouter require configuration before they can be enabled.

## Example Usage

```hcl
resource "cloudstack_physicalnetwork" "default" {
  name = "test-physical-network"
  zone = "zone-name"
}

resource "cloudstack_network_service_provider" "virtualrouter" {
  name                = "VirtualRouter"
  physical_network_id = cloudstack_physicalnetwork.default.id
  service_list        = ["Dhcp", "Dns", "Firewall", "LoadBalancer", "SourceNat", "StaticNat", "PortForwarding", "Vpn"]
  state               = "Enabled"
}

resource "cloudstack_network_service_provider" "vpcvirtualrouter" {
  name                = "VpcVirtualRouter"
  physical_network_id = cloudstack_physicalnetwork.default.id
  service_list        = ["Dhcp", "Dns", "SourceNat", "StaticNat", "NetworkACL", "PortForwarding", "Lb", "UserData", "Vpn"]
}

resource "cloudstack_network_service_provider" "securitygroup" {
  name                = "SecurityGroupProvider"
  physical_network_id = cloudstack_physicalnetwork.default.id
  # Note: service_list is predefined for SecurityGroupProvider
  state               = "Enabled"  # Optional: providers are created in "Disabled" state by default
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network service provider. Possible values include: VirtualRouter, VpcVirtualRouter, InternalLbVm, ConfigDrive, etc.
* `physical_network_id` - (Required) The ID of the physical network to which to add the network service provider.
* `destination_physical_network_id` - (Optional) The destination physical network ID.
* `service_list` - (Optional) The list of services to be enabled for this service provider. Possible values include: Dhcp, Dns, Firewall, Gateway, LoadBalancer, NetworkACL, PortForwarding, SourceNat, StaticNat, UserData, Vpn, etc.
* `state` - (Optional) The state of the network service provider. Possible values are "Enabled" and "Disabled". This can be used to enable or disable the provider.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network service provider.
* `state` - The state of the network service provider.

## Import

Network service providers can be imported using the network service provider ID, e.g.

```shell
terraform import cloudstack_network_service_provider.virtualrouter 5fb307e2-0e11-11ee-be56-0242ac120002
```
