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
			State: importStatePassthrough,
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

			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"accountid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"userid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackProjectCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the name and display_text
	name := d.Get("name").(string)
	displaytext := name
	if v, ok := d.GetOk("display_text"); ok {
		displaytext = v.(string)
	}

	// The CloudStack API expects displaytext as the first parameter and name as the second
	p := cs.Project.NewCreateProjectParams(name, displaytext)

	// Set the domain if provided
	if domain, ok := d.GetOk("domain"); ok {
		domainid, e := retrieveID(cs, "domain", domain.(string))
		if e != nil {
			return e.Error()
		}
		p.SetDomainid(domainid)
	}

	// Set the account if provided
	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	// Set the accountid if provided
	if accountid, ok := d.GetOk("accountid"); ok {
		p.SetAccountid(accountid.(string))
	}

	// Set the userid if provided
	if userid, ok := d.GetOk("userid"); ok {
		p.SetUserid(userid.(string))
	}

	log.Printf("[DEBUG] Creating project %s", name)
	r, err := cs.Project.CreateProject(p)
	if err != nil {
		return fmt.Errorf("Error creating project %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCloudStackProjectRead(d, meta)
}

func resourceCloudStackProjectRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Retrieving project %s", d.Id())

	// Get the project details
	p := cs.Project.NewListProjectsParams()
	p.SetId(d.Id())

	l, err := cs.Project.ListProjects(p)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			log.Printf("[DEBUG] Project %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	if l.Count == 0 {
		log.Printf("[DEBUG] Project %s does no longer exist", d.Id())
		d.SetId("")
		return nil
	}

	project := l.Projects[0]
	// The CloudStack API seems to swap name and display_text, so we need to swap them back
	d.Set("name", project.Displaytext)
	d.Set("display_text", project.Name)
	d.Set("domain", project.Domain)

	// Only set the account, accountid, and userid if they were explicitly set in the configuration
	if _, ok := d.GetOk("account"); ok && len(project.Owner) > 0 {
		for _, owner := range project.Owner {
			if account, ok := owner["account"]; ok {
				d.Set("account", account)
			}
		}
	}

	if _, ok := d.GetOk("accountid"); ok && len(project.Owner) > 0 {
		for _, owner := range project.Owner {
			if accountid, ok := owner["accountid"]; ok {
				d.Set("accountid", accountid)
			}
		}
	}

	if _, ok := d.GetOk("userid"); ok && len(project.Owner) > 0 {
		for _, owner := range project.Owner {
			if userid, ok := owner["userid"]; ok {
				d.Set("userid", userid)
			}
		}
	}

	return nil
}

func resourceCloudStackProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Check if the name or display text is changed
	if d.HasChange("name") || d.HasChange("display_text") {
		// Create a new parameter struct
		p := cs.Project.NewUpdateProjectParams(d.Id())

		// The CloudStack API seems to swap name and display_text, so we need to swap them here
		if d.HasChange("name") {
			p.SetDisplaytext(d.Get("name").(string))
		}

		if d.HasChange("display_text") {
			p.SetName(d.Get("display_text").(string))
		}

		log.Printf("[DEBUG] Updating project %s", d.Id())
		_, err := cs.Project.UpdateProject(p)
		if err != nil {
			return fmt.Errorf("Error updating project %s: %s", d.Id(), err)
		}
	}

	return resourceCloudStackProjectRead(d, meta)
}

func resourceCloudStackProjectDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Project.NewDeleteProjectParams(d.Id())

	log.Printf("[INFO] Deleting project: %s", d.Id())
	_, err := cs.Project.DeleteProject(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting project %s: %s", d.Id(), err)
	}

	return nil
}
