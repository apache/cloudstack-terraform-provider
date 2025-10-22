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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackUserData() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackUserDataCreate,
		Read:   resourceCloudStackUserDataRead,
		Delete: resourceCloudStackUserDataDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the user data",
			},

			"userdata": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user data content to be registered",
			},

			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional account for the user data. Must be used with domain_id.",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional domain ID for the user data. If the account parameter is used, domain_id must also be used.",
			},

			"params": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional comma separated list of variables declared in user data content.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An optional project for the user data.",
			},
		},
	}
}

func resourceCloudStackUserDataCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.User.NewRegisterUserDataParams(d.Get("name").(string), d.Get("userdata").(string))
	if v, ok := d.GetOk("account"); ok {
		p.SetAccount(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(v.(string))
	}
	if v, ok := d.GetOk("params"); ok {
		paramsList := v.(*schema.Set).List()
		var params []string
		for _, param := range paramsList {
			params = append(params, param.(string))
		}
		p.SetParams(strings.Join(params, ","))
	}

	userdata, err := cs.User.RegisterUserData(p)
	if err != nil {
		return fmt.Errorf("Error registering user data: %s", err)
	}

	d.SetId(userdata.Id)

	return resourceCloudStackUserDataRead(d, meta)
}

func resourceCloudStackUserDataRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	id := d.Id()

	p := cs.User.NewListUserDataParams()
	p.SetId(id)

	userdata, err := cs.User.ListUserData(p)
	if err != nil {
		return fmt.Errorf("Error retrieving user data with ID %s: %s", id, err)
	}

	d.Set("name", userdata.UserData[0].Name)
	d.Set("userdata", userdata.UserData[0].Userdata)
	if d.Get("account").(string) != "" {
		d.Set("account", userdata.UserData[0].Account)
	}
	if d.Get("domain_id").(string) != "" {
		d.Set("domain_id", userdata.UserData[0].Domainid)
	}
	if userdata.UserData[0].Params != "" {
		paramsList := strings.Split(userdata.UserData[0].Params, ",")
		var paramsSet []interface{}
		for _, param := range paramsList {
			paramsSet = append(paramsSet, param)
		}
		d.Set("params", schema.NewSet(schema.HashString, paramsSet))
	}
	if userdata.UserData[0].Projectid != "" {
		d.Set("project_id", userdata.UserData[0].Projectid)
	}

	return nil
}

func resourceCloudStackUserDataDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.User.NewDeleteUserDataParams(d.Id())
	_, err := cs.User.DeleteUserData(p)
	if err != nil {
		return fmt.Errorf("Error deleting user data with ID %s: %s", d.Id(), err)
	}

	return nil
}
