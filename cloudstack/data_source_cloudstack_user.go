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
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudstackUser() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackUserRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceCloudStackUserRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.User.NewListUsersParams()
	csUsers, err := cs.User.ListUsers(p)

	if err != nil {
		return fmt.Errorf("Failed to list users: %s", err)
	}

	filters := d.Get("filter")
	var users []*cloudstack.User

	for _, u := range csUsers.Users {
		match, err := applyUserFilters(u, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			users = append(users, u)
		}
	}

	if len(users) == 0 {
		return fmt.Errorf("No user is matching with the specified regex")
	}
	//return the latest user from the list of filtered userss according
	//to its creation date
	user, err := latestUser(users)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected users: %s\n", user.Username)

	return userDescriptionAttributes(d, user)
}

func userDescriptionAttributes(d *schema.ResourceData, user *cloudstack.User) error {
	d.SetId(user.Id)
	d.Set("account", user.Account)
	d.Set("email", user.Email)
	d.Set("first_name", user.Firstname)
	d.Set("last_name", user.Lastname)
	d.Set("username", user.Username)

	return nil
}

func latestUser(users []*cloudstack.User) (*cloudstack.User, error) {
	var latest time.Time
	var user *cloudstack.User

	for _, u := range users {
		created, err := time.Parse("2006-01-02T15:04:05-0700", u.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a user: %s", err)
		}

		if created.After(latest) {
			latest = created
			user = u
		}
	}

	return user, nil
}

func applyUserFilters(user *cloudstack.User, filters *schema.Set) (bool, error) {
	var userJSON map[string]interface{}
	k, _ := json.Marshal(user)
	err := json.Unmarshal(k, &userJSON)
	if err != nil {
		return false, err
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
		userField := userJSON[updatedName].(string)
		if !r.MatchString(userField) {
			return false, nil
		}
	}
	return true, nil
}
