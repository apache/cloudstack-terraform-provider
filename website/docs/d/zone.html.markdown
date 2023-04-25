---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_zone"
sidebar_current: "docs-cloudstack-cloudstack_zone"
description: |-
  Gets information about cloudstack zone.
---

# cloudstack_zone

Use this datasource to get information about a zone for use in other resources.

### Example Usage

```hcl
  data "cloudstack_zone" "zone-data-source"{
    filter{
    name = "name"
    value="TestZone"
    }
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the zone.
* `dns1` - The first DNS for the Zone.
* `internal_dns1` - The first internal DNS for the Zone.
* `network_type` - The network type of the zone; can be Basic or Advanced.