---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_role_permission"
description: |-
  Creates a role permission (rule) for a role.
---

# cloudstack_role_permission

Creates a role permission. A role permission is a single rule that allows or
denies a role access to an API (or a wildcard set of APIs).

Rules belonging to the same role are evaluated in the order in which they are
created, and the first matching rule wins. Order the corresponding
`cloudstack_role_permission` resources accordingly (for example with
`depends_on`) when precedence matters.

## Example Usage

```hcl
resource "cloudstack_role" "custom" {
  name = "custom-role"
  type = "User"
}

# Allow listing virtual machines
resource "cloudstack_role_permission" "list_vms" {
  role_id     = cloudstack_role.custom.id
  rule        = "listVirtualMachines"
  permission  = "allow"
  description = "Allow listing virtual machines"
}

# Deny every other API using a wildcard
resource "cloudstack_role_permission" "deny_all" {
  role_id    = cloudstack_role.custom.id
  rule       = "*"
  permission = "deny"

  depends_on = [cloudstack_role_permission.list_vms]
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) ID of the role the permission belongs to. Changing this
  forces a new resource to be created.
* `rule` - (Required) The API name or a wildcard (e.g. `list*` or `*`) the rule
  applies to. Changing this forces a new resource to be created.
* `permission` - (Required) Whether the rule is allowed or denied. Valid options
  are: `allow`, `deny`.
* `description` - (Optional) A description for the role permission. Changing this
  forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role permission.
