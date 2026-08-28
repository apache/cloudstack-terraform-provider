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
	"sync"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var rolePermissionLocks sync.Map

type rolePermissionSpec struct {
	ID          string
	Rule        string
	Permission  string
	Description string
}

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
				Description: "ID of the role the permissions belong to.",
			},
			"authoritative": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether permissions not declared in this resource should be deleted.",
			},
			"managed_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"permission": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ordered list of role permission rules. Rules are evaluated from top to bottom.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the role permission.",
						},
						"rule": {
							Type:        schema.TypeString,
							Required:    true,
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
							Description: "A description for the role permission.",
						},
					},
				},
			},
		},
	}
}

func resourceCloudStackRolePermissionCreate(d *schema.ResourceData, meta interface{}) error {
	roleID := d.Get("role_id").(string)
	d.SetId(roleID)

	roleLock := rolePermissionLock(roleID)
	roleLock.Lock()
	defer roleLock.Unlock()

	reconcileErr := reconcileCloudStackRolePermissions(d, meta, nil)
	if readErr := resourceCloudStackRolePermissionRead(d, meta); readErr != nil {
		if reconcileErr != nil {
			return reconcileErr
		}
		return readErr
	}

	return reconcileErr
}

func resourceCloudStackRolePermissionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	roleID := d.Get("role_id").(string)
	if roleID == "" {
		roleID = d.Id()
	}

	rolePermissions, err := listCloudStackRolePermissions(cs, roleID)
	if err != nil {
		return fmt.Errorf("Error listing Role Permissions: %s", err)
	}

	permissionsByID := make(map[string]*cloudstack.RolePermission)
	for _, rp := range rolePermissions {
		permissionsByID[rp.Id] = rp
	}

	var missing bool
	var readPermissions []interface{}
	used := make(map[string]bool)

	for _, desired := range rolePermissionSpecs(d.Get("permission").([]interface{})) {
		var rp *cloudstack.RolePermission
		if desired.ID != "" {
			rp = permissionsByID[desired.ID]
		}
		if rp == nil {
			rp = findMatchingRolePermission(rolePermissions, desired, used)
		}
		if rp == nil {
			missing = true
			readPermissions = append(readPermissions, rolePermissionState(desired))
			continue
		}

		used[rp.Id] = true
		readPermissions = append(readPermissions, rolePermissionState(rolePermissionSpec{
			ID:          rp.Id,
			Rule:        rp.Rule,
			Permission:  rp.Permission,
			Description: rp.Description,
		}))
	}

	managedIDs := make([]interface{}, 0, len(used))
	for id := range used {
		managedIDs = append(managedIDs, id)
	}
	if err := d.Set("managed_ids", managedIDs); err != nil {
		return fmt.Errorf("Error setting Role Permission managed IDs: %s", err)
	}

	if d.Get("authoritative").(bool) {
		for _, rp := range rolePermissions {
			if used[rp.Id] {
				continue
			}
			readPermissions = append(readPermissions, rolePermissionState(rolePermissionSpec{
				ID:          rp.Id,
				Rule:        rp.Rule,
				Permission:  rp.Permission,
				Description: rp.Description,
			}))
		}
	}

	if err := d.Set("permission", readPermissions); err != nil {
		return fmt.Errorf("Error setting Role Permissions: %s", err)
	}

	if missing {
		log.Printf("[DEBUG] One or more Role Permissions for role %s no longer exist", roleID)
		d.SetId("")
	}

	return nil
}

func resourceCloudStackRolePermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	roleID := d.Get("role_id").(string)
	roleLock := rolePermissionLock(roleID)
	roleLock.Lock()
	defer roleLock.Unlock()

	oldManagedIDs := make(map[string]bool)
	if d.HasChange("managed_ids") {
		oldRaw, _ := d.GetChange("managed_ids")
		for _, id := range oldRaw.(*schema.Set).List() {
			oldManagedIDs[id.(string)] = true
		}
	}

	reconcileErr := reconcileCloudStackRolePermissions(d, meta, oldManagedIDs)
	if readErr := resourceCloudStackRolePermissionRead(d, meta); readErr != nil {
		if reconcileErr != nil {
			return reconcileErr
		}
		return readErr
	}

	return reconcileErr
}

func resourceCloudStackRolePermissionDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	roleID := d.Get("role_id").(string)

	roleLock := rolePermissionLock(roleID)
	roleLock.Lock()
	defer roleLock.Unlock()

	rolePermissions, err := listCloudStackRolePermissions(cs, roleID)
	if err != nil {
		return fmt.Errorf("Error listing Role Permissions: %s", err)
	}

	if d.Get("authoritative").(bool) {
		for _, rp := range rolePermissions {
			if err := deleteCloudStackRolePermission(cs, rp.Id); err != nil {
				return err
			}
		}
		return nil
	}

	rolePermissionsByID := make(map[string]*cloudstack.RolePermission)
	for _, rp := range rolePermissions {
		rolePermissionsByID[rp.Id] = rp
	}

	used := make(map[string]bool)
	for _, permission := range rolePermissionSpecs(d.Get("permission").([]interface{})) {
		ruleID := permission.ID
		if ruleID == "" {
			if rp := findMatchingRolePermission(rolePermissions, permission, used); rp != nil {
				ruleID = rp.Id
			}
		}
		if ruleID == "" || rolePermissionsByID[ruleID] == nil {
			continue
		}
		used[ruleID] = true
		if err := deleteCloudStackRolePermission(cs, ruleID); err != nil {
			return err
		}
	}

	return nil
}

func reconcileCloudStackRolePermissions(d *schema.ResourceData, meta interface{}, oldManagedIDs map[string]bool) error {
	cs := meta.(*cloudstack.CloudStackClient)
	roleID := d.Get("role_id").(string)

	rolePermissions, err := listCloudStackRolePermissions(cs, roleID)
	if err != nil {
		return fmt.Errorf("Error listing Role Permissions: %s", err)
	}

	rolePermissionsByID := make(map[string]*cloudstack.RolePermission)
	for _, rp := range rolePermissions {
		rolePermissionsByID[rp.Id] = rp
	}

	desiredPermissions := rolePermissionSpecs(d.Get("permission").([]interface{}))

	used := make(map[string]bool)
	deleted := make(map[string]bool)
	matches := make([]*cloudstack.RolePermission, len(desiredPermissions))
	toDelete := make([]string, 0)

	for i, desired := range desiredPermissions {
		rp := rolePermissionsByID[desired.ID]
		if rp != nil && (rp.Rule != desired.Rule || rp.Description != desired.Description) {
			rp = nil
		}
		if rp == nil {
			rp = findMatchingRolePermission(rolePermissions, desired, used)
		}
		if rp == nil {
			if stale := rolePermissionsByID[desired.ID]; stale != nil && !used[stale.Id] {
				toDelete = append(toDelete, stale.Id)
				used[stale.Id] = true
			}
			continue
		}
		used[rp.Id] = true
		matches[i] = rp
	}

	for _, id := range toDelete {
		if err := deleteCloudStackRolePermission(cs, id); err != nil {
			return err
		}
		deleted[id] = true
	}

	managedIDs := make([]string, 0, len(desiredPermissions))
	managedIDSet := make(map[string]bool)

	for i, desired := range desiredPermissions {
		rp := matches[i]
		if rp == nil {
			rp, err = createCloudStackRolePermission(cs, roleID, desired)
			if err != nil {
				return err
			}
		} else if rp.Permission != desired.Permission {
			if err := updateCloudStackRolePermission(cs, roleID, rp.Id, desired.Permission); err != nil {
				return err
			}
		}

		managedIDs = append(managedIDs, rp.Id)
		managedIDSet[rp.Id] = true
	}

	if d.Get("authoritative").(bool) {
		for _, rp := range rolePermissions {
			if managedIDSet[rp.Id] || deleted[rp.Id] {
				continue
			}
			if err := deleteCloudStackRolePermission(cs, rp.Id); err != nil {
				return err
			}
			deleted[rp.Id] = true
		}
	} else {
		for oldID := range oldManagedIDs {
			if managedIDSet[oldID] || deleted[oldID] {
				continue
			}
			if rolePermissionsByID[oldID] == nil {
				continue
			}
			if err := deleteCloudStackRolePermission(cs, oldID); err != nil {
				return err
			}
			deleted[oldID] = true
		}
	}

	rolePermissions, err = listCloudStackRolePermissions(cs, roleID)
	if err != nil {
		return fmt.Errorf("Error listing Role Permissions: %s", err)
	}

	ruleOrder := append([]string{}, managedIDs...)
	for _, rp := range rolePermissions {
		if !managedIDSet[rp.Id] {
			ruleOrder = append(ruleOrder, rp.Id)
		}
	}

	if err := orderCloudStackRolePermissions(cs, roleID, ruleOrder); err != nil {
		return err
	}

	return nil
}

func listCloudStackRolePermissions(cs *cloudstack.CloudStackClient, roleID string) ([]*cloudstack.RolePermission, error) {
	p := cs.Role.NewListRolePermissionsParams()
	p.SetRoleid(roleID)

	l, err := cs.Role.ListRolePermissions(p)
	if err != nil {
		return nil, err
	}

	return l.RolePermissions, nil
}

func createCloudStackRolePermission(cs *cloudstack.CloudStackClient, roleID string, permission rolePermissionSpec) (*cloudstack.RolePermission, error) {
	p := cs.Role.NewCreateRolePermissionParams(permission.Permission, roleID, permission.Rule)
	if permission.Description != "" {
		p.SetDescription(permission.Description)
	}

	log.Printf("[DEBUG] Creating Role Permission %s (%s) for role %s", permission.Rule, permission.Permission, roleID)
	r, err := cs.Role.CreateRolePermission(p)
	if err != nil {
		return nil, fmt.Errorf("Error creating Role Permission: %s", err)
	}

	return &cloudstack.RolePermission{
		Id:          r.Id,
		Roleid:      roleID,
		Rule:        permission.Rule,
		Permission:  permission.Permission,
		Description: permission.Description,
	}, nil
}

func updateCloudStackRolePermission(cs *cloudstack.CloudStackClient, roleID, ruleID, permission string) error {
	p := cs.Role.NewUpdateRolePermissionParams(roleID)
	p.SetRuleid(ruleID)
	p.SetPermission(permission)

	log.Printf("[DEBUG] Updating Role Permission %s", ruleID)
	if _, err := cs.Role.UpdateRolePermission(p); err != nil {
		return fmt.Errorf("Error updating Role Permission: %s", err)
	}

	return nil
}

func deleteCloudStackRolePermission(cs *cloudstack.CloudStackClient, ruleID string) error {
	p := cs.Role.NewDeleteRolePermissionParams(ruleID)

	log.Printf("[DEBUG] Deleting Role Permission %s", ruleID)
	if _, err := cs.Role.DeleteRolePermission(p); err != nil {
		return fmt.Errorf("Error deleting Role Permission: %s", err)
	}

	return nil
}

func orderCloudStackRolePermissions(cs *cloudstack.CloudStackClient, roleID string, ruleIDs []string) error {
	if len(ruleIDs) == 0 {
		return nil
	}

	p := cs.Role.NewUpdateRolePermissionParams(roleID)
	p.SetRuleorder(ruleIDs)

	log.Printf("[DEBUG] Ordering Role Permissions for role %s: %v", roleID, ruleIDs)
	if _, err := cs.Role.UpdateRolePermission(p); err != nil {
		return fmt.Errorf("Error ordering Role Permissions: %s", err)
	}

	return nil
}

func rolePermissionSpecs(raw []interface{}) []rolePermissionSpec {
	permissions := make([]rolePermissionSpec, 0, len(raw))
	for _, item := range raw {
		permissionMap := item.(map[string]interface{})
		permissions = append(permissions, rolePermissionSpec{
			ID:          rolePermissionString(permissionMap, "id"),
			Rule:        rolePermissionString(permissionMap, "rule"),
			Permission:  rolePermissionString(permissionMap, "permission"),
			Description: rolePermissionString(permissionMap, "description"),
		})
	}

	return permissions
}

func rolePermissionState(permission rolePermissionSpec) map[string]interface{} {
	return map[string]interface{}{
		"id":          permission.ID,
		"rule":        permission.Rule,
		"permission":  permission.Permission,
		"description": permission.Description,
	}
}

func findMatchingRolePermission(rolePermissions []*cloudstack.RolePermission, desired rolePermissionSpec, used map[string]bool) *cloudstack.RolePermission {
	for _, rp := range rolePermissions {
		if used[rp.Id] {
			continue
		}
		if rp.Rule == desired.Rule && rp.Description == desired.Description {
			return rp
		}
	}

	return nil
}

func rolePermissionLock(roleID string) *sync.Mutex {
	lock, _ := rolePermissionLocks.LoadOrStore(roleID, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func rolePermissionString(permission map[string]interface{}, key string) string {
	if value, ok := permission[key].(string); ok {
		return value
	}

	return ""
}
