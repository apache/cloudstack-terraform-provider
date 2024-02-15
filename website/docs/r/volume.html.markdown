---
layout: default
page_title: "CloudStack: cloudstack_volume"
sidebar_current: "docs-cloudstack-resource-volume"
description: |-
    Creates a Volume
---
# CloudStack: cloudstack_volume

A `cloudstack_volume` resource manages a volume within CloudStack.

## Example Usage

```hcl
resource "cloudstack_volume" "example" {
    name = "example-volume"
    disk_offering_id = "a6f7e5fb-1b9a-417e-a46e-7e3d715f34d3"
    zone_id = "b0fcd7cc-5e14-499d-a2ff-ecf49840f1ab"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume. Forces new resource.
* `disk_offering_id` - (Required) The ID of the disk offering for the volume. Forces new resource.
* `zone_id` - (Required) The ID of the zone where the volume will be created. Forces new resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the volume.
* `name` - The name of the volume.
* `disk_offering_id` - The ID of the disk offering for the volume.
* `zone_id` - The ID of the zone where the volume resides.

## Import

Volumes can be imported; use `<VOLUMEID>` as the import ID. For example:

```shell
$ terraform import cloudstack_volume.example <VOLUMEID>
```
