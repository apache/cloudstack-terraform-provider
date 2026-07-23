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
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackVPCOffering() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackVPCOfferingRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"for_nsx": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"internet_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specify_as_number": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"supported_services": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_provider_list": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_capability_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capability_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capability_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceCloudStackVPCOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.VPC.NewListVPCOfferingsParams()
	csVPCOfferings, err := cs.VPC.ListVPCOfferings(p)

	if err != nil {
		return fmt.Errorf("Failed to list VPC offerings: %s", err)
	}

	filters := d.Get("filter")
	var vpcOfferings []*cloudstack.VPCOffering

	for _, o := range csVPCOfferings.VPCOfferings {
		match, err := applyVPCOfferingFilters(o, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			vpcOfferings = append(vpcOfferings, o)
		}
	}

	if len(vpcOfferings) == 0 {
		return fmt.Errorf("No VPC offering is matching with the specified regex")
	}
	//return the latest VPC offering from the list of filtered VPC offerings according
	//to its creation date
	vpcOffering, err := latestVPCOffering(vpcOfferings)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected VPC offering: %s\n", vpcOffering.Displaytext)

	fullVPCOffering, _, err := cs.VPC.GetVPCOfferingByID(vpcOffering.Id)
	if err != nil {
		return fmt.Errorf("Error retrieving full VPC offering details: %s", err)
	}

	return vpcOfferingDescriptionAttributes(d, fullVPCOffering)
}

func vpcOfferingDescriptionAttributes(d *schema.ResourceData, vpcOffering *cloudstack.VPCOffering) error {
	d.SetId(vpcOffering.Id)
	d.Set("name", vpcOffering.Name)
	d.Set("display_text", vpcOffering.Displaytext)
	d.Set("enable", vpcOffering.State == "Enabled")
	d.Set("for_nsx", vpcOffering.Fornsx)
	d.Set("internet_protocol", vpcOffering.Internetprotocol)
	d.Set("network_mode", vpcOffering.Networkmode)
	d.Set("routing_mode", vpcOffering.Routingmode)
	d.Set("specify_as_number", vpcOffering.Specifyasnumber)
	d.Set("is_default", vpcOffering.Isdefault)

	if len(vpcOffering.Service) > 0 {
		services := make([]string, len(vpcOffering.Service))
		serviceProviders := make(map[string]string)
		var serviceCapabilities []map[string]interface{}
		for i, service := range vpcOffering.Service {
			services[i] = service.Name
			if len(service.Provider) > 0 {
				serviceProviders[service.Name] = service.Provider[0].Name
			}
			for _, capability := range service.Capability {
				serviceCapabilities = append(serviceCapabilities, map[string]interface{}{
					"service":          service.Name,
					"capability_type":  capability.Name,
					"capability_value": capability.Value,
				})
			}
		}
		d.Set("supported_services", services)
		d.Set("service_provider_list", serviceProviders)
		d.Set("service_capability_list", serviceCapabilities)
	}

	return nil
}

func latestVPCOffering(vpcOfferings []*cloudstack.VPCOffering) (*cloudstack.VPCOffering, error) {
	var latest time.Time
	var vpcOffering *cloudstack.VPCOffering

	for _, o := range vpcOfferings {
		created, err := time.Parse("2006-01-02T15:04:05-0700", o.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a VPC offering: %s", err)
		}

		if created.After(latest) {
			latest = created
			vpcOffering = o
		}
	}

	return vpcOffering, nil
}

func applyVPCOfferingFilters(vpcOffering *cloudstack.VPCOffering, filters *schema.Set) (bool, error) {
	var vpcOfferingJSON map[string]interface{}
	k, _ := json.Marshal(vpcOffering)
	err := json.Unmarshal(k, &vpcOfferingJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		vpcOfferingField, ok := vpcOfferingJSON[updatedName].(string)
		if !ok {
			continue
		}
		if !r.MatchString(vpcOfferingField) {
			return false, nil
		}

	}
	return true, nil
}
