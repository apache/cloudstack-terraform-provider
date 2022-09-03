---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_user"
sidebar_current: "docs-cloudstack-cloudstack_user"
description: |-
  Gets information about cloudstack user.
---

# cloudstack_user

Use this datasource to get information about a cloudstack user for use in other resources.

### Example Usage

```hcl
data "cloudstack_user" "user-data-source"{
    filter{
    name = "first_name"
    value= "jon"
    }
  }
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `account` - The account name of the userg.
* `email` - The user email address.
* `first_name` - The user firstname.
* `last_name` - The user lastname.
* `username` - The user name