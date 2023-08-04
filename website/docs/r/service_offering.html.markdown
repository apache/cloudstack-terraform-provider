---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_service_offering"
sidebar_current: "docs-cloudstack-resource-service-offering"
description: |-
  Creates a service offering.
---

# cloudstack_service_offering

Creates a service offering.

## Example Usage

Basic usage:

```hcl
resource "cloudstack_service_offering" "default" {
    name       = "Small"
    cpu_number = 2
    cpu_speed  = 1000
    memory     = 4096
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the service offering.
    Changing this forces a new resource to be created.

* `display_text` - (Optional) The display text of the service offering.

* `cpu_number` - (Optional) The number of CPU cores.
    Changing this forces a new resource to be created.

* `cpu_speed` - (Optional) The speed of the CPU in Mhz.
    Changing this forces a new resource to be created.

* `memory` - (Optional) Memory reserved by the VM in MB.
    Changing this forces a new resource to be created.

* `host_tags` - (Optional) The host tags for the service offering.

* `limit_cpu_use` - (Optional) Restrict the CPU usage to committed service offering.
    Changing this forces a new resource to be created.

* `offer_ha` - (Optional) The HA for the service offering.
    Changing this forces a new resource to be created.

* `storage_type` - (Optional) The storage type of the service offering. Values are `local` and `shared`.
    Changing this forces a new resource to be created.
    

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service offering.
