---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_affinity_group"
sidebar_current: "docs-cloudstack-resource-affinity-group"
description: |-
  Creates an affinity group.
---

# cloudstack_affinity_group

Creates an affinity group.

## Example Usage

```hcl
resource "cloudstack_affinity_group" "default" {
  name = "test-affinity-group"
  type = "host anti-affinity"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the affinity group. Changing this
    forces a new resource to be created.

* `description` - (Optional) The description of the affinity group.

* `type` - (Required) The affinity group type. Changing this
    forces a new resource to be created.

* `project` - (Optional) The name or ID of the project to register this
    affinity group to. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the affinity group.
* `description` - The description of the affinity group.

## Import

Affinity groups can be imported; use `<AFFINITY GROUP ID>` as the import ID. For
example:

```shell
terraform import cloudstack_affinity_group.default 6226ea4d-9cbe-4cc9-b30c-b9532146da5b
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_affinity_group.default my-project/6226ea4d-9cbe-4cc9-b30c-b9532146da5b
```

