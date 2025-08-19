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

func dataSourceCloudstackProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackProjectRead,

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

			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cpu_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cpu_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cpu_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"memory_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"memory_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"memory_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_storage_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_storage_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_storage_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_storage_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_storage_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_storage_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vm_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vm_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vm_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"volume_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"volume_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"volume_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_available": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_total": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceCloudstackProjectRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Project.NewListProjectsParams()

	csProjects, err := cs.Project.ListProjects(p)
	if err != nil {
		return fmt.Errorf("Failed to list projects: %s", err)
	}

	filters := d.Get("filter")
	var projects []*cloudstack.Project

	for _, project := range csProjects.Projects {
		match, err := applyProjectFilters(project, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			projects = append(projects, project)
		}
	}

	if len(projects) == 0 {
		return fmt.Errorf("No project is matching with the specified regex")
	}

	//return the latest project from the list of projects according to its creation date
	project := projects[0]
	log.Printf("[DEBUG] Selected project: %s\n", project.Displaytext)

	return projectDescriptionAttributes(d, project)
}

func projectDescriptionAttributes(d *schema.ResourceData, project *cloudstack.Project) error {
	d.SetId(project.Id)
	d.Set("name", project.Name)
	d.Set("display_text", project.Displaytext)

	// Extract account information from Owner slice
	if len(project.Owner) > 0 {
		d.Set("account", project.Owner[0]["account"])
		d.Set("account_id", project.Owner[0]["accountid"])
	}

	d.Set("domain", project.Domain)
	d.Set("domain_id", project.Domainid)
	d.Set("state", project.State)
	d.Set("cpu_available", project.Cpuavailable)
	d.Set("cpu_limit", project.Cpulimit)
	d.Set("cpu_total", project.Cputotal)
	d.Set("memory_available", project.Memoryavailable)
	d.Set("memory_limit", project.Memorylimit)
	d.Set("memory_total", project.Memorytotal)
	d.Set("network_available", project.Networkavailable)
	d.Set("network_limit", project.Networklimit)
	d.Set("network_total", project.Networktotal)
	d.Set("primary_storage_available", project.Primarystorageavailable)
	d.Set("primary_storage_limit", project.Primarystoragelimit)
	d.Set("primary_storage_total", project.Primarystoragetotal)
	d.Set("secondary_storage_available", project.Secondarystorageavailable)
	d.Set("secondary_storage_limit", project.Secondarystoragelimit)
	d.Set("secondary_storage_total", project.Secondarystoragetotal)
	d.Set("vm_available", project.Vmavailable)
	d.Set("vm_limit", project.Vmlimit)
	d.Set("vm_total", project.Vmtotal)
	d.Set("volume_available", project.Volumeavailable)
	d.Set("volume_limit", project.Volumelimit)
	d.Set("volume_total", project.Volumetotal)
	d.Set("vpc_available", project.Vpcavailable)
	d.Set("vpc_limit", project.Vpclimit)
	d.Set("vpc_total", project.Vpctotal)

	// Set tags
	tags := make(map[string]interface{})
	for _, tag := range project.Tags {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func applyProjectFilters(project *cloudstack.Project, filters *schema.Set) (bool, error) {
	var projectJSON map[string]interface{}
	t, _ := json.Marshal(project)
	err := json.Unmarshal(t, &projectJSON)
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
		if projectField, exists := projectJSON[updatedName]; exists {
			if projectFieldStr, ok := projectField.(string); ok {
				if !r.MatchString(projectFieldStr) {
					return false, nil
				}
			}
		} else {
			return false, fmt.Errorf("Field %s does not exist in project", updatedName)
		}
	}

	return true, nil
}
