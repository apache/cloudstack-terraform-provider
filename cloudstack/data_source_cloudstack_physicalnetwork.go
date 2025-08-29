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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudStackPhysicalNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudStackPhysicalNetworkRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"broadcast_domain_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isolation_methods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_speed": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudStackPhysicalNetworkRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Network.NewListPhysicalNetworksParams()
	physicalNetworks, err := cs.Network.ListPhysicalNetworks(p)

	if err != nil {
		return fmt.Errorf("Failed to list physical networks: %s", err)
	}
	filters := d.Get("filter")
	var physicalNetwork *cloudstack.PhysicalNetwork

	for _, pn := range physicalNetworks.PhysicalNetworks {
		match, err := applyPhysicalNetworkFilters(pn, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			physicalNetwork = pn
		}
	}

	if physicalNetwork == nil {
		return fmt.Errorf("No physical network is matching with the specified regex")
	}
	log.Printf("[DEBUG] Selected physical network: %s\n", physicalNetwork.Name)

	return physicalNetworkDescriptionAttributes(d, physicalNetwork)
}

func physicalNetworkDescriptionAttributes(d *schema.ResourceData, physicalNetwork *cloudstack.PhysicalNetwork) error {
	d.SetId(physicalNetwork.Id)
	d.Set("name", physicalNetwork.Name)
	d.Set("broadcast_domain_range", physicalNetwork.Broadcastdomainrange)
	d.Set("network_speed", physicalNetwork.Networkspeed)
	d.Set("vlan", physicalNetwork.Vlan)

	// Set isolation methods
	if physicalNetwork.Isolationmethods != "" {
		methods := strings.Split(physicalNetwork.Isolationmethods, ",")
		d.Set("isolation_methods", methods)
	}

	// Set the zone
	d.Set("zone", physicalNetwork.Zonename)

	// Physical networks don't support tags in CloudStack API

	return nil
}

func applyPhysicalNetworkFilters(physicalNetwork *cloudstack.PhysicalNetwork, filters *schema.Set) (bool, error) {
	var physicalNetworkJSON map[string]interface{}
	k, _ := json.Marshal(physicalNetwork)
	err := json.Unmarshal(k, &physicalNetworkJSON)
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
		physicalNetworkField := physicalNetworkJSON[updatedName].(string)
		if !r.MatchString(physicalNetworkField) {
			return false, nil
		}
	}
	return true, nil
}
