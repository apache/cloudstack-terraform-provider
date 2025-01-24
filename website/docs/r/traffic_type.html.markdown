---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_traffic_type"
sidebar_current: "docs-cloudstack-resource-traffic-type"
description: |-
  Adds traffic type to a physical network.
---

# cloudstack_traffic_type

Adds traffic type to a physical network.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_traffic_type" "example" {
	physical_network_id = cloudstack_physical_network.example.id
	traffic_type        = "Management"
	kvm_network_label   = "example"
}
```

## Argument Reference

The following arguments are supported:

* `hyperv_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a Hyperv host.
* `isolation_method` - (Optional) Used if physical network has multiple isolation types and traffic type is public. Choose which isolation method. Valid options currently 'vlan' or 'vxlan', defaults to 'vlan'..
* `kvm_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a KVM host.
* `ovm3_network_label` - (Optional) The network name of the physical device dedicated to this traffic on an OVM3 host.
* `physical_network_id` - (Required) The Physical Network ID.
* `traffic_type` - (Required) The trafficType to be added to the physical network.
* `vlan` - (Optional) The VLAN id to be used for Management traffic by VMware host.
* `vmware_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a VMware host.
* `xen_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a XenServer host.


## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.

