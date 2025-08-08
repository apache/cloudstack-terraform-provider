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

func dataSourceCloudstackRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackRoleRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackRoleRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Role.NewListRolesParams()

	csRoles, err := cs.Role.ListRoles(p)
	if err != nil {
		return fmt.Errorf("failed to list roles: %s", err)
	}

	filters := d.Get("filter")
	var role *cloudstack.Role

	for _, r := range csRoles.Roles {
		match, err := applyRoleFilters(r, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			role = r
			break
		}
	}

	if role == nil {
		return fmt.Errorf("no role is matching with the specified criteria")
	}
	log.Printf("[DEBUG] Selected role: %s\n", role.Name)

	return roleDescriptionAttributes(d, role)
}

func roleDescriptionAttributes(d *schema.ResourceData, role *cloudstack.Role) error {
	d.SetId(role.Id)
	d.Set("name", role.Name)
	d.Set("type", role.Type)
	d.Set("description", role.Description)
	d.Set("is_public", role.Ispublic)

	return nil
}

func latestRole(roles []*cloudstack.Role) (*cloudstack.Role, error) {
	// Since the Role struct doesn't have a Created field,
	// we'll just return the first role in the list
	if len(roles) > 0 {
		return roles[0], nil
	}
	return nil, fmt.Errorf("no roles found")
}

func applyRoleFilters(role *cloudstack.Role, filters *schema.Set) (bool, error) {
	var roleJSON map[string]interface{}
	k, _ := json.Marshal(role)
	err := json.Unmarshal(k, &roleJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")

		// Check if the field exists in the role JSON
		roleField, ok := roleJSON[updatedName]
		if !ok {
			return false, fmt.Errorf("field %s does not exist in role", updatedName)
		}

		// Convert the field to string for regex matching
		var roleFieldStr string
		switch v := roleField.(type) {
		case string:
			roleFieldStr = v
		case bool:
			roleFieldStr = fmt.Sprintf("%t", v)
		case float64:
			roleFieldStr = fmt.Sprintf("%g", v)
		default:
			roleFieldStr = fmt.Sprintf("%v", v)
		}

		if !r.MatchString(roleFieldStr) {
			return false, nil
		}
	}

	return true, nil
}
