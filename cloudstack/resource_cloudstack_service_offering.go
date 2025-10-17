package cloudstack

import (
	"fmt"

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
				Description: "Whether service offering allows custom CPU/memory or not",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
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
				ForceNew:    true,
			},

			"max_iops": {
				Description: "Maximum IOPS",
				Type:        schema.TypeInt,
				Optional:    true,
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

	if v, ok := d.GetOk("customized"); ok {
		p.SetCustomized(v.(bool))
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

	if v, ok := d.GetOk("domain_id"); ok {
		domains := make([]string, 0)
		for _, d := range v.([]interface{}) {
			domains = append(domains, d.(string))
		}
		p.SetDomainid(domains)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		zones := make([]string, 0)
		for _, z := range v.([]interface{}) {
			zones = append(zones, z.(string))
		}
		p.SetZoneid(zones)
	}

	// Handle service offering details (custom configurations only, GPU is separate)
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
	d.Set("customized", so.Iscustomized)

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

	// IOPS limits (min/max CPU and memory are write-only, not returned by API)
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

	// Tags field is write-only, not returned by API - skip setting

	// Domain and Zone IDs
	if so.Domainid != "" {
		d.Set("domain_id", []string{so.Domainid})
	}

	if so.Zoneid != "" {
		d.Set("zone_id", []string{so.Zoneid})
	}

	// Set service offering details (excluding system-managed keys)
	if so.Serviceofferingdetails != nil {
		details := make(map[string]string)
		// List of keys that are set via dedicated schema fields and should not appear in serviceofferingdetails
		systemKeys := map[string]bool{
			"pciDevice":    true, // GPU card
			"vgpuType":     true, // vGPU profile
			"mincpunumber": true, // min_cpu_number parameter
			"maxcpunumber": true, // max_cpu_number parameter
			"minmemory":    true, // min_memory parameter
			"maxmemory":    true, // max_memory parameter
		}
		for k, v := range so.Serviceofferingdetails {
			if !systemKeys[k] {
				details[k] = v
			}
		}
		// Only set if user originally configured service_offering_details
		// This prevents drift from CloudStack's internal detail keys
		if _, ok := d.GetOk("service_offering_details"); ok && len(details) > 0 {
			d.Set("service_offering_details", details)
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
