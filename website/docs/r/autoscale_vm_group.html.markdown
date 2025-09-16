---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_autoscale_vm_group"
sidebar_current: "docs-cloudstack-autoscale-vm-group"
description: |-
  Creates an autoscale VM group.
---

# cloudstack_autoscale_vm_group

Creates an autoscale VM group that automatically scales virtual machines based on policies and load balancer rules.

## Example Usage

```hcl
# Basic autoscale VM group
resource "cloudstack_autoscale_vm_group" "vm_group" {
  name             = "web-server-autoscale"
  lbrule_id        = cloudstack_loadbalancer_rule.lb.id
  min_members      = 1
  max_members      = 5
  vm_profile_id    = cloudstack_autoscale_vm_profile.profile.id
  state            = "enable"  # or "disable"
  cleanup          = true
  
  scaleup_policy_ids = [
    cloudstack_autoscale_policy.scale_up_policy.id
  ]
  
  scaledown_policy_ids = [
    cloudstack_autoscale_policy.scale_down_policy.id
  ]
}

# Autoscale VM group with optional parameters
resource "cloudstack_autoscale_vm_group" "advanced_vm_group" {
  name             = "advanced-autoscale-group"
  lbrule_id        = cloudstack_loadbalancer_rule.lb.id
  min_members      = 2
  max_members      = 10
  vm_profile_id    = cloudstack_autoscale_vm_profile.profile.id
  interval         = 30    # Monitor every 30 seconds
  display          = true
  state            = "enable"
  cleanup          = false  # Keep VMs when deleting group
  
  scaleup_policy_ids = [
    cloudstack_autoscale_policy.cpu_scale_up.id,
    cloudstack_autoscale_policy.memory_scale_up.id
  ]
  
  scaledown_policy_ids = [
    cloudstack_autoscale_policy.scale_down.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `lbrule_id` - (Required) The ID of the load balancer rule. Changing this forces a new resource to be created.

* `min_members` - (Required) The minimum number of members in the VM group. The number of instances will be equal to or more than this number.

* `max_members` - (Required) The maximum number of members in the VM group. The number of instances will be equal to or less than this number.

* `vm_profile_id` - (Required) The ID of the autoscale VM profile that contains information about the VMs in the group. Changing this forces a new resource to be created.

* `scaleup_policy_ids` - (Required) A list of scale-up autoscale policy IDs.

* `scaledown_policy_ids` - (Required) A list of scale-down autoscale policy IDs.

* `name` - (Optional) The name of the autoscale VM group.

* `interval` - (Optional) The frequency in seconds at which performance counters are collected. Defaults to CloudStack's default interval.

* `display` - (Optional) Whether to display the group to the end user. Defaults to `true`.

* `state` - (Optional) The state of the autoscale VM group. Valid values are:
  * `"enable"` - Enable the autoscale group (default)
  * `"disable"` - Disable the autoscale group
  
  **Note**: When set to `"disable"`, the autoscale group stops monitoring and scaling, but existing VMs remain running.

* `cleanup` - (Optional) Whether all members of the autoscale VM group should be cleaned up when the group is deleted. Defaults to `false`.
  * `true` - Destroy all VMs when deleting the autoscale group
  * `false` - Leave VMs running when deleting the autoscale group

## Attributes Reference

The following attributes are exported:

* `id` - The autoscale VM group ID.

## State Management

The `state` parameter allows you to enable or disable the autoscale group without destroying it:

* **Enabling**: Changes from `"disable"` to `"enable"` will activate autoscale monitoring and allow automatic scaling.
* **Disabling**: Changes from `"enable"` to `"disable"` will stop autoscale monitoring but keep all existing VMs running.

This is useful for temporarily pausing autoscale behavior during maintenance or testing.

## Import

Autoscale VM groups can be imported using the `id`, e.g.

```
$ terraform import cloudstack_autoscale_vm_group.vm_group eb22f91-7454-4107-89f4-36afcdf33021
```
