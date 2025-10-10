---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_autoscale_vm_group"
sidebar_current: "docs-cloudstack-data-source-autoscale-vm-group"
description: |-
  Gets information about a CloudStack autoscale VM group.
---

# cloudstack_autoscale_vm_group

Use this data source to get information about a CloudStack autoscale VM group.

## Example Usage

```hcl
# Get autoscale VM group by ID
data "cloudstack_autoscale_vm_group" "existing_group" {
  id = "156a819a-dec1-4166-aab3-657c271fa4a3"
}

# Get autoscale VM group by name
data "cloudstack_autoscale_vm_group" "web_group" {
  filter {
    name  = "name"
    value = "web-server-autoscale"
  }
}

# Output information about the group
output "autoscale_group_state" {
  value = data.cloudstack_autoscale_vm_group.existing_group.state
}

output "current_members" {
  value = "Min: ${data.cloudstack_autoscale_vm_group.existing_group.min_members}, Max: ${data.cloudstack_autoscale_vm_group.existing_group.max_members}"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the autoscale VM group.

* `filter` - (Optional) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `id` - The autoscale VM group ID.

* `name` - The name of the autoscale VM group.

* `lbrule_id` - The load balancer rule ID.

* `min_members` - The minimum number of members.

* `max_members` - The maximum number of members.

* `vm_profile_id` - The VM profile ID.

* `interval` - The monitoring interval in seconds.

* `display` - Whether the group is displayed to end users.

* `state` - The current state of the group (enable or disable).

* `scaleup_policy_ids` - The list of scale-up policy IDs.

* `scaledown_policy_ids` - The list of scale-down policy IDs.
