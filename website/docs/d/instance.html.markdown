---
layout: "cloudstack"
page_title: "Cloudstack: cloudstack_instance"
sidebar_current: "docs-cloudstack-datasource-instance"
description: |-
  Gets information about cloudstack instance.
---

# cloudstack_instance

Use this datasource to get information about an instance for use in other resources.

### Example Usage

```hcl
data "cloudstack_instance" "my_instance" {
  filter {
    name = "name" 
    value = "server-a"
  }
  
  nic {
    ip_address="10.1.1.37"
  }
}
```

### Argument Reference

* `filter` - (Required) One or more name/value pairs to filter off of. You can apply filters on any exported attributes.

## Attributes Reference

The following attributes are exported:

* `instance_id` - The ID of the virtual machine.
* `account` - The account associated with the virtual machine.
* `display_name` - The user generated name. The name of the virtual machine is returned if no displayname exists.
* `state` - The state of the virtual machine.
* `host_id` - The ID of the host for the virtual machine.
* `zone_id` - The ID of the availability zone for the virtual machine.
* `created` - The date when this virtual machine was created.
* `nic` - The list of nics associated with vm.
