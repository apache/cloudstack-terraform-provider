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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackPhysicalNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackPhysicalNetworkCreate,
		Read:   resourceCloudStackPhysicalNetworkRead,
		Update: resourceCloudStackPhysicalNetworkUpdate,
		Delete: resourceCloudStackPhysicalNetworkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"broadcast_domain_range": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validOptions := []string{
						"ZONE",
						"POD",
					}
					err := validateOptions(validOptions, v.(string), k)
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isolation_methods": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validOptions := []string{
						"VLAN",
						"VXLAN",
						"GRE",
						"SST",
						"BCF_SEGMENT",
						"SSP",
						"ODL",
						"L3VPN",
						"VCS",
					}
					err := validateOptions(validOptions, v.(string), k)
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
				ForceNew: true,
			},
			"network_speed": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					vnetRange, _ := v.(string)

					ranges := strings.Split(vnetRange, ",")
					for _, r := range ranges {
						// Split the range into start and end
						parts := strings.Split(r, "-")
						if len(parts) != 2 {
							errors = append(errors, fmt.Errorf("%q must consist of a range defined by two numbers separated by a dash, got %s", k, r))
							continue
						}

						start, errStart := strconv.Atoi(parts[0])
						end, errEnd := strconv.Atoi(parts[1])
						if errStart != nil || errEnd != nil {
							errors = append(errors, fmt.Errorf("%q contains non-numeric values in the range: %s", k, r))
							continue
						}

						if start < 0 || start > 4095 || end < 0 || end > 4095 {
							errors = append(errors, fmt.Errorf("%q numbers must be between 0 and 4095, got range: %s", k, r))
							continue
						}

						if start > end {
							errors = append(errors, fmt.Errorf("%q range start must be less than or equal to range end, got start: %d and end: %d", k, start, end))
						}
					}

					return
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Enabled",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validOptions := []string{
						"Enabled",
						"Disabled",
					}
					err := validateOptions(validOptions, v.(string), k)
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
		},
	}
}

func resourceCloudStackPhysicalNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	zoneID := d.Get("zone_id").(string)

	// Create a new parameter struct
	p := cs.Network.NewCreatePhysicalNetworkParams(
		name,
		zoneID,
	)

	if broadcastDomainRange, ok := d.GetOk("broadcast_domain_range"); ok {
		p.SetBroadcastdomainrange(broadcastDomainRange.(string))
	}

	if domainID, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(domainID.(string))
	}

	if isolationMethods, ok := d.GetOk("isolation_methods"); ok {
		p.SetIsolationmethods([]string{isolationMethods.(string)})
	}

	if tags, ok := d.GetOk("tags"); ok {
		p.SetTags(convertToStringArray(tags.([]interface{})))
	}

	if networkSpeed, ok := d.GetOk("network_speed"); ok {
		p.SetNetworkspeed(networkSpeed.(string))
	}

	if vlan, ok := d.GetOk("vlan"); ok {
		p.SetVlan(vlan.(string))
	}

	log.Printf("[DEBUG] Creating Physical Network %s", name)

	n, err := cs.Network.CreatePhysicalNetwork(p)

	if err != nil {
		return err
	}

	if state, ok := d.GetOk("state"); ok {
		p := cs.Network.NewUpdatePhysicalNetworkParams(n.Id)
		log.Printf("[DEBUG] state changed for %s, starting update", name)
		p.SetState(state.(string))
		_, err := cs.Network.UpdatePhysicalNetwork(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the state %s: %s", name, err)
		}
	}

	log.Printf("[DEBUG] Physical Network %s successfully created", name)

	d.SetId(n.Id)

	return resourceCloudStackPhysicalNetworkRead(d, meta)
}

func resourceCloudStackPhysicalNetworkRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Physical Network %s", d.Get("name").(string))

	// Get the Physical Network details
	p, count, err := cs.Network.GetPhysicalNetworkByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Physical Network %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(p.Id)
	d.Set("name", p.Name)
	d.Set("zone_id", p.Zoneid)
	d.Set("broadcast_domain_range", p.Broadcastdomainrange)
	d.Set("domain_id", p.Domainid)
	d.Set("isolation_methods", p.Isolationmethods)
	d.Set("network_speed", p.Networkspeed)
	d.Set("vlan", p.Vlan)
	d.Set("state", p.State)
	d.Set("zone_name", p.Zonename)
	d.Set("tags", p.Tags)
	return nil
}

func resourceCloudStackPhysicalNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	if d.HasChange("network_speed") {
		p := cs.Network.NewUpdatePhysicalNetworkParams(d.Id())
		log.Printf("[DEBUG] network_speed changed for %s, starting update", name)
		p.SetNetworkspeed(d.Get("network_speed").(string))
		_, err := cs.Network.UpdatePhysicalNetwork(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the network_speed %s: %s", name, err)
		}
	}
	if d.HasChange("state") {
		p := cs.Network.NewUpdatePhysicalNetworkParams(d.Id())
		log.Printf("[DEBUG] state changed for %s, starting update", name)
		p.SetState(d.Get("state").(string))
		_, err := cs.Network.UpdatePhysicalNetwork(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the state %s: %s", name, err)
		}
	}
	if d.HasChange("tags") {
		p := cs.Network.NewUpdatePhysicalNetworkParams(d.Id())
		log.Printf("[DEBUG] tags changed for %s, starting update", name)
		p.SetTags(convertToStringArray(d.Get("tags").([]interface{})))
		_, err := cs.Network.UpdatePhysicalNetwork(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the tags %s: %s", name, err)
		}
	}
	if d.HasChange("vlan") {
		p := cs.Network.NewUpdatePhysicalNetworkParams(d.Id())
		log.Printf("[DEBUG] vlan changed for %s, starting update", name)
		p.SetVlan(d.Get("vlan").(string))
		_, err := cs.Network.UpdatePhysicalNetwork(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the vlan %s: %s", name, err)
		}
	}

	return resourceCloudStackPhysicalNetworkRead(d, meta)
}

func resourceCloudStackPhysicalNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Network.NewDeletePhysicalNetworkParams(d.Id())
	_, err := cs.Network.DeletePhysicalNetwork(p)

	if err != nil {
		return fmt.Errorf("Error deleting Physical Network: %s", err)
	}
	return nil
}

func convertToStringArray(interfaces []interface{}) []string {
	strings := make([]string, len(interfaces))
	for i, v := range interfaces {
		strings[i] = v.(string)
	}
	return strings
}
