---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_network_acl"
sidebar_current: "docs-cloudstack-resource-network-acl"
description: |-
  Creates a Network ACL for the given VPC.
---

# cloudstack_network_acl

Creates a Network ACL for the given VPC.

## Example Usage

```hcl
resource "cloudstack_network_acl" "default" {
  name   = "test-acl"
  vpc_id = "76f6e8dc-07e3-4971-b2a2-8831b0cc4cb4"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ACL. Changing this forces a new resource
    to be created.

* `description` - (Optional) The description of the ACL. Changing this forces a
    new resource to be created.

* `project` - (Optional) The name or ID of the project to deploy this
    instance to. Changing this forces a new resource to be created.

* `vpc_id` - (Required) The ID of the VPC to create this ACL for. Changing this
   forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Network ACL

## Import

Network ACLs can be imported; use `<NETWORK ACL ID>` as the import ID. For
example:

```shell
terraform import cloudstack_network_acl.default e8b5982a-1b50-4ea9-9920-6ea2290c7359
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_network_acl.default my-project/e8b5982a-1b50-4ea9-9920-6ea2290c7359
```
