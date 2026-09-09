---
layout: default
page_title: "CloudStack: cloudstack_disk_offering"
sidebar_current: "docs-cloudstack-resource-disk_offering"
description: |-
    Creates a Disk Offering
---

# CloudStack: cloudstack_disk_offering

A `cloudstack_disk_offering` resource manages a disk offering within CloudStack.

## Example Usage

```hcl
resource "cloudstack_disk_offering" "example" {
    name         = "example-disk-offering"
    display_text = "Example Disk Offering"
    disk_size    = 100

    storage {
        min_iops                    = 1000
        max_iops                    = 5000
        customized_iops             = false
        hypervisor_snapshot_reserve = 25
    }

    hypervisor {
        bytes_read_rate  = 1048576
        bytes_write_rate = 1048576
    }
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the disk offering.
* `display_text` - (Required) The display text of the disk offering.
* `cache_mode` - (Optional) The cache mode to use for this disk offering. Values
    are `none`, `writeback` or `writethrough`. Computed when not set.
* `disk_size` - (Optional) The size of the disk offering in GB. When set, the
    offering uses a fixed size; when omitted, the offering allows a custom disk
    size to be specified at deployment time. Changing this forces a new resource
    to be created.
* `disk_offering_strictness` - (Optional) Whether the disk offering size is
    strictly enforced and cannot be resized at deployment time. When `true`,
    resize is not allowed. Changing this forces a new resource to be created.
* `domain_id` - (Optional) The list of domain IDs that can use this disk
    offering. Leave empty to make the offering public.
* `iops_read_rate` - (Optional) The IO requests read rate of the disk offering.
* `iops_read_rate_max` - (Optional) The burst IO requests read rate of the disk
    offering.
* `iops_read_rate_max_length` - (Optional) The length (in seconds) of the IOPS
    read rate burst.
* `iops_write_rate` - (Optional) The IO requests write rate of the disk offering.
* `iops_write_rate_max` - (Optional) The burst IO requests write rate of the disk
    offering.
* `iops_write_rate_max_length` - (Optional) The length (in seconds) of the IOPS
    write rate burst.
* `provisioning_type` - (Optional) The provisioning type used to create volumes.
    Values are `thin`, `sparse` and `fat`. Computed when not set. Changing this
    forces a new resource to be created.
* `storage_type` - (Optional) The storage type of the disk offering. Values are
    `local` and `shared`. Computed when not set. Changing this forces a new
    resource to be created.
* `tags` - (Optional) The storage tags for the disk offering.
* `display_offering` - (Optional) Whether the disk offering is displayed to the
    end user. Defaults to `true`.
* `zone_id` - (Optional) The list of zone IDs the disk offering is restricted to.
    Leave empty to make the offering available in all zones.
* `hypervisor` - (Optional) A hypervisor-side QoS block (throughput limits).
    Only one block is supported. The structure is documented below.
* `storage` - (Optional) A storage-side QoS block (IOPS limits). Only one block
    is supported. The structure is documented below.

The `hypervisor` block supports the following (all optional; changing any forces
a new resource to be created):

* `bytes_read_rate` - The bytes read rate of the disk offering.
* `bytes_read_rate_max` - The burst bytes read rate of the disk offering.
* `bytes_read_rate_max_length` - The length (in seconds) of the read rate burst.
* `bytes_write_rate` - The bytes write rate of the disk offering.
* `bytes_write_rate_max` - The burst bytes write rate of the disk offering.
* `bytes_write_rate_max_length` - The length (in seconds) of the write rate burst.

The `storage` block supports the following (all optional; changing any forces a
new resource to be created):

* `min_iops` - The minimum IOPS of the disk offering.
* `max_iops` - The maximum IOPS of the disk offering.
* `customized_iops` - Whether the disk offering IOPS are customizable at
    deployment time.
* `hypervisor_snapshot_reserve` - Hypervisor snapshot reserve space as a percent
    of a volume (for managed storage using Xen or VMware).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk offering.
* `name` - The name of the disk offering.
* `display_text` - The display text of the disk offering.
* `cache_mode` - The cache mode of the disk offering.
* `disk_size` - The size of the disk offering in GB.
* `disk_offering_strictness` - Whether the disk offering size is strictly
    enforced.
* `domain_id` - The list of domain IDs that can use this disk offering.
* `iops_read_rate` - The IO requests read rate of the disk offering.
* `iops_read_rate_max` - The burst IO requests read rate of the disk offering.
* `iops_read_rate_max_length` - The length (in seconds) of the IOPS read burst.
* `iops_write_rate` - The IO requests write rate of the disk offering.
* `iops_write_rate_max` - The burst IO requests write rate of the disk offering.
* `iops_write_rate_max_length` - The length (in seconds) of the IOPS write burst.
* `provisioning_type` - The provisioning type of the disk offering.
* `storage_type` - The storage type of the disk offering.
* `tags` - The storage tags for the disk offering.
* `display_offering` - Whether the disk offering is displayed to the end user.
* `zone_id` - The list of zone IDs the disk offering is restricted to.
* `hypervisor` - The hypervisor QoS block, as documented above.
* `storage` - The storage QoS block, as documented above.

## Import

Disk offerings can be imported; use `<DISKOFFERINGID>` as the import ID. For example:

```shell
$ terraform import cloudstack_disk_offering.example <DISKOFFERINGID>
```
