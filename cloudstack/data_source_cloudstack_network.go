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

func dataSourceCloudstackNetwork() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackNetworkRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			// Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_offering_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_offering_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"acl_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"traffic_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func datasourceCloudStackNetworkRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Network.NewListNetworksParams()

	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	csNetworks, err := cs.Network.ListNetworks(p)
	if err != nil {
		return fmt.Errorf("Failed to list networks: %s", err)
	}

	filters := d.Get("filter")
	var networks []*cloudstack.Network

	for _, n := range csNetworks.Networks {
		match, err := applyNetworkFilters(n, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			networks = append(networks, n)
		}
	}

	if len(networks) == 0 {
		return fmt.Errorf("No network is matching with the specified regex")
	}

	network, err := latestNetwork(networks)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected network: %s\n", network.Name)

	return networkDescriptionAttributes(d, network)
}

func networkDescriptionAttributes(d *schema.ResourceData, network *cloudstack.Network) error {
	d.SetId(network.Id)
	d.Set("name", network.Name)
	d.Set("display_text", network.Displaytext)
	d.Set("cidr", network.Cidr)
	d.Set("gateway", network.Gateway)
	d.Set("network_domain", network.Networkdomain)
	d.Set("network_offering_id", network.Networkofferingid)
	d.Set("network_offering_name", network.Networkofferingname)
	d.Set("project", network.Project)
	d.Set("project_id", network.Projectid)
	d.Set("vpc_id", network.Vpcid)
	d.Set("acl_id", network.Aclid)
	d.Set("zone_id", network.Zoneid)
	d.Set("zone_name", network.Zonename)
	d.Set("state", network.State)
	d.Set("type", network.Type)
	d.Set("traffic_type", network.Traffictype)
	d.Set("tags", tagsToMap(network.Tags))

	return nil
}

func latestNetwork(networks []*cloudstack.Network) (*cloudstack.Network, error) {
	var latest time.Time
	var network *cloudstack.Network

	for _, n := range networks {
		created, err := time.Parse("2006-01-02T15:04:05-0700", n.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a network: %s", err)
		}

		if created.After(latest) {
			latest = created
			network = n
		}
	}

	return network, nil
}

func applyNetworkFilters(network *cloudstack.Network, filters *schema.Set) (bool, error) {
	var networkJSON map[string]interface{}
	k, _ := json.Marshal(network)
	err := json.Unmarshal(k, &networkJSON)
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
		networkField := fmt.Sprintf("%v", networkJSON[updatedName])
		if !r.MatchString(networkField) {
			return false, nil
		}
	}

	return true, nil
}
