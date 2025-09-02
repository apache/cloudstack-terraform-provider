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

func dataSourceCloudstackAutoscaleVMProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackAutoscaleVMProfileRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"service_offering": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"template": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"destroy_vm_grace_period": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"counter_param_list": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"user_data": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data_details": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"other_deploy_params": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCloudstackAutoscaleVMProfileRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	id, idOk := d.GetOk("id")

	if !idOk {
		return fmt.Errorf("'id' must be specified")
	}

	p := cs.AutoScale.NewListAutoScaleVmProfilesParams()
	p.SetId(id.(string))

	resp, err := cs.AutoScale.ListAutoScaleVmProfiles(p)
	if err != nil {
		return fmt.Errorf("failed to list autoscale VM profiles: %s", err)
	}

	if resp.Count == 0 {
		return fmt.Errorf("autoscale VM profile with ID %s not found", id.(string))
	}

	profile := resp.AutoScaleVmProfiles[0]

	log.Printf("[DEBUG] Found autoscale VM profile: %s", profile.Id)

	d.SetId(profile.Id)
	d.Set("service_offering", profile.Serviceofferingid)
	d.Set("template", profile.Templateid)
	d.Set("zone", profile.Zoneid)
	d.Set("account_name", profile.Account)
	d.Set("domain_id", profile.Domainid)
	if profile.Projectid != "" {
		d.Set("project_id", profile.Projectid)
	}
	d.Set("display", profile.Fordisplay)

	if profile.Userdata != "" {
		d.Set("user_data", profile.Userdata)
	}

	return nil
}
