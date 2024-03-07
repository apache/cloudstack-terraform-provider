---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_disk"
sidebar_current: "docs-cloudstack-resource-disk"
description: |-
  Creates a disk volume from a disk offering. This disk volume will be attached to a virtual machine if the optional parameters are configured.
---

# cloudstack_disk

Creates a disk volume from a disk offering. This disk volume will be attached to
a virtual machine if the optional parameters are configured.

## Example Usage

```hcl
resource "cloudstack_disk" "default" {
  name               = "test-disk"
  attach             = "true"
  disk_offering      = "custom"
  size               = 50
  virtual_machine_id = "server-1"
  zone               = "zone-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the disk volume. Changing this forces a new
    resource to be created.

* `attach` - (Optional) Determines whether or not to attach the disk volume to a
    virtual machine (defaults false).

* `device_id` - (Optional) The device ID to map the disk volume to within the guest OS.

* `disk_offering` - (Required) The name or ID of the disk offering to use for
    this disk volume.

* `size` - (Optional) The size of the disk volume in gigabytes.

* `shrink_ok` - (Optional) Verifies if the disk volume is allowed to shrink when
    resizing (defaults false).

* `virtual_machine_id` - (Optional) The ID of the virtual machine to which you want
    to attach the disk volume.

* `project` - (Optional) The name or ID of the project to deploy this
    instance to. Changing this forces a new resource to be created.

* `zone` - (Required) The name or ID of the zone where this disk volume will be available.
    Changing this forces a new resource to be created.

* `reattach_on_change` - (Optional) Determines whether or not to detach the disk volume
    from the virtual machine on disk offering or size change.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk volume.
* `device_id` - The device ID the disk volume is mapped to within the guest OS.

## Import

Disks can be imported; use `<DISK ID>` as the import ID. For
example:

```shell
terraform import cloudstack_disk.default 6f3ee798-d417-4e7a-92bc-95ad41cf1244
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_disk.default my-project/6f3ee798-d417-4e7a-92bc-95ad41cf1244
```
