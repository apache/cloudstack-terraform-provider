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

func resourceCloudStackUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackUserCreate,
		Read:   resourceCloudStackUserRead,
		Update: resourceCloudStackUserUpdate,
		Delete: resourceCloudStackUserDelete,
		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackUserCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	account := d.Get("account").(string)
	email := d.Get("email").(string)
	first_name := d.Get("first_name").(string)
	last_name := d.Get("last_name").(string)
	password := d.Get("password").(string)
	username := d.Get("username").(string)

	// Create a new parameter struct
	p := cs.User.NewCreateUserParams(account, email, first_name, last_name, password, username)

	log.Printf("[DEBUG] Creating User %s", username)
	u, err := cs.User.CreateUser(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] User %s successfully created", username)
	d.SetId(u.Id)

	return resourceCloudStackUserRead(d, meta)
}

func resourceCloudStackUserUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.User.NewUpdateUserParams(d.Id())

	if d.HasChange("email") {
		p.SetEmail(d.Get("email").(string))
	}
	if d.HasChange("first_name") {
		p.SetFirstname(d.Get("first_name").(string))
	}
	if d.HasChange("last_name") {
		p.SetLastname(d.Get("last_name").(string))
	}
	if d.HasChange("password") {
		p.SetPassword(d.Get("password").(string))
	}

	_, err := cs.User.UpdateUser(p)
	if err != nil {
		return fmt.Errorf("Error updating User %s: %s", d.Id(), err)
	}

	return resourceCloudStackUserRead(d, meta)
}

func resourceCloudStackUserDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.User.NewDeleteUserParams(d.Id())
	_, err := cs.User.DeleteUser(p)

	if err != nil {
		return fmt.Errorf("Error deleting User: %s", err)
	}

	return nil
}

func resourceCloudStackUserRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	user, count, err := cs.User.GetUserByID(d.Id())
	if err != nil {
		if count == 0 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading User %s: %s", d.Id(), err)
	}

	d.Set("account", user.Account)
	d.Set("email", user.Email)
	d.Set("first_name", user.Firstname)
	d.Set("last_name", user.Lastname)
	d.Set("username", user.Username)

	return nil
}
