---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_vgpu_profile"
description: |-
  Gets information about a vGPU profile.
---

# cloudstack_vgpu_profile

Use this data source to get information about a vGPU profile for use in other resources.

## Example Usage

```hcl
data "cloudstack_vgpu_profile" "profile" {
  filter {
    name = "name"
    value = "passthrough"
  }
}

output "vgpu_profile_id" {
  value = data.cloudstack_vgpu_profile.profile.id
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) One or more name/value pairs to filter off of. See detailed documentation below.

### Filter Arguments

* `name` - (Required) The name of the field to filter on. This can be any of the fields returned by the CloudStack API.
* `value` - (Required) The value to filter on. This should be a regular expression.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vGPU profile.
* `name` - The name of the vGPU profile.
* `description` - The description of the vGPU profile.
* `device_id` - The device id of the GPU card.
* `device_name` - The device name of the GPU card.
* `gpu_card_id` - The GPU card id of the vGPU profile.
* `gpu_card_name` - The GPU card name of the vGPU profile.
* `max_heads` - The maximum displays per vGPU instance.
* `max_resolution_x` - The maximum X resolution per display.
* `max_resolution_y` - The maximum Y resolution per display.
* `max_vgpu_per_physical_gpu` - The maximum number of vGPU instances per physical GPU.
* `vendor_id` - The vendor id of the GPU card.
* `vendor_name` - The vendor name of the GPU card.
* `video_ram` - The video RAM size in MB for the vGPU profile.
