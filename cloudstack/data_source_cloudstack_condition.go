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

func dataSourceCloudstackCondition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackConditionRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"counter_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"relational_operator": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"threshold": {
				Type:     schema.TypeFloat,
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

func dataSourceCloudstackConditionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	id, idOk := d.GetOk("id")

	if !idOk {
		return fmt.Errorf("'id' must be specified")
	}

	p := cs.AutoScale.NewListConditionsParams()
	p.SetId(id.(string))

	resp, err := cs.AutoScale.ListConditions(p)
	if err != nil {
		return fmt.Errorf("failed to list conditions: %s", err)
	}

	if resp.Count == 0 {
		return fmt.Errorf("condition with ID %s not found", id.(string))
	}

	condition := resp.Conditions[0]

	log.Printf("[DEBUG] Found condition: %s", condition.Id)

	d.SetId(condition.Id)
	d.Set("counter_id", condition.Counterid)
	d.Set("relational_operator", condition.Relationaloperator)
	d.Set("threshold", condition.Threshold)
	d.Set("account_name", condition.Account)
	d.Set("domain_id", condition.Domainid)
	if condition.Projectid != "" {
		d.Set("project_id", condition.Projectid)
	}

	return nil
}
