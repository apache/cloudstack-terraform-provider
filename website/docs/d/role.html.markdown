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
  filter {
    name = "name"
    value = "Admin"
  }
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

* `filter` - (Required) One or more name/value pairs to filter off of. See the example below for usage.

## Filter Example

```hcl
data "cloudstack_role" "admin" {
  filter {
    name = "name"
    value = "Admin"
  }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the role.
* `name` - The name of the role.
* `type` - The type of the role.
* `description` - The description of the role.
* `is_public` - Whether the role is public or not.
