---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_disk_offering"
sidebar_current: "docs-cloudstack-cloudstack_disk_offering"
description: |-
  Gets information about cloudstack disk offering.
---

# cloudstack_disk_offering

Use this datasource to get information about a disk offering for use in other resources.

### Example Usage

```hcl
data "cloudstack_disk_offering" "disk-offering-data-source" {
  filter {
    name  = "name"
    value = "custom"
  }
}
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the disk offering.
* `display_text` - The display text of the disk offering.
* `disk_size` - The size of the disk offering in GB.
* `customized` - Whether the disk offering allows a custom disk size.
* `storage_type` - The storage type of the disk offering.
* `provisioning_type` - The provisioning type of the disk offering.
* `tags` - The storage tags for the disk offering.
* `display_offering` - Whether the disk offering is displayed to the end user.
