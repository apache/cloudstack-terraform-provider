---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_role_permission"
description: |-
  Manages ordered role permissions for a role.
---

# cloudstack_role_permission

Manages an ordered list of role permissions for a role. A role permission is a
single rule that allows or denies access to an API or wildcard set of APIs.

Rules belonging to the same role are evaluated in order, and the first matching
rule wins. This resource stores that order explicitly and reapplies it after
rules are added, removed, or recreated.

By default, only the permissions declared in this resource are managed.
Undeclared permissions on the role are preserved and ordered after the declared
permissions. Set `authoritative = true` to delete undeclared permissions and
make the CloudStack role permission list exactly match this resource.

## Example Usage

```hcl
resource "cloudstack_role" "custom" {
  name = "custom-role"
  type = "User"
}

resource "cloudstack_role_permission" "custom" {
  role_id = cloudstack_role.custom.id

  permission {
    rule        = "listVirtualMachines"
    permission  = "allow"
    description = "Allow listing virtual machines"
  }

  permission {
    rule       = "*"
    permission = "deny"
  }
}
```

## Argument Reference

The following arguments are supported:

* `role_id` - (Required) ID of the role the permissions belong to. Changing this
  forces a new resource to be created.
* `authoritative` - (Optional) Whether permissions not declared in this resource
  should be deleted. Defaults to `false`.
* `permission` - (Optional) Ordered list of role permission rules. Each block
  supports the following:
  * `rule` - (Required) The API name or a wildcard (e.g. `list*` or `*`) the
    rule applies to.
  * `permission` - (Required) Whether the rule is allowed or denied. Valid
    options are: `allow`, `deny`.
  * `description` - (Optional) A description for the role permission.

## Attributes Reference

The following attributes are exported:

* `id` - The role ID.
* `permission.*.id` - The ID of each role permission rule.
