---
layout: "cloudstack"
page_title: "CloudStack: cloudstack_service_offering_fixed"
sidebar_current: "docs-cloudstack-resource-service-offering-fixed"
description: |-
  Creates a service offering.
---

# cloudstack_service_offering_fixed

Creates a service offering.

## Example Usage

```hcl
resource "cloudstack_service_offering_fixed" "example" {
	display_text = "example"
	name         = "example"
	
	cpu_number     = 2
	cpu_speed      = 2500
	memory         = 2048

	host_tags          = "test0101, test0202"
	network_rate       = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false
}
```

## Argument Reference

The following arguments are supported:

* `cpu_number` - (Required) the CPU number of the service offering

* `cpu_speed` - (Required) the CPU speed of the service offering in MHz.

* `display_text` - (Required) alternate display text of the disk offering

* `memory` - (Required) the total memory of the service offering in MB

* `name` - (Required) name of the disk offering

---


* `deployment_planner` - (Optional) The deployment planner heuristics used to deploy a VM of this offering. If null, value of global config vm.deployment.planner is used

* `disk_offering_id` - (Optional) the ID of the disk offering to which service offering should be mapped

* `disk_hypervisor` - (Optional) block as defined below.

* `disk_offering` - (Optional) block as defined below.

* `disk_storage` - (Optional) block as defined below.

* `domain_id` - (Optional) the ID of the containing domain(s), null for public offerings

* `dynamic_scaling_enabled` - (Optional) true if virtual machine needs to be dynamically scalable of cpu or memory

* `host_tags` - (Optional) the host tag for this service offering.

* `is_volatile` - (Optional) true if the virtual machine needs to be volatile so that on every reboot of VM, original root disk is dettached then destroyed and a fresh root disk is created and attached to VM

* `limit_cpu_use` - (Optional) restrict the CPU usage to committed service offering

* `network_rate` - (Optional) The maximum number of CPUs to be set with Custom Computer Offering

* `offer_ha` - (Optional) the HA for the service offering

* `zone_id` - (Optional) the ID of the containing zone(s), null for public offerings

---
A `disk_hypervisor` block supports the following:

* `bytes_read_rate` - (Optional) bytes read rate of the disk offering

* `bytes_read_rate_max` - (Optional) burst bytes read rate of the disk offering

* `bytes_read_rate_max_length` - (Optional) length (in seconds) of the burst

* `bytes_write_rate` - (Optional) bytes write rate of the disk offering

* `bytes_write_rate_max` - (Optional) burst bytes write rate of the disk offering

* `bytes_write_rate_max_length` - (Optional) length (in seconds) of the burst

---
A `disk_offering` block supports the following:

* `cache_mode` - (Optional) the cache mode to use for this disk offering. none, writeback or writethrough

* `disk_offering_strictness` - (Optional) True/False to indicate the strictness of the disk offering association with the compute offering. When set to true, override of disk offering is not allowed when VM is deployed and change disk offering is not allowed for the ROOT disk after the VM is deployed

* `provisioning_type` - (Optional) provisioning type used to create volumes. Valid values are thin, sparse, fat.

* `root_disk_size` - (Optional) the Root disk size in GB.

* `storage_type` - (Optional) the storage type of the service offering. Values are local and shared.

* `tags` - (Optional) the tags for this service offering.

---
A `disk_storage` block supports the following:

* `customized_iops` - (Optional) whether compute offering iops is custom or not

* `hypervisor_snapshot_reserve` - (Optional) Hypervisor snapshot reserve space as a percent of a volume (for managed storage using Xen or VMware)

* `max_iops` - (Optional) max iops of the compute offering

* `min_iops` - (Optional) min iops of the compute offering



## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The id of the affinity group.


## Import

Service offerings can be imported; use `<ID>` as the import ID. For
example:

```shell
terraform import cloudstack_service_offering_fixed.example 6226ea4d-9cbe-4cc9-b30c-b9532146da5b
```
