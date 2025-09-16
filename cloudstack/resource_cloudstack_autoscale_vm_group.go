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
	"time"

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

			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enable",
				Description: "the state of the autoscale vm group (enable or disable)",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "enable" && value != "disable" {
						errors = append(errors, fmt.Errorf("state must be either 'enable' or 'disable'"))
					}
					return
				},
			},

			"cleanup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "true if all the members of autoscale vm group has to be cleaned up, false otherwise",
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

	requestedState := d.Get("state").(string)
	if requestedState == "disable" {
		log.Printf("[DEBUG] Disabling autoscale VM group as requested: %s", resp.Id)
		disableParams := cs.AutoScale.NewDisableAutoScaleVmGroupParams(resp.Id)
		_, err = cs.AutoScale.DisableAutoScaleVmGroup(disableParams)
		if err != nil {
			return fmt.Errorf("Error disabling autoscale VM group after creation: %s", err)
		}
	}

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

	d.Set("state", requestedState)

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

	terraformState := "enable"
	if strings.ToLower(group.State) == "disabled" {
		terraformState = "disable"
	}

	currentConfigState := d.Get("state").(string)
	log.Printf("[DEBUG] CloudStack API state: %s, mapped to Terraform state: %s", group.State, terraformState)
	log.Printf("[DEBUG] Current config state: %s, will set state to: %s", currentConfigState, terraformState)

	d.Set("state", terraformState)

	scaleUpPolicyIds := schema.NewSet(schema.HashString, []interface{}{})
	if group.Scaleuppolicies != nil {
		for _, policy := range group.Scaleuppolicies {
			// Extract the ID from the AutoScalePolicy object
			if policy != nil {
				scaleUpPolicyIds.Add(policy.Id)
			}
		}
	}
	d.Set("scaleup_policy_ids", scaleUpPolicyIds)

	scaleDownPolicyIds := schema.NewSet(schema.HashString, []interface{}{})
	if group.Scaledownpolicies != nil {
		for _, policy := range group.Scaledownpolicies {
			// Extract the ID from the AutoScalePolicy object
			if policy != nil {
				scaleDownPolicyIds.Add(policy.Id)
			}
		}
	}
	d.Set("scaledown_policy_ids", scaleDownPolicyIds)

	return nil
}

func resourceCloudStackAutoScaleVMGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	oldState, newState := d.GetChange("state")
	log.Printf("[DEBUG] State in terraform config: %s", d.Get("state").(string))
	log.Printf("[DEBUG] Old state: %s, New state: %s", oldState.(string), newState.(string))

	changes := []string{}
	if d.HasChange("name") {
		changes = append(changes, "name")
	}
	if d.HasChange("min_members") {
		changes = append(changes, "min_members")
	}
	if d.HasChange("max_members") {
		changes = append(changes, "max_members")
	}
	if d.HasChange("interval") {
		changes = append(changes, "interval")
	}
	if d.HasChange("scaleup_policy_ids") {
		changes = append(changes, "scaleup_policy_ids")
	}
	if d.HasChange("scaledown_policy_ids") {
		changes = append(changes, "scaledown_policy_ids")
	}
	if d.HasChange("display") {
		changes = append(changes, "display")
	}
	if d.HasChange("state") {
		changes = append(changes, "state")
	}
	log.Printf("[DEBUG] Detected changes in autoscale VM group: %v", changes)

	if d.HasChange("name") || d.HasChange("min_members") || d.HasChange("max_members") ||
		d.HasChange("interval") || d.HasChange("scaleup_policy_ids") ||
		d.HasChange("scaledown_policy_ids") || d.HasChange("display") || d.HasChange("state") {

		log.Printf("[DEBUG] Updating autoscale VM group: %s", d.Id())

		// Check current state to determine operation order
		currentState := "enable"
		if oldState, newState := d.GetChange("state"); d.HasChange("state") {
			currentState = oldState.(string)
			log.Printf("[DEBUG] State change detected: %s -> %s", currentState, newState.(string))
		}

		if d.HasChange("state") {
			newState := d.Get("state").(string)
			if newState == "disable" && currentState == "enable" {
				log.Printf("[DEBUG] Disabling autoscale VM group before other updates: %s", d.Id())
				disableParams := cs.AutoScale.NewDisableAutoScaleVmGroupParams(d.Id())
				_, err := cs.AutoScale.DisableAutoScaleVmGroup(disableParams)
				if err != nil {
					return fmt.Errorf("Error disabling autoscale VM group: %s", err)
				}
				// Wait a moment for disable to take effect
				time.Sleep(1 * time.Second)
			}
		}

		if d.HasChange("name") || d.HasChange("min_members") || d.HasChange("max_members") ||
			d.HasChange("interval") || d.HasChange("scaleup_policy_ids") ||
			d.HasChange("scaledown_policy_ids") || d.HasChange("display") {

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

			log.Printf("[DEBUG] Applying parameter updates to autoscale VM group: %s", d.Id())
			_, err := cs.AutoScale.UpdateAutoScaleVmGroup(p)
			if err != nil {
				return fmt.Errorf("Error updating autoscale VM group parameters: %s", err)
			}
		}

		if d.HasChange("state") {
			newState := d.Get("state").(string)
			if newState == "enable" {
				log.Printf("[DEBUG] Enabling autoscale VM group after updates: %s", d.Id())
				enableParams := cs.AutoScale.NewEnableAutoScaleVmGroupParams(d.Id())
				_, err := cs.AutoScale.EnableAutoScaleVmGroup(enableParams)
				if err != nil {
					return fmt.Errorf("Error enabling autoscale VM group: %s", err)
				}
			}
		}

		log.Printf("[DEBUG] Autoscale VM group updated successfully: %s", d.Id())
	}

	return resourceCloudStackAutoScaleVMGroupRead(d, meta)
}

func resourceCloudStackAutoScaleVMGroupDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewDeleteAutoScaleVmGroupParams(d.Id())

	cleanup := d.Get("cleanup").(bool)
	p.SetCleanup(cleanup)
	log.Printf("[DEBUG] Deleting autoscale VM group with cleanup=%t", cleanup)
	log.Printf("[DEBUG] This deletion was triggered by Terraform - check if this should be an update instead")

	log.Printf("[DEBUG] Disabling autoscale VM group before deletion: %s", d.Id())
	disableParams := cs.AutoScale.NewDisableAutoScaleVmGroupParams(d.Id())
	_, err := cs.AutoScale.DisableAutoScaleVmGroup(disableParams)
	if err != nil {
		if !strings.Contains(err.Error(), "Invalid parameter id value") &&
			!strings.Contains(err.Error(), "entity does not exist") &&
			!strings.Contains(err.Error(), "already disabled") {
			return fmt.Errorf("Error disabling autoscale VM group: %s", err)
		}
	}

	time.Sleep(2 * time.Second)
	log.Printf("[DEBUG] Autoscale VM group disabled, proceeding with deletion: %s", d.Id())

	log.Printf("[DEBUG] Deleting autoscale VM group: %s", d.Id())
	_, err = cs.AutoScale.DeleteAutoScaleVmGroup(p)
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
