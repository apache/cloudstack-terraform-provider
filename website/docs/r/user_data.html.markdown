---
layout: "cloudstack"
page_title: "Cloud* `user_data` - (Required) The user data script content. Can be plain text or base64 encoded.

**Note:** UserData content has size limitations:

* Base64 encoded content up to 4KB when using HTTP GET requests
* Base64 encoded content up to 1MB when using HTTP POST requests (used by default)
* The CloudStack global setting `vm.userdata.max.length` may impose additional limits

* `account` - (Optional) The account to register the UserData under. If not specified, uses the current account.oudstack_user_data"
sidebar_current: "docs-cloudstack-resource-user-data"
description: |-
    Creates a UserData resource
---

# CloudStack: cloudstack_user_data

A `cloudstack_user_data` resource manages a registered UserData object within CloudStack. This allows you to store and reuse user data scripts that can be referenced by virtual machines.

## Example Usage

```hcl
resource "cloudstack_user_data" "web_server_init" {
  name      = "web-server-bootstrap"
  user_data = <<-EOF
    #!/bin/bash
    yum update -y
    yum install -y httpd
    systemctl start httpd
    systemctl enable httpd
    echo "<h1>Hello from CloudStack!</h1>" > /var/www/html/index.html
  EOF
}

# Use the UserData in an instance
resource "cloudstack_instance" "web_server" {
  name             = "web-server"
  service_offering = "Medium Instance"
  network_id       = cloudstack_network.default.id
  template         = "CentOS 7"
  zone             = data.cloudstack_zone.zone.name
  user_data_id     = cloudstack_user_data.web_server_init.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the UserData.
* `user_data` - (Required) The user data script content. Can be plain text or base64 encoded.

**Note:** UserData content has size limitations:

- Base64 encoded content up to 4KB when using HTTP GET requests
- Base64 encoded content up to 1MB when using HTTP POST requests (used by default)
- The CloudStack global setting `vm.userdata.max.length` may impose additional limits

* `account` - (Optional) The account to register the UserData under. If not specified, uses the current account.
* `project` - (Optional) The project to register the UserData under.
* `params` - (Optional) Additional parameters for the UserData.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the UserData.
* `domain` - The domain of the account that owns the UserData.
* `domain_id` - The ID of the domain.
* `account_id` - The ID of the account that owns the UserData.

## Import

UserData can be imported using the UserData ID:

```shell
$ terraform import cloudstack_user_data.example <userdata-id>
```
