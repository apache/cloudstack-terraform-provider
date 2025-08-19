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

func dataSourceCloudStackDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudStackDomainRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			// Computed attributes
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the domain",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the domain",
			},

			"network_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network domain",
			},

			"level": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The level of the domain",
			},

			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path of the domain",
			},

			"parent_domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain ID of the parent domain",
			},

			"parent_domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain name of the parent domain",
			},

			"has_child": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the domain has one or more sub-domains",
			},

			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when this domain was created",
			},

			// Resource limits
			"cpu_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of cpu cores the domain can own",
			},

			"memory_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total memory (in MB) the domain can own",
			},

			"network_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of networks the domain can own",
			},

			"primary_storage_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total primary storage space (in GiB) the domain can own",
			},

			"secondary_storage_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total secondary storage space (in GiB) the domain can own",
			},

			"snapshot_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of snapshots the domain can own",
			},

			"ip_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of public ip addresses this domain can acquire",
			},

			"project_limit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of projects the domain can own",
			},

			// Resource usage
			"cpu_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of cpu cores owned by domain",
			},

			"memory_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total memory (in MB) owned by domain",
			},

			"network_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of networks owned by domain",
			},

			"primary_storage_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total primary storage space (in GiB) owned by domain",
			},

			"secondary_storage_total": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The total secondary storage space (in GiB) owned by domain",
			},

			"snapshot_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of snapshots owned by domain",
			},

			"ip_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of public ip addresses allocated for this domain",
			},

			"project_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of projects owned by domain",
			},

			// Resource availability
			"cpu_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of cpu cores available to be created for this domain",
			},

			"memory_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total memory (in MB) available to be created for this domain",
			},

			"network_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of networks available to be created for this domain",
			},

			"primary_storage_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total primary storage space (in GiB) available to be used for this domain",
			},

			"secondary_storage_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total secondary storage space (in GiB) available to be used for this domain",
			},

			"snapshot_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of snapshots available to be created for this domain",
			},

			"ip_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of public ip addresses available for this domain to acquire",
			},

			"project_available": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The total number of projects available to be created for this domain",
			},
		},
	}
}

func dataSourceCloudStackDomainRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Domain.NewListDomainsParams()

	csDomains, err := cs.Domain.ListDomains(p)
	if err != nil {
		return fmt.Errorf("Failed to list domains: %s", err)
	}

	filters := d.Get("filter")
	var domains []*cloudstack.Domain

	for _, domain := range csDomains.Domains {
		match, err := applyDomainFilters(domain, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			domains = append(domains, domain)
		}
	}

	if len(domains) == 0 {
		return fmt.Errorf("No domain is matching with the specified regex")
	}

	// Return the first matching domain
	domain := domains[0]
	log.Printf("[DEBUG] Selected domain: %s\n", domain.Name)

	return domainDescriptionAttributes(d, domain)
}

func domainDescriptionAttributes(d *schema.ResourceData, domain *cloudstack.Domain) error {
	d.SetId(domain.Id)
	d.Set("domain_id", domain.Id)
	d.Set("name", domain.Name)
	d.Set("network_domain", domain.Networkdomain)
	d.Set("level", domain.Level)
	d.Set("path", domain.Path)
	d.Set("parent_domain_id", domain.Parentdomainid)
	d.Set("parent_domain_name", domain.Parentdomainname)
	d.Set("has_child", domain.Haschild)
	d.Set("created", domain.Created)

	// Set resource limits
	d.Set("cpu_limit", domain.Cpulimit)
	d.Set("memory_limit", domain.Memorylimit)
	d.Set("network_limit", domain.Networklimit)
	d.Set("primary_storage_limit", domain.Primarystoragelimit)
	d.Set("secondary_storage_limit", domain.Secondarystoragelimit)
	d.Set("snapshot_limit", domain.Snapshotlimit)
	d.Set("ip_limit", domain.Iplimit)
	d.Set("project_limit", domain.Projectlimit)

	// Set resource usage
	d.Set("cpu_total", domain.Cputotal)
	d.Set("memory_total", domain.Memorytotal)
	d.Set("network_total", domain.Networktotal)
	d.Set("primary_storage_total", domain.Primarystoragetotal)
	d.Set("secondary_storage_total", domain.Secondarystoragetotal)
	d.Set("snapshot_total", domain.Snapshottotal)
	d.Set("ip_total", domain.Iptotal)
	d.Set("project_total", domain.Projecttotal)

	// Set resource availability
	d.Set("cpu_available", domain.Cpuavailable)
	d.Set("memory_available", domain.Memoryavailable)
	d.Set("network_available", domain.Networkavailable)
	d.Set("primary_storage_available", domain.Primarystorageavailable)
	d.Set("secondary_storage_available", domain.Secondarystorageavailable)
	d.Set("snapshot_available", domain.Snapshotavailable)
	d.Set("ip_available", domain.Ipavailable)
	d.Set("project_available", domain.Projectavailable)

	return nil
}

func applyDomainFilters(domain *cloudstack.Domain, filters *schema.Set) (bool, error) {
	var domainJSON map[string]interface{}
	t, _ := json.Marshal(domain)
	err := json.Unmarshal(t, &domainJSON)
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

		// Check if the field exists in the JSON structure
		if domainField, exists := domainJSON[updatedName]; exists {
			if domainFieldStr, ok := domainField.(string); ok {
				if !r.MatchString(domainFieldStr) {
					return false, nil
				}
			}
		} else {
			return false, fmt.Errorf("Field %s does not exist in domain", updatedName)
		}
	}

	return true, nil
}
