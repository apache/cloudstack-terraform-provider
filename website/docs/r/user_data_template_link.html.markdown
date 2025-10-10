---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_user_data_template_link"
sidebar_current: "docs-cloudstack-resource-user-data-template-link"
description: |-
  Attaches an existing user data object to a template or ISO.
---

# cloudstack_user_data_template_link

Manages the association between a CloudStack template (or ISO) and a stored
user data object. Linking user data allows VMs created from the template to
receive the script automatically.

## Example Usage

### Link user data to a template

```hcl
resource "cloudstack_user_data" "bootstrap" {
  name      = "bootstrap"
  user_data = "#!/bin/bash\necho bootstrap > /var/tmp/bootstrap.log"
}

resource "cloudstack_template" "base" {
  name       = "base-template"
  format     = "QCOW2"
  hypervisor = "KVM"
  os_type    = "Ubuntu 22.04 (64-bit)"
  url        = "http://example.com/images/ubuntu-2204.qcow2"
}

resource "cloudstack_user_data_template_link" "base_userdata" {
  template_id      = cloudstack_template.base.id
  user_data_id     = cloudstack_user_data.bootstrap.id
  user_data_policy = "ALLOWOVERRIDE"
}
```

### Link user data to an ISO

```hcl
resource "cloudstack_user_data_template_link" "iso_userdata" {
  iso_id           = cloudstack_template.iso_image.id
  user_data_id     = cloudstack_user_data.bootstrap.id
  user_data_policy = "APPEND"
}
```

## Argument Reference

The following arguments are supported:

* `template_id` - (Optional) The ID of the template to link. Conflicts with
  `iso_id`. One of `template_id` or `iso_id` must be supplied. Changing the
  target template forces a new resource to be created.

* `iso_id` - (Optional) The ID of the ISO to link. Conflicts with `template_id`.
  Changing the target ISO forces a new resource to be created.

* `user_data_id` - (Optional) The ID of the user data object to associate. When
  omitted, the resource only ensures the template is present and can be used to
  remove an existing link via deletion.

* `user_data_policy` - (Optional) The policy that defines how the linked user
  data interacts with any user data provided during VM deployment. Valid values
  are `ALLOWOVERRIDE`, `APPEND`, and `DENYOVERRIDE`. Defaults to `ALLOWOVERRIDE`.

## Attributes Reference

The following attributes are exported:

* `id` - The template or ISO ID reported by CloudStack.
* `name` - The name of the template or ISO.
* `display_text` - The display text of the template or ISO.
* `is_ready` - Whether the template or ISO is ready for use.
* `template_type` - The CloudStack template type.
* `user_data_name` - The name of the linked user data object, if any.
* `user_data_params` - The parameters stored with the linked user data.

## Import

User data template links can be imported using the CloudStack identifier of the
linked template or ISO, e.g.

```shell
terraform import cloudstack_user_data_template_link.base_userdata template-8a7d5c64-0605-4ed6-b3a2-7c91c7f5d4bb
```
