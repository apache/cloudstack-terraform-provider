---
subcategory: "Pod"
layout: "cloudstack"
page_title: "CloudStack: cloudstack_pod"
description: |-
  Creates a pod.
---

# cloudstack_pod

Creates a pod.

## Example Usage

```hcl
resource "cloudstack_pod" "default" {
  name = "pod-1"
  zone_id = "1"
  gateway = "10.1.1.1"
  netmask = "255.255.255.0"
  start_ip = "10.1.1.100"
  end_ip = "10.1.1.200"
}
```

## Argument Reference

The following arguments are supported:

<<<<<<< HEAD
* `name` - (Required) The name of the pod.
* `zone_id` - (Required) The Zone ID in which the pod will be created.
* `gateway` - (Required) The gateway for the pod.
* `netmask` - (Required) The netmask for the pod.
* `start_ip` - (Required) The starting IP address for the pod.
* `end_ip` - (Required) The ending IP address for the pod.
* `allocation_state` - (Optional) Allocation state of this pod for allocation of new resources.
=======
* `allocation_state` - (Optional) allocation state of this Pod for allocation of new resources.
* `end_ip` - (Optional) the ending IP address for the Pod.
* `gateway` - (Required) the gateway for the Pod.
* `name` - (Required) the name of the Pod.
* `netmask` - (Required) the netmask for the Pod.
* `start_ip` - (Required) the starting IP address for the Pod.
* `zone_id` - (Required) the Zone ID in which the Pod will be created.

>>>>>>> apache/main

## Attributes Reference

The following attributes are exported:

<<<<<<< HEAD
* `id` - The ID of the pod.
* `allocation_state` - The allocation state of the pod.
* `zone_name` - The name of the zone where the pod is created.
* `vlan_id` - The VLAN ID associated with the pod.

## Import

Pods can be imported; use `<POD ID>` as the import ID. For example:

```shell
terraform import cloudstack_pod.default 5fb02d7f-9513-4f96-9fbe-b5d167f4e90b
=======
* `id` - The instance ID.



## Import

A pod can be imported; use `<POD ID>` as the import ID. For
example:

```shell
terraform import cloudstack_pod.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
>>>>>>> apache/main
