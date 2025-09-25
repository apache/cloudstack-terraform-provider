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

func dataSourceCloudstackAutoscalePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackAutoscalePolicyRead,

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

			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"quiet_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"condition_ids": {
				Type:     schema.TypeSet,
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
		},
	}
}

func dataSourceCloudstackAutoscalePolicyRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}

	var policy *cloudstack.AutoScalePolicy

	if idOk {
		p := cs.AutoScale.NewListAutoScalePoliciesParams()
		p.SetId(id.(string))

		resp, err := cs.AutoScale.ListAutoScalePolicies(p)
		if err != nil {
			return fmt.Errorf("failed to list autoscale policies: %s", err)
		}

		if resp.Count == 0 {
			return fmt.Errorf("autoscale policy with ID %s not found", id.(string))
		}

		policy = resp.AutoScalePolicies[0]
	} else {
		p := cs.AutoScale.NewListAutoScalePoliciesParams()

		resp, err := cs.AutoScale.ListAutoScalePolicies(p)
		if err != nil {
			return fmt.Errorf("failed to list autoscale policies: %s", err)
		}

		for _, pol := range resp.AutoScalePolicies {
			if pol.Name == name.(string) {
				policy = pol
				break
			}
		}

		if policy == nil {
			return fmt.Errorf("autoscale policy with name %s not found", name.(string))
		}
	}

	log.Printf("[DEBUG] Found autoscale policy: %s", policy.Name)

	d.SetId(policy.Id)
	d.Set("name", policy.Name)
	d.Set("action", policy.Action)
	d.Set("duration", policy.Duration)
	d.Set("quiet_time", policy.Quiettime)
	d.Set("account_name", policy.Account)
	d.Set("domain_id", policy.Domainid)
	if policy.Projectid != "" {
		d.Set("project_id", policy.Projectid)
	}

	conditionIds := make([]string, len(policy.Conditions))
	for i, condition := range policy.Conditions {
		conditionIds[i] = condition.Id
	}
	d.Set("condition_ids", conditionIds)

	return nil
}
