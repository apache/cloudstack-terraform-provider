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
	"log"
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
				Description: "Speed of CPU in Mhz",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"host_tags": {
				Description: "The host tag for this service offering",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"limit_cpu_use": {
				Description: "Restrict the CPU usage to committed service offering",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"memory": {
				Description: "The total memory of the service offering in MB",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"offer_ha": {
				Description: "The HA for the service offering",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"storage_type": {
				Description: "The storage type of the service offering. Values are local and shared",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "shared",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)

					if v == "local" || v == "shared" {
						return
					}

					errs = append(errs, fmt.Errorf("storage type should be either local or shared, got %s", v))

					return
				},
			},
			"customized": {
				Description: "Whether service offering allows custom CPU/memory or not",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"min_cpu_number": {
				Description: "Minimum number of CPU cores allowed",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"max_cpu_number": {
				Description: "Maximum number of CPU cores allowed",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"min_memory": {
				Description: "Minimum memory allowed (MB)",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"max_memory": {
				Description: "Maximum memory allowed (MB)",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"encrypt_root": {
				Description: "Encrypt the root disk for VMs using this service offering",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"storage_tags": {
				Description: "Storage tags to associate with the service offering",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"service_offering_details": {
				Description: "Service offering details for GPU configuration and other advanced settings",
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		DeprecationMessage: "cloudstack_service_offering is deprecated and will be removed in a future release. Please use the following resources instead: " +
			"cloudstack_service_offering_constrained, cloudstack_service_offering_unconstrained or cloudstack_service_offering_fixed.",
	}
}

func resourceCloudStackServiceOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	display_text := d.Get("display_text").(string)

	// Create a new parameter struct
	p := cs.ServiceOffering.NewCreateServiceOfferingParams(display_text, name)

	cpuNumber, cpuNumberOk := d.GetOk("cpu_number")
	if cpuNumberOk {
		p.SetCpunumber(cpuNumber.(int))
	}

	cpuSpeed, cpuSpeedOk := d.GetOk("cpu_speed")
	if cpuSpeedOk {
		p.SetCpuspeed(cpuSpeed.(int))
	}

	if v, ok := d.GetOk("host_tags"); ok {
		p.SetHosttags(v.(string))
	}

	if v, ok := d.GetOk("limit_cpu_use"); ok {
		p.SetLimitcpuuse(v.(bool))
	}

	memory, memoryOk := d.GetOk("memory")
	if memoryOk {
		p.SetMemory(memory.(int))
	}

	if v, ok := d.GetOk("offer_ha"); ok {
		p.SetOfferha(v.(bool))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		p.SetStoragetype(v.(string))
	}

	customized := false
	if v, ok := d.GetOk("customized"); ok {
		customized = v.(bool)
	}
	if !cpuNumberOk && !cpuSpeedOk && !memoryOk {
		customized = true
	}
	p.SetCustomized(customized)

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

	if v, ok := d.GetOk("encrypt_root"); ok {
		p.SetEncryptroot(v.(bool))
	}

	if v, ok := d.GetOk("storage_tags"); ok {
		p.SetTags(v.(string))
	}

	if details, ok := d.GetOk("service_offering_details"); ok {
		serviceOfferingDetails := make(map[string]string)
		for k, v := range details.(map[string]interface{}) {
			serviceOfferingDetails[k] = v.(string)
		}
		p.SetServiceofferingdetails(serviceOfferingDetails)
	}

	log.Printf("[DEBUG] Creating Service Offering %s", name)
	s, err := cs.ServiceOffering.CreateServiceOffering(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Service Offering %s successfully created", name)
	d.SetId(s.Id)

	return resourceCloudStackServiceOfferingRead(d, meta)
}

func resourceCloudStackServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Service Offering %s", d.Get("name").(string))

	// Get the Service Offering details
	s, count, err := cs.ServiceOffering.GetServiceOfferingByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Service Offering %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(s.Id)

	fields := map[string]interface{}{
		"name":                     s.Name,
		"display_text":             s.Displaytext,
		"cpu_number":               s.Cpunumber,
		"cpu_speed":                s.Cpuspeed,
		"host_tags":                s.Hosttags,
		"limit_cpu_use":            s.Limitcpuuse,
		"memory":                   s.Memory,
		"offer_ha":                 s.Offerha,
		"storage_type":             s.Storagetype,
		"customized":               s.Iscustomized,
		"min_cpu_number":           getIntFromDetails(s.Serviceofferingdetails, "mincpunumber"),
		"max_cpu_number":           getIntFromDetails(s.Serviceofferingdetails, "maxcpunumber"),
		"min_memory":               getIntFromDetails(s.Serviceofferingdetails, "minmemory"),
		"max_memory":               getIntFromDetails(s.Serviceofferingdetails, "maxmemory"),
		"encrypt_root":             s.Encryptroot,
		"storage_tags":             s.Storagetags,
		"service_offering_details": getServiceOfferingDetails(s.Serviceofferingdetails),
	}

	for k, v := range fields {
		d.Set(k, v)
	}

	return nil
}

func resourceCloudStackServiceOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	// Check if the name is changed and if so, update the service offering
	if d.HasChange("name") {
		log.Printf("[DEBUG] Name changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())

		// Set the new name
		p.SetName(d.Get("name").(string))

		// Update the name
		_, err := cs.ServiceOffering.UpdateServiceOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the name for service offering %s: %s", name, err)
		}

	}

	// Check if the display text is changed and if so, update seervice offering
	if d.HasChange("display_text") {
		log.Printf("[DEBUG] Display text changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())

		// Set the new display text
		p.SetName(d.Get("display_text").(string))

		// Update the display text
		_, err := cs.ServiceOffering.UpdateServiceOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the display text for service offering %s: %s", name, err)
		}

	}

	if d.HasChange("host_tags") {
		log.Printf("[DEBUG] Host tags changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())

		// Set the new host tags
		p.SetHosttags(d.Get("host_tags").(string))

		// Update the host tags
		_, err := cs.ServiceOffering.UpdateServiceOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the host tags for service offering %s: %s", name, err)
		}

	}

	if d.HasChange("storage_tags") {
		log.Printf("[DEBUG] Tags changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())

		// Set the new tags
		p.SetStoragetags(d.Get("storage_tags").(string))

		// Update the host tags
		_, err := cs.ServiceOffering.UpdateServiceOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the storage tags for service offering %s: %s", name, err)
		}

	}

	return resourceCloudStackServiceOfferingRead(d, meta)
}

func resourceCloudStackServiceOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.ServiceOffering.NewDeleteServiceOfferingParams(d.Id())
	_, err := cs.ServiceOffering.DeleteServiceOffering(p)

	if err != nil {
		return fmt.Errorf("Error deleting Service Offering: %s", err)
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
