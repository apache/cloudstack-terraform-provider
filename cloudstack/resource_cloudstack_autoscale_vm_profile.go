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

func resourceCloudStackAutoScaleVMProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackAutoScaleVMProfileCreate,
		Read:   resourceCloudStackAutoScaleVMProfileRead,
		Update: resourceCloudStackAutoScaleVMProfileUpdate,
		Delete: resourceCloudStackAutoScaleVMProfileDelete,

		Schema: map[string]*schema.Schema{
			"service_offering": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the service offering of the auto deployed virtual machine",
			},

			"template": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the template of the auto deployed virtual machine",
			},

			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "availability zone for the auto deployed virtual machine",
			},

			"destroy_vm_grace_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "the time allowed for existing connections to get closed before a vm is expunged",
			},

			"other_deploy_params": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "parameters other than zoneId/serviceOfferringId/templateId of the auto deployed virtual machine",
			},

			"counter_param_list": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "counterparam list. Example: counterparam[0].name=snmpcommunity&counterparam[0].value=public&counterparam[1].name=snmpport&counterparam[1].value=161",
			},

			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "an optional binary data that can be sent to the virtual machine upon a successful deployment. This binary data must be base64 encoded before adding it to the request.",
			},

			"user_data_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the ID of the Userdata",
			},

			"user_data_details": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "used to specify the parameters values for the variables in userdata",
			},

			"autoscale_user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the ID of the user used to launch and destroy the VMs",
			},

			"display": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "an optional field, whether to the display the profile to the end user or not",
			},

			"account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "account that will own the autoscale VM profile",
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "an optional project for the autoscale VM profile",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "domain ID of the account owning a autoscale VM profile",
			},

			"metadata": metadataSchema(),
		},
	}
}

func resourceCloudStackAutoScaleVMProfileCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	serviceofferingid, e := retrieveID(cs, "service_offering", d.Get("service_offering").(string))
	if e != nil {
		return e.Error()
	}

	zoneid, e := retrieveID(cs, "zone", d.Get("zone").(string))
	if e != nil {
		return e.Error()
	}

	templateid, e := retrieveTemplateID(cs, zoneid, d.Get("template").(string))
	if e != nil {
		return e.Error()
	}

	p := cs.AutoScale.NewCreateAutoScaleVmProfileParams(serviceofferingid, templateid, zoneid)

	if v, ok := d.GetOk("destroy_vm_grace_period"); ok {
		duration, err := time.ParseDuration(v.(string))
		if err != nil {
			return err
		}
		p.SetExpungevmgraceperiod(int(duration.Seconds()))
	}

	if v, ok := d.GetOk("other_deploy_params"); ok {
		nv := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			nv[k] = v.(string)
		}
		p.SetOtherdeployparams(nv)
	}

	if v, ok := d.GetOk("counter_param_list"); ok {
		nv := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			nv[k] = v.(string)
		}
		p.SetCounterparam(nv)
	}

	if v, ok := d.GetOk("user_data"); ok {
		p.SetUserdata(v.(string))
	}

	if v, ok := d.GetOk("user_data_id"); ok {
		p.SetUserdataid(v.(string))
	}

	if v, ok := d.GetOk("user_data_details"); ok {
		nv := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			nv[k] = v.(string)
		}
		p.SetUserdatadetails(nv)
	}

	if v, ok := d.GetOk("autoscale_user_id"); ok {
		p.SetAutoscaleuserid(v.(string))
	}

	if v, ok := d.GetOk("display"); ok {
		p.SetFordisplay(v.(bool))
	}

	if v, ok := d.GetOk("account_name"); ok {
		p.SetAccount(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(v.(string))
	}

	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}

	r, err := cs.AutoScale.CreateAutoScaleVmProfile(p)
	if err != nil {
		return fmt.Errorf("Error creating AutoScaleVmProfile %s: %s", d.Id(), err)
	}

	d.SetId(r.Id)

	if err = setMetadata(cs, d, "AutoScaleVmProfile"); err != nil {
		return fmt.Errorf("Error setting metadata on the AutoScaleVmProfile %s: %s", d.Id(), err)
	}

	return resourceCloudStackAutoScaleVMProfileRead(d, meta)
}

func resourceCloudStackAutoScaleVMProfileRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p, count, err := cs.AutoScale.GetAutoScaleVmProfileByID(d.Id())

	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] AutoScaleVmProfile %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	zone, _, err := cs.Zone.GetZoneByID(p.Zoneid)
	if err != nil {
		return err
	}

	offering, _, err := cs.ServiceOffering.GetServiceOfferingByID(p.Serviceofferingid)
	if err != nil {
		return err
	}

	template, _, err := cs.Template.GetTemplateByID(p.Templateid, "executable", cloudstack.WithZone(p.Zoneid))
	if err != nil {
		return err
	}

	setValueOrID(d, "service_offering", offering.Name, p.Serviceofferingid)
	setValueOrID(d, "template", template.Name, p.Templateid)
	setValueOrID(d, "zone", zone.Name, p.Zoneid)

	d.Set("destroy_vm_grace_period", (time.Duration(p.Expungevmgraceperiod) * time.Second).String())

	if p.Otherdeployparams != nil {
		d.Set("other_deploy_params", p.Otherdeployparams)
	}

	if p.Userdata != "" {
		d.Set("user_data", p.Userdata)
	}

	if p.Userdataid != "" {
		d.Set("user_data_id", p.Userdataid)
	}

	if p.Userdatadetails != "" {
		if _, ok := d.GetOk("user_data_details"); !ok {
			d.Set("user_data_details", map[string]interface{}{})
		}
	}

	if p.Autoscaleuserid != "" {
		d.Set("autoscale_user_id", p.Autoscaleuserid)
	}

	d.Set("display", p.Fordisplay)

	if p.Account != "" {
		d.Set("account_name", p.Account)
	}

	if p.Projectid != "" {
		d.Set("project_id", p.Projectid)
	}

	if p.Domainid != "" {
		d.Set("domain_id", p.Domainid)
	}

	metadata, err := getMetadata(cs, d, "AutoScaleVmProfile")
	if err != nil {
		return err
	}
	d.Set("metadata", metadata)

	return nil
}

// waitForVMGroupsState waits for the specified VM groups to reach the desired state
func waitForVMGroupsState(cs *cloudstack.CloudStackClient, groupIDs []string, desiredState string) error {
	maxRetries := 30 // 30 * 2 seconds = 60 seconds max wait
	for i := 0; i < maxRetries; i++ {
		allInDesiredState := true

		for _, groupID := range groupIDs {
			group, _, err := cs.AutoScale.GetAutoScaleVmGroupByID(groupID)
			if err != nil {
				return fmt.Errorf("Error checking state of VM group %s: %s", groupID, err)
			}

			groupInDesiredState := false
			if desiredState == "disabled" {
				groupInDesiredState = (group.State == "disabled")
			} else if desiredState == "enabled" {
				groupInDesiredState = (group.State == "enabled")
			} else {
				groupInDesiredState = (group.State == desiredState)
			}

			if !groupInDesiredState {
				allInDesiredState = false
				log.Printf("[DEBUG] VM group %s is in state '%s', waiting for '%s'", groupID, group.State, desiredState)
				break
			}
		}

		if allInDesiredState {
			log.Printf("[INFO] All VM groups have reached desired state: %s", desiredState)
			return nil
		}

		if i < maxRetries-1 {
			log.Printf("[INFO] Waiting for VM groups to reach state '%s' (attempt %d/%d)", desiredState, i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
	}

	return fmt.Errorf("Timeout waiting for VM groups to reach state '%s' after %d seconds", desiredState, maxRetries*2)
}

func waitForVMGroupsToBeDisabled(cs *cloudstack.CloudStackClient, profileID string) error {
	log.Printf("[DEBUG] Waiting for VM groups using profile %s to be disabled", profileID)
	listParams := cs.AutoScale.NewListAutoScaleVmGroupsParams()
	listParams.SetVmprofileid(profileID)

	groups, err := cs.AutoScale.ListAutoScaleVmGroups(listParams)
	if err != nil {
		log.Printf("[ERROR] Failed to list VM groups for profile %s: %s", profileID, err)
		return fmt.Errorf("Error listing autoscale VM groups: %s", err)
	}

	log.Printf("[DEBUG] Found %d VM groups using profile %s", len(groups.AutoScaleVmGroups), profileID)

	var groupIDs []string
	for _, group := range groups.AutoScaleVmGroups {
		log.Printf("[DEBUG] VM group %s (%s) current state: %s", group.Name, group.Id, group.State)
		groupIDs = append(groupIDs, group.Id)
	}

	if len(groupIDs) == 0 {
		log.Printf("[DEBUG] No VM groups found using profile %s", profileID)
		return nil
	}

	log.Printf("[INFO] Waiting for %d VM groups to be disabled for profile update", len(groupIDs))
	if err := waitForVMGroupsState(cs, groupIDs, "disabled"); err != nil {
		return fmt.Errorf("Autoscale VM groups must be disabled before updating profile: %s", err)
	}

	log.Printf("[DEBUG] All VM groups are now disabled for profile %s", profileID)
	return nil
}

func resourceCloudStackAutoScaleVMProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Profile update requested for ID: %s", d.Id())
	for _, key := range []string{"template", "destroy_vm_grace_period", "counter_param_list", "user_data", "user_data_id", "user_data_details", "autoscale_user_id", "display", "metadata"} {
		if d.HasChange(key) {
			old, new := d.GetChange(key)
			log.Printf("[DEBUG] Field '%s' changed from %v to %v", key, old, new)
		}
	}

	// Check if we only have metadata changes (which don't require CloudStack API update)
	onlyMetadataChanges := d.HasChange("metadata") &&
		!d.HasChange("template") &&
		!d.HasChange("destroy_vm_grace_period") &&
		!d.HasChange("counter_param_list") &&
		!d.HasChange("user_data") &&
		!d.HasChange("user_data_id") &&
		!d.HasChange("user_data_details") &&
		!d.HasChange("autoscale_user_id") &&
		!d.HasChange("display")

	if !onlyMetadataChanges {
		if err := waitForVMGroupsToBeDisabled(cs, d.Id()); err != nil {
			return fmt.Errorf("Autoscale VM groups must be disabled before updating profile: %s", err)
		}

		p := cs.AutoScale.NewUpdateAutoScaleVmProfileParams(d.Id())

		if d.HasChange("template") {
			zoneid, e := retrieveID(cs, "zone", d.Get("zone").(string))
			if e != nil {
				return e.Error()
			}
			templateid, e := retrieveTemplateID(cs, zoneid, d.Get("template").(string))
			if e != nil {
				return e.Error()
			}
			p.SetTemplateid(templateid)
		}

		if d.HasChange("destroy_vm_grace_period") {
			if v, ok := d.GetOk("destroy_vm_grace_period"); ok {
				duration, err := time.ParseDuration(v.(string))
				if err != nil {
					return err
				}
				p.SetExpungevmgraceperiod(int(duration.Seconds()))
			}
		}

		if d.HasChange("counter_param_list") {
			if v, ok := d.GetOk("counter_param_list"); ok {
				nv := make(map[string]string)
				for k, v := range v.(map[string]interface{}) {
					nv[k] = v.(string)
				}
				p.SetCounterparam(nv)
			}
		}

		if d.HasChange("user_data") {
			if v, ok := d.GetOk("user_data"); ok {
				p.SetUserdata(v.(string))
			}
		}

		if d.HasChange("user_data_id") {
			if v, ok := d.GetOk("user_data_id"); ok {
				p.SetUserdataid(v.(string))
			}
		}

		if d.HasChange("user_data_details") {
			if v, ok := d.GetOk("user_data_details"); ok {
				nv := make(map[string]string)
				for k, v := range v.(map[string]interface{}) {
					nv[k] = v.(string)
				}
				p.SetUserdatadetails(nv)
			}
		}

		if d.HasChange("autoscale_user_id") {
			if v, ok := d.GetOk("autoscale_user_id"); ok {
				p.SetAutoscaleuserid(v.(string))
			}
		}

		if d.HasChange("display") {
			if v, ok := d.GetOk("display"); ok {
				p.SetFordisplay(v.(bool))
			}
		}

		log.Printf("[DEBUG] Performing CloudStack API update for profile %s", d.Id())
		_, updateErr := cs.AutoScale.UpdateAutoScaleVmProfile(p)
		if updateErr != nil {
			return fmt.Errorf("Error updating AutoScaleVmProfile %s: %s", d.Id(), updateErr)
		}
	}

	if d.HasChange("metadata") {
		if metadataErr := updateMetadata(cs, d, "AutoScaleVmProfile"); metadataErr != nil {
			return fmt.Errorf("Error updating tags on AutoScaleVmProfile %s: %s", d.Id(), metadataErr)
		}
	}

	return resourceCloudStackAutoScaleVMProfileRead(d, meta)
}

func resourceCloudStackAutoScaleVMProfileDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewDeleteAutoScaleVmProfileParams(d.Id())

	log.Printf("[INFO] Deleting AutoScaleVmProfile: %s", d.Id())
	_, err := cs.AutoScale.DeleteAutoScaleVmProfile(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting AutoScaleVmProfile %s: %s", d.Id(), err)
	}
	return nil
}
