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

	// Set service offering details (excluding GPU-related keys)
	if so.Serviceofferingdetails != nil {
		details := make(map[string]string)
		for k, v := range so.Serviceofferingdetails {
			// Skip GPU-related keys as they're handled above
			if k != "pciDevice" && k != "vgpuType" {
				details[k] = v
			}
		}
		if len(details) > 0 {
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
