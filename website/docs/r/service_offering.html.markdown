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

### Basic Service Offering

```hcl
resource "cloudstack_service_offering" "example" {
    name = "example-service-offering"
    display_text = "Example Service Offering"
    cpu_number = 2
    memory = 4096
}
```

### GPU Service Offering

```hcl
resource "cloudstack_service_offering" "gpu_offering" {
    name = "gpu-a6000"
    display_text = "GPU A6000 Instance"
    cpu_number = 4
    memory = 16384
    
    service_offering_details = {
        pciDevice = "Group of NVIDIA A6000 GPUs"
        vgpuType  = "A6000-8A"
    }
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

* `customized` - (Optional) Whether the service offering allows custom CPU and memory values. Set to `true` to enable users to specify CPU/memory within the min/max constraints for constrained offerings and any value for unconstrained offerings.
    Changing this forces a new resource to be created.

* `min_cpu_number` - (Optional) Minimum number of CPU cores allowed for customized offerings.
    Changing this forces a new resource to be created.

* `max_cpu_number` - (Optional) Maximum number of CPU cores allowed for customized offerings.
    Changing this forces a new resource to be created.

* `min_memory` - (Optional) Minimum memory (in MB) allowed for customized offerings.
    Changing this forces a new resource to be created.

* `max_memory` - (Optional) Maximum memory (in MB) allowed for customized offerings.
    Changing this forces a new resource to be created.

* `encrypt_root` - (Optional) Whether to encrypt the root disk for VMs using this service offering.
    Changing this forces a new resource to be created.

* `storage_tags` - (Optional) Storage tags to associate with the service offering.

* `service_offering_details` - (Optional) A map of service offering details for GPU configuration and other advanced settings. Common keys include `pciDevice` and `vgpuType` for GPU offerings.
    Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service offering.

## Import

Service offerings can be imported; use `<SERVICEOFFERINGID>` as the import ID. For example:

```shell
$ terraform import cloudstack_service_offering.example <SERVICEOFFERINGID>
```
