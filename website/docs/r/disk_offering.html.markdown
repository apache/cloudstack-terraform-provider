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
    name = "example-disk-offering"
    display_text = "Example Disk Offering"
    disk_size = 100
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the disk offering.
* `display_text` - (Required) The display text of the disk offering.
* `disk_size` - (Optional) The size of the disk offering in GB. Conflicts with
    `customized`. Changing this forces a new resource to be created.
* `customized` - (Optional) Whether the disk offering allows a custom disk size
    to be specified at deployment time. Conflicts with `disk_size`. Defaults to
    `false`. Changing this forces a new resource to be created.
* `storage_type` - (Optional) The storage type of the disk offering. Values are
    `local` and `shared`. Defaults to `shared`. Changing this forces a new
    resource to be created.
* `provisioning_type` - (Optional) The provisioning type used to create volumes.
    Values are `thin`, `sparse` and `fat`. Defaults to `thin`. Changing this
    forces a new resource to be created.
* `tags` - (Optional) The storage tags for the disk offering.
* `display_offering` - (Optional) Whether the disk offering is displayed to the
    end user. Defaults to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk offering.
* `name` - The name of the disk offering.
* `display_text` - The display text of the disk offering.
* `disk_size` - The size of the disk offering in GB.
* `customized` - Whether the disk offering allows a custom disk size.
* `storage_type` - The storage type of the disk offering.
* `provisioning_type` - The provisioning type of the disk offering.
* `tags` - The storage tags for the disk offering.
* `display_offering` - Whether the disk offering is displayed to the end user.

## Import

Disk offerings can be imported; use `<DISKOFFERINGID>` as the import ID. For example:

```shell
$ terraform import cloudstack_disk_offering.example <DISKOFFERINGID>
```
