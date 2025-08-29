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
		Read: dataSourceCloudstackProjectRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceCloudstackProjectRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Project Data Source Read Started")

	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Project.NewListProjectsParams()
	csProjects, err := cs.Project.ListProjects(p)

	if err != nil {
		return fmt.Errorf("Failed to list projects: %s", err)
	}

	filters := d.Get("filter")
	var projects []*cloudstack.Project

	for _, i := range csProjects.Projects {
		match, err := applyProjectFilters(i, filters.(*schema.Set))
		if err != nil {
			return err
		}

		if match {
			projects = append(projects, i)
		}
	}

	if len(projects) == 0 {
		return fmt.Errorf("No project is matching with the specified regex")
	}
	//return the latest project from the list of filtered projects according
	//to its creation date
	project, err := latestProject(projects)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected projects: %s\n", project.Displaytext)

	return projectDescriptionAttributes(d, project)
}

func projectDescriptionAttributes(d *schema.ResourceData, project *cloudstack.Project) error {
	d.SetId(project.Id)
	d.Set("project_id", project.Id)
	d.Set("created", project.Created)
	d.Set("display_text", project.Displaytext)
	d.Set("state", project.State)

	d.Set("tags", tagsToMap(project.Tags))

	return nil
}

func latestProject(projects []*cloudstack.Project) (*cloudstack.Project, error) {
	var latest time.Time
	var project *cloudstack.Project

	for _, i := range projects {
		created, err := time.Parse("2006-01-02T15:04:05-0700", i.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of an project: %s", err)
		}

		if created.After(latest) {
			latest = created
			project = i
		}
	}

	return project, nil
}

func applyProjectFilters(project *cloudstack.Project, filters *schema.Set) (bool, error) {
	var projectJSON map[string]interface{}
	i, _ := json.Marshal(project)
	err := json.Unmarshal(i, &projectJSON)
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
		projectField := projectJSON[updatedName].(string)
		if !r.MatchString(projectField) {
			return false, nil
		}

	}
	return true, nil
}
