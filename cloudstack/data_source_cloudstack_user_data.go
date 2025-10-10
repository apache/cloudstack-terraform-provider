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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackUserData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackUserDataRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"account": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"params": {
				Type:     schema.TypeString,
				Computed: true,
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

			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackUserDataRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.User.NewListUserDataParams()
	p.SetName(d.Get("name").(string))

	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	if project, ok := d.GetOk("project"); ok {
		if project.(string) != "" {
			projectid, retrieveErr := retrieveID(cs, "project", project.(string))
			if retrieveErr != nil {
				return retrieveErr.Error()
			}
			p.SetProjectid(projectid)
		}
	}

	resp, err := cs.User.ListUserData(p)
	if err != nil {
		return fmt.Errorf("Error listing UserData: %s", err)
	}

	if resp.Count == 0 || len(resp.UserData) == 0 {
		return fmt.Errorf("UserData %s not found", d.Get("name").(string))
	}

	if resp.Count > 1 && len(resp.UserData) > 1 {
		return fmt.Errorf("Multiple UserData entries found for name %s", d.Get("name").(string))
	}

	userData := resp.UserData[0]

	d.SetId(userData.Id)
	d.Set("name", userData.Name)
	d.Set("account", userData.Account)
	d.Set("account_id", userData.Accountid)
	d.Set("domain", userData.Domain)
	d.Set("domain_id", userData.Domainid)
	d.Set("params", userData.Params)

	if userData.Project != "" {
		d.Set("project", userData.Project)
	}

	if userData.Userdata != "" {
		decoded, err := base64.StdEncoding.DecodeString(userData.Userdata)
		if err != nil {
			d.Set("user_data", userData.Userdata)
		} else {
			d.Set("user_data", string(decoded))
		}
	}

	return nil
}
