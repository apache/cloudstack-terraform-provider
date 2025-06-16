---
subcategory: ""
layout: "cloudstack"
page_title: "CloudStack: cloudstack_role"
description: |-
  Creates a role.
---

# cloudstack_role

Creates a role.

## Example Usage

```hcl
# Create a role with a specific type
resource "cloudstack_role" "admin" {
  name        = "Admin"
  type        = "Admin"
  description = "Administrator role"
  is_public   = true
}

# Create a role by cloning an existing role
resource "cloudstack_role" "custom_admin" {
  name        = "CustomAdmin"
  role_id     = "12345678-1234-1234-1234-123456789012"
  description = "Custom administrator role cloned from an existing role"
  is_public   = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the role.
* `type` - (Optional) The type of the role. Valid options are: Admin, ResourceAdmin, DomainAdmin, User. Either `type` or `role_id` must be specified.
* `description` - (Optional) The description of the role.
* `is_public` - (Optional) Whether the role is public or not. Defaults to `true`.
* `role_id` - (Optional) ID of the role to be cloned from. Either `role_id` or `type` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role.

## Import

Roles can be imported using the role ID, e.g.

```
terraform import cloudstack_role.admin 12345678-1234-1234-1234-123456789012
