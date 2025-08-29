---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_limits"
sidebar_current: "docs-cloudstack-datasource-limits"
description: |-
  Gets information about CloudStack resource limits.
---

# cloudstack_limits

Use this data source to retrieve information about CloudStack resource limits for accounts, domains, and projects.

## Example Usage

```hcl
# Get all resource limits for a specific domain
data "cloudstack_limits" "domain_limits" {
  domainid = "domain-uuid"
}

# Get instance limits for a specific account
data "cloudstack_limits" "account_instance_limits" {
  type         = "instance"
  account      = "acct1"
  domainid     = "domain-uuid"
}

# Get primary storage limits for a project
data "cloudstack_limits" "project_storage_limits" {
  type         = "primarystorage"
  projectid    = "project-uuid"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Optional) The type of resource to list the limits. Available types are:
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
* `account` - (Optional) List resources by account. Must be used with the `domainid` parameter.
* `domainid` - (Optional) List only resources belonging to the domain specified.
* `projectid` - (Optional) List resource limits by project.

## Attributes Reference

The following attributes are exported:

* `limits` - A list of resource limits. Each limit has the following attributes:
  * `resourcetype` - The type of resource.
  * `resourcetypename` - The name of the resource type.
  * `max` - The maximum number of the resource. A value of `-1` indicates unlimited resources. A value of `0` means zero resources are allowed, though the CloudStack API may return `-1` for a limit set to `0`.
  * `account` - The account of the resource limit.
  * `domain` - The domain name of the resource limit.
  * `domainid` - The domain ID of the resource limit.
  * `project` - The project name of the resource limit.
  * `projectid` - The project ID of the resource limit.
