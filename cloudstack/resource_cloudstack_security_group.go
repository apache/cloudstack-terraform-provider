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

func resourceCloudStackSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackSecurityGroupCreate,
		Read:   resourceCloudStackSecurityGroupRead,
		Delete: resourceCloudStackSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
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

			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	// Validate that account is used with domain
	if account, ok := d.GetOk("account"); ok {
		if _, domainOk := d.GetOk("domain"); !domainOk {
			return fmt.Errorf("account parameter requires domain to be set")
		}
		// Account and project are mutually exclusive
		if _, projectOk := d.GetOk("project_id"); projectOk {
			return fmt.Errorf("account and project parameters are mutually exclusive")
		}
		log.Printf("[DEBUG] Creating security group %s for account %s", name, account)
	}

	// Create a new parameter struct
	p := cs.SecurityGroup.NewCreateSecurityGroupParams(name)

	// Set the description
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	} else {
		p.SetDescription(name)
	}

	// Set the account if provided
	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	// If there is a domain supplied, we retrieve and set the domain id
	if err := setDomainid(p, cs, d); err != nil {
		return err
	}

	// If there is a project_id supplied, we set it directly
	if projectID, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(projectID.(string))
	}

	r, err := cs.SecurityGroup.CreateSecurityGroup(p)
	if err != nil {
		return fmt.Errorf("Error creating security group %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCloudStackSecurityGroupRead(d, meta)
}

func resourceCloudStackSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the security group details
	sg, count, err := cs.SecurityGroup.GetSecurityGroupByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project_id").(string)),
	)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Security group %s does not longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	// Update the config
	d.Set("name", sg.Name)
	d.Set("description", sg.Description)

	// Only set account if it was explicitly configured
	if _, ok := d.GetOk("account"); ok {
		d.Set("account", sg.Account)
	}

	setValueOrID(d, "domain", sg.Domain, sg.Domainid)
	setValueOrID(d, "project_id", sg.Project, sg.Projectid)

	return nil
}

func resourceCloudStackSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.SecurityGroup.NewDeleteSecurityGroupParams()
	p.SetId(d.Id())

	// Set the account if provided
	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	// If there is a domain supplied, we retrieve and set the domain id
	if err := setDomainid(p, cs, d); err != nil {
		return err
	}

	// If there is a project_id supplied, we set it directly
	if projectID, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(projectID.(string))
	}

	// Delete the security group
	_, err := cs.SecurityGroup.DeleteSecurityGroup(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting security group: %s", err)
	}

	return nil
}
