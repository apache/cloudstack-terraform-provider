---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_service_offering"
sidebar_current: "docs-cloudstack-cloudstack_service_offering"
description: |-
  Gets information about cloudstack service offering.
---

# cloudstack_service_offering

Use this datasource to get information about a service offering for use in other resources.

### Example Usage

```hcl
    data "cloudstack_service_offering" "service-offering-data-source"{
    filter{
    name = "name"
    value = "TestServiceUpdate"
    }  
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the service offering.
* `display_text` - An alternate display text of the service offering.