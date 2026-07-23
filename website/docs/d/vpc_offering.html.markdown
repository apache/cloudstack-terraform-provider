---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_vpc_offering"
sidebar_current: "docs-cloudstack-cloudstack_vpc_offering"
description: |-
  Gets information about cloudstack VPC offering.
---

# cloudstack_vpc_offering

Use this datasource to get information about a VPC offering for use in other resources.

### Example Usage

```hcl
  data "cloudstack_vpc_offering" "vpc-off-data-source"{
    filter{
    name = "name"
    value="Default VPC offering"
    }
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the VPC offering.
* `display_text` - An alternate display text of the VPC offering.
* `enable` - Whether the VPC offering is enabled.
* `for_nsx` - Whether this VPC offering is meant to be used for NSX.
* `internet_protocol` - The internet protocol supported by the VPC offering.
* `network_mode` - The mode with which the VPC will operate.
* `routing_mode` - The routing mode for the VPC offering.
* `specify_as_number` - Whether the VPC offering supports choosing an AS number.
* `is_default` - `true` if this is the default VPC offering.
* `supported_services` - The list of supported services for this VPC offering.
* `service_provider_list` - The map of service providers for the supported services.
* `service_capability_list` - The map of service capabilities for this VPC offering.
