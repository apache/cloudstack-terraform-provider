---
layout: default
title: "CloudStack: cloudstack_user"
sidebar_current: "docs-cloudstack-resource-user"
description: |-
    Creates a User
---

# CloudStack: cloudstack_user

A `cloudstack_user` resource manages a user within CloudStack.

## Example Usage

```hcl
resource "cloudstack_user" "example" {
    account = "example-account"
    email = "user@example.com"
    first_name = "John"
    last_name = "Doe"
    password = "securepassword"
    username = "jdoe"
}
```


## Argument Reference

The following arguments are supported:

* `account` - (Optional) The account the user belongs to.
* `email` - (Required) The email address of the user.
* `first_name` - (Required) The first name of the user.
* `last_name` - (Required) The last name of the user.
* `password` - (Required) The password for the user.
* `username` - (Required) The username of the user.

## Attributes Reference

No attributes are exported.

## Import

Users can be imported; use `<USERID>` as the import ID. For example:

```shell
$ terraform import cloudstack_user.example <USERID>
```
