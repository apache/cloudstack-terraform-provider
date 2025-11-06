---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_security_group"
sidebar_current: "docs-cloudstack-resource-security-group"
description: |-
  Creates a security group.
---

# cloudstack_security_group

Creates a security group.

## Example Usage

```hcl
resource "cloudstack_security_group" "default" {
  name        = "allow_web"
  description = "Allow access to HTTP and HTTPS"
}
```

### With Account and Domain

```hcl
resource "cloudstack_security_group" "account_sg" {
  name        = "allow_web"
  description = "Allow access to HTTP and HTTPS"
  account     = "my-account"
  domain      = "example-domain"
}
```

### With Project

```hcl
resource "cloudstack_project" "my_project" {
  name        = "my-project"
  displaytext = "My Project"
}

resource "cloudstack_security_group" "project_sg" {
  name        = "allow_web"
  description = "Allow access to HTTP and HTTPS"
  project_id  = cloudstack_project.my_project.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security group. Changing this forces a
    new resource to be created.

* `description` - (Optional) The description of the security group. Changing
    this forces a new resource to be created.

* `account` - (Optional) The account name to create the security group for.
    Must be used with `domain`. Cannot be used with `project_id`. Changing this
    forces a new resource to be created.

* `domain` - (Optional) The name or ID of the domain to create this security
    group in. Changing this forces a new resource to be created.

* `project_id` - (Optional) The ID of the project to create this security
    group in. Cannot be used with `account`. Changing this forces a new
    resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group.

## Import

Security groups can be imported; use `<SECURITY GROUP ID>` as the import ID. For
example:

```shell
terraform import cloudstack_security_group.default e54970f1-f563-46dd-a365-2b2e9b78c54b
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_security_group.default my-project/e54970f1-f563-46dd-a365-2b2e9b78c54b
```
