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
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloudStackSystemServiceOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackSystemServiceOfferingCreate,
		Read:   resourceCloudStackSystemServiceOfferingRead,
		Update: resourceCloudStackSystemServiceOfferingUpdate,
		Delete: resourceCloudStackSystemServiceOfferingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceCloudStackSystemServiceOfferingCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"display_text": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"system_vm_type": {
				Description: "The system VM type that uses this offering",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
				ValidateFunc: validation.StringInSlice([]string{
					"domainrouter",
					"consoleproxy",
					"secondarystoragevm",
				}, true),
			},
			"cpu_number": {
				Description:  "Number of CPU cores",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"cpu_speed": {
				Description:  "Speed of each CPU core in MHz",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"memory": {
				Description:  "Memory reserved by the system VM in MB",
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(32),
			},
			"storage_type": {
				Description: "The storage type of the offering",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "shared",
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
				ValidateFunc: validation.StringInSlice([]string{"local", "shared"}, true),
			},
			"network_rate": {
				Description:  "Network rate in Mbps; valid only for domain router offerings",
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"offer_ha": {
				Description: "Whether the system offering supports HA",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"limit_cpu_use": {
				Description: "Whether CPU usage is limited to the offering's committed resources",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"host_tags": {
				Description: "Host tags associated with the offering",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"storage_tags": {
				Description: "Storage tags associated with the offering",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain_ids": {
				Description: "IDs of the domains that can use the offering; omit for a public offering",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceCloudStackSystemServiceOfferingCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	_, hasNetworkRate := d.GetOk("network_rate")
	return validateSystemServiceOfferingConfiguration(
		d.Get("system_vm_type").(string),
		hasNetworkRate,
		d.Get("storage_type").(string),
		d.Get("offer_ha").(bool),
	)
}

func validateSystemServiceOfferingConfiguration(systemVMType string, hasNetworkRate bool, storageType string, offerHA bool) error {
	if hasNetworkRate && !strings.EqualFold(systemVMType, "domainrouter") {
		return fmt.Errorf("network_rate can only be set when system_vm_type is domainrouter")
	}

	if offerHA && strings.EqualFold(storageType, "local") {
		return fmt.Errorf("offer_ha cannot be enabled when storage_type is local")
	}

	return nil
}

func resourceCloudStackSystemServiceOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	p := cs.ServiceOffering.NewCreateServiceOfferingParams(d.Get("display_text").(string), name)

	p.SetIssystem(true)
	p.SetSystemvmtype(d.Get("system_vm_type").(string))
	p.SetCpunumber(d.Get("cpu_number").(int))
	p.SetCpuspeed(d.Get("cpu_speed").(int))
	p.SetMemory(d.Get("memory").(int))
	p.SetCustomized(false)
	p.SetStoragetype(d.Get("storage_type").(string))
	p.SetOfferha(d.Get("offer_ha").(bool))
	p.SetLimitcpuuse(d.Get("limit_cpu_use").(bool))

	if v, ok := d.GetOk("network_rate"); ok {
		p.SetNetworkrate(v.(int))
	}
	if v, ok := d.GetOk("host_tags"); ok {
		p.SetHosttags(v.(string))
	}
	if v, ok := d.GetOk("storage_tags"); ok {
		p.SetTags(v.(string))
	}
	if domainIDs := expandSystemServiceOfferingDomainIDs(d.Get("domain_ids")); len(domainIDs) > 0 {
		p.SetDomainid(domainIDs)
	}

	log.Printf("[DEBUG] Creating System Service Offering %s", name)
	offering, err := cs.ServiceOffering.CreateServiceOffering(p)
	if err != nil {
		return fmt.Errorf("error creating System Service Offering %s: %s", name, err)
	}

	d.SetId(offering.Id)
	log.Printf("[DEBUG] System Service Offering %s successfully created", name)

	return resourceCloudStackSystemServiceOfferingRead(d, meta)
}

func resourceCloudStackSystemServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving System Service Offering %s", d.Id())

	offering, count, err := getSystemServiceOfferingByID(cs, d.Id())
	if err != nil {
		return err
	}
	if count == 0 {
		log.Printf("[DEBUG] System Service Offering %s no longer exists", d.Id())
		d.SetId("")
		return nil
	}

	fields := map[string]interface{}{
		"name":           offering.Name,
		"display_text":   offering.Displaytext,
		"system_vm_type": strings.ToLower(offering.Systemvmtype),
		"cpu_number":     offering.Cpunumber,
		"cpu_speed":      offering.Cpuspeed,
		"memory":         offering.Memory,
		"storage_type":   strings.ToLower(offering.Storagetype),
		"network_rate":   offering.Networkrate,
		"offer_ha":       offering.Offerha,
		"limit_cpu_use":  offering.Limitcpuuse,
		"host_tags":      offering.Hosttags,
		"storage_tags":   offering.Storagetags,
		"domain_ids":     flattenSystemServiceOfferingDomainIDs(offering.Domainid),
	}

	for key, value := range fields {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("error setting %s for System Service Offering %s: %s", key, d.Id(), err)
		}
	}

	return nil
}

func resourceCloudStackSystemServiceOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())
	hasChanges := false

	if d.HasChange("name") {
		p.SetName(d.Get("name").(string))
		hasChanges = true
	}
	if d.HasChange("display_text") {
		p.SetDisplaytext(d.Get("display_text").(string))
		hasChanges = true
	}
	if d.HasChange("host_tags") {
		p.SetHosttags(d.Get("host_tags").(string))
		hasChanges = true
	}
	if d.HasChange("storage_tags") {
		p.SetStoragetags(d.Get("storage_tags").(string))
		hasChanges = true
	}
	if d.HasChange("domain_ids") {
		domainIDs := expandSystemServiceOfferingDomainIDs(d.Get("domain_ids"))
		if len(domainIDs) == 0 {
			p.SetDomainid("public")
		} else {
			p.SetDomainid(strings.Join(domainIDs, ","))
		}
		hasChanges = true
	}

	if hasChanges {
		log.Printf("[DEBUG] Updating System Service Offering %s", d.Id())
		if _, err := cs.ServiceOffering.UpdateServiceOffering(p); err != nil {
			return fmt.Errorf("error updating System Service Offering %s: %s", d.Id(), err)
		}
	}

	return resourceCloudStackSystemServiceOfferingRead(d, meta)
}

func resourceCloudStackSystemServiceOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.ServiceOffering.NewDeleteServiceOfferingParams(d.Id())

	log.Printf("[DEBUG] Deleting System Service Offering %s", d.Id())
	if _, err := cs.ServiceOffering.DeleteServiceOffering(p); err != nil {
		return fmt.Errorf("error deleting System Service Offering %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func getSystemServiceOfferingByID(cs *cloudstack.CloudStackClient, id string) (*cloudstack.ServiceOffering, int, error) {
	p := cs.ServiceOffering.NewListServiceOfferingsParams()
	p.SetId(id)
	p.SetIssystem(true)

	response, err := cs.ServiceOffering.ListServiceOfferings(p)
	if err != nil {
		return nil, -1, fmt.Errorf("error retrieving System Service Offering %s: %s", id, err)
	}
	if response.Count == 0 {
		return nil, 0, nil
	}
	if response.Count != 1 {
		return nil, response.Count, fmt.Errorf("expected one System Service Offering for ID %s, got %d", id, response.Count)
	}

	offering := response.ServiceOfferings[0]
	if !offering.Issystem {
		return nil, 1, fmt.Errorf("service offering %s is not a system offering", id)
	}

	return offering, 1, nil
}

func expandSystemServiceOfferingDomainIDs(value interface{}) []string {
	set, ok := value.(*schema.Set)
	if !ok || set.Len() == 0 {
		return nil
	}

	domainIDs := make([]string, 0, set.Len())
	for _, value := range set.List() {
		domainIDs = append(domainIDs, value.(string))
	}

	return domainIDs
}

func flattenSystemServiceOfferingDomainIDs(domainID string) []string {
	if domainID == "" || strings.EqualFold(domainID, "public") {
		return nil
	}

	values := strings.Split(domainID, ",")
	domainIDs := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			domainIDs = append(domainIDs, value)
		}
	}

	return domainIDs
}
