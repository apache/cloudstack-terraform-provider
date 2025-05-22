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

func dataSourceCloudstackRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackRoleRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
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

	var err error
	var role *cloudstack.Role

	if id, ok := d.GetOk("id"); ok {
		log.Printf("[DEBUG] Getting Role by ID: %s", id.(string))
		role, _, err = cs.Role.GetRoleByID(id.(string))
	} else if name, ok := d.GetOk("name"); ok {
		log.Printf("[DEBUG] Getting Role by name: %s", name.(string))
		role, _, err = cs.Role.GetRoleByName(name.(string))
	} else {
		return fmt.Errorf("Either 'id' or 'name' must be specified")
	}

	if err != nil {
		return err
	}

	d.SetId(role.Id)
	d.Set("name", role.Name)
	d.Set("type", role.Type)
	d.Set("description", role.Description)
	d.Set("is_public", role.Ispublic)

	return nil
}
