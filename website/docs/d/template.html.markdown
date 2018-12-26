---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_template"
sidebar_current: "docs-cloudstack-datasource-template"
description: |-
  Get informations on a Cloudstack template.
---

# cloudstack_template

Use this datasource to get the ID of a template for use in other resources.

### Example Usage

```hcl
data "cloudstack_template" "my_template" {
  template_filter = "featured"

  filter {
    name = "name"
    value = "CentOS 7\\.1"
  }

  filter {
    name = "hypervisor"
    value = "KVM"
  }
}
```

### Argument Reference

* `template_filter` - (Required) The template filter. Possible values are `featured`, `self`, `selfexecutable`, `sharedexecutable`, `executable` and `community` (see the Cloudstack API *listTemplate* command documentation).

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `id` - The template ID.
* `account` - The account name to which the template belongs.
* `created` - The date this template was created.
* `display_text` - The template display text.
* `format` - The format of the template.
* `hypervisor` - The hypervisor on which the templates runs.
* `name` - The template name.
* `size` - The size of the template.
