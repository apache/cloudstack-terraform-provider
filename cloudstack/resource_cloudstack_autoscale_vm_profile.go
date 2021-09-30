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
	"net/url"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackAutoScaleVMProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackAutoScaleVMProfileCreate,
		Read:   resourceCloudStackAutoScaleVMProfileRead,
		Update: resourceCloudStackAutoScaleVMProfileUpdate,
		Delete: resourceCloudStackAutoScaleVMProfileDelete,

		Schema: map[string]*schema.Schema{
			"service_offering": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"template": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destroy_vm_grace_period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"other_deploy_params": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"metadata": metadataSchema(),
		},
	}
}

func resourceCloudStackAutoScaleVMProfileCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Retrieve the service_offering ID
	serviceofferingid, e := retrieveID(cs, "service_offering", d.Get("service_offering").(string))
	if e != nil {
		return e.Error()
	}

	// Retrieve the zone ID
	zoneid, e := retrieveID(cs, "zone", d.Get("zone").(string))
	if e != nil {
		return e.Error()
	}

	// Retrieve the template ID
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
		p.SetDestroyvmgraceperiod(int(duration.Seconds()))
	}

	if v, ok := d.GetOk("other_deploy_params"); ok {
		otherMap := v.(map[string]interface{})
		result := url.Values{}
		for k, v := range otherMap {
			result.Set(k, fmt.Sprint(v))
		}
		p.SetOtherdeployparams(result.Encode())
	}

	// Create the new vm profile
	r, err := cs.AutoScale.CreateAutoScaleVmProfile(p)
	if err != nil {
		return fmt.Errorf("Error creating AutoScaleVmProfile %s: %s", d.Id(), err)
	}

	d.SetId(r.Id)

	// Set metadata if necessary
	if err = setMetadata(cs, d, "AutoScaleVmProfile"); err != nil {
		return fmt.Errorf("Error setting metadata on the AutoScaleVmProfile %s: %s", d.Id(), err)
	}

	return nil
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

	d.Set("destroy_vm_grace_period", (time.Duration(p.Destroyvmgraceperiod) * time.Second).String())

	if p.Otherdeployparams != "" {
		var values url.Values
		values, err = url.ParseQuery(p.Otherdeployparams)
		if err != nil {
			return err
		}
		otherParams := make(map[string]interface{}, len(values))
		for key := range values {
			otherParams[key] = values.Get(key)
		}
		d.Set("other_deploy_params", otherParams)
	}

	metadata, err := getMetadata(cs, d, "AutoScaleVmProfile")
	if err != nil {
		return err
	}
	d.Set("metadata", metadata)

	return nil
}

func resourceCloudStackAutoScaleVMProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
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
		duration, err := time.ParseDuration(d.Get("destroy_vm_grace_period").(string))
		if err != nil {
			return err
		}
		p.SetDestroyvmgraceperiod(int(duration.Seconds()))
	}

	_, err := cs.AutoScale.UpdateAutoScaleVmProfile(p)
	if err != nil {
		return fmt.Errorf("Error updating AutoScaleVmProfile %s: %s", d.Id(), err)
	}

	if d.HasChange("metadata") {
		if err := updateMetadata(cs, d, "AutoScaleVmProfile"); err != nil {
			return fmt.Errorf("Error updating tags on AutoScaleVmProfile %s: %s", d.Id(), err)
		}
	}

	return resourceCloudStackAutoScaleVMProfileRead(d, meta)
}

func resourceCloudStackAutoScaleVMProfileDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.AutoScale.NewDeleteAutoScaleVmProfileParams(d.Id())

	// Delete the template
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
