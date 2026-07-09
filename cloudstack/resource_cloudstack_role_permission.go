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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloudStackRolePermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackRolePermissionCreate,
		Read:   resourceCloudStackRolePermissionRead,
		Update: resourceCloudStackRolePermissionUpdate,
		Delete: resourceCloudStackRolePermissionDelete,
		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the role the permission (rule) belongs to.",
			},
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The API name or wildcard (e.g. 'list*') the permission applies to.",
			},
			"permission": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"allow", "deny"}, false),
				Description:  "Whether the rule is allowed or denied. Valid options are: allow, deny.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description for the role permission.",
			},
		},
	}
}

func resourceCloudStackRolePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	roleID := d.Get("role_id").(string)
	rule := d.Get("rule").(string)
	permission := d.Get("permission").(string)

	// Create a new parameter struct
	p := cs.Role.NewCreateRolePermissionParams(permission, roleID, rule)

	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	}

	log.Printf("[DEBUG] Creating Role Permission %s (%s) for role %s", rule, permission, roleID)
	r, err := cs.Role.CreateRolePermission(p)

	if err != nil {
		return fmt.Errorf("Error creating Role Permission: %s", err)
	}

	log.Printf("[DEBUG] Role Permission %s successfully created", rule)
	d.SetId(r.Id)

	return resourceCloudStackRolePermissionRead(d, meta)
}

func resourceCloudStackRolePermissionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	roleID := d.Get("role_id").(string)

	// The API only supports listing permissions by role, so fetch them all
	// and locate the one matching this resource's ID.
	p := cs.Role.NewListRolePermissionsParams()
	p.SetRoleid(roleID)

	l, err := cs.Role.ListRolePermissions(p)
	if err != nil {
		return fmt.Errorf("Error listing Role Permissions: %s", err)
	}

	for _, rp := range l.RolePermissions {
		if rp.Id == d.Id() {
			d.Set("role_id", rp.Roleid)
			d.Set("rule", rp.Rule)
			d.Set("permission", rp.Permission)
			d.Set("description", rp.Description)
			return nil
		}
	}

	log.Printf("[DEBUG] Role Permission %s no longer exists", d.Id())
	d.SetId("")

	return nil
}

func resourceCloudStackRolePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Only the permission (allow/deny) can be changed in place; the role_id,
	// rule and description are all ForceNew.
	p := cs.Role.NewUpdateRolePermissionParams(d.Get("role_id").(string))
	p.SetRuleid(d.Id())
	p.SetPermission(d.Get("permission").(string))

	log.Printf("[DEBUG] Updating Role Permission %s", d.Id())
	_, err := cs.Role.UpdateRolePermission(p)

	if err != nil {
		return fmt.Errorf("Error updating Role Permission: %s", err)
	}

	return resourceCloudStackRolePermissionRead(d, meta)
}

func resourceCloudStackRolePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Role.NewDeleteRolePermissionParams(d.Id())

	log.Printf("[DEBUG] Deleting Role Permission %s", d.Id())
	_, err := cs.Role.DeleteRolePermission(p)

	if err != nil {
		return fmt.Errorf("Error deleting Role Permission: %s", err)
	}

	return nil
}
