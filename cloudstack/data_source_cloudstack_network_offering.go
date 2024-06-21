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

func dataSourceCloudstackNetworkOffering() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackNetworkOfferingRead,
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
			"guest_ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceCloudStackNetworkOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.NetworkOffering.NewListNetworkOfferingsParams()
	csNetworkOfferings, err := cs.NetworkOffering.ListNetworkOfferings(p)

	if err != nil {
		return fmt.Errorf("Failed to list network offerings: %s", err)
	}

	filters := d.Get("filter")
	var networkOfferings []*cloudstack.NetworkOffering

	for _, n := range csNetworkOfferings.NetworkOfferings {
		match, err := applyNetworkOfferingFilters(n, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			networkOfferings = append(networkOfferings, n)
		}
	}

	if len(networkOfferings) == 0 {
		return fmt.Errorf("No network offering is matching with the specified regex")
	}
	//return the latest network offering from the list of filtered network offerings according
	//to its creation date
	networkOffering, err := latestNetworkOffering(networkOfferings)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected network offerings: %s\n", networkOffering.Displaytext)

	return networkOfferingDescriptionAttributes(d, networkOffering)
}

func networkOfferingDescriptionAttributes(d *schema.ResourceData, networkOffering *cloudstack.NetworkOffering) error {
	d.SetId(networkOffering.Id)
	d.Set("name", networkOffering.Name)
	d.Set("display_text", networkOffering.Displaytext)
	d.Set("guest_ip_type", networkOffering.Guestiptype)
	d.Set("traffic_type", networkOffering.Traffictype)

	return nil
}

func latestNetworkOffering(networkOfferings []*cloudstack.NetworkOffering) (*cloudstack.NetworkOffering, error) {
	var latest time.Time
	var networkOffering *cloudstack.NetworkOffering

	for _, n := range networkOfferings {
		created, err := time.Parse("2006-01-02T15:04:05-0700", n.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a network offering: %s", err)
		}

		if created.After(latest) {
			latest = created
			networkOffering = n
		}
	}

	return networkOffering, nil
}

func applyNetworkOfferingFilters(networkOffering *cloudstack.NetworkOffering, filters *schema.Set) (bool, error) {
	var networkOfferingJSON map[string]interface{}
	k, _ := json.Marshal(networkOffering)
	err := json.Unmarshal(k, &networkOfferingJSON)
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
		networkOfferingField := networkOfferingJSON[updatedName].(string)
		if !r.MatchString(networkOfferingField) {
			return false, nil
		}

	}
	return true, nil
}
