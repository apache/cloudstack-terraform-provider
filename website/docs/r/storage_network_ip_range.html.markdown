---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_storage_network_ip_range"
sidebar_current: "docs-cloudstack-resource-storage-network-ip-range"
description: |-
  Creates a storage network IP range.
---

# cloudstack_storage_network_ip_range

Creates a storage network IP range for a pod.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_storage_network_ip_range" "default" {
    pod_id   = cloudstack_pod.default.id
    gateway  = "10.1.1.1"
    netmask  = "255.255.255.0"
    start_ip = "10.1.1.2"
    end_ip   = "10.1.1.10"
    vlan     = 100
}
```

## Argument Reference

The following arguments are supported:

- `pod_id` - (Required) The Pod ID for the storage network IP range. Changing
  this forces a new resource to be created.
- `gateway` - (Required) The gateway for the storage network IP range. Changing
  this forces a new resource to be created.
- `netmask` - (Required) The netmask for the storage network IP range. Changing
  this forces a new resource to be created.
- `start_ip` - (Required) The beginning IP address in the storage network IP range.
  Changing this forces a new resource to be created.
- `end_ip` - (Optional) The ending IP address in the storage network IP range.
  Changing this forces a new resource to be created.
- `vlan` - (Optional) The optional VLAN of the storage network IP range.

`netmask`, `start_ip`, and `end_ip` force replacement rather than an in-place
update: CloudStack's `updateStorageNetworkIpRange` API validates a new range
against the record's own current IPs without excluding the record being
updated. An edit that keeps any overlap with the current range (for example
shrinking or extending it) fails with an IP overlap error against itself;
only moving to a range fully disjoint from the current one would succeed via
the API. Since Terraform can't tell which case applies before attempting the
update, these fields always force a replacement.

Replacement is more destructive than the update it replaces: the existing
range is deleted before the new one is created (the two would overlap, so
`create_before_destroy` can't help), so the storage network briefly has no
IP range configured for this pod during the apply.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the storage network IP range.
- `end_ip` - The ending IP address in the storage network IP range.
- `network_id` - The network ID of the storage network IP range.

## Import

Storage network IP ranges can be imported; use `<ID>` as the import ID. For
example:

```shell
terraform import cloudstack_storage_network_ip_range.default 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
