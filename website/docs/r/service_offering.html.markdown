---
layout: default
title: "CloudStack: cloudstack_service_offering"
sidebar_current: "docs-cloudstack-resource-service_offering"
description: |-
    Creates a Service Offering
---

# CloudStack: cloudstack_service_offering

A `cloudstack_service_offering` resource manages a service offering within CloudStack.

## Example Usage

```hcl
resource "cloudstack_service_offering" "example" {
    name = "example-service-offering"
    display_text = "Example Service Offering"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service offering.
* `display_text` - (Required) The display text of the service offering.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service offering.
* `name` - The name of the service offering.
* `display_text` - The display text of the service offering.

## Import

Service offerings can be imported; use `<SERVICEOFFERINGID>` as the import ID. For example:

```shell
$ terraform import cloudstack_service_offering.example <SERVICEOFFERINGID>
```
