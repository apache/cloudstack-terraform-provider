---
subcategory: ""
layout: "cloudstack"
page_title: "CloudStack: cloudstack_role"
description: |-
  Gets information about a role.
---

# cloudstack_role

Use this data source to get information about a role for use in other resources.

## Example Usage

```hcl
data "cloudstack_role" "admin" {
  name = "Admin"
}

resource "cloudstack_account" "example" {
  email       = "example@example.com"
  first_name  = "John"
  last_name   = "Doe"
  password    = "password"
  username    = "johndoe"
  account_type = 1
  role_id     = data.cloudstack_role.admin.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the role.
* `name` - (Optional) The name of the role.

At least one of the above arguments is required.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role.
* `name` - The name of the role.
* `type` - The type of the role.
* `description` - The description of the role.
* `is_public` - Whether the role is public or not.
