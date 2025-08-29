---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_project"
sidebar_current: "docs-cloudstack-cloudstack_project"
description: |-
  Gets information about CloudStack projects.
---

# cloudstack_project

Use this datasource to get information about a CloudStack project for use in other resources.

## Example Usage

### Basic Usage

```hcl
data "cloudstack_project" "my_project" {
  filter {
    name = "name"
    value = "my-project"
  }
}
```

### With Multiple Filters

```hcl
data "cloudstack_project" "admin_project" {
  filter {
    name = "name"
    value = "admin-project"
  }
  filter {
    name = "domain"
    value = "ROOT"
  }
  filter {
    name = "account"
    value = "admin"
  }
}
```

## Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the project.
* `name` - The name of the project.
* `display_text` - The display text of the project.
* `domain` - The domain where the project belongs.
* `account` - The account who is the admin for the project.
* `state` - The current state of the project.
* `tags` - A map of tags assigned to the project.