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

func resourceCloudStackAccount() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCloudStackAccountRead,
		Update: resourceCloudStackAccountUpdate,
		Create: resourceCloudStackAccountCreate,
		Delete: resourceCloudStackAccountDelete,
		Schema: map[string]*schema.Schema{
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
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCloudStackAccountCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	email := d.Get("email").(string)
	first_name := d.Get("first_name").(string)
	last_name := d.Get("last_name").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	role_id := d.Get("role_id").(string)
	account_type := d.Get("account_type").(int)
	account := d.Get("account").(string)
	domain_id := d.Get("domain_id").(string)

	// Create a new parameter struct
	p := cs.Account.NewCreateAccountParams(email, first_name, last_name, password, username)
	p.SetAccounttype(int(account_type))
	p.SetRoleid(role_id)
	if account != "" {
		p.SetAccount(account)
	} else {
		p.SetAccount(username)
	}
	p.SetDomainid(domain_id)

	log.Printf("[DEBUG] Creating Account %s", account)
	a, err := cs.Account.CreateAccount(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Account %s successfully created", account)
	d.SetId(a.Id)

	return resourceCloudStackAccountRead(d, meta)
}

func resourceCloudStackAccountRead(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceCloudStackAccountUpdate(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceCloudStackAccountDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Account.NewDeleteAccountParams(d.Id())
	_, err := cs.Account.DeleteAccount(p)

	if err != nil {
		return fmt.Errorf("Error deleting Account: %s", err)
	}

	return nil
}
