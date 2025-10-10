---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_user_data"
sidebar_current: "docs-cloudstack-cloudstack_user_data"
description: |-
  Retrieves information about an existing CloudStack UserData definition.
---

# cloudstack_user_data

Use this data source to look up a registered CloudStack UserData definition so that its contents can be re-used across resources.

## Example Usage

```hcl
data "cloudstack_user_data" "cloudinit" {
  name = "bootstrap-userdata"

  # Optional filters
  account = "devops"
  project = "automation"
}

resource "cloudstack_instance" "vm" {
  # ... other arguments ...
  user_data = data.cloudstack_user_data.cloudinit.user_data
}
```

## Argument Reference

The following arguments are supported:

* `name` – (Required) The name of the UserData definition to retrieve.
* `account` – (Optional) Limits the lookup to a specific account.
* `project` – (Optional) Limits the lookup to a specific project. You may supply either the project name or ID.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` – The ID of the UserData definition.
* `account` – The owning account name.
* `account_id` – The owning account ID.
* `domain` – The domain name in which the UserData resides.
* `domain_id` – The domain ID in which the UserData resides.
* `params` – The optional parameters string associated with the UserData.
* `project` – The project name, when applicable.
* `user_data` – The decoded UserData contents.
