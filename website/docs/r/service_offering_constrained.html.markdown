---
page_title: "cloudstack_service_offering_constrained Resource"
description: |-
	Provides a CloudStack Constrained Service Offering resource. This allows you to create and manage constrained compute offerings in CloudStack.
---

# cloudstack_service_offering_constrained

Provides a CloudStack Constrained Service Offering resource. This resource allows you to create and manage service offerings with constrained CPU and memory parameters.

## Example Usage

```hcl
resource "cloudstack_service_offering_constrained" "example" {
	display_text = "Example Constrained Offering"
	name         = "example_constrained"

	cpu_speed      = 2500
	max_cpu_number = 10
	min_cpu_number = 2
	max_memory     = 4096
	min_memory     = 1024

	host_tags    = "tag1,tag2"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false
	zone_ids = ["zone-uuid"]

	disk_offering {
		cache_mode                = "none"
		disk_offering_strictness  = true
		provisioning_type         = "thin"
		storage_type              = "local"
	}
}
```

## Argument Reference

The following arguments are supported:

- `display_text` (String, Required) - The display text of the service offering.
- `name` (String, Required) - The name of the service offering.
- `cpu_speed` (Int, Required) - CPU speed in MHz.
- `max_cpu_number` (Int, Required) - Maximum number of CPUs.
- `min_cpu_number` (Int, Required) - Minimum number of CPUs.
- `max_memory` (Int, Required) - Maximum memory in MB.
- `min_memory` (Int, Required) - Minimum memory in MB.
- `deployment_planner` (String, Optional) - The deployment planner for the service offering.
- `disk_offering_id` (String, Optional) - The ID of the disk offering.
- `domain_ids` (Set of String, Optional) - The ID(s) of the containing domain(s), null for public offerings.
- `dynamic_scaling_enabled` (Bool, Optional, Default: false) - Enable dynamic scaling of the service offering.
- `host_tags` (String, Optional) - The host tag for this service offering.
- `is_volatile` (Bool, Optional, Default: false) - Service offering is volatile.
- `limit_cpu_use` (Bool, Optional, Default: false) - Restrict the CPU usage to committed service offering.
- `network_rate` (Int, Optional) - Data transfer rate in megabits per second.
- `offer_ha` (Bool, Optional, Default: false) - Enable HA for the service offering.
- `zone_ids` (Set of String, Optional) - The ID(s) of the zone(s).

### Nested Blocks

#### `disk_offering` (Block, Optional)

- `cache_mode` (String, Required) - The cache mode to use for this disk offering. One of `none`, `writeback`, or `writethrough`.
- `disk_offering_strictness` (Bool, Required) - True/False to indicate the strictness of the disk offering association with the compute offering.
- `provisioning_type` (String, Required) - Provisioning type used to create volumes. Valid values are `thin`, `sparse`, `fat`.
- `root_disk_size` (Int, Optional) - The root disk size in GB.
- `storage_type` (String, Required) - The storage type of the service offering. Values are `local` and `shared`.
- `storage_tags` (String, Optional) - The tags for the service offering.

#### `disk_hypervisor` (Block, Optional)

- `bytes_read_rate` (Int, Required) - IO requests read rate of the disk offering.
- `bytes_read_rate_max` (Int, Required) - Burst requests read rate of the disk offering.
- `bytes_read_rate_max_length` (Int, Required) - Length (in seconds) of the burst.
- `bytes_write_rate` (Int, Required) - IO requests write rate of the disk offering.
- `bytes_write_rate_max` (Int, Required) - Burst IO requests write rate of the disk offering.
- `bytes_write_rate_max_length` (Int, Required) - Length (in seconds) of the burst.

#### `disk_storage` (Block, Optional)

- `customized_iops` (Bool, Optional) - True if disk offering uses custom IOPS, false otherwise.
- `hypervisor_snapshot_reserve` (Int, Optional) - Hypervisor snapshot reserve space as a percent of a volume (for managed storage using Xen or VMware).
- `max_iops` (Int, Optional) - Max IOPS of the compute offering.
- `min_iops` (Int, Optional) - Min IOPS of the compute offering.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The UUID of the service offering.
