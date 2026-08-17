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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackRoleCreate,
		Read:   resourceCloudStackRoleRead,
		Update: resourceCloudStackRoleUpdate,
		Delete: resourceCloudStackRoleDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of the role, valid options are: Admin, ResourceAdmin, DomainAdmin, User",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether the role will be visible to all users (public) or only to root admins (private). Default is true.",
			},
			"role_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the role to be cloned from. Either role_id or type must be passed in.",
			},
		},
	}
}

func resourceCloudStackRoleCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)

	// Create a new parameter struct
	p := cs.Role.NewCreateRoleParams(name)

	// Check if either role_id or type is provided
	roleID, roleIDOk := d.GetOk("role_id")
	roleType, roleTypeOk := d.GetOk("type")

	if roleIDOk {
		p.SetRoleid(roleID.(string))
	} else if roleTypeOk {
		p.SetType(roleType.(string))
	} else {
		// According to the API, either roleid or type must be passed in
		return fmt.Errorf("either role_id or type must be specified")
	}

	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	}

	if isPublic, ok := d.GetOk("is_public"); ok {
		p.SetIspublic(isPublic.(bool))
	}

	log.Printf("[DEBUG] Creating Role %s", name)
	r, err := cs.Role.CreateRole(p)

	if err != nil {
		return fmt.Errorf("Error creating Role: %s", err)
	}

	log.Printf("[DEBUG] Role %s successfully created", name)
	d.SetId(r.Id)

	return resourceCloudStackRoleRead(d, meta)
}

func resourceCloudStackRoleRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the Role details
	r, count, err := cs.Role.GetRoleByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Role %s does not exist", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Role: %s", err)
	}

	d.Set("name", r.Name)
	d.Set("type", r.Type)
	d.Set("description", r.Description)
	d.Set("is_public", r.Ispublic)

	return nil
}

func resourceCloudStackRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Role.NewUpdateRoleParams(d.Id())

	if d.HasChange("name") {
		p.SetName(d.Get("name").(string))
	}

	if d.HasChange("type") {
		p.SetType(d.Get("type").(string))
	}

	if d.HasChange("description") {
		p.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("is_public") {
		p.SetIspublic(d.Get("is_public").(bool))
	}

	log.Printf("[DEBUG] Updating Role %s", d.Get("name").(string))
	_, err := cs.Role.UpdateRole(p)

	if err != nil {
		return fmt.Errorf("Error updating Role: %s", err)
	}

	return resourceCloudStackRoleRead(d, meta)
}

func resourceCloudStackRoleDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Role.NewDeleteRoleParams(d.Id())

	log.Printf("[DEBUG] Deleting Role %s", d.Get("name").(string))
	_, err := cs.Role.DeleteRole(p)

	if err != nil {
		return fmt.Errorf("Error deleting Role: %s", err)
	}

	return nil
}
