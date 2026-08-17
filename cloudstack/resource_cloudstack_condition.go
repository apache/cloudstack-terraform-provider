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

func resourceCloudStackCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackConditionCreate,
		Read:   resourceCloudStackConditionRead,
		Update: resourceCloudStackConditionUpdate,
		Delete: resourceCloudStackConditionDelete,

		Schema: map[string]*schema.Schema{
			"counter_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the counter to be used in the condition.",
				ForceNew:    true,
			},
			"relational_operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Relational Operator to be used with threshold. Valid values are EQ, GT, LT, GE, LE.",
			},
			"threshold": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Value for which the Counter will be evaluated with the Operator selected.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the account of the condition. Must be used with the domainId parameter.",
				ForceNew:    true,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the domain ID of the account.",
				ForceNew:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "optional project for the condition",
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudStackConditionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewListConditionsParams()
	p.SetId(d.Id())

	resp, err := cs.AutoScale.ListConditions(p)
	if err != nil {
		return fmt.Errorf("Error retrieving condition: %s", err)
	}

	if resp.Count == 0 {
		log.Printf("[DEBUG] Condition %s no longer exists", d.Id())
		d.SetId("")
		return nil
	}

	condition := resp.Conditions[0]
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

func resourceCloudStackConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.HasChange("relational_operator") || d.HasChange("threshold") {
		log.Printf("[DEBUG] Updating condition: %s", d.Id())

		relationaloperator := d.Get("relational_operator").(string)
		threshold := d.Get("threshold").(float64)

		p := cs.AutoScale.NewUpdateConditionParams(d.Id(), relationaloperator, int64(threshold))

		_, err := cs.AutoScale.UpdateCondition(p)
		if err != nil {
			return fmt.Errorf("Error updating condition: %s", err)
		}

		log.Printf("[DEBUG] Condition updated successfully: %s", d.Id())
	}

	return resourceCloudStackConditionRead(d, meta)
}

func resourceCloudStackConditionDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewDeleteConditionParams(d.Id())

	log.Printf("[DEBUG] Deleting condition: %s", d.Id())
	_, err := cs.AutoScale.DeleteCondition(p)
	if err != nil {
		return fmt.Errorf("Error deleting condition: %s", err)
	}

	return nil
}
func resourceCloudStackConditionCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	counterid := d.Get("counter_id")
	relationaloperator := d.Get("relational_operator").(string)
	threshold := d.Get("threshold").(float64)

	account, accountOk := d.GetOk("account_name")
	domainid, domainOk := d.GetOk("domain_id")

	if !accountOk || !domainOk {
		return fmt.Errorf("account_name and domain_id are required fields")
	}

	p := cs.AutoScale.NewCreateConditionParams(counterid.(string), relationaloperator, int64(threshold))
	p.SetAccount(account.(string))
	p.SetDomainid(domainid.(string))

	if v, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(v.(string))
	}

	log.Printf("[DEBUG] Creating condition")
	resp, err := cs.AutoScale.CreateCondition(p)
	if err != nil {
		return fmt.Errorf("Error creating condition: %s", err)
	}

	d.SetId(resp.Id)
	log.Printf("[DEBUG] Condition created with ID: %s", resp.Id)

	// Set the values directly instead of calling read to avoid JSON unmarshaling issues
	d.Set("counter_id", counterid.(string))
	d.Set("relational_operator", relationaloperator)
	d.Set("threshold", threshold)
	d.Set("account_name", account.(string))
	d.Set("domain_id", domainid.(string))
	if v, ok := d.GetOk("project_id"); ok {
		d.Set("project_id", v.(string))
	}

	return nil
}
