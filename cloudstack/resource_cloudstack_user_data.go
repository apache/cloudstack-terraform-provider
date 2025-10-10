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
	"encoding/base64"
	"fmt"
	"log"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackUserData() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackUserDataCreate,
		Read:   resourceCloudStackUserDataRead,
		Update: resourceCloudStackUserDataUpdate,
		Delete: resourceCloudStackUserDataDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Required: true,
			},

			"account": {
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

			"params": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackUserDataCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	userData := d.Get("user_data").(string)

	// Encode user data as base64 if not already encoded
	ud, err := getUserData(userData)
	if err != nil {
		return fmt.Errorf("Error encoding user data: %s", err)
	}

	// Validate user data size (CloudStack API limitation)
	if len(ud) > 1048576 { // 1MB in bytes for base64 encoded content
		return fmt.Errorf("UserData is too large: %d bytes (max 1MB for base64 encoded content). Consider reducing content or using CloudStack global setting vm.userdata.max.length", len(ud))
	}

	// Create a new parameter struct
	p := cs.User.NewRegisterUserDataParams(name, ud)

	// Set optional parameters
	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	if params, ok := d.GetOk("params"); ok {
		p.SetParams(params.(string))
	}

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	log.Printf("[DEBUG] Registering UserData %s", name)
	r, err := cs.User.RegisterUserData(p)
	if err != nil {
		return fmt.Errorf("Error registering UserData %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCloudStackUserDataRead(d, meta)
}

func resourceCloudStackUserDataRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the UserData details
	userData, count, err := cs.User.GetUserDataByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project").(string)),
	)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] UserData %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	// Update the config
	d.Set("name", userData.Name)
	d.Set("account", userData.Account)
	d.Set("account_id", userData.Accountid)
	d.Set("domain", userData.Domain)
	d.Set("domain_id", userData.Domainid)
	d.Set("params", userData.Params)

	// Decode and set the user data
	if userData.Userdata != "" {
		decoded, err := base64.StdEncoding.DecodeString(userData.Userdata)
		if err != nil {
			// If decoding fails, assume it's already plain text
			d.Set("user_data", userData.Userdata)
		} else {
			d.Set("user_data", string(decoded))
		}
	}

	return nil
}

func resourceCloudStackUserDataUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	if d.HasChange("user_data") || d.HasChange("params") {
		// For updates, we need to delete and recreate as CloudStack doesn't have an update API
		log.Printf("[DEBUG] UserData %s has changes, recreating", name)

		// Validate user data size before proceeding with update
		if d.HasChange("user_data") {
			userData := d.Get("user_data").(string)
			ud, err := getUserData(userData)
			if err != nil {
				return fmt.Errorf("Error encoding user data: %s", err)
			}
			if len(ud) > 1048576 { // 1MB in bytes for base64 encoded content
				return fmt.Errorf("UserData is too large: %d bytes (max 1MB for base64 encoded content). Consider reducing content or using CloudStack global setting vm.userdata.max.length", len(ud))
			}
		}

		// Delete the old UserData
		if err := resourceCloudStackUserDataDelete(d, meta); err != nil {
			return err
		}

		// Create new UserData
		return resourceCloudStackUserDataCreate(d, meta)
	}

	return resourceCloudStackUserDataRead(d, meta)
}

func resourceCloudStackUserDataDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.User.NewDeleteUserDataParams(d.Id())

	// Set optional parameters if they were used during creation
	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	if project, ok := d.GetOk("project"); ok {
		if !cloudstack.IsID(project.(string)) {
			id, _, err := cs.Project.GetProjectID(project.(string))
			if err != nil {
				return err
			}
			p.SetProjectid(id)
		} else {
			p.SetProjectid(project.(string))
		}
	}

	log.Printf("[DEBUG] Deleting UserData %s", d.Get("name").(string))
	_, err := cs.User.DeleteUserData(p)
	if err != nil {
		return fmt.Errorf("Error deleting UserData: %s", err)
	}

	return nil
}
