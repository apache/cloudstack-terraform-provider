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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
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

func resourceCloudStackAccountRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Reading Account %s", d.Id())

	p := cs.Account.NewListAccountsParams()
	p.SetId(d.Id())

	accounts, err := cs.Account.ListAccounts(p)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %s: %s", d.Id(), err)
	}

	if accounts.Count == 0 {
		log.Printf("[DEBUG] Account %s does no longer exist", d.Id())
		d.SetId("")
		return nil
	}

	account := accounts.Accounts[0]

	d.Set("account_type", account.Accounttype)
	d.Set("role_id", account.Roleid)
	d.Set("account", account.Name)
	d.Set("domain_id", account.Domainid)

	if len(account.User) > 0 {
		user := account.User[0]
		d.Set("email", user.Email)
		d.Set("first_name", user.Firstname)
		d.Set("last_name", user.Lastname)
		d.Set("username", user.Username)
	}

	log.Printf("[DEBUG] Account %s successfully read", d.Id())
	return nil
}

func resourceCloudStackAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Updating Account %s", d.Id())

	// Handle account-level changes
	if d.HasChange("role_id") || d.HasChange("account") || d.HasChange("domain_id") {
		p := cs.Account.NewUpdateAccountParams()
		p.SetId(d.Id())

		if d.HasChange("role_id") {
			p.SetRoleid(d.Get("role_id").(string))
		}
		if d.HasChange("account") {
			p.SetNewname(d.Get("account").(string))
		}
		if d.HasChange("domain_id") {
			p.SetDomainid(d.Get("domain_id").(string))
		}

		_, err := cs.Account.UpdateAccount(p)
		if err != nil {
			return fmt.Errorf("Error updating Account %s: %s", d.Id(), err)
		}
	}

	// Handle user-level changes via updateUser API
	if d.HasChange("email") || d.HasChange("first_name") || d.HasChange("last_name") || d.HasChange("password") {
		lp := cs.Account.NewListAccountsParams()
		lp.SetId(d.Id())
		accounts, err := cs.Account.ListAccounts(lp)
		if err != nil {
			return fmt.Errorf("Error retrieving Account %s for user update: %s", d.Id(), err)
		}
		if accounts.Count == 0 || len(accounts.Accounts[0].User) == 0 {
			return fmt.Errorf("Account %s has no users to update", d.Id())
		}

		userID := accounts.Accounts[0].User[0].Id
		up := cs.User.NewUpdateUserParams(userID)

		if d.HasChange("email") {
			up.SetEmail(d.Get("email").(string))
		}
		if d.HasChange("first_name") {
			up.SetFirstname(d.Get("first_name").(string))
		}
		if d.HasChange("last_name") {
			up.SetLastname(d.Get("last_name").(string))
		}
		if d.HasChange("password") {
			up.SetPassword(d.Get("password").(string))
		}

		_, err = cs.User.UpdateUser(up)
		if err != nil {
			return fmt.Errorf("Error updating user for Account %s: %s", d.Id(), err)
		}
	}

	log.Printf("[DEBUG] Account %s successfully updated", d.Id())
	return resourceCloudStackAccountRead(d, meta)
}

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
