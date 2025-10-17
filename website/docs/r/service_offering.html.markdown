---
layout: default
page_title: "CloudStack: cloudstack_service_offering"
sidebar_current: "docs-cloudstack-resource-service_offering"
description: |-
    Creates and manages a Service Offering in CloudStack
---

# cloudstack_service_offering

Provides a CloudStack Service Offering resource. This resource can be used to create, modify, and delete service offerings that define the compute resources (CPU, memory, storage) available to virtual machines.

## Understanding CloudStack Service Offering Types

CloudStack supports **three types** of service offerings based on the `customized` parameter and how CPU/memory are configured:

### 1. Fixed Offering (Fixed CPU and Memory)

Users **cannot** change CPU or memory when deploying VMs. The values are fixed.

**How to create:** Specify both `cpu_number` AND `memory` (do NOT set `customized`).

```hcl
resource "cloudstack_service_offering" "fixed" {
  name         = "small-fixed"
  display_text = "Small Fixed Instance - 2 CPU, 4GB RAM"
  cpu_number   = 2
  cpu_speed    = 2000
  memory       = 4096
  # customized is automatically set to false
}
```

### 2. Custom Constrained Offering (User Choice with Limits)

Users **can** choose CPU and memory **within** the min/max limits you define.

**How to create:** Set `customized = true` AND specify min/max constraints.

```hcl
resource "cloudstack_service_offering" "constrained" {
  name           = "custom-constrained"
  display_text   = "Custom Constrained - Choose between 2-8 CPU, 2-16GB RAM"
  customized     = true
  min_cpu_number = 2
  max_cpu_number = 8
  min_memory     = 2048  # 2 GB
  max_memory     = 16384 # 16 GB
}
```

### 3. Custom Unconstrained Offering (User Choice without Limits)

Users **can** choose **any** CPU and memory values (no restrictions).

**How to create:** Set `customized = true` WITHOUT min/max constraints.

```hcl
resource "cloudstack_service_offering" "unconstrained" {
  name         = "custom-unlimited"
  display_text = "Custom Unlimited - Choose any CPU/RAM"
  customized   = true
  # No min/max limits - users have complete freedom
}
```

## Example Usage

### Basic Fixed Service Offering

```hcl
resource "cloudstack_service_offering" "basic" {
  name         = "basic-offering"
  display_text = "Basic Service Offering"
  cpu_number   = 2
  memory       = 4096
}
```

### GPU Service Offering

```hcl
resource "cloudstack_service_offering" "gpu" {
  name         = "gpu-a6000"
  display_text = "GPU A6000 Instance"
  cpu_number   = 8
  memory       = 32768
  
  service_offering_details = {
    pciDevice = "Group of NVIDIA A6000 GPUs"
    vgpuType  = "A6000-8A"
  }
}
```

### GPU Service Offering with Direct GPU Parameters

```hcl
resource "cloudstack_service_offering" "tesla_p100_passthrough" {
  name         = "tesla-p100-passthrough"
  display_text = "Tesla P100 GPU - 2 vCPU, 1GB RAM (Passthrough Mode)"
  
  # Fixed CPU and Memory configuration
  cpu_number   = 2
  cpu_speed    = 1000
  memory       = 1024  # 1 GB
  
  # Storage configuration
  storage_type        = "shared"
  provisioning_type   = "thin"
  
  # Performance and HA settings
  offer_ha                = false
  limit_cpu_use           = false
  is_volatile             = false
  dynamic_scaling_enabled = true
  encrypt_root            = false
  
  # Cache mode
  cache_mode = "none"
  
  # GPU Configuration - Direct API mapping
  # gpu_card -> serviceofferingdetails["pciDevice"] -> resolves to gpucardid
  # gpu_type -> SetVgpuprofileid() + serviceofferingdetails["vgpuType"] -> resolves to vgpuprofileid/vgpuprofilename
  # gpu_count -> SetGpucount() -> number of GPUs
  
  gpu_card  = "Tesla P100 Auto Created"                      # GPU Card Name
  gpu_type  = "8cb04cca-8395-44f8-9d1b-48eb08d48bed"         # vGPU Profile UUID
  gpu_count = 1                                               # Number of GPUs
}
```

### High Performance with IOPS Limits

```hcl
resource "cloudstack_service_offering" "high_performance" {
  name                = "high-performance"
  display_text        = "High Performance Instance"
  cpu_number          = 4
  memory              = 8192
  storage_type        = "shared"
  
  # IOPS configuration
  bytes_read_rate     = 10485760  # 10 MB/s
  bytes_write_rate    = 10485760  # 10 MB/s
  iops_read_rate      = 1000
  iops_write_rate     = 1000
  
  # Additional settings
  offer_ha            = true
  limit_cpu_use       = true
  host_tags           = "high-performance"
  storage_tags        = "ssd"
}
```

### Dynamic Scaling Enabled

```hcl
resource "cloudstack_service_offering" "scalable" {
  name                    = "scalable-offering"
  display_text            = "Dynamically Scalable Instance"
  cpu_number              = 4
  memory                  = 8192
  dynamic_scaling_enabled = true
  disk_iops_min           = 500
  disk_iops_max           = 2000
}
```

## Argument Reference

The following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the service offering. Changing this forces a new resource to be created.

### Basic Configuration

* `display_text` - (Optional) The display text of the service offering. If not provided, defaults to the name.

### CPU and Memory Configuration

* `cpu_number` - (Optional) The number of CPU cores. **Note:** When specified together with `memory`, creates a **Fixed Offering** (users cannot change CPU/memory). Changing this forces a new resource to be created.

* `cpu_speed` - (Optional) The clock rate of the CPU cores in MHz. Changing this forces a new resource to be created.

* `memory` - (Optional) The total memory for the service offering in MB. **Note:** When specified together with `cpu_number`, creates a **Fixed Offering** (users cannot change CPU/memory). Changing this forces a new resource to be created.

### Customization Options (Choose Offering Type)

* `customized` - (Optional) Controls whether users can choose CPU and memory when deploying VMs:
  - **Not set** + `cpu_number` + `memory` specified = **Fixed Offering** (no user choice)
  - **`true`** + min/max limits = **Custom Constrained** (user choice within limits)
  - **`true`** + no limits = **Custom Unconstrained** (any user choice)
  
  Changing this forces a new resource to be created.

* `customized_iops` - (Optional) Whether compute offering IOPS should be customizable. Changing this forces a new resource to be created.

### Custom Constrained Limits

Use these **only** when `customized = true` to set boundaries:

* `min_cpu_number` - (Optional) Minimum number of CPU cores users can select. Changing this forces a new resource to be created.

* `max_cpu_number` - (Optional) Maximum number of CPU cores users can select. Changing this forces a new resource to be created.

* `min_memory` - (Optional) Minimum memory in MB users can select. Changing this forces a new resource to be created.

* `max_memory` - (Optional) Maximum memory in MB users can select. Changing this forces a new resource to be created.

### Storage Configuration

* `storage_type` - (Optional) The storage type of the service offering. Valid values are `local` and `shared`. Changing this forces a new resource to be created.

* `storage_tags` - (Optional) Comma-separated list of tags for matching storage pools. Changing this forces a new resource to be created.

* `host_tags` - (Optional) Comma-separated list of tags for matching hosts. Changing this forces a new resource to be created.

* `root_disk_size` - (Optional) The root disk size in GB. Changing this forces a new resource to be created.

* `encrypt_root` - (Optional) Whether to encrypt the root disk. Changing this forces a new resource to be created.

### GPU Configuration

* `gpu_card` - (Optional) The GPU card name for GPU-enabled service offerings. This maps to `serviceofferingdetails["pciDevice"]` and CloudStack automatically resolves it to `gpucardid`. Example: `"Tesla P100 Auto Created"`. Changing this forces a new resource to be created.

* `gpu_type` - (Optional) The vGPU profile UUID or type for GPU-enabled service offerings. This parameter serves dual purposes:
  - Sets the vGPU profile ID via `SetVgpuprofileid()` API parameter
  - Populates `serviceofferingdetails["vgpuType"]`
  - CloudStack uses this to determine both `vgpuprofileid` and `vgpuprofilename`
  
  Example: `"8cb04cca-8395-44f8-9d1b-48eb08d48bed"` for passthrough mode. Changing this forces a new resource to be created.

* `gpu_count` - (Optional) The number of GPUs to allocate for this service offering. Maps directly to `SetGpucount()` API parameter. Default is `1` if not specified. Changing this forces a new resource to be created.

~> **Note:** For GPU offerings, ensure your CloudStack hosts are properly configured with GPU passthrough and that the GPU card names and vGPU profile UUIDs match your physical GPU configuration.

### IOPS and Bandwidth Limits

~> **Note:** All IOPS and bandwidth parameters are immutable after creation. Any changes will force recreation of the service offering.

* `bytes_read_rate` - (Optional) Bytes read rate in bytes per second. Changing this forces a new resource to be created.

* `bytes_read_rate_max` - (Optional) Burst bytes read rate in bytes per second. Changing this forces a new resource to be created.

* `bytes_read_rate_max_length` - (Optional) Length of the burst bytes read rate in seconds. Changing this forces a new resource to be created.

* `bytes_write_rate` - (Optional) Bytes write rate in bytes per second. Changing this forces a new resource to be created.

* `bytes_write_rate_max` - (Optional) Burst bytes write rate in bytes per second. Changing this forces a new resource to be created.

* `bytes_write_rate_max_length` - (Optional) Length of the burst bytes write rate in seconds. Changing this forces a new resource to be created.

* `iops_read_rate` - (Optional) IO requests read rate in IOPS. Changing this forces a new resource to be created.

* `iops_read_rate_max` - (Optional) Burst IO requests read rate in IOPS. Changing this forces a new resource to be created.

* `iops_read_rate_max_length` - (Optional) Length of the burst IO requests read rate in seconds. Changing this forces a new resource to be created.

* `iops_write_rate` - (Optional) IO requests write rate in IOPS. Changing this forces a new resource to be created.

* `iops_write_rate_max` - (Optional) Burst IO requests write rate in IOPS. Changing this forces a new resource to be created.

* `iops_write_rate_max_length` - (Optional) Length of the burst IO requests write rate in seconds. Changing this forces a new resource to be created.

### Dynamic Scaling and IOPS

* `disk_iops_min` - (Optional) Minimum IOPS for the disk. Changing this forces a new resource to be created.

* `disk_iops_max` - (Optional) Maximum IOPS for the disk. Changing this forces a new resource to be created.

* `hypervisor_snapshot_reserve` - (Optional) Hypervisor snapshot reserve space as a percent of a volume (for managed storage using Xen or VMware). Changing this forces a new resource to be created.

### High Availability and Performance

* `offer_ha` - (Optional) Whether HA (High Availability) is enabled for the service offering. Changing this forces a new resource to be created.

* `limit_cpu_use` - (Optional) Restrict the CPU usage to the committed service offering percentage. Changing this forces a new resource to be created.

* `dynamic_scaling_enabled` - (Optional) Enable dynamic scaling of VM CPU and memory. Changing this forces a new resource to be created.

### Cache and Deployment

* `cache_mode` - (Optional) The cache mode to use. Valid values: `none`, `writeback`, `writethrough`. Changing this forces a new resource to be created.

* `deployment_planner` - (Optional) The deployment planner heuristics used to deploy VMs. Changing this forces a new resource to be created.

### System VM Configuration

* `is_system` - (Optional) Whether the offering is for system VMs. Changing this forces a new resource to be created.

* `system_vm_type` - (Optional) The system VM type. Required when `is_system` is true. Valid values: `domainrouter`, `consoleproxy`, `secondarystoragevm`. Changing this forces a new resource to be created.

### Zone and Domain

* `zone_id` - (Optional) The ID of the zone for which the service offering is available. Changing this forces a new resource to be created.

* `domain_id` - (Optional) The ID of the domain for which the service offering is available. Changing this forces a new resource to be created.

### Provisioning

* `provisioning_type` - (Optional) Provisioning type used to create volumes. Valid values: `thin`, `sparse`, `fat`. Changing this forces a new resource to be created.

* `is_volatile` - (Optional) If true, the VM's original root disk is destroyed and recreated on every reboot. Changing this forces a new resource to be created.

### Advanced Settings

* `service_offering_details` - (Optional) A map of service offering details for GPU configuration and other advanced settings. Common keys include:
  - `pciDevice` - PCI device for GPU passthrough
  - `vgpuType` - vGPU type for GPU offerings
  
  ~> **Note:** CloudStack may add system keys (like `External:key`, `purge.db.entities`) to this map. Terraform automatically filters these out to prevent state drift.
  
  Changing this forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the service offering.
* `created` - The date when the service offering was created.

## Import

Service offerings can be imported using the service offering ID:

```shell
terraform import cloudstack_service_offering.example 550e8400-e29b-41d4-a716-446655440000
```

When importing resources associated with a project, use the format `project_name/offering_id`:

```shell
terraform import cloudstack_service_offering.example my-project/550e8400-e29b-41d4-a716-446655440000
```

## Important Notes

### Understanding Offering Types Decision Tree

```
Need to create a Service Offering?
│
├─ Users should choose CPU/RAM? ──NO──> Fixed Offering
│                                       (set cpu_number + memory)
│                                       
└─ YES
   │
   ├─ Limit their choices? ──NO──> Custom Unconstrained
   │                                (customized = true, no min/max)
   │
   └─ YES ──> Custom Constrained
              (customized = true, set min/max limits)
```

### ForceNew Behavior

⚠️ **Most parameters are immutable** due to CloudStack API limitations. Changes to the following will destroy and recreate the service offering:

- All CPU, memory, and customization parameters
- All storage and IOPS/bandwidth parameters
- All HA, scaling, and system VM parameters
- `service_offering_details` map

**Only these can be updated in-place:**
- `display_text`
- `host_tags`
- `storage_tags`

### CloudStack API Behaviors

1. **`customized` parameter logic:**
   - If you provide both `cpu_number` AND `memory`: Terraform automatically sets `customized = false` (Fixed Offering)
   - If you set `customized = true`: Users can choose CPU/RAM at VM deployment time
   - If neither are set: CloudStack creates a Custom Unconstrained offering

2. **`service_offering_details` filtering:**
   - CloudStack automatically adds system-managed keys: `External:key`, `External:value`, `purge.db.entities`
   - Terraform filters these out automatically to prevent state drift
   - Only the keys YOU configure are tracked in Terraform state

3. **Write-only parameters:**
   - Some parameters (`tags`, `lease_duration`) are not returned by CloudStack's Read API
   - These are write-only and won't appear in state after creation

### Performance Configuration Examples

**Example: 10 MB/s read rate**
```hcl
bytes_read_rate = 10485760  # 10 * 1024 * 1024 = 10 MB/s in bytes
```

**Example: 1000 IOPS**
```hcl
iops_read_rate = 1000  # Direct IOPS value
```

**Example: Burst configuration**
```hcl
bytes_write_rate            = 5242880   # 5 MB/s baseline
bytes_write_rate_max        = 10485760  # 10 MB/s burst
bytes_write_rate_max_length = 60        # Burst for 60 seconds
```
