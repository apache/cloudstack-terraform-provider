---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_network"
sidebar_current: "docs-cloudstack-resource-network"
description: |-
  Creates a network.
---

# cloudstack_network

Creates a network.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_network" "default" {
  name             = "test-network"
  cidr             = "10.0.0.0/16"
  network_offering = "Default Network"
  zone             = "zone-1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network.

* `display_text` - (Optional) The display text of the network.

* `cidr` - (Required) The CIDR block for the network. Changing this forces a new
    resource to be created.

* `gateway` - (Optional) Gateway that will be provided to the instances in this
    network. Defaults to the first usable IP in the range.

* `startip` - (Optional) Start of the IP block that will be available on the
    network. Defaults to the second available IP in the range.

* `endip` - (Optional) End of the IP block that will be available on the
    network. Defaults to the last available IP in the range.

* `network_domain` - (Optional) DNS domain for the network.

* `network_offering` - (Required) The name or ID of the network offering to use
    for this network.

* `vlan` - (Optional) The VLAN number (1-4095) the network will use. This might be
    required by the Network Offering if specifyVlan=true is set. Only the ROOT
    admin can set this value.

* `vpc_id` - (Optional) The VPC ID in which to create this network. Changing
    this forces a new resource to be created.

* `acl_id` - (Optional) The ACL ID that should be attached to the network or
    `none` if you do not want to attach an ACL. You can dynamically attach and
    swap ACL's, but if you want to detach an attached ACL and revert to using
    `none`, this will force a new resource to be created. (defaults `none`)

* `project` - (Optional) The name or ID of the project to deploy this
    instance to. Changing this forces a new resource to be created.

* `source_nat_ip` - (Optional) If set to `true` a public IP will be associated
    with the network. This is mainly used when the network supports the source
    NAT service which claims the first associated IP address. This prevents the
    ability to manage the IP address as an independent entity.

* `zone` - (Required) The name or ID of the zone where this network will be
    available. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network.
* `display_text` - The display text of the network.
* `network_domain` - DNS domain for the network.
* `source_nat_ip_id` - The ID of the associated source NAT IP.

## Import

Networks can be imported; use `<NETWORK ID>` as the import ID. For
example:

```shell
terraform import cloudstack_network.default 36619b20-5584-43bf-9a84-e242bacd5582
```

When importing into a project you need to prefix the import ID with the project name:

```shell
terraform import cloudstack_network.default my-project/36619b20-5584-43bf-9a84-e242bacd5582
```
