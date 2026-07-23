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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackVPCOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackVPCOfferingCreate,
		Read:   resourceCloudStackVPCOfferingRead,
		Update: resourceCloudStackVPCOfferingUpdate,
		Delete: resourceCloudStackVPCOfferingDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "the ID of the containing domain(s), null for public offerings",
				ForceNew:    true,
			},
			"zone_id": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "the ID of the containing zone(s), null for public offerings",
				ForceNew:    true,
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set to true if the offering is to be enabled during creation. Default is false",
			},
			"for_nsx": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if the VPC offering is meant to be used for NSX, false otherwise",
				ForceNew:    true,
			},
			"nsx_support_lb": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if the NSX supports Lb service, false otherwise",
				ForceNew:    true,
			},
			"internet_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The internet protocol of the VPC offering. Options are ipv4 and dualstack. Default is ipv4. dualstack will create a VPC offering that supports both IPv4 and IPv6",
				ForceNew:    true,
			},
			"network_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Indicates the mode with which the VPC will operate. Valid option: NATTED or ROUTED",
				ForceNew:    true,
			},
			"routing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "the routing mode for the VPC offering. Supported types are: Static or Dynamic.",
				ForceNew:    true,
			},
			"specify_as_number": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if the VPC offering supports choosing AS number",
				ForceNew:    true,
			},
			"network_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "network provider for the VPC offering",
				ForceNew:    true,
			},
			"service_offering_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the ID of the service offering for the VPC router appliance",
				ForceNew:    true,
			},
			"supported_services": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "the list of supported services",
				ForceNew:    true,
			},
			"service_provider_list": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "provider to service mapping. If not specified, the provider for the service will be mapped to the default provider on the physical network",
				ForceNew:    true,
			},
			"service_capability_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "desired service capabilities as part of the VPC offering",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "the service the capability applies to, e.g. SourceNat",
						},
						"capability_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "the capability type, e.g. RedundantRouter",
						},
						"capability_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "the capability value, e.g. true",
						},
					},
				},
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "true if this is the default VPC offering",
			},
		},
	}
}

func resourceCloudStackVPCOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	displayText := d.Get("display_text").(string)

	// Create a new parameter struct
	p := cs.VPC.NewCreateVPCOfferingParams(displayText, name)

	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(expandStringList(v.([]interface{})))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		p.SetZoneid(expandStringList(v.([]interface{})))
	}

	if v, ok := d.GetOk("enable"); ok {
		p.SetEnable(v.(bool))
	}

	if v, ok := d.GetOk("for_nsx"); ok {
		p.SetFornsx(v.(bool))
	}

	if v, ok := d.GetOk("nsx_support_lb"); ok {
		p.SetNsxsupportlb(v.(bool))
	}

	if v, ok := d.GetOk("internet_protocol"); ok {
		p.SetInternetprotocol(v.(string))
	}

	if v, ok := d.GetOk("network_mode"); ok {
		p.SetNetworkmode(v.(string))
	}

	if v, ok := d.GetOk("routing_mode"); ok {
		p.SetRoutingmode(v.(string))
	}

	if v, ok := d.GetOk("specify_as_number"); ok {
		p.SetSpecifyasnumber(v.(bool))
	}

	if v, ok := d.GetOk("network_provider"); ok {
		p.SetProvider(v.(string))
	}

	if v, ok := d.GetOk("service_offering_id"); ok {
		serviceOfferingID, e := retrieveID(cs, "service_offering", v.(string))
		if e != nil {
			return e.Error()
		}
		p.SetServiceofferingid(serviceOfferingID)
	}

	if v, ok := d.GetOk("supported_services"); ok {
		var supportedServices []string
		for _, service := range v.(*schema.Set).List() {
			supportedServices = append(supportedServices, service.(string))
		}
		p.SetSupportedservices(supportedServices)
	}

	if v, ok := d.GetOk("service_provider_list"); ok {
		m := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			m[key] = value.(string)
		}
		p.SetServiceproviderlist(m)
	}

	if v, ok := d.GetOk("service_capability_list"); ok {
		var capabilities []map[string]string
		for _, item := range v.([]interface{}) {
			c := item.(map[string]interface{})
			capabilities = append(capabilities, map[string]string{
				"service":         c["service"].(string),
				"capabilitytype":  c["capability_type"].(string),
				"capabilityvalue": c["capability_value"].(string),
			})
		}
		p.SetServicecapabilitylist(capabilities)
	}

	log.Printf("[DEBUG] Creating VPC Offering %s", name)
	o, err := cs.VPC.CreateVPCOffering(p)
	if err != nil {
		return fmt.Errorf("Error creating VPC Offering %s: %s", name, err)
	}

	d.SetId(o.Id)

	log.Printf("[DEBUG] VPC Offering %s successfully created", name)
	return resourceCloudStackVPCOfferingRead(d, meta)
}

func resourceCloudStackVPCOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving VPC Offering %s", d.Id())

	// Get the VPC Offering details
	o, count, err := cs.VPC.GetVPCOfferingByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] VPC Offering %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId(o.Id)
	d.Set("name", o.Name)
	d.Set("display_text", o.Displaytext)
	d.Set("enable", o.State == "Enabled")
	d.Set("is_default", o.Isdefault)

	d.Set("for_nsx", o.Fornsx)
	d.Set("specify_as_number", o.Specifyasnumber)

	if o.Internetprotocol != "" {
		d.Set("internet_protocol", o.Internetprotocol)
	}

	if o.Networkmode != "" {
		d.Set("network_mode", o.Networkmode)
	}

	if o.Routingmode != "" {
		d.Set("routing_mode", o.Routingmode)
	}

	if o.Domainid != "" {
		d.Set("domain_id", []string{o.Domainid})
	}

	if o.Zoneid != "" {
		d.Set("zone_id", []string{o.Zoneid})
	}

	if len(o.Service) > 0 {
		services := make([]string, len(o.Service))
		serviceProviders := make(map[string]string)
		var capabilities []map[string]interface{}
		for i, service := range o.Service {
			services[i] = service.Name
			if len(service.Provider) > 0 {
				serviceProviders[service.Name] = service.Provider[0].Name
			}
			for _, capability := range service.Capability {
				capabilities = append(capabilities, map[string]interface{}{
					"service":          service.Name,
					"capability_type":  capability.Name,
					"capability_value": capability.Value,
				})
			}
		}
		d.Set("supported_services", services)
		d.Set("service_provider_list", serviceProviders)
		if len(capabilities) > 0 {
			d.Set("service_capability_list", capabilities)
		}
	}

	return nil
}

func resourceCloudStackVPCOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)

	// Check if the name is changed and if so, update the VPC offering
	if d.HasChange("name") {
		log.Printf("[DEBUG] Name changed for %s, starting update", name)

		p := cs.VPC.NewUpdateVPCOfferingParams(d.Id())
		p.SetName(name)

		_, err := cs.VPC.UpdateVPCOffering(p)
		if err != nil {
			return fmt.Errorf("Error updating the name for VPC offering %s: %s", name, err)
		}
	}

	// Check if the display text is changed and if so, update the VPC offering
	if d.HasChange("display_text") {
		log.Printf("[DEBUG] Display text changed for %s, starting update", name)

		p := cs.VPC.NewUpdateVPCOfferingParams(d.Id())
		p.SetDisplaytext(d.Get("display_text").(string))

		_, err := cs.VPC.UpdateVPCOffering(p)
		if err != nil {
			return fmt.Errorf("Error updating the display text for VPC offering %s: %s", name, err)
		}
	}

	// Check if the enabled state is changed and if so, update the VPC offering
	if d.HasChange("enable") {
		log.Printf("[DEBUG] State changed for %s, starting update", name)

		state := "Disabled"
		if d.Get("enable").(bool) {
			state = "Enabled"
		}

		p := cs.VPC.NewUpdateVPCOfferingParams(d.Id())
		p.SetState(state)

		_, err := cs.VPC.UpdateVPCOffering(p)
		if err != nil {
			return fmt.Errorf("Error updating the state for VPC offering %s: %s", name, err)
		}
	}

	return resourceCloudStackVPCOfferingRead(d, meta)
}

func resourceCloudStackVPCOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.VPC.NewDeleteVPCOfferingParams(d.Id())
	_, err := cs.VPC.DeleteVPCOffering(p)
	if err != nil {
		return fmt.Errorf("Error deleting VPC Offering: %s", err)
	}

	return nil
}

func expandStringList(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}
