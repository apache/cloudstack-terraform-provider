---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_traffic_type"
sidebar_current: "docs-cloudstack-resource-traffic-type"
description: |-
  Adds a traffic type to a physical network.
---

# cloudstack_traffic_type

Adds a traffic type to a physical network.

## Example Usage

```hcl
resource "cloudstack_physicalnetwork" "default" {
  name = "test-physical-network"
  zone = "zone-name"
}

resource "cloudstack_traffic_type" "management" {
  physical_network_id = cloudstack_physicalnetwork.default.id
  type                = "Management"
  
  kvm_network_label    = "cloudbr0"
  xen_network_label    = "xenbr0"
  vmware_network_label = "VM Network"
}

resource "cloudstack_traffic_type" "guest" {
  physical_network_id = cloudstack_physicalnetwork.default.id
  type                = "Guest"
  
  kvm_network_label    = "cloudbr1"
  xen_network_label    = "xenbr1"
  vmware_network_label = "VM Guest Network"
}
```

## Argument Reference

The following arguments are supported:

* `physical_network_id` - (Required) The ID of the physical network to which the traffic type is being added.
* `type` - (Required) The type of traffic (e.g., Management, Guest, Public, Storage).
* `kvm_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a KVM host.
* `vlan` - (Optional) The VLAN ID to be used for Management traffic by VMware host.
* `xen_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a XenServer host.
* `vmware_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a VMware host.
* `hyperv_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on a HyperV host.
* `ovm3_network_label` - (Optional) The network name label of the physical device dedicated to this traffic on an OVM3 host.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the traffic type.

## Import

Traffic types can be imported using the traffic type ID, e.g.

```shell
terraform import cloudstack_traffic_type.management 5fb307e2-0e11-11ee-be56-0242ac120002
```
