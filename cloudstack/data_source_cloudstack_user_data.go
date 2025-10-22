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

func dataSourceCloudstackUserData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackUserDataRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
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
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"userdata_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"userdata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"params": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackUserDataRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	p := cs.User.NewListUserDataParams()
	p.SetName(name)

	if v, ok := d.GetOk("account"); ok {
		p.SetAccount(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}

	log.Printf("[DEBUG] Listing user data with name: %s", name)
	userdataList, err := cs.User.ListUserData(p)
	if err != nil {
		return fmt.Errorf("Error listing user data with name %s: %s", name, err)
	}

	if len(userdataList.UserData) == 0 {
		return fmt.Errorf("No user data found with name: %s", name)
	}
	if len(userdataList.UserData) > 1 {
		return fmt.Errorf("Multiple user data entries found with name: %s", name)
	}

	userdata := userdataList.UserData[0]

	d.SetId(userdata.Id)
	d.Set("name", userdata.Name)
	d.Set("account", userdata.Account)
	d.Set("account_id", userdata.Accountid)
	d.Set("domain", userdata.Domain)
	d.Set("domain_id", userdata.Domainid)
	d.Set("userdata_id", userdata.Id)
	d.Set("params", userdata.Params)

	if userdata.Project != "" {
		d.Set("project", userdata.Project)
		d.Set("project_id", userdata.Projectid)
	}

	if userdata.Userdata != "" {
		decoded, err := base64.StdEncoding.DecodeString(userdata.Userdata)
		if err != nil {
			d.Set("userdata", userdata.Userdata) // Fallback: use raw data
		} else {
			d.Set("userdata", string(decoded))
		}
	}
	return nil
}
