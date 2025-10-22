---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_user_data"
sidebar_current: "docs-cloudstack-datasource-user-data"
description: |-
  Get information about a CloudStack user data.
---

# cloudstack_user_data

Use this data source to retrieve information about a CloudStack user data by either its name or ID.

## Example Usage

### Find User Data by Name

```hcl
data "cloudstack_user_data" "web_init" {
  filter {
    name  = "name"
    value = "web-server-init"
  }
}

# Use the user data in an instance
resource "cloudstack_instance" "web" {
  name        = "web-server"
  userdata_id = data.cloudstack_user_data.web_init.id
  # ... other arguments ...
}
```

### Find User Data by ID

```hcl
data "cloudstack_user_data" "app_init" {
  filter {
    name  = "id"
    value = "12345678-1234-1234-1234-123456789012"
  }
}
```

### Find Project-Scoped User Data

```hcl
data "cloudstack_user_data" "project_init" {
  project = "my-project"
  
  filter {
    name  = "name"
    value = "project-specific-init"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply multiple filters to narrow down the results. See [Filters](#filters) below for more details.

* `project` - (Optional) The name or ID of the project to search in.

### Filters

The `filter` block supports the following arguments:

* `name` - (Required) The name of the filter. Valid filter names are:
  * `id` - Filter by user data ID
  * `name` - Filter by user data name
  * `account` - Filter by account name
  * `domainid` - Filter by domain ID

* `value` - (Required) The value to filter by.

## Attributes Reference

The following attributes are exported:

* `id` - The user data ID.
* `name` - The name of the user data.
* `userdata` - The user data content.
* `account` - The account name owning the user data.
* `domain_id` - The domain ID where the user data belongs.
* `project_id` - The project ID if the user data is project-scoped.
* `params` - The list of parameter names defined in the user data (comma-separated string).

## Example with Template Integration

```hcl
# Find existing user data
data "cloudstack_user_data" "app_bootstrap" {
  filter {
    name  = "name"
    value = "application-bootstrap"
  }
}

# Use with template
resource "cloudstack_template" "app_template" {
  name         = "application-template"
  display_text = "Application Template with Bootstrap"
  # ... other template arguments ...
  
  userdata_link {
    userdata_id     = data.cloudstack_user_data.app_bootstrap.id
    userdata_policy = "ALLOWOVERRIDE"
  }
}

# Deploy instance with parameterized user data
resource "cloudstack_instance" "app_server" {
  name     = "app-server-01"
  template = cloudstack_template.app_template.id
  # ... other instance arguments ...
  
  userdata_id = data.cloudstack_user_data.app_bootstrap.id
  
  userdata_details = {
    "environment" = "production"
    "app_version" = "v2.1.0"
    "debug_enabled" = "false"
  }
}
```
