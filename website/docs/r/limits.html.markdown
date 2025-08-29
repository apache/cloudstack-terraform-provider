---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_limits"
sidebar_current: "docs-cloudstack-limits"
description: |-
  Provides a CloudStack limits resource.
---

# cloudstack_limits

Provides a CloudStack limits resource. This can be used to manage resource limits for accounts, domains, and projects within CloudStack.

## Example Usage

```hcl
# Set instance limit for the root domain
resource "cloudstack_limits" "instance_limit" {
  type         = "instance"
  max          = 20
}

# Set volume limit for a specific account in a domain
resource "cloudstack_limits" "volume_limit" {
  type         = "volume"
  max          = 50
  account      = "acct1"
  domainid     = "domain-uuid"
}

# Set primary storage limit for a project
resource "cloudstack_limits" "storage_limit" {
  type         = "primarystorage"
  max          = 1000  # GB
  projectid    = "project-uuid"
}

# Set unlimited CPU limit
resource "cloudstack_limits" "cpu_unlimited" {
  type         = "cpu"
  max          = -1  # Unlimited
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, ForceNew) The type of resource to update. Available types are:
  * `instance`
  * `ip`
  * `volume`
  * `snapshot`
  * `template`
  * `project`
  * `network`
  * `vpc`
  * `cpu`
  * `memory`
  * `primarystorage`
  * `secondarystorage`

* `account` - (Optional, ForceNew) Update resource for a specified account. Must be used with the `domainid` parameter.
* `domainid` - (Optional, ForceNew) Update resource limits for all accounts in specified domain. If used with the `account` parameter, updates resource limits for a specified account in specified domain.
* `max` - (Optional) Maximum resource limit. Use `-1` for unlimited resource limit. A value of `0` means zero resources are allowed, though the CloudStack API may return `-1` for a limit set to `0`.
* `projectid` - (Optional, ForceNew) Update resource limits for project.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `type` - The type of resource.
* `max` - The maximum number of the resource.
* `account` - The account of the resource limit.
* `domainid` - The domain ID of the resource limit.
* `projectid` - The project ID of the resource limit.

## Import

Resource limits can be imported using the resource type (numeric), account, domain ID, and project ID, e.g.

```bash
terraform import cloudstack_limits.instance_limit 0
terraform import cloudstack_limits.volume_limit 2-acct1-domain-uuid
terraform import cloudstack_limits.storage_limit 10-project-uuid
```

When importing, the numeric resource type is used in the import ID. The provider will automatically convert the numeric type to the corresponding string type after import.
