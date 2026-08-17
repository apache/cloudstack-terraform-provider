---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_autoscale_policy"
sidebar_current: "docs-cloudstack-autoscale-policy"
description: |-
  Creates an autoscale policy.
---

# cloudstack_autoscale_policy

Creates an autoscale policy that defines when and how to scale virtual machines based on conditions.

## Example Usage

```hcl
resource "cloudstack_condition" "scale_up_condition" {
  counter_id          = data.cloudstack_counter.cpu_counter.id
  relational_operator = "GT"
  threshold           = 80.0
  account_name        = "admin"
  domain_id           = "67bc8dbe-8416-11f0-9a72-1e001b000238"
}

resource "cloudstack_autoscale_policy" "scale_up_policy" {
  name         = "scale-up-policy"
  action       = "scaleup"  # Case insensitive: scaleup/SCALEUP
  duration     = 300        # 5 minutes
  quiet_time   = 300        # 5 minutes
  condition_ids = [cloudstack_condition.scale_up_condition.id]
}

resource "cloudstack_autoscale_policy" "scale_down_policy" {
  name         = "scale-down-policy"
  action       = "scaledown"  # Case insensitive: scaledown/SCALEDOWN
  duration     = 300
  quiet_time   = 600          # 10 minutes quiet time
  condition_ids = [cloudstack_condition.scale_down_condition.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the autoscale policy.

* `action` - (Required) The action to be executed when conditions are met. Valid values are:
  * `"scaleup"` or `"SCALEUP"` - Scale up (add instances)
  * `"scaledown"` or `"SCALEDOWN"` - Scale down (remove instances)
  
  **Note**: The action field is case-insensitive.

* `duration` - (Required) The duration in seconds for which the conditions must be true before the action is taken.

* `quiet_time` - (Optional) The cool down period in seconds during which the policy should not be evaluated after the action has been taken.

* `condition_ids` - (Required) A list of condition IDs that must all evaluate to true for the policy to trigger.

## Attributes Reference

The following attributes are exported:

* `id` - The autoscale policy ID.
