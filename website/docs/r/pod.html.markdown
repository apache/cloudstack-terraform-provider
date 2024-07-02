---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_pod"
sidebar_current: "docs-cloudstack-resource-pod"
description: |-
  Creates a new Pod.
---

# cloudstack_pod

Creates a new Pod.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_pod" "example" {
	allocation_state = "Disabled"
	gateway          = "172.29.0.1"
	name             = "example"
	netmask          = "255.255.240.0"
	start_ip         =  "172.29.0.2"
	zone_id          =  cloudstack_zone.example.id
}
```

## Argument Reference

The following arguments are supported:

* `allocation_state` - (Optional) allocation state of this Pod for allocation of new resources.
* `end_ip` - (Optional) the ending IP address for the Pod.
* `gateway` - (Required) the gateway for the Pod.
* `name` - (Required) the name of the Pod.
* `netmask` - (Required) the netmask for the Pod.
* `start_ip` - (Required) the starting IP address for the Pod.
* `zone_id` - (Required) the Zone ID in which the Pod will be created.


## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.



## Import

A pod can be imported; use `<POD ID>` as the import ID. For
example:

```shell
terraform import cloudstack_pod.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
