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
* `disk_size` - (Required) The size of the disk offering in GB.
* `encrypt` - (Optional)  Volumes using this offering should be encrypted
* `tags` - (Optional)  tags for the disk offering

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk offering.
* `name` - The name of the disk offering.
* `display_text` - The display text of the disk offering.
* `disk_size` - The size of the disk offering in GB.

## Import

Disk offerings can be imported; use `<DISKOFFERINGID>` as the import ID. For example:

```shell
$ terraform import cloudstack_disk_offering.example <DISKOFFERINGID>
```
