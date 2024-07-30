---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_storage_pool"
sidebar_current: "docs-cloudstack-resource-storage-pool"
description: |-
  Creates a storage pool.
---

# cloudstack_storage_pool

Creates a storage pool.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_storage_pool" "example" {
	name         = "example"
	url          = "nfs://10.147.28.6/export/home/sandbox/primary11"
	zone_id      = "0ed38eb3-f279-4951-ac20-fef39ebab20c"
	cluster_id   = "9daeeb36-d8b7-497a-9b53-bbebba88c817"
	pod_id       = "2ff52b73-139e-4c40-a0a3-5b7d87d8e3c4"
	scope        = "CLUSTER"
	hypervisor   = "Simulator"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) the cluster ID for the storage pool.
* `hypervisor` - (Optional) hypervisor type of the hosts in zone that will be attached to this storage pool. KVM, VMware supported as of now.
* `name` - (Required) the name for the storage pool.
* `pod_id` - (Optional) the Pod ID for the storage pool.
* `storage_provider` - (Optional) the storage provider name.
* `scope` - (Optional) the scope of the storage: cluster or zone.
* `state` - (Optional) the state of the storage pool.
* `tags` - (Optional) the tags for the storage pool.
* `url` - (Required) the URL of the storage pool.
* `zone_id` - (Required) the Zone ID for the storage pool.


## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.



## Import

Storage pools can be imported; use `<STORAGE POOL ID>` as the import ID. For
example:

```shell
terraform import cloudstack_storage_pool.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
