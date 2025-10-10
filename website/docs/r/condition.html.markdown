---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_condition"
sidebar_current: "docs-cloudstack-condition"
description: |-
  Creates a condition for autoscale policies.
---

# cloudstack_condition

Creates a condition that evaluates performance metrics against thresholds for autoscale policies.

## Example Usage

```hcl
# Reference an existing counter
data "cloudstack_counter" "cpu_counter" {
  id = "959e11c0-8416-11f0-9a72-1e001b000238"
}

resource "cloudstack_condition" "scale_up_condition" {
  counter_id          = data.cloudstack_counter.cpu_counter.id
  relational_operator = "GT"
  threshold           = 80.0
  account_name        = "admin"
  domain_id           = "1"
}

resource "cloudstack_condition" "scale_down_condition" {
  counter_id          = data.cloudstack_counter.cpu_counter.id
  relational_operator = "LT"
  threshold           = 20.0
  account_name        = "admin"
  domain_id           = "1"
}
```

## Argument Reference

The following arguments are supported:

* `counter_id` - (Required) The ID of the counter to monitor.

* `relational_operator` - (Required) The relational operator for the condition. Valid values are:
  * `"GT"` - Greater than
  * `"LT"` - Less than
  * `"EQ"` - Equal to
  * `"GE"` - Greater than or equal to
  * `"LE"` - Less than or equal to

* `threshold` - (Required) The threshold value to compare against.

* `account_name` - (Required) The account name that owns this condition.

* `domain_id` - (Required) The domain ID where the condition will be created.

## Attributes Reference

The following attributes are exported:

* `id` - The condition ID.
