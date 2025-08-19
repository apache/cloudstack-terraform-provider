---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_user_data"
sidebar_current: "docs-cloudstack-datasource-user-data"
description: |-
  Gets information about a CloudStack UserData.
---

# cloudstack_user_data

Use this data source to get information about a CloudStack UserData for use in other resources.

## Example Usage

```hcl
data "cloudstack_user_data" "bootstrap" {
  filter {
    name  = "name"
    value = "web-server-bootstrap"
  }
}

resource "cloudstack_instance" "web" {
  name             = "web-server"
  service_offering = "Medium Instance"
  network_id       = cloudstack_network.default.id
  template         = "CentOS 7"
  zone             = data.cloudstack_zone.zone.name
  user_data_id     = data.cloudstack_user_data.bootstrap.id
}
```

**Note:** UserData content has size limitations:

- Base64 encoded content up to 4KB when using HTTP GET requests
- Base64 encoded content up to 1MB when using HTTP POST requests (default)
- The CloudStack global setting `vm.userdata.max.length` may impose additional limits

## Argument Reference

- `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the UserData.
- `name` - The name of the UserData.
- `user_data` - The user data script content (decoded if base64 encoded).
- `account` - The account that owns the UserData.
- `account_id` - The ID of the account that owns the UserData.
- `domain` - The domain of the account.
- `domain_id` - The ID of the domain.
- `params` - Additional parameters of the UserData.
- `project` - The project that owns the UserData.
- `project_id` - The ID of the project.
