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

func resourceCloudStackNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNetworkACLCreate,
		Read:   resourceCloudStackNetworkACLRead,
		Delete: resourceCloudStackNetworkACLDelete,
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

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	vpcID := d.Get("vpc_id").(string)

	// If no project is explicitly set, try to inherit it from the VPC
	// and set it in the state so the Read function can use it
	if _, ok := d.GetOk("project"); !ok {
		// Get the VPC to retrieve its project
		// Use projectid=-1 to search across all projects
		vpc, count, err := cs.VPC.GetVPCByID(vpcID, cloudstack.WithProject("-1"))
		if err == nil && count > 0 && vpc.Projectid != "" {
			log.Printf("[DEBUG] Inheriting project %s from VPC %s", vpc.Projectid, vpcID)
			// Set the project in the resource data for state management
			d.Set("project", vpc.Project)
		}
	}

	// Create a new parameter struct
	// Note: CreateNetworkACLListParams doesn't support SetProjectid
	// The ACL will be created in the same project as the VPC automatically
	p := cs.NetworkACL.NewCreateNetworkACLListParams(name, vpcID)

	// Set the description
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	} else {
		p.SetDescription(name)
	}

	// Create the new network ACL list
	r, err := cs.NetworkACL.CreateNetworkACLList(p)
	if err != nil {
		return fmt.Errorf("Error creating network ACL list %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCloudStackNetworkACLRead(d, meta)
}

func resourceCloudStackNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the network ACL list details
	f, count, err := cs.NetworkACL.GetNetworkACLListByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project").(string)),
	)

	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] Network ACL list %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", f.Name)
	d.Set("description", f.Description)
	d.Set("vpc_id", f.Vpcid)

	// If project is not already set in state, try to get it from the VPC
	if d.Get("project").(string) == "" {
		// Get the VPC to retrieve its project
		vpc, vpcCount, vpcErr := cs.VPC.GetVPCByID(f.Vpcid, cloudstack.WithProject("-1"))
		if vpcErr == nil && vpcCount > 0 && vpc.Project != "" {
			log.Printf("[DEBUG] Setting project %s from VPC %s for ACL %s", vpc.Project, f.Vpcid, f.Name)
			setValueOrID(d, "project", vpc.Project, vpc.Projectid)
		}
	}

	return nil
}

func resourceCloudStackNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.NetworkACL.NewDeleteNetworkACLListParams(d.Id())

	// Delete the network ACL list
	_, err := Retry(3, func() (interface{}, error) {
		return cs.NetworkACL.DeleteNetworkACLList(p)
	})
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting network ACL list %s: %s", d.Get("name").(string), err)
	}

	return nil
}
