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

func resourceCloudStackNetworkOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNetworkOfferingCreate,
		Read:   resourceCloudStackNetworkOfferingRead,
		Update: resourceCloudStackNetworkOfferingUpdate,
		Delete: resourceCloudStackNetworkOfferingDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"guest_ip_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "the ID of the containing domain(s), null for public offerings",
			},
			"network_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "data transfer rate in megabits per second allowed",
			},
			"network_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates the mode with which the network will operate. Valid option: NATTED or ROUTED",
			},
			"max_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "maximum number of concurrent connections supported by the network offering",
			},
			"conserve_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if the network offering is IP conserve mode enabled",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "set to true if the offering is to be enabled during creation. Default is false",
			},
			"for_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if network offering is meant to be used for VPC, false otherwise.",
			},
			"for_nsx": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if network offering is meant to be used for NSX, false otherwise",
			},
			"internet_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The internet protocol of network offering. Options are ipv4 and dualstack. Default is ipv4. dualstack will create a network offering that supports both IPv4 and IPv6",
			},
			"routing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the routing mode for the network offering. Supported types are: Static or Dynamic.",
			},
			"specify_vlan": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if network offering supports vlans, false otherwise",
			},
			"supported_services": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "the list of supported services",
			},
			"service_provider_list": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "provider to service mapping. If not specified, the provider for the service will be mapped to the default provider on the physical network",
			},
			"specify_ip_ranges": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if network offering supports specifying ip ranges; defaulted to false if not specified",
			},
			"specify_as_number": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "true if network offering supports choosing AS number",
			},
		},
	}
}

func resourceCloudStackNetworkOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	display_text := d.Get("display_text").(string)
	guest_ip_type := d.Get("guest_ip_type").(string)
	traffic_type := d.Get("traffic_type").(string)

	// Create a new parameter struct
	p := cs.NetworkOffering.NewCreateNetworkOfferingParams(display_text, guest_ip_type, name, traffic_type)

	if guest_ip_type == "Shared" {
		p.SetSpecifyvlan(true)
		p.SetSpecifyipranges(true)
	}

	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.([]string))
	}

	if v, ok := d.GetOk("network_rate"); ok {
		p.SetNetworkrate(v.(int))
	}

	if v, ok := d.GetOk("network_mode"); ok {
		p.SetNetworkmode(v.(string))
	}

	if v, ok := d.GetOk("max_connections"); ok {
		p.SetMaxconnections(v.(int))
	}

	if v, ok := d.GetOk("conserve_mode"); ok {
		p.SetConservemode(v.(bool))
	}

	if v, ok := d.GetOk("enable"); ok {
		p.SetEnable(v.(bool))
	}

	if v, ok := d.GetOk("for_vpc"); ok {
		p.SetForvpc(v.(bool))
	}

	if v, ok := d.GetOk("for_nsx"); ok {
		p.SetFornsx(v.(bool))
	}

	if v, ok := d.GetOk("internet_protocol"); ok {
		p.SetInternetprotocol(v.(string))
	}

	if v, ok := d.GetOk("routing_mode"); ok {
		p.SetRoutingmode(v.(string))
	}

	if v, ok := d.GetOk("specify_vlan"); ok {
		p.SetSpecifyvlan(v.(bool))
	}

	var supported_services []string
	if v, ok := d.GetOk("supported_services"); ok {
		for _, supported_service := range v.(*schema.Set).List() {
			supported_services = append(supported_services, supported_service.(string))
		}
	}
	p.SetSupportedservices(supported_services)

	if v, ok := d.GetOk("service_provider_list"); ok {
		m := make(map[string]string)
		for key, value := range v.(map[string]interface{}) {
			m[key] = value.(string)
		}
		p.SetServiceproviderlist(m)
	}

	if v, ok := d.GetOk("specify_ip_ranges"); ok {
		p.SetSpecifyipranges(v.(bool))
	}

	if v, ok := d.GetOk("specify_as_number"); ok {
		p.SetSpecifyasnumber(v.(bool))
	}

	log.Printf("[DEBUG] Creating Network Offering %s", name)
	n, err := cs.NetworkOffering.CreateNetworkOffering(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Network Offering %s successfully created", name)
	d.SetId(n.Id)

	return resourceCloudStackNetworkOfferingRead(d, meta)
}

func resourceCloudStackNetworkOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	// Check if the name is changed and if so, update the network offering
	if d.HasChange("name") {
		log.Printf("[DEBUG] Name changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new name
		p.SetName(d.Get("name").(string))

		// Update the name
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the name for network offering %s: %s", name, err)
		}

	}

	// Check if the display text is changed and if so, update the virtual machine
	if d.HasChange("display_text") {
		log.Printf("[DEBUG] Display text changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new display text
		p.SetName(d.Get("display_text").(string))

		// Update the display text
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the display text for network offering %s: %s", name, err)
		}

	}

	// Check if the guest ip type is changed and if so, update the virtual machine
	if d.HasChange("guest_ip_type") {
		log.Printf("[DEBUG] Guest ip type changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new guest ip type
		p.SetName(d.Get("guest_ip_type").(string))

		// Update the guest ip type
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the guest ip type for network offering %s: %s", name, err)
		}

	}

	// Check if the traffic type is changed and if so, update the virtual machine
	if d.HasChange("traffic_type") {
		log.Printf("[DEBUG] Traffic type changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new traffic type
		p.SetName(d.Get("traffic_type").(string))

		// Update the traffic type
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the traffic type for network offering %s: %s", name, err)
		}

	}

	return resourceCloudStackInstanceRead(d, meta)
}

func resourceCloudStackNetworkOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.NetworkOffering.NewDeleteNetworkOfferingParams(d.Id())
	_, err := cs.NetworkOffering.DeleteNetworkOffering(p)

	if err != nil {
		return fmt.Errorf("Error deleting Network Offering: %s", err)
	}

	return nil
}

func resourceCloudStackNetworkOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Network Offering %s", d.Get("name").(string))

	// Get the Network Offering details
	n, count, err := cs.NetworkOffering.GetNetworkOfferingByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Network Offering %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId(n.Id)
	d.Set("name", n.Name)
	d.Set("display_text", n.Displaytext)
	d.Set("guest_ip_type", n.Guestiptype)
	d.Set("traffic_type", n.Traffictype)
	d.Set("network_rate", n.Networkrate)
	d.Set("network_mode", n.Networkmode)
	d.Set("conserve_mode", n.Conservemode)
	d.Set("enable", n.State == "Enabled")
	d.Set("for_vpc", n.Forvpc)
	d.Set("for_nsx", n.Fornsx)
	d.Set("specify_vlan", n.Specifyvlan)
	d.Set("specify_ip_ranges", n.Specifyipranges)
	d.Set("specify_as_number", n.Specifyasnumber)
	d.Set("internet_protocol", n.Internetprotocol)
	d.Set("routing_mode", n.Routingmode)
	d.Set("max_connections", n.Maxconnections)

	// Set supported services
	if len(n.Service) > 0 {
		services := make([]string, len(n.Service))
		for i, service := range n.Service {
			services[i] = service.Name
		}
		d.Set("supported_services", services)
	}

	// Set service provider list
	if len(n.Service) > 0 {
		serviceProviders := make(map[string]string)
		for _, service := range n.Service {
			if len(service.Provider) > 0 {
				serviceProviders[service.Name] = service.Provider[0].Name
			}
		}
		d.Set("service_provider_list", serviceProviders)
	}

	return nil
}
