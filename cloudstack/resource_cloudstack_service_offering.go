//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package cloudstack

import (
	"fmt"
	"strconv"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackServiceOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackServiceOfferingCreate,
		Read:   resourceCloudStackServiceOfferingRead,
		Update: resourceCloudStackServiceOfferingUpdate,
		Delete: resourceCloudStackServiceOfferingDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu_number": {
				Description: "Number of CPU cores",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"cpu_speed": {
				Description: "Speed of CPU in MHz",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"memory": {
				Description: "Amount of memory in MB",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"host_tags": {
				Description: "The host tags for this service offering",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"storage_type": {
				Description: "The storage type of the service offering",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "shared",
			},
			"service_offering_details": {
				Description: "Service offering details for custom configuration",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"customized": {
				Description: "Whether service offering allows custom CPU/memory or not. If not specified, CloudStack automatically determines based on cpu_number and memory presence: creates customizable offering (true) when cpu_number/memory are omitted, or fixed offering (false) when they are provided.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"gpu_card": {
				Description: "GPU card name (e.g., 'Tesla P100 Auto Created')",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"gpu_type": {
				Description: "GPU profile/type (e.g., 'passthrough', 'GRID V100-8Q')",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"gpu_count": {
				Description: "Number of GPUs",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},

			// Behavior settings
			"offer_ha": {
				Description: "Whether to offer HA to the VMs created with this offering",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			},

			"dynamic_scaling_enabled": {
				Description: "Whether to enable dynamic scaling",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			},

			// Disk configuration
			"root_disk_size": {
				Description: "Root disk size in GB",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},

			"provisioning_type": {
				Description: "Provisioning type: thin, sparse, or fat",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			},

			// Security
			"encrypt_root": {
				Description: "Whether to encrypt the root disk or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			}, // Customized Offering Limits
			"min_cpu_number": {
				Description: "Minimum number of CPUs for customized offerings",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},

			"max_cpu_number": {
				Description: "Maximum number of CPUs for customized offerings",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},

			"min_memory": {
				Description: "Minimum memory in MB for customized offerings",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},

			"max_memory": {
				Description: "Maximum memory in MB for customized offerings",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},

			// IOPS Limits
			"min_iops": {
				Description: "Minimum IOPS",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},

			"max_iops": {
				Description: "Maximum IOPS",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},

			// GPU Display
			"gpu_display": {
				Description: "Whether to display GPU in UI or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			},

			// High Priority Parameters
			"limit_cpu_use": {
				Description: "Restrict CPU usage to the service offering value",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			},

			"is_volatile": {
				Description: "True if the virtual machine needs to be volatile (root disk destroyed on stop)",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			},

			"customized_iops": {
				Description: "Whether compute offering iops is custom or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true, // CloudStack returns default value
				ForceNew:    true,
			},

			"tags": {
				Description: "Comma-separated list of tags for the service offering",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"domain_id": {
				Description: "The ID(s) of the domain(s) to which the service offering belongs",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"zone_id": {
				Description: "The ID(s) of the zone(s) this service offering belongs to",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			// IOPS/Bandwidth parameters (Phase 2)
			// Note: CloudStack API does not support updating these after creation
			"disk_iops_read_rate": {
				Description: "IO requests read rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_iops_write_rate": {
				Description: "IO requests write rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_iops_read_rate_max": {
				Description: "IO requests read rate max of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_iops_write_rate_max": {
				Description: "IO requests write rate max of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_iops_read_rate_max_length": {
				Description: "Burst duration in seconds for read rate max",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_iops_write_rate_max_length": {
				Description: "Burst duration in seconds for write rate max",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_bytes_read_rate": {
				Description: "Bytes read rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_bytes_write_rate": {
				Description: "Bytes write rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_bytes_read_rate_max": {
				Description: "Bytes read rate max of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"disk_bytes_write_rate_max": {
				Description: "Bytes write rate max of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"bytes_read_rate_max_length": {
				Description: "Burst duration in seconds for bytes read rate max",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"bytes_write_rate_max_length": {
				Description: "Burst duration in seconds for bytes write rate max",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			}, // Hypervisor Parameters (Phase 3)
			"hypervisor_snapshot_reserve": {
				Description: "Hypervisor snapshot reserve space as a percent of a volume (for managed storage using Xen or VMware)",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},

			"cache_mode": {
				Description: "The cache mode to use for the disk offering. Valid values: none, writeback, writethrough",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"deployment_planner": {
				Description: "Deployment planner heuristics to use for the service offering",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"storage_policy": {
				Description: "Name of the storage policy (for VMware)",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			// Low Priority Parameters (Phase 4)
			"network_rate": {
				Description: "Data transfer rate in megabits per second allowed",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			},

			"purge_resources": {
				Description: "Whether to purge resources on deletion",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated after creation
			}, "system_vm_type": {
				Description: "The system VM type. Possible values: domainrouter, consoleproxy, secondarystoragevm",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Cannot be updated - defines fundamental VM type
			},

			"disk_offering_id": {
				Description: "The ID of the disk offering to associate with this service offering",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"disk_offering_strictness": {
				Description: "Whether to strictly enforce the disk offering (requires disk_offering_id)",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},

			"external_details": {
				Description: "External system metadata (CMDB, billing, etc.)",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"is_system": {
				Description: "Whether this is a system VM offering (WARNING: For CloudStack internal use)",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},

			"lease_duration": {
				Description: "Lease duration in seconds (for temporary offerings with auto-cleanup)",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},

			"lease_expiry_action": {
				Description: "Action when lease expires. Possible values: destroy, stop",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceCloudStackServiceOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create parameters structure
	p := cs.ServiceOffering.NewCreateServiceOfferingParams(
		d.Get("display_text").(string),
		d.Get("name").(string),
	)

	// Set optional parameters
	if v, ok := d.GetOk("cpu_number"); ok {
		p.SetCpunumber(v.(int))
	}

	if v, ok := d.GetOk("cpu_speed"); ok {
		p.SetCpuspeed(v.(int))
	}

	if v, ok := d.GetOk("memory"); ok {
		p.SetMemory(v.(int))
	}

	if v, ok := d.GetOk("host_tags"); ok {
		p.SetHosttags(v.(string))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		p.SetStoragetype(v.(string))
	}

	// Handle customized parameter with CloudStack UI logic:
	// 1. If user explicitly sets customized, use that value
	// 2. If user provides cpu_number AND memory without customized, default to false (Fixed Offering)
	// 3. If none specified, CloudStack creates Custom unconstrained (customized=true)
	if v, ok := d.GetOkExists("customized"); ok {
		// User explicitly configured customized
		p.SetCustomized(v.(bool))
	} else {
		// User didn't specify customized - check if cpu/memory are provided
		_, hasCpuNumber := d.GetOk("cpu_number")
		_, hasMemory := d.GetOk("memory")

		if hasCpuNumber && hasMemory {
			// Both cpu and memory provided → Fixed Offering (customized=false)
			p.SetCustomized(false)
		}
		// If neither provided → CloudStack will default to Custom unconstrained (customized=true)
		// Don't send customized parameter, let CloudStack decide
	}

	// Handle GPU parameters
	// GPU configuration uses dedicated API parameters (not serviceofferingdetails)

	// Set vGPU profile ID (UUID)
	if v, ok := d.GetOk("gpu_type"); ok {
		p.SetVgpuprofileid(v.(string))
	}

	// Set GPU count
	if v, ok := d.GetOk("gpu_count"); ok {
		p.SetGpucount(v.(int))
	}

	// Set GPU display
	if v, ok := d.GetOk("gpu_display"); ok {
		p.SetGpudisplay(v.(bool))
	}

	// High Availability
	if v, ok := d.GetOk("offer_ha"); ok {
		p.SetOfferha(v.(bool))
	}

	// Dynamic Scaling
	if v, ok := d.GetOk("dynamic_scaling_enabled"); ok {
		p.SetDynamicscalingenabled(v.(bool))
	}

	// Disk Configuration
	if v, ok := d.GetOk("root_disk_size"); ok {
		p.SetRootdisksize(int64(v.(int)))
	}

	if v, ok := d.GetOk("provisioning_type"); ok {
		p.SetProvisioningtype(v.(string))
	}

	// Security
	if v, ok := d.GetOk("encrypt_root"); ok {
		p.SetEncryptroot(v.(bool))
	}

	// Customized Offering Limits
	if v, ok := d.GetOk("min_cpu_number"); ok {
		p.SetMincpunumber(v.(int))
	}

	if v, ok := d.GetOk("max_cpu_number"); ok {
		p.SetMaxcpunumber(v.(int))
	}

	if v, ok := d.GetOk("min_memory"); ok {
		p.SetMinmemory(v.(int))
	}

	if v, ok := d.GetOk("max_memory"); ok {
		p.SetMaxmemory(v.(int))
	}

	// IOPS Limits
	if v, ok := d.GetOk("min_iops"); ok {
		p.SetMiniops(int64(v.(int)))
	}

	if v, ok := d.GetOk("max_iops"); ok {
		p.SetMaxiops(int64(v.(int)))
	}

	// High Priority Parameters
	if v, ok := d.GetOk("limit_cpu_use"); ok {
		p.SetLimitcpuuse(v.(bool))
	}

	if v, ok := d.GetOk("is_volatile"); ok {
		p.SetIsvolatile(v.(bool))
	}

	if v, ok := d.GetOk("customized_iops"); ok {
		p.SetCustomizediops(v.(bool))
	}

	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(v.(string))
	}

	if details, ok := d.GetOk("service_offering_details"); ok {
		serviceOfferingDetails := make(map[string]string)
		for k, v := range details.(map[string]interface{}) {
			serviceOfferingDetails[k] = v.(string)
		}
		p.SetServiceofferingdetails(serviceOfferingDetails)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		zones := make([]string, 0)
		for _, z := range v.([]interface{}) {
			zones = append(zones, z.(string))
		}
		p.SetZoneid(zones)
	}

	if v, ok := d.GetOk("disk_iops_read_rate"); ok {
		p.SetIopsreadrate(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_iops_write_rate"); ok {
		p.SetIopswriterate(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_iops_read_rate_max"); ok {
		p.SetIopsreadratemax(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_iops_write_rate_max"); ok {
		p.SetIopswriteratemax(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_iops_read_rate_max_length"); ok {
		p.SetIopsreadratemaxlength(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_iops_write_rate_max_length"); ok {
		p.SetIopswriteratemaxlength(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_bytes_read_rate"); ok {
		p.SetBytesreadrate(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_bytes_write_rate"); ok {
		p.SetByteswriterate(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_bytes_read_rate_max"); ok {
		p.SetBytesreadratemax(int64(v.(int)))
	}

	if v, ok := d.GetOk("disk_bytes_write_rate_max"); ok {
		p.SetByteswriteratemax(int64(v.(int)))
	}

	if v, ok := d.GetOk("bytes_read_rate_max_length"); ok {
		p.SetBytesreadratemaxlength(int64(v.(int)))
	}

	if v, ok := d.GetOk("bytes_write_rate_max_length"); ok {
		p.SetByteswriteratemaxlength(int64(v.(int)))
	}

	if v, ok := d.GetOk("hypervisor_snapshot_reserve"); ok {
		p.SetHypervisorsnapshotreserve(v.(int))
	}

	if v, ok := d.GetOk("cache_mode"); ok {
		p.SetCachemode(v.(string))
	}

	if v, ok := d.GetOk("deployment_planner"); ok {
		p.SetDeploymentplanner(v.(string))
	}

	if v, ok := d.GetOk("storage_policy"); ok {
		p.SetStoragepolicy(v.(string))
	}

	if v, ok := d.GetOk("network_rate"); ok {
		p.SetNetworkrate(v.(int))
	}

	if v, ok := d.GetOk("purge_resources"); ok {
		p.SetPurgeresources(v.(bool))
	}

	if v, ok := d.GetOk("system_vm_type"); ok {
		p.SetSystemvmtype(v.(string))
	}

	if v, ok := d.GetOk("disk_offering_id"); ok {
		p.SetDiskofferingid(v.(string))
	}

	if v, ok := d.GetOk("disk_offering_strictness"); ok {
		p.SetDiskofferingstrictness(v.(bool))
	}

	if v, ok := d.GetOk("external_details"); ok {
		extDetails := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			extDetails[key] = value.(string)
		}
		p.SetExternaldetails(extDetails)
	}

	if v, ok := d.GetOk("is_system"); ok {
		p.SetIssystem(v.(bool))
	}

	if v, ok := d.GetOk("lease_duration"); ok {
		p.SetLeaseduration(v.(int))
	}

	if v, ok := d.GetOk("lease_expiry_action"); ok {
		p.SetLeaseexpiryaction(v.(string))
	}

	if v, ok := d.GetOk("service_offering_details"); ok {
		details := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			details[key] = value.(string)
		}
		p.SetServiceofferingdetails(details)
	}

	// Create the service offering
	so, err := cs.ServiceOffering.CreateServiceOffering(p)
	if err != nil {
		return fmt.Errorf("Error creating service offering %s: %s", d.Get("name").(string), err)
	}

	d.SetId(so.Id)

	return resourceCloudStackServiceOfferingRead(d, meta)
}

func resourceCloudStackServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	so, count, err := cs.ServiceOffering.GetServiceOfferingByID(d.Id())
	if err != nil {
		if count == 0 {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", so.Name)
	d.Set("display_text", so.Displaytext)
	d.Set("cpu_number", so.Cpunumber)
	d.Set("cpu_speed", so.Cpuspeed)
	d.Set("memory", so.Memory)
	d.Set("host_tags", so.Hosttags)
	d.Set("storage_type", so.Storagetype)

	// Set GPU fields from dedicated response fields
	// Use gpucardname (not gpucardid) to match what user provides in terraform config
	if so.Gpucardname != "" {
		d.Set("gpu_card", so.Gpucardname)
	}

	// Use vgpuprofileid (UUID) as configured
	if so.Vgpuprofileid != "" {
		d.Set("gpu_type", so.Vgpuprofileid)
	}

	if so.Gpucount > 0 {
		d.Set("gpu_count", so.Gpucount)
	}

	// Set computed fields (CloudStack returns default values)
	d.Set("gpu_display", so.Gpudisplay)
	d.Set("offer_ha", so.Offerha)
	d.Set("dynamic_scaling_enabled", so.Dynamicscalingenabled)
	d.Set("encrypt_root", so.Encryptroot)
	d.Set("provisioning_type", so.Provisioningtype)

	if so.Rootdisksize > 0 {
		d.Set("root_disk_size", int(so.Rootdisksize))
	}

	// IOPS limits - only set if returned by API (> 0)
	if so.Miniops > 0 {
		d.Set("min_iops", int(so.Miniops))
	}
	if so.Maxiops > 0 {
		d.Set("max_iops", int(so.Maxiops))
	}

	// High Priority Parameters
	d.Set("limit_cpu_use", so.Limitcpuuse)
	d.Set("is_volatile", so.Isvolatile)
	d.Set("customized_iops", so.Iscustomizediops)
	d.Set("customized", so.Iscustomized)

	// Tags field is write-only, not returned by API - skip setting

	// Domain and Zone IDs
	if so.Domainid != "" {
		d.Set("domain_id", []string{so.Domainid})
	}

	if so.Zoneid != "" {
		d.Set("zone_id", []string{so.Zoneid})
	}

	// IOPS/Bandwidth Parameters (Phase 2)
	if so.DiskIopsReadRate > 0 {
		d.Set("disk_iops_read_rate", int(so.DiskIopsReadRate))
	}

	if so.DiskIopsWriteRate > 0 {
		d.Set("disk_iops_write_rate", int(so.DiskIopsWriteRate))
	}

	if so.DiskIopsReadRateMax > 0 {
		d.Set("disk_iops_read_rate_max", int(so.DiskIopsReadRateMax))
	}

	if so.DiskIopsWriteRateMax > 0 {
		d.Set("disk_iops_write_rate_max", int(so.DiskIopsWriteRateMax))
	}

	if so.DiskIopsReadRateMaxLength > 0 {
		d.Set("disk_iops_read_rate_max_length", int(so.DiskIopsReadRateMaxLength))
	}

	if so.DiskIopsWriteRateMaxLength > 0 {
		d.Set("disk_iops_write_rate_max_length", int(so.DiskIopsWriteRateMaxLength))
	}

	if so.DiskBytesReadRate > 0 {
		d.Set("disk_bytes_read_rate", int(so.DiskBytesReadRate))
	}

	if so.DiskBytesWriteRate > 0 {
		d.Set("disk_bytes_write_rate", int(so.DiskBytesWriteRate))
	}

	if so.DiskBytesReadRateMax > 0 {
		d.Set("disk_bytes_read_rate_max", int(so.DiskBytesReadRateMax))
	}

	if so.DiskBytesWriteRateMax > 0 {
		d.Set("disk_bytes_write_rate_max", int(so.DiskBytesWriteRateMax))
	}

	if so.DiskBytesReadRateMaxLength > 0 {
		d.Set("bytes_read_rate_max_length", int(so.DiskBytesReadRateMaxLength))
	}

	if so.DiskBytesWriteRateMaxLength > 0 {
		d.Set("bytes_write_rate_max_length", int(so.DiskBytesWriteRateMaxLength))
	}

	// Hypervisor Parameters (Phase 3)
	if so.Hypervisorsnapshotreserve > 0 {
		d.Set("hypervisor_snapshot_reserve", so.Hypervisorsnapshotreserve)
	}

	if so.CacheMode != "" {
		d.Set("cache_mode", so.CacheMode)
	}

	if so.Deploymentplanner != "" {
		d.Set("deployment_planner", so.Deploymentplanner)
	}

	// Note: storage_policy field doesn't exist in ServiceOffering response
	// This is a write-only parameter for VMware environments

	// Low Priority Parameters (Phase 4)
	if so.Networkrate > 0 {
		d.Set("network_rate", so.Networkrate)
	}

	// Note: purge_resources is write-only, not returned by API
	// Note: system_vm_type is write-only, not returned by API

	// Final Parameters - Complete SDK Coverage (Phase 5)
	// Only set disk_offering_id if it was explicitly configured by user
	if _, ok := d.GetOk("disk_offering_id"); ok {
		d.Set("disk_offering_id", so.Diskofferingid)
	}

	// Only set disk_offering_strictness if it was explicitly configured
	if _, ok := d.GetOk("disk_offering_strictness"); ok {
		d.Set("disk_offering_strictness", so.Diskofferingstrictness)
	}

	// Note: external_details is write-only, not returned by API

	// Only set is_system if it was explicitly configured
	if _, ok := d.GetOk("is_system"); ok {
		d.Set("is_system", so.Issystem)
	}

	if so.Leaseduration > 0 {
		d.Set("lease_duration", so.Leaseduration)
	}

	if so.Leaseexpiryaction != "" {
		d.Set("lease_expiry_action", so.Leaseexpiryaction)
	}

	// Set service offering details (only user-configured keys)
	if so.Serviceofferingdetails != nil {
		// Only process if user originally configured service_offering_details
		if configuredDetails, ok := d.GetOk("service_offering_details"); ok {
			details := make(map[string]string)
			configuredMap := configuredDetails.(map[string]interface{})

			// Only include keys that the user explicitly configured
			// This prevents drift from CloudStack-managed keys like "External:*", "purge.db.entities", etc.
			for userKey := range configuredMap {
				if cloudValue, exists := so.Serviceofferingdetails[userKey]; exists {
					details[userKey] = cloudValue
				}
			}

			if len(details) > 0 {
				d.Set("service_offering_details", details)
			}
		}
	}

	return nil
}

func resourceCloudStackServiceOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Check if name, display_text or host_tags changed
	if d.HasChange("name") || d.HasChange("display_text") || d.HasChange("host_tags") {
		// Create parameters structure
		p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())

		if d.HasChange("name") {
			p.SetName(d.Get("name").(string))
		}

		if d.HasChange("display_text") {
			p.SetDisplaytext(d.Get("display_text").(string))
		}

		if d.HasChange("host_tags") {
			p.SetHosttags(d.Get("host_tags").(string))
		}

		// Update the service offering
		_, err := cs.ServiceOffering.UpdateServiceOffering(p)
		if err != nil {
			return fmt.Errorf("Error updating service offering %s: %s", d.Get("name").(string), err)
		}
	}

	return resourceCloudStackServiceOfferingRead(d, meta)
}

func resourceCloudStackServiceOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create parameters structure
	p := cs.ServiceOffering.NewDeleteServiceOfferingParams(d.Id())

	// Delete the service offering
	_, err := cs.ServiceOffering.DeleteServiceOffering(p)
	if err != nil {
		return fmt.Errorf("Error deleting service offering: %s", err)
	}

	return nil
}

// getIntFromDetails extracts an integer value from the service offering details map.
func getIntFromDetails(details map[string]string, key string) interface{} {
	if details == nil {
		return nil
	}
	if val, ok := details[key]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return nil
}

// getServiceOfferingDetails extracts custom service offering details while excluding
// built-in details that are handled as separate schema fields
func getServiceOfferingDetails(details map[string]string) map[string]interface{} {
	if details == nil {
		return make(map[string]interface{})
	}

	// List of built-in details that are handled as separate schema fields
	builtInKeys := map[string]bool{
		"mincpunumber": true,
		"maxcpunumber": true,
		"minmemory":    true,
		"maxmemory":    true,
	}

	result := make(map[string]interface{})
	for k, v := range details {
		if !builtInKeys[k] {
			result[k] = v
		}
	}

	return result
}
