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
	"regexp"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackUserData() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudStackUserDataRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},

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

			"params": {
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
		},
	}
}

func dataSourceCloudStackUserDataRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.User.NewListUserDataParams()
	csUserData, err := cs.User.ListUserData(p)

	if err != nil {
		return fmt.Errorf("Failed to list UserData: %s", err)
	}

	filters := d.Get("filter")
	var userDataList []*cloudstack.UserData

	for _, ud := range csUserData.UserData {
		match, err := applyUserDataFilters(ud, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			userDataList = append(userDataList, ud)
		}
	}

	if len(userDataList) == 0 {
		return fmt.Errorf("No UserData is matching with the specified regex")
	}

	// Return the first match
	userData := userDataList[0]
	log.Printf("[DEBUG] Selected UserData: %s\n", userData.Name)

	return userDataDescriptionAttributes(d, userData)
}

func userDataDescriptionAttributes(d *schema.ResourceData, userData *cloudstack.UserData) error {
	d.SetId(userData.Id)
	d.Set("name", userData.Name)
	d.Set("account", userData.Account)
	d.Set("account_id", userData.Accountid)
	d.Set("domain", userData.Domain)
	d.Set("domain_id", userData.Domainid)
	d.Set("params", userData.Params)
	d.Set("project", userData.Project)
	d.Set("project_id", userData.Projectid)

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

func applyUserDataFilters(userData *cloudstack.UserData, filters *schema.Set) (bool, error) {
	var userDataMap map[string]interface{}
	userDataMap = map[string]interface{}{
		"name":      userData.Name,
		"account":   userData.Account,
		"accountid": userData.Accountid,
		"domain":    userData.Domain,
		"domainid":  userData.Domainid,
		"params":    userData.Params,
		"project":   userData.Project,
		"projectid": userData.Projectid,
		"userdata":  userData.Userdata,
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		log.Print(m)
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		log.Print(updatedName)
		if userDataField, ok := userDataMap[updatedName].(string); ok {
			if !r.MatchString(userDataField) {
				return false, nil
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}
