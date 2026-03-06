---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_static_route"
sidebar_current: "docs-cloudstack-resource-static-route"
description: |-
  Creates a static route.
---

# cloudstack_static_route

Creates a static route for the given private gateway or VPC.

## Example Usage

Using a private gateway:

```hcl
resource "cloudstack_static_route" "default" {
  cidr       = "10.0.0.0/16"
  gateway_id = "76f607e3-e8dc-4971-8831-b2a2b0cc4cb4"
}
```

Using a nexthop IP address:

```hcl
resource "cloudstack_static_route" "with_nexthop" {
  cidr    = "10.0.0.0/16"
  nexthop = "192.168.1.1"
  vpc_id  = "76f607e3-e8dc-4971-8831-b2a2b0cc4cb4"
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The CIDR for the static route. Changing this forces
    a new resource to be created.

* `gateway_id` - (Optional) The ID of the Private gateway. Changing this forces
    a new resource to be created. Conflicts with `nexthop` and `vpc_id`.

* `nexthop` - (Optional) The IP address of the nexthop for the static route.
    Changing this forces a new resource to be created. Conflicts with `gateway_id`.
    Must be used together with `vpc_id`.

* `vpc_id` - (Optional) The ID of the VPC. Required when using `nexthop`.
    Changing this forces a new resource to be created. Conflicts with `gateway_id`.

**Note:** Either `gateway_id` or (`nexthop` + `vpc_id`) must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the static route.
