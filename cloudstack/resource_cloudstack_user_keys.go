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

func resourceCloudStackUserKeys() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackUserKeysCreate,
		Read:   resourceCloudStackUserKeysRead,
		Delete: resourceCloudStackUserKeysDelete,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the user for which to register API keys",
			},

			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Account name (required for non-admin users)",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Domain ID (required for non-admin users)",
			},

			// Computed attributes
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated API key for the user",
			},

			"secret_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated secret key for the user",
			},
		},
	}
}

func resourceCloudStackUserKeysCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	userID := d.Get("user_id").(string)

	log.Printf("[DEBUG] Registering API keys for user %s", userID)

	// Create new register user keys params
	p := cs.User.NewRegisterUserKeysParams(userID)

	// Register the user keys
	r, err := cs.User.RegisterUserKeys(p)
	if err != nil {
		return fmt.Errorf("Error registering user keys for user %s: %s", userID, err)
	}

	log.Printf("[DEBUG] API keys successfully registered for user %s", userID)

	// Set the resource ID to the user ID since keys are tied to the user
	d.SetId(userID)

	// Set the computed attributes
	if err := d.Set("api_key", r.Apikey); err != nil {
		return fmt.Errorf("Error setting api_key: %s", err)
	}

	if err := d.Set("secret_key", r.Secretkey); err != nil {
		return fmt.Errorf("Error setting secret_key: %s", err)
	}

	return resourceCloudStackUserKeysRead(d, meta)
}

func resourceCloudStackUserKeysRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	userID := d.Id()

	log.Printf("[DEBUG] Reading user keys for user %s", userID)

	// List users to check if the user still exists
	p := cs.User.NewListUsersParams()
	p.SetId(userID)

	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	if domainID, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(domainID.(string))
	}

	r, err := cs.User.ListUsers(p)
	if err != nil {
		return fmt.Errorf("Error reading user %s: %s", userID, err)
	}

	if r.Count == 0 {
		log.Printf("[DEBUG] User %s no longer exists, removing from state", userID)
		d.SetId("")
		return nil
	}

	user := r.Users[0]

	// Verify this is the correct user
	if user.Id != userID {
		log.Printf("[DEBUG] User ID mismatch, removing from state")
		d.SetId("")
		return nil
	}

	// Set user_id attribute
	if err := d.Set("user_id", user.Id); err != nil {
		return fmt.Errorf("Error setting user_id: %s", err)
	}

	if err := d.Set("account", user.Account); err != nil {
		return fmt.Errorf("Error setting account: %s", err)
	}

	if err := d.Set("domain_id", user.Domainid); err != nil {
		return fmt.Errorf("Error setting domain_id: %s", err)
	}

	// Note: We cannot read back the actual API keys from CloudStack for security reasons
	// The keys are only available at creation time
	// We keep the existing values in state if they exist

	return nil
}

func resourceCloudStackUserKeysDelete(d *schema.ResourceData, meta interface{}) error {
	// Note: CloudStack doesn't provide an API to delete user keys directly
	// Keys are typically disabled by disabling the user account
	// For this resource, we'll just remove it from state
	log.Printf("[DEBUG] Removing user keys resource from state for user %s", d.Id())
	d.SetId("")
	return nil
}
