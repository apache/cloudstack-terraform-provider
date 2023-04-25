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
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudstackServiceOffering() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackServiceOfferingRead,
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
		},
	}
}

func datasourceCloudStackServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.ServiceOffering.NewListServiceOfferingsParams()
	csServiceOfferings, err := cs.ServiceOffering.ListServiceOfferings(p)

	if err != nil {
		return fmt.Errorf("Failed to list service offerings: %s", err)
	}

	filters := d.Get("filter")
	var serviceOfferings []*cloudstack.ServiceOffering

	for _, s := range csServiceOfferings.ServiceOfferings {
		match, err := applyServiceOfferingFilters(s, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			serviceOfferings = append(serviceOfferings, s)
		}
	}

	if len(serviceOfferings) == 0 {
		return fmt.Errorf("No service offering is matching with the specified regex")
	}
	//return the latest service offering from the list of filtered service according
	//to its creation date
	serviceOffering, err := latestServiceOffering(serviceOfferings)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected service offerings: %s\n", serviceOffering.Displaytext)

	return serviceOfferingDescriptionAttributes(d, serviceOffering)
}

func serviceOfferingDescriptionAttributes(d *schema.ResourceData, serviceOffering *cloudstack.ServiceOffering) error {
	d.SetId(serviceOffering.Id)
	d.Set("name", serviceOffering.Name)
	d.Set("display_text", serviceOffering.Displaytext)

	return nil
}

func latestServiceOffering(serviceOfferings []*cloudstack.ServiceOffering) (*cloudstack.ServiceOffering, error) {
	var latest time.Time
	var serviceOffering *cloudstack.ServiceOffering

	for _, s := range serviceOfferings {
		created, err := time.Parse("2006-01-02T15:04:05-0700", s.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of an service offering: %s", err)
		}

		if created.After(latest) {
			latest = created
			serviceOffering = s
		}
	}

	return serviceOffering, nil
}

func applyServiceOfferingFilters(serviceOffering *cloudstack.ServiceOffering, filters *schema.Set) (bool, error) {
	var serviceOfferingJSON map[string]interface{}
	k, _ := json.Marshal(serviceOffering)
	err := json.Unmarshal(k, &serviceOfferingJSON)
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
		serviceOfferingField := serviceOfferingJSON[updatedName].(string)
		if !r.MatchString(serviceOfferingField) {
			return false, nil
		}

	}
	return true, nil
}
