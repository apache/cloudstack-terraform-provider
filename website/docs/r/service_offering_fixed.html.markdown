---
page_title: "cloudstack_service_offering_fixed Resource"
sidebar_current: terraform-resource-service_offering_fixed
description: |-
  Provides a CloudStack Service Offering (Fixed) resource. This resource allows you to create and manage fixed compute service offerings in CloudStack.
---

# cloudstack_service_offering_fixed

Provides a CloudStack Service Offering (Fixed) resource. This resource allows you to create and manage fixed compute service offerings in CloudStack.

## Example Usage

```hcl
resource "cloudstack_service_offering_fixed" "fixed1" {
  name         = "fixed1"
  display_text = "fixed1"
  cpu_number   = 2
  cpu_speed    = 2000
  memory       = 4096
  # Optional common attributes:
  # deployment_planner = "FirstFit"
  # disk_offering_id   = "..."
  # domain_ids         = ["...", "..."]
  # dynamic_scaling_enabled = false
  # host_tags          = "..."
  # is_volatile        = false
  # limit_cpu_use      = false
  # network_rate       = 1000
  # offer_ha           = false
  # zone_ids           = ["..."]
  # disk_offering { ... }
  # disk_hypervisor { ... }
  # disk_storage { ... }
}
```

## Argument Reference

The following arguments are supported:

- `name` (Required) - The name of the service offering.
- `display_text` (Required) - The display text of the service offering.
- `cpu_number` (Required) - Number of CPU cores.
- `cpu_speed` (Required) - The CPU speed in MHz (not applicable to KVM).
- `memory` (Required) - The total memory in MB.

### Common Attributes

- `deployment_planner` (Optional) - The deployment planner for the service offering.
- `disk_offering_id` (Optional) - The ID of the disk offering.
- `domain_ids` (Optional) - The ID(s) of the containing domain(s), null for public offerings.
- `dynamic_scaling_enabled` (Optional, Computed) - Enable dynamic scaling of the service offering. Defaults to `false`.
- `host_tags` (Optional) - The host tag for this service offering.
- `id` (Computed) - The UUID of the service offering.
- `is_volatile` (Optional, Computed) - Service offering is volatile. Defaults to `false`.
- `limit_cpu_use` (Optional, Computed) - Restrict the CPU usage to committed service offering. Defaults to `false`.
- `network_rate` (Optional) - Data transfer rate in megabits per second.
- `offer_ha` (Optional, Computed) - The HA for the service offering. Defaults to `false`.
- `zone_ids` (Optional) - The ID(s) of the zone(s).

### Nested Blocks

#### `disk_offering` (Optional)

- `cache_mode` (Required) - The cache mode to use for this disk offering. One of `none`, `writeback`, or `writethrough`.
- `disk_offering_strictness` (Required) - True/False to indicate the strictness of the disk offering association.
- `provisioning_type` (Required) - Provisioning type used to create volumes. Valid values: `thin`, `sparse`, `fat`.
- `root_disk_size` (Optional) - The root disk size in GB.
- `storage_type` (Required) - The storage type. Values: `local`, `shared`.
- `storage_tags` (Optional) - The tags for the service offering.

#### `disk_hypervisor` (Optional)

- `bytes_read_rate` (Required) - IO requests read rate.
- `bytes_read_rate_max` (Required) - Burst requests read rate.
- `bytes_read_rate_max_length` (Required) - Length (in seconds) of the burst.
- `bytes_write_rate` (Required) - IO requests write rate.
- `bytes_write_rate_max` (Required) - Burst IO requests write rate.
- `bytes_write_rate_max_length` (Required) - Length (in seconds) of the burst.

#### `disk_storage` (Optional)

- `customized_iops` (Optional) - True if disk offering uses custom IOPS.
- `hypervisor_snapshot_reserve` (Optional) - Hypervisor snapshot reserve space as a percent of a volume.
- `max_iops` (Optional) - Max IOPS of the compute offering.
- `min_iops` (Optional) - Min IOPS of the compute offering.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the service offering.

## Import

Service offerings can be imported using the ID:

```sh
terraform import cloudstack_service_offering_fixed.example <id>
```
