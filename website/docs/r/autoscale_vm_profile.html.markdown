---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_autoscale_vm_profile"
sidebar_current: "docs-cloudstack-autoscale-vm-profile"
description: |-
  Creates an autoscale VM profile.
---

# cloudstack_autoscale_vm_profile

Creates an autoscale VM profile.

## Example Usage

```hcl
resource "cloudstack_autoscale_vm_profile" "profile1" {
  service_offering        = "small"
  template                = "CentOS 6.5"
  zone                    = "zone-1"
  destroy_vm_grace_period = "45s"
  
  other_deploy_params = {
    networkids  = "6eb22f91-7454-4107-89f4-36afcdf33021"
    displayname = "profile1vm"
  }

  metadata = {
    mydata = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `service_offering` - (Required) The name or ID of the service offering used
    for instances. Changing this forces a new resource to be created.

* `template` - (Required) The name or ID of the template used for instances.

* `zone` - (Required) The name or ID of the zone where instances will be
    created. Changing this forces a new resource to be created.

* `destroy_vm_grace_period` - (Optional) A time interval to wait for graceful
    shutdown of instances.

* `other_deploy_params` - (Optional) A mapping of additional params used when
    creating new instances.

* `metadata` - (Optional) A mapping of metadata key/values to assign to the
    resource.

## Attributes Reference

The following attributes are exported:

* `id` - The autoscale VM profile ID.
