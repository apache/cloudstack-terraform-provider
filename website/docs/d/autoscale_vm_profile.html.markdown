---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_autoscale_vm_profile"
sidebar_current: "docs-cloudstack-data-source-autoscale-vm-profile"
description: |-
  Gets information about a CloudStack autoscale VM profile.
---

# cloudstack_autoscale_vm_profile

Use this data source to get information about a CloudStack autoscale VM profile.

## Example Usage

```hcl
# Get VM profile by ID
data "cloudstack_autoscale_vm_profile" "existing_profile" {
  id = "a596f7a2-95b8-4f0e-9f15-88f4091f18fe"
}

# Get VM profile by filter
data "cloudstack_autoscale_vm_profile" "web_profile" {
  filter {
    name  = "service_offering"
    value = "Small Instance"
  }
}

# Use in an autoscale VM group
resource "cloudstack_autoscale_vm_group" "vm_group" {
  name             = "web-autoscale"
  lbrule_id        = cloudstack_loadbalancer_rule.lb.id
  min_members      = 1
  max_members      = 5
  vm_profile_id    = data.cloudstack_autoscale_vm_profile.existing_profile.id
  
  scaleup_policy_ids = [cloudstack_autoscale_policy.scale_up.id]
  scaledown_policy_ids = [cloudstack_autoscale_policy.scale_down.id]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the autoscale VM profile.

* `filter` - (Optional) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `id` - The autoscale VM profile ID.

* `service_offering` - The service offering name or ID.

* `template` - The template name or ID.

* `zone` - The zone name or ID.

* `destroy_vm_grace_period` - The grace period for VM destruction.

* `counter_param_list` - Counter parameters for monitoring.

* `user_data` - User data for VM initialization.

* `user_data_details` - Additional user data details.

* `account_name` - The account name that owns the profile.

* `domain_id` - The domain ID where the profile exists.

* `display` - Whether the profile is displayed to end users.

* `other_deploy_params` - Additional deployment parameters.
