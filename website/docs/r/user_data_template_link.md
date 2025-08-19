---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_user_data_template_link"
sidebar_current: "docs-cloudstack-resource-user-data-template-link"
description: |-
  Links UserData to a Template or ISO in CloudStack.
---

# cloudstack_user_data_template_link

Links UserData to a Template or ISO in CloudStack. This allows you to associate
user data with templates or ISOs, which will be used when instances are created
from those templates/ISOs.

## Example Usage

### Link UserData to Template

```hcl
resource "cloudstack_user_data" "web_config" {
  name      = "web-server-config"
  user_data = base64encode(<<-EOF
    #!/bin/bash
    yum update -y
    yum install -y httpd
    systemctl start httpd
    systemctl enable httpd
    echo "<h1>Hello from CloudStack</h1>" > /var/www/html/index.html
  EOF
  )
}

resource "cloudstack_template" "web_template" {
  name              = "web-server-template"
  format            = "QCOW2"
  hypervisor        = "KVM"
  os_type           = "CentOS 7"
  url               = "https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img"
  is_extractable    = true
  is_featured       = false
  is_public         = false
  password_enabled  = false
}

resource "cloudstack_user_data_template_link" "web_link" {
  template_id      = cloudstack_template.web_template.id
  user_data_id     = cloudstack_user_data.web_config.id
  user_data_policy = "ALLOWOVERRIDE"
}
```

### Link UserData to ISO

```hcl
resource "cloudstack_user_data" "iso_config" {
  name      = "iso-config"
  user_data = base64encode("#!/bin/bash\necho 'ISO Configuration' > /tmp/iso-config.txt")
}

resource "cloudstack_template" "custom_iso" {
  name              = "custom-iso"
  format            = "ISO"
  hypervisor        = "KVM"
  os_type           = "Other"
  url               = "http://releases.ubuntu.com/20.04/ubuntu-20.04.3-live-server-amd64.iso"
  is_extractable    = true
  is_featured       = false
  is_public         = false
  password_enabled  = false
}

resource "cloudstack_user_data_template_link" "iso_link" {
  iso_id           = cloudstack_template.custom_iso.id
  user_data_id     = cloudstack_user_data.iso_config.id
  user_data_policy = "APPEND"
}
```

## Argument Reference

The following arguments are supported:

* `template_id` - (Optional) The ID of the template to link UserData to.
  Conflicts with `iso_id`. Either `template_id` or `iso_id` must be specified.

* `iso_id` - (Optional) The ID of the ISO to link UserData to.
  Conflicts with `template_id`. Either `template_id` or `iso_id` must be specified.

* `user_data_id` - (Optional) The ID of the UserData to link. If not specified,
  any existing UserData link will be removed.

* `user_data_policy` - (Optional) The policy for how UserData should be handled.
  Valid values are:
  * `ALLOWOVERRIDE` - (Default) Allow instance-level user data to override template user data
  * `APPEND` - Append instance-level user data to template user data
  * `DENYOVERRIDE` - Deny instance-level user data override; only use template user data

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the template or ISO.

* `name` - The name of the template or ISO.

* `display_text` - The display text of the template or ISO.

* `is_ready` - Whether the template or ISO is ready for use.

* `template_type` - The type of the template (e.g., "USER", "BUILTIN").

* `user_data_name` - The name of the linked UserData.

* `user_data_params` - The parameters of the linked UserData.

## Import

UserData Template Links can be imported using the template or ISO ID:

```shell
$ terraform import cloudstack_user_data_template_link.example template-12345678
```

Or for ISOs:

```shell
$ terraform import cloudstack_user_data_template_link.example iso-87654321
```
