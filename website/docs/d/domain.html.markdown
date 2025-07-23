---
layout: default
page_title: "CloudStack: cloudstack_domain Data Source"
sidebar_current: "docs-cloudstack-datasource-domain"
description: |-
    Retrieves information about a Domain
---

# CloudStack: cloudstack_domain Data Source

A `cloudstack_domain` data source retrieves information about a domain within CloudStack.

## Example Usage

```hcl
data "cloudstack_domain" "my_domain" {
  filter {
    name = "name"
    value = "ROOT"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) A block to filter the domains. The filter block supports the following:
  * `name` - (Required) The name of the filter.
  * `value` - (Required) The value of the filter.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the domain.
* `name` - The name of the domain.
* `network_domain` - The network domain for the domain.
* `parent_domain_id` - The ID of the parent domain.
