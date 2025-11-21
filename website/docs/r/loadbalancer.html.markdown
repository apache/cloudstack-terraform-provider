---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_loadbalancer"
sidebar_current: "docs-cloudstack-resource-loadbalancer"
description: |-
  Creates an internal load balancer.
---

# cloudstack_loadbalancer

Creates an internal load balancer.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_loadbalancer" "example" {
  algorithm                = "Source"
  description              = "Example Load Balancer"
  instanceport             = "8081"
  name                     = "internal-lb-example"
  networkid                = "0ae8fa84-c78e-441a-8628-917d276c7d5c"
  scheme                   = "Internal"
  sourceipaddressnetworkid = "0ae8fa84-c78e-441a-8628-917d276c7d5c"
  sourceport               = "8081"
  virtualmachineids        = [ "48600f99-d890-472c-bc1a-7379af22727c", "749485a1-4081-49cd-9668-1d160ac94488", "9bbf3c6d-c8b8-42e8-ab37-e4f60a71f61d" ]
}
```

## Argument Reference

The following arguments are supported:

* `algorithm` - (Required) load balancer algorithm (source, roundrobin, leastconn).
* `description` - (Optional) the description of the load balancer.
* `instanceport` - (Required) the TCP port of the virtual machine where the network traffic will be load balanced to.
* `name` - (Required) name of the load balancer.
* `networkid` - (Required) The guest network the load balancer will be created for.
* `scheme` - (Required) the load balancer scheme. Supported value in this release is Internal.
* `sourceipaddressnetworkid` - (Required) the network id of the source ip address.
* `sourceport` - (Required) the source port the network traffic will be load balanced from.
* `sourceipaddress` - (Optional) the source IP address the network traffic will be load balanced from.
* `virtualmachineids` - (Optional) the list of IDs of the virtual machine that are being assigned to the load balancer rule(i.e. virtualMachineIds=1,2,3).


## Attributes Reference

The following attributes are exported:

* `id` - the Load Balancer ID.



## Import

A load balancer can be imported; use `<LOADBALANCER ID>` as the import ID. For
example:

```shell
terraform import cloudstack_loadbalancer.example eefce154-e759-45a4-989a-9a432792801b
```
