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

func dataSourceCloudstackProject() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackProjectRead,
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

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func datasourceCloudStackProjectRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Project.NewListProjectsParams()
	csProjects, err := cs.Project.ListProjects(p)

	if err != nil {
		return fmt.Errorf("failed to list projects: %s", err)
	}

	filters := d.Get("filter")
	var projects []*cloudstack.Project

	for _, v := range csProjects.Projects {
		match, err := applyProjectFilters(v, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			projects = append(projects, v)
		}
	}

	if len(projects) == 0 {
		return fmt.Errorf("no project matches the specified filters")
	}

	// Return the latest project from the list of filtered projects according
	// to its creation date
	project, err := latestProject(projects)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected project: %s\n", project.Name)

	return projectDescriptionAttributes(d, project)
}

func projectDescriptionAttributes(d *schema.ResourceData, project *cloudstack.Project) error {
	d.SetId(project.Id)
	d.Set("name", project.Name)
	d.Set("display_text", project.Displaytext)
	d.Set("domain", project.Domain)
	d.Set("state", project.State)

	// Handle account information safely
	if len(project.Owner) > 0 {
		for _, owner := range project.Owner {
			if account, ok := owner["account"]; ok {
				d.Set("account", account)
				break
			}
		}
	}

	d.Set("tags", tagsToMap(project.Tags))

	return nil
}

func latestProject(projects []*cloudstack.Project) (*cloudstack.Project, error) {
	var latest time.Time
	var project *cloudstack.Project

	for _, v := range projects {
		created, err := time.Parse("2006-01-02T15:04:05-0700", v.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to parse creation date of a project: %s", err)
		}

		if created.After(latest) {
			latest = created
			project = v
		}
	}

	return project, nil
}

func applyProjectFilters(project *cloudstack.Project, filters *schema.Set) (bool, error) {
	var projectJSON map[string]interface{}
	k, _ := json.Marshal(project)
	err := json.Unmarshal(k, &projectJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("invalid regex: %s", err)
		}

		// Handle special case for owner/account
		if m["name"].(string) == "account" {
			if len(project.Owner) == 0 {
				return false, nil
			}

			found := false
			for _, owner := range project.Owner {
				if account, ok := owner["account"]; ok {
					if r.MatchString(fmt.Sprintf("%v", account)) {
						found = true
						break
					}
				}
			}

			if !found {
				return false, nil
			}
			continue
		}

		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")

		// Handle fields that might not exist in the JSON
		fieldValue, exists := projectJSON[updatedName]
		if !exists {
			return false, nil
		}

		// Handle different types of fields
		switch v := fieldValue.(type) {
		case string:
			if !r.MatchString(v) {
				return false, nil
			}
		case float64:
			if !r.MatchString(fmt.Sprintf("%v", v)) {
				return false, nil
			}
		case bool:
			if !r.MatchString(fmt.Sprintf("%v", v)) {
				return false, nil
			}
		default:
			// Skip fields that aren't simple types
			continue
		}
	}

	return true, nil
}
