---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_counter"
sidebar_current: "docs-cloudstack-data-source-counter"
description: |-
  Gets information about a CloudStack counter.
---

# cloudstack_counter

Use this data source to get information about a CloudStack counter for use in autoscale conditions.

## Example Usage

```hcl
# Get counter by ID
data "cloudstack_counter" "cpu_counter" {
  id = "959e11c0-8416-11f0-9a72-1e001b000238"
}

# Get counter by name
data "cloudstack_counter" "memory_counter" {
  filter {
    name  = "name" 
    value = "VM CPU - average percentage"
  }
}

# Use in a condition
resource "cloudstack_condition" "scale_up" {
  counter_id          = data.cloudstack_counter.cpu_counter.id
  relational_operator = "GT"
  threshold           = 80.0
  account_name        = "admin"
  domain_id           = "1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the counter.

* `filter` - (Optional) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `id` - The counter ID.

* `name` - The name of the counter.

* `source` - The source of the counter.

* `value` - The metric value monitored by the counter.
