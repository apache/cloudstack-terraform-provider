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

func dataSourceCloudstackDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackDomainRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			// Computed values
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackDomainRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Domain Data Source Read Started")

	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Domain.NewListDomainsParams()

	var filterName, filterValue string
	var filterByName, filterByID bool

	// Apply filters if provided
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		for _, f := range filters.(*schema.Set).List() {
			m := f.(map[string]interface{})
			name := m["name"].(string)
			value := m["value"].(string)

			switch name {
			case "name":
				p.SetName(value)
				filterName = value
				filterByName = true
				log.Printf("[DEBUG] Filtering by name: %s", value)
			case "id":
				p.SetId(value)
				filterValue = value
				filterByID = true
				log.Printf("[DEBUG] Filtering by ID: %s", value)
			}
		}
	}

	csDomains, err := cs.Domain.ListDomains(p)
	if err != nil {
		return fmt.Errorf("failed to list domains: %s", err)
	}

	log.Printf("[DEBUG] Found %d domains from CloudStack API", len(csDomains.Domains))

	var domain *cloudstack.Domain

	// If we have results from the API call, select the appropriate domain
	if len(csDomains.Domains) > 0 {
		// If we filtered by ID or name through the API, we should have a specific result
		if filterByID || filterByName {
			// Since we used API filtering, the first result should be our match
			domain = csDomains.Domains[0]
			log.Printf("[DEBUG] Using API-filtered domain: %s", domain.Name)
		} else {
			// If no filters were applied, we need to handle this case
			// This shouldn't happen with the current schema as filters are required
			return fmt.Errorf("no filter criteria specified")
		}
	}

	if domain == nil {
		if filterByName {
			return fmt.Errorf("no domain found with name: %s", filterName)
		} else if filterByID {
			return fmt.Errorf("no domain found with ID: %s", filterValue)
		} else {
			return fmt.Errorf("no domain found matching the specified criteria")
		}
	}

	log.Printf("[DEBUG] Selected domain: %s (ID: %s)", domain.Name, domain.Id)

	return domainDescriptionAttributes(d, domain)
}

func domainDescriptionAttributes(d *schema.ResourceData, domain *cloudstack.Domain) error {
	d.SetId(domain.Id)
	d.Set("domain_id", domain.Id)
	d.Set("name", domain.Name)
	d.Set("network_domain", domain.Networkdomain)
	d.Set("parent_domain_id", domain.Parentdomainid)

	return nil
}
