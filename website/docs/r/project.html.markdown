---
subcategory: "CloudStack"
layout: "cloudstack"
page_title: "CloudStack: cloudstack_project"
description: |-
  Creates a project.
---

# cloudstack_project

Creates a project.

## Example Usage

```hcl
resource "cloudstack_project" "myproject" {
  name         = "terraform-project"
  display_text = "Terraform Managed Project"
  domain       = "root"
}
```

### With Account and User ID

```hcl
resource "cloudstack_project" "myproject" {
  name         = "terraform-project"
  display_text = "Terraform Managed Project"
  domain       = "root"
  account      = "admin"
  userid       = "b0afc3ca-a99c-4fb4-98ad-8564acab10a4"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project.
* `display_text` - (Required) The display text of the project. Required for API version 4.18 and lower compatibility. This requirement will be removed when support for API versions older than 4.18 is dropped.
* `domain` - (Optional) The domain where the project will be created. This cannot be changed after the project is created.
* `account` - (Optional) The account who will be Admin for the project. Requires `domain` to be set. This can be updated after the project is created.
* `accountid` - (Optional) The ID of the account owning the project. This can be updated after the project is created.
* `userid` - (Optional) The user ID of the account to be assigned as owner of the project (Project Admin). This can be updated after the project is created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the project.
* `name` - The name of the project.
* `display_text` - The display text of the project.
* `domain` - The domain where the project was created.

## Import

Projects can be imported using the project ID, e.g.

```sh
terraform import cloudstack_project.myproject 5cf69677-7e4b-4bf4-b868-f0b02bb72ee0
```
