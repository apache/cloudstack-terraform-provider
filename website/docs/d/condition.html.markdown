---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_condition"
sidebar_current: "docs-cloudstack-data-source-condition"
description: |-
  Gets information about a CloudStack autoscale condition.
---

# cloudstack_condition

Use this data source to get information about a CloudStack autoscale condition.

## Example Usage

```hcl
# Get condition by ID
data "cloudstack_condition" "existing_condition" {
  id = "c2f0591b-ce9b-499a-81f2-8fc6318b0c72"
}

# Get condition by filter
data "cloudstack_condition" "cpu_condition" {
  filter {
    name  = "threshold"
    value = "80"
  }
}

# Use in a policy
resource "cloudstack_autoscale_policy" "scale_policy" {
  name         = "scale-up-policy"
  action       = "scaleup"
  duration     = 300
  quiet_time   = 300
  condition_ids = [data.cloudstack_condition.existing_condition.id]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the condition.

* `filter` - (Optional) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `id` - The condition ID.

* `counter_id` - The counter ID being monitored.

* `relational_operator` - The relational operator for the condition.

* `threshold` - The threshold value.

* `account_name` - The account name that owns the condition.

* `domain_id` - The domain ID where the condition exists.
