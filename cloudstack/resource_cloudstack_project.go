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
	"fmt"
	"log"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackProjectCreate,
		Read:   resourceCloudStackProjectRead,
		Update: resourceCloudStackProjectUpdate,
		Delete: resourceCloudStackProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed attributes
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

func resourceCloudStackProjectCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	// Create new project params with both displaytext and name
	p := cs.Project.NewCreateProjectParams(d.Get("display_text").(string), d.Get("name").(string)) // Set optional parameters
	if displayText, ok := d.GetOk("display_text"); ok {
		p.SetDisplaytext(displayText.(string))
	} else {
		p.SetDisplaytext(name)
	}

	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	if accountID, ok := d.GetOk("account_id"); ok {
		p.SetAccountid(accountID.(string))
	}

	if domainID, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(domainID.(string))
	}

	if userID, ok := d.GetOk("user_id"); ok {
		p.SetUserid(userID.(string))
	}

	log.Printf("[DEBUG] Creating project %s", name)

	r, err := cs.Project.CreateProject(p)
	if err != nil {
		return fmt.Errorf("Error creating project %s: %s", name, err)
	}

	log.Printf("[DEBUG] Project %s successfully created", name)
	d.SetId(r.Id)

	return resourceCloudStackProjectRead(d, meta)
}

func resourceCloudStackProjectRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Reading project %s", d.Id())

	// Get the project details
	p := cs.Project.NewListProjectsParams()
	p.SetId(d.Id())

	list, err := cs.Project.ListProjects(p)
	if err != nil {
		return fmt.Errorf("Failed to find project: %s", err)
	}

	if list.Count == 0 {
		log.Printf("[DEBUG] Project %s does no longer exist", d.Id())
		d.SetId("")
		return nil
	}

	if list.Count > 1 {
		return fmt.Errorf("Found more than one project with ID: %s", d.Id())
	}

	project := list.Projects[0]

	d.Set("name", project.Name)
	d.Set("display_text", project.Displaytext)
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

	// Set account information from Owner slice
	if len(project.Owner) > 0 {
		d.Set("account", project.Owner[0]["account"])
		d.Set("account_id", project.Owner[0]["accountid"])
	}

	setValueOrID(d, "domain_id", project.Domain, project.Domainid)

	// Set tags
	tags := make(map[string]interface{})
	for _, tag := range project.Tags {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func resourceCloudStackProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	// Create a new parameter struct
	p := cs.Project.NewUpdateProjectParams(d.Id())

	if d.HasChange("name") {
		p.SetName(name)
	}

	if d.HasChange("display_text") {
		p.SetDisplaytext(d.Get("display_text").(string))
	}

	if d.HasChange("account") {
		p.SetAccount(d.Get("account").(string))
	}

	log.Printf("[DEBUG] Updating project %s", name)

	_, err := cs.Project.UpdateProject(p)
	if err != nil {
		return fmt.Errorf("Error updating project %s: %s", name, err)
	}

	return resourceCloudStackProjectRead(d, meta)
}

func resourceCloudStackProjectDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Project.NewDeleteProjectParams(d.Id())
	p.SetCleanup(true) // Clean up all project resources

	log.Printf("[INFO] Deleting project: %s", d.Get("name").(string))

	_, err := cs.Project.DeleteProject(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting project %s: %s", d.Get("name").(string), err)
	}

	return nil
}
