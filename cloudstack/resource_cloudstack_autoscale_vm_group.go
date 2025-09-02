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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackAutoScaleVMGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackAutoScaleVMGroupCreate,
		Read:   resourceCloudStackAutoScaleVMGroupRead,
		Update: resourceCloudStackAutoScaleVMGroupUpdate,
		Delete: resourceCloudStackAutoScaleVMGroupDelete,

		Schema: map[string]*schema.Schema{
			"lbrule_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the ID of the load balancer rule",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the name of the autoscale vmgroup",
			},

			"min_members": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "the minimum number of members in the vmgroup, the number of instances in the vm group will be equal to or more than this number",
			},

			"max_members": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "the maximum number of members in the vmgroup, The number of instances in the vm group will be equal to or less than this number",
			},

			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "the frequency in which the performance counters to be collected",
			},

			"scaleup_policy_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "list of scaleup autoscale policies",
			},

			"scaledown_policy_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "list of scaledown autoscale policies",
			},

			"vm_profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the autoscale profile that contains information about the vms in the vm group",
			},

			"display": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "an optional field, whether to the display the group to the end user or not",
			},
		},
	}
}

func resourceCloudStackAutoScaleVMGroupCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	lbruleid := d.Get("lbrule_id").(string)
	minmembers := d.Get("min_members").(int)
	maxmembers := d.Get("max_members").(int)
	vmprofileid := d.Get("vm_profile_id").(string)

	scaleUpPolicyIds := []string{}
	if v, ok := d.GetOk("scaleup_policy_ids"); ok {
		scaleUpSet := v.(*schema.Set)
		for _, id := range scaleUpSet.List() {
			scaleUpPolicyIds = append(scaleUpPolicyIds, id.(string))
		}
	}

	scaleDownPolicyIds := []string{}
	if v, ok := d.GetOk("scaledown_policy_ids"); ok {
		scaleDownSet := v.(*schema.Set)
		for _, id := range scaleDownSet.List() {
			scaleDownPolicyIds = append(scaleDownPolicyIds, id.(string))
		}
	}

	p := cs.AutoScale.NewCreateAutoScaleVmGroupParams(lbruleid, maxmembers, minmembers, scaleDownPolicyIds, scaleUpPolicyIds, vmprofileid)

	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}

	if v, ok := d.GetOk("interval"); ok {
		p.SetInterval(v.(int))
	}

	if v, ok := d.GetOk("display"); ok {
		p.SetFordisplay(v.(bool))
	}

	log.Printf("[DEBUG] Creating autoscale VM group")
	resp, err := cs.AutoScale.CreateAutoScaleVmGroup(p)
	if err != nil {
		return fmt.Errorf("Error creating autoscale VM group: %s", err)
	}

	d.SetId(resp.Id)
	log.Printf("[DEBUG] Autoscale VM group created with ID: %s", resp.Id)

	if v, ok := d.GetOk("name"); ok {
		d.Set("name", v.(string))
	}
	d.Set("lbrule_id", lbruleid)
	d.Set("min_members", minmembers)
	d.Set("max_members", maxmembers)
	d.Set("vm_profile_id", vmprofileid)
	if v, ok := d.GetOk("interval"); ok {
		d.Set("interval", v.(int))
	}
	if v, ok := d.GetOk("display"); ok {
		d.Set("display", v.(bool))
	}

	scaleUpSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, id := range scaleUpPolicyIds {
		scaleUpSet.Add(id)
	}
	d.Set("scaleup_policy_ids", scaleUpSet)

	scaleDownSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, id := range scaleDownPolicyIds {
		scaleDownSet.Add(id)
	}
	d.Set("scaledown_policy_ids", scaleDownSet)

	return nil
}

func resourceCloudStackAutoScaleVMGroupRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewListAutoScaleVmGroupsParams()
	p.SetId(d.Id())

	resp, err := cs.AutoScale.ListAutoScaleVmGroups(p)
	if err != nil {
		return fmt.Errorf("Error retrieving autoscale VM group: %s", err)
	}

	if resp.Count == 0 {
		log.Printf("[DEBUG] Autoscale VM group %s no longer exists", d.Id())
		d.SetId("")
		return nil
	}

	group := resp.AutoScaleVmGroups[0]
	d.Set("name", group.Name)
	d.Set("lbrule_id", group.Lbruleid)
	d.Set("min_members", group.Minmembers)
	d.Set("max_members", group.Maxmembers)
	d.Set("interval", group.Interval)
	d.Set("vm_profile_id", group.Vmprofileid)
	d.Set("display", group.Fordisplay)

	scaleUpPolicyIds := schema.NewSet(schema.HashString, []interface{}{})
	if group.Scaleuppolicies != nil {
		for _, policyId := range group.Scaleuppolicies {
			scaleUpPolicyIds.Add(policyId)
		}
	}
	d.Set("scaleup_policy_ids", scaleUpPolicyIds)

	scaleDownPolicyIds := schema.NewSet(schema.HashString, []interface{}{})
	if group.Scaledownpolicies != nil {
		for _, policyId := range group.Scaledownpolicies {
			scaleDownPolicyIds.Add(policyId)
		}
	}
	d.Set("scaledown_policy_ids", scaleDownPolicyIds)

	return nil
}

func resourceCloudStackAutoScaleVMGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if d.HasChange("name") || d.HasChange("min_members") || d.HasChange("max_members") ||
		d.HasChange("interval") || d.HasChange("scaleup_policy_ids") ||
		d.HasChange("scaledown_policy_ids") || d.HasChange("display") {

		log.Printf("[DEBUG] Updating autoscale VM group: %s", d.Id())

		p := cs.AutoScale.NewUpdateAutoScaleVmGroupParams(d.Id())

		if d.HasChange("name") {
			if v, ok := d.GetOk("name"); ok {
				p.SetName(v.(string))
			}
		}

		if d.HasChange("min_members") {
			minmembers := d.Get("min_members").(int)
			p.SetMinmembers(minmembers)
		}

		if d.HasChange("max_members") {
			maxmembers := d.Get("max_members").(int)
			p.SetMaxmembers(maxmembers)
		}

		if d.HasChange("interval") {
			if v, ok := d.GetOk("interval"); ok {
				p.SetInterval(v.(int))
			}
		}

		if d.HasChange("scaleup_policy_ids") {
			scaleUpPolicyIds := []string{}
			if v, ok := d.GetOk("scaleup_policy_ids"); ok {
				scaleUpSet := v.(*schema.Set)
				for _, id := range scaleUpSet.List() {
					scaleUpPolicyIds = append(scaleUpPolicyIds, id.(string))
				}
			}
			p.SetScaleuppolicyids(scaleUpPolicyIds)
		}

		if d.HasChange("scaledown_policy_ids") {
			scaleDownPolicyIds := []string{}
			if v, ok := d.GetOk("scaledown_policy_ids"); ok {
				scaleDownSet := v.(*schema.Set)
				for _, id := range scaleDownSet.List() {
					scaleDownPolicyIds = append(scaleDownPolicyIds, id.(string))
				}
			}
			p.SetScaledownpolicyids(scaleDownPolicyIds)
		}

		if d.HasChange("display") {
			if v, ok := d.GetOk("display"); ok {
				p.SetFordisplay(v.(bool))
			}
		}

		_, err := cs.AutoScale.UpdateAutoScaleVmGroup(p)
		if err != nil {
			return fmt.Errorf("Error updating autoscale VM group: %s", err)
		}

		log.Printf("[DEBUG] Autoscale VM group updated successfully: %s", d.Id())
	}

	return resourceCloudStackAutoScaleVMGroupRead(d, meta)
}

func resourceCloudStackAutoScaleVMGroupDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewDeleteAutoScaleVmGroupParams(d.Id())

	log.Printf("[DEBUG] Deleting autoscale VM group: %s", d.Id())
	_, err := cs.AutoScale.DeleteAutoScaleVmGroup(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting autoscale VM group: %s", err)
	}

	return nil
}
