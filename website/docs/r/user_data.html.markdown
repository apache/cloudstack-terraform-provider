---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_user_data"
sidebar_current: "docs-cloudstack-resource-user-data"
description: |-
  Manages reusable user data scripts that can be linked to templates or instances.
---

# cloudstack_user_data

Registers a reusable piece of user data in CloudStack. The stored script can be
linked to templates or referenced by instances that support user data.

## Example Usage

```hcl
resource "cloudstack_user_data" "bootstrap" {
  name      = "bootstrap-script"
  user_data = <<-EOF
    #!/bin/bash
    echo "Hello from Terraform" > /var/tmp/hello.txt
  EOF

  params  = "key=value"
  project = "infra-project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The display name for the user data object.

* `user_data` - (Required) The script or payload to store. The provider handles
  Base64 encoding and validates the CloudStack size limit (1 MB encoded).

* `account` - (Optional) The name of the account that owns the user data. Changing
  this forces a new resource to be created.

* `project` - (Optional) The name or ID of the project that owns the user data.
  Changing this forces a new resource to be created.

* `params` - (Optional) Additional parameters that will be passed to the user
  data when it is executed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the user data object.
* `account_id` - The ID of the owning account.
* `domain` - The name of the domain that owns the user data.
* `domain_id` - The ID of the domain that owns the user data.

## Import

User data can be imported using the `id`, e.g.

```shell
terraform import cloudstack_user_data.bootstrap 2c1bab14-5fcb-4b52-bdba-2f7d4a4fb916
```
