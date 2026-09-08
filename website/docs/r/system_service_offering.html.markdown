---
layout: default
title: "CloudStack: cloudstack_system_service_offering"
sidebar_current: "docs-cloudstack-resource-system-service-offering"
description: |-
    Creates a System Service Offering
---

# CloudStack: cloudstack_system_service_offering

A `cloudstack_system_service_offering` resource manages a service offering for
CloudStack system VMs.

## Example Usage

```hcl
resource "cloudstack_system_service_offering" "router" {
  name           = "redundant-router-offering"
  display_text   = "Redundant router offering"
  system_vm_type = "domainrouter"
  cpu_number     = 2
  cpu_speed      = 1000
  memory         = 2048
  network_rate   = 200
  offer_ha       = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the system service offering.

* `display_text` - (Required) Display text of the system service offering.

* `system_vm_type` - (Required) Type of system VM that uses the offering. Valid
  values are `domainrouter`, `consoleproxy`, and `secondarystoragevm`. Changing
  this forces a new resource to be created.

* `cpu_number` - (Required) Number of CPU cores. Changing this forces a new
  resource to be created.

* `cpu_speed` - (Required) Speed of each CPU core in MHz. Changing this forces
  a new resource to be created.

* `memory` - (Required) Memory reserved by the system VM in MB. Changing this
  forces a new resource to be created.

* `storage_type` - (Optional) Storage type of the offering. Valid values are
  `local` and `shared`. Defaults to `shared`. Changing this forces a new resource
  to be created. `local` storage cannot be combined with `offer_ha = true`.

* `network_rate` - (Optional) Network rate in Mbps. This can only be set for a
  `domainrouter` offering. Changing this forces a new resource to be created.

* `offer_ha` - (Optional) Whether the offering supports HA. Defaults to `false`.
  Changing this forces a new resource to be created.

* `limit_cpu_use` - (Optional) Whether CPU usage is limited to the offering's
  committed resources. Defaults to `false`. Changing this forces a new resource
  to be created.

* `host_tags` - (Optional) Host tags associated with the offering.

* `storage_tags` - (Optional) Storage tags associated with the offering.

* `domain_ids` - (Optional) Set of domain IDs that can use the offering. Omit
  this argument to make the offering public.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the system service offering.

## Import

System service offerings can be imported using their ID. For example:

```shell
$ terraform import cloudstack_system_service_offering.router <SYSTEMSERVICEOFFERINGID>
```
