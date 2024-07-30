---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_disk_offering"
sidebar_current: "docs-cloudstack-resource-disk-offering"
description: |-
  Creates a disk offering.
---

# cloudstack_disk_offering

Creates a disk offering.

## Example Usage

```hcl
resource "cloudstack_disk_offering" "example" {
	display_text = "example"
	name         = "example"

	storage_type = "shared"
	provisioning_type = "fat"
	cache_mode        = "writeback"
	tags              = "test1,test2"

	disk_size = 7
}
```

## Argument Reference

The following arguments are supported:

* `display_text` - (Required) alternate display text of the disk offering
* `name` - (Required) name of the disk offering

---

* `cache_mode` - (Optional) the cache mode to use for this disk offering. none, writeback or writethrough
* 
* `disk_size` - (Optional) size of the disk offering in GB (1GB = 1,073,741,824 bytes)
* 
* `disk_offering_strictness` - (Optional) To allow or disallow the resize operation on the disks created from this disk offering, if the flag is true then resize is not allowed

* `domain_id` - (Optional) the ID of the containing domain(s), null for public offerings

* `iops_read_rate` - (Optional) io requests read rate of the disk offering

* `iops_read_rate_max` - (Optional) burst requests read rate of the disk offering

* `iops_read_rate_max_length` - (Optional) length (in seconds) of the burst

* `iops_write_rate` - (Optional) io requests write rate of the disk offering

* `iops_write_rate_max` - (Optional) burst io requests write rate of the disk offering

* `iops_write_rate_max_length` - (Optional) length (in seconds) of the burst

* `provisioning_type` - (Optional) provisioning type used to create volumes. Valid values are thin, sparse, fat.

* `storage_type` - (Optional) the storage type of the disk offering. Values are local and shared.

* `tags` - (Optional) tags for the disk offering

* `zone_id` - (Optional) the ID of the containing zone(s), null for public offerings

* `hypervisor` - (Optional) A `backend_request` block as defined below.

* `storage` - (Optional) A `backend_request` block as defined below.

---
A `hypervisor` block supports the following:

* `bytes_read_rate` - (Optional) bytes read rate of the disk offering

* `bytes_read_rate_max` - (Optional) burst bytes read rate of the disk offering

* `bytes_read_rate_max_length` - (Optional) length (in seconds) of the burst

* `bytes_write_rate` - (Optional) bytes write rate of the disk offering

* `bytes_write_rate_max` - (Optional) burst bytes write rate of the disk offering

* `bytes_write_rate_max_length` - (Optional) length (in seconds) of the burst

---
A `storage` block supports the following:

* `max_iops` - (Optional) max iops of the disk offering

* `min_iops` - (Optional) min iops of the disk offering

* `customized_iops` - (Optional) whether disk offering iops is custom or not

* `hypervisor_snapshot_reserve` - (Optional) Hypervisor snapshot reserve space as a percent of a volume (for managed storage using Xen or VMware)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The id of the affinity group.


## Import

Disk offerings can be imported; use `<ID>` as the import ID. For
example:

```shell
terraform import cloudstack_disk_offering.example 6226ea4d-9cbe-4cc9-b30c-b9532146da5b
```