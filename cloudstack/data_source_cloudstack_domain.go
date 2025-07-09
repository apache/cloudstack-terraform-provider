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
	csDomains, err := cs.Domain.ListDomains(p)

	if err != nil {
		return fmt.Errorf("failed to list domains: %s", err)
	}

	var domain *cloudstack.Domain

	for _, d := range csDomains.Domains {
		if d.Name == "ROOT" {
			domain = d
			break
		}
	}

	if domain == nil {
		return fmt.Errorf("no domain is matching with the specified name")
	}

	log.Printf("[DEBUG] Selected domain: %s\n", domain.Name)

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
