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

* `name` - (Required) The name of the pod.
* `zone_id` - (Required) The Zone ID in which the pod will be created.
* `gateway` - (Required) The gateway for the pod.
* `netmask` - (Required) The netmask for the pod.
* `start_ip` - (Required) The starting IP address for the pod.
* `end_ip` - (Required) The ending IP address for the pod.
* `allocation_state` - (Optional) Allocation state of this pod for allocation of new resources.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the pod.
* `allocation_state` - The allocation state of the pod.
* `zone_name` - The name of the zone where the pod is created.
* `vlan_id` - The VLAN ID associated with the pod.

## Import

Pods can be imported; use `<POD ID>` as the import ID. For example:

```shell
terraform import cloudstack_pod.default 5fb02d7f-9513-4f96-9fbe-b5d167f4e90b