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

func dataSourceCloudstackAutoscaleVMGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackAutoscaleVMGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"lbrule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"min_members": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_members": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"vm_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"scaleup_policy_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"scaledown_policy_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"available_virtual_machine_count": {
				Type:     schema.TypeInt,
				Computed: true,
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
		},
	}
}

func dataSourceCloudstackAutoscaleVMGroupRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}

	var group *cloudstack.AutoScaleVmGroup

	if idOk {
		p := cs.AutoScale.NewListAutoScaleVmGroupsParams()
		p.SetId(id.(string))

		resp, err := cs.AutoScale.ListAutoScaleVmGroups(p)
		if err != nil {
			return fmt.Errorf("failed to list autoscale VM groups: %s", err)
		}

		if resp.Count == 0 {
			return fmt.Errorf("autoscale VM group with ID %s not found", id.(string))
		}

		group = resp.AutoScaleVmGroups[0]
	} else {
		p := cs.AutoScale.NewListAutoScaleVmGroupsParams()

		resp, err := cs.AutoScale.ListAutoScaleVmGroups(p)
		if err != nil {
			return fmt.Errorf("failed to list autoscale VM groups: %s", err)
		}

		for _, grp := range resp.AutoScaleVmGroups {
			if grp.Name == name.(string) {
				group = grp
				break
			}
		}

		if group == nil {
			return fmt.Errorf("autoscale VM group with name %s not found", name.(string))
		}
	}

	log.Printf("[DEBUG] Found autoscale VM group: %s", group.Name)

	d.SetId(group.Id)
	d.Set("name", group.Name)
	d.Set("lbrule_id", group.Lbruleid)
	d.Set("min_members", group.Minmembers)
	d.Set("max_members", group.Maxmembers)
	d.Set("vm_profile_id", group.Vmprofileid)
	d.Set("state", group.State)
	d.Set("interval", group.Interval)
	d.Set("available_virtual_machine_count", group.Availablevirtualmachinecount)
	d.Set("account_name", group.Account)
	d.Set("domain_id", group.Domainid)
	if group.Projectid != "" {
		d.Set("project_id", group.Projectid)
	}

	scaleupPolicyIds := make([]string, len(group.Scaleuppolicies))
	for i, policy := range group.Scaleuppolicies {
		scaleupPolicyIds[i] = policy.Id
	}
	d.Set("scaleup_policy_ids", scaleupPolicyIds)

	scaledownPolicyIds := make([]string, len(group.Scaledownpolicies))
	for i, policy := range group.Scaledownpolicies {
		scaledownPolicyIds[i] = policy.Id
	}
	d.Set("scaledown_policy_ids", scaledownPolicyIds)

	return nil
}
