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
resource "cloudstack_role" "admin" {
  name        = "Admin"
  type        = "Admin"
  description = "Administrator role"
  is_public   = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the role.
* `type` - (Optional) The type of the role. Defaults to the CloudStack default.
* `description` - (Optional) The description of the role.
* `is_public` - (Optional) Whether the role is public or not. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role.

## Import

Roles can be imported using the role ID, e.g.

```
terraform import cloudstack_role.admin 12345678-1234-1234-1234-123456789012
