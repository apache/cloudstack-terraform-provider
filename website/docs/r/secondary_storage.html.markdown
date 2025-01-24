---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_secondary_storage"
sidebar_current: "docs-cloudstack-resource-secondary-storage"
description: |-
  Creates a changeme.
---

# cloudstack_secondary_storage

Create secondary storage

## Example Usage

Basic usage:

```hcl
resource "cloudstack_secondary_storage" "example" {
	name             = "example"
	storage_provider = "NFS"
	url              = "nfs://10.147.28.6:/export/home/sandbox/secondary"
	zone_id          = data.cloudstack_zone.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) the name for the image store.
* `storage_provider` - (Required) the image store provider name.
* `url` - (Optional) the URL for the image store.
* `zone_id` - (Optional) the Zone ID for the image store.


## Attributes Reference

The following attributes are exported:

* `id` - The instance ID.



## Import

changeme can be imported; use `<ZONE ID>` as the import ID. For
example:

```shell
terraform import cloudstack_secondary_storage.example 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
