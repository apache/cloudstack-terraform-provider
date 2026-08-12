---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_gpu_card"
description: |-
  Gets information about a GPU card.
---

# cloudstack_gpu_card

Use this data source to get information about a GPU card for use in other resources.

## Example Usage

```hcl
data "cloudstack_gpu_card" "card" {
  filter {
    name = "name"
    value = "NVIDIA.*"
  }
}

output "gpu_card_id" {
  value = data.cloudstack_gpu_card.card.id
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

* `id` - The ID of the GPU card.
* `name` - The name of the GPU card.
* `device_id` - The device id of the GPU card.
* `device_name` - The device name of the GPU card.
* `vendor_id` - The vendor id of the GPU card.
* `vendor_name` - The vendor name of the GPU card.
