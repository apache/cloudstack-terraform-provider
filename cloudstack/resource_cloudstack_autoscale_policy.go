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

func resourceCloudStackAutoScalePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackAutoScalePolicyCreate,
		Read:   resourceCloudStackAutoScalePolicyRead,
		Update: resourceCloudStackAutoScalePolicyUpdate,
		Delete: resourceCloudStackAutoScalePolicyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the name of the autoscale policy",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the action to be executed if all the conditions evaluate to true for the specified duration",
				ForceNew:    true,
			},
			"duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "the duration in which the conditions have to be true before action is taken",
			},
			"quiet_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "the cool down period in which the policy should not be evaluated after the action has been taken",
			},
			"condition_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "the list of IDs of the conditions that are being evaluated on every interval",
			},
		},
	}
}

func resourceCloudStackAutoScalePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	action := d.Get("action").(string)
	duration := d.Get("duration").(int)

	conditionIds := []string{}
	if v, ok := d.GetOk("condition_ids"); ok {
		conditionSet := v.(*schema.Set)
		for _, id := range conditionSet.List() {
			conditionIds = append(conditionIds, id.(string))
		}
	}

	p := cs.AutoScale.NewCreateAutoScalePolicyParams(action, conditionIds, duration)

	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}
	if v, ok := d.GetOk("quiet_time"); ok {
		p.SetQuiettime(v.(int))
	}

	log.Printf("[DEBUG] Creating autoscale policy")
	resp, err := cs.AutoScale.CreateAutoScalePolicy(p)
	if err != nil {
		return fmt.Errorf("Error creating autoscale policy: %s", err)
	}

	d.SetId(resp.Id)
	log.Printf("[DEBUG] Autoscale policy created with ID: %s", resp.Id)

	return resourceCloudStackAutoScalePolicyRead(d, meta)
}

func resourceCloudStackAutoScalePolicyRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewListAutoScalePoliciesParams()
	p.SetId(d.Id())

	resp, err := cs.AutoScale.ListAutoScalePolicies(p)
	if err != nil {
		return fmt.Errorf("Error retrieving autoscale policy: %s", err)
	}

	if resp.Count == 0 {
		log.Printf("[DEBUG] Autoscale policy %s no longer exists", d.Id())
		d.SetId("")
		return nil
	}

	policy := resp.AutoScalePolicies[0]
	d.Set("name", policy.Name)
	d.Set("action", policy.Action)
	d.Set("duration", policy.Duration)
	d.Set("quiet_time", policy.Quiettime)

	conditionIds := schema.NewSet(schema.HashString, []interface{}{})
	for _, condition := range policy.Conditions {
		var conditionInterface interface{} = condition
		switch v := conditionInterface.(type) {
		case string:
			conditionIds.Add(v)
		case map[string]interface{}:
			if id, ok := v["id"].(string); ok {
				conditionIds.Add(id)
			}
		default:
			log.Printf("[DEBUG] Unexpected condition type: %T, value: %+v", condition, condition)
			conditionIds.Add(fmt.Sprintf("%v", condition))
		}
	}
	d.Set("condition_ids", conditionIds)

	return nil
}

func resourceCloudStackAutoScalePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.HasChange("name") || d.HasChange("condition_ids") || d.HasChange("duration") || d.HasChange("quiet_time") {
		log.Printf("[DEBUG] Updating autoscale policy: %s", d.Id())

		p := cs.AutoScale.NewUpdateAutoScalePolicyParams(d.Id())

		if d.HasChange("name") {
			if v, ok := d.GetOk("name"); ok {
				p.SetName(v.(string))
			}
		}

		if d.HasChange("duration") {
			duration := d.Get("duration").(int)
			p.SetDuration(duration)
		}

		if d.HasChange("quiet_time") {
			if v, ok := d.GetOk("quiet_time"); ok {
				p.SetQuiettime(v.(int))
			}
		}

		if d.HasChange("condition_ids") {
			conditionIds := []string{}
			if v, ok := d.GetOk("condition_ids"); ok {
				conditionSet := v.(*schema.Set)
				for _, id := range conditionSet.List() {
					conditionIds = append(conditionIds, id.(string))
				}
			}
			p.SetConditionids(conditionIds)
		}

		_, err := cs.AutoScale.UpdateAutoScalePolicy(p)
		if err != nil {
			return fmt.Errorf("Error updating autoscale policy: %s", err)
		}

		log.Printf("[DEBUG] Autoscale policy updated successfully: %s", d.Id())
	}

	return resourceCloudStackAutoScalePolicyRead(d, meta)
}

func resourceCloudStackAutoScalePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewDeleteAutoScalePolicyParams(d.Id())

	log.Printf("[DEBUG] Deleting autoscale policy: %s", d.Id())
	_, err := cs.AutoScale.DeleteAutoScalePolicy(p)
	if err != nil {
		return fmt.Errorf("Error deleting autoscale policy: %s", err)
	}

	return nil
}
