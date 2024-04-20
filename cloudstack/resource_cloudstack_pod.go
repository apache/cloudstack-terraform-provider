// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackPod() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackPodCreate,
		Read:   resourceCloudStackPodRead,
		Update: resourceCloudStackPodUpdate,
		Delete: resourceCloudStackPodDelete,
		Schema: map[string]*schema.Schema{
			"gateway": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"netmask": {
				Type:     schema.TypeString,
				Required: true,
			},

			"start_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Enabled",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validOptions := []string{
						"Enabled",
						"Disabled",
					}
					err := validateOptions(validOptions, v.(string), k)
					if err != nil {
						errors = append(errors, err)
					}
					return
				},
			},
			"has_annotations": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackPodCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	gateway := d.Get("gateway").(string)
	name := d.Get("name").(string)
	netmask := d.Get("netmask").(string)
	startIP := d.Get("start_ip").(string)
	zoneID := d.Get("zone_id").(string)

	// Create a new parameter struct
	p := cs.Pod.NewCreatePodParams(
		gateway,
		name,
		netmask,
		startIP,
		zoneID,
	)

	// If there is a end_ip supplied, add it to the parameter struct
	if endIP, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(endIP.(string))
	}

	// If there is a end_ip supplied, add it to the parameter struct
	if allocationState, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(allocationState.(string))
	}

	log.Printf("[DEBUG] Creating Pod %s", name)

	n, err := cs.Pod.CreatePod(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Pod %s successfully created", name)

	d.SetId(n.Id)

	return resourceCloudStackPodRead(d, meta)
}

func resourceCloudStackPodRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Pod %s", d.Get("name").(string))

	// Get the Pod details
	p, count, err := cs.Pod.GetPodByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Pod %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(p.Id)
	d.Set("allocation_state", p.Allocationstate)
	d.Set("end_ip", p.Endip)
	d.Set("gateway", p.Gateway)
	d.Set("has_annotations", p.Hasannotations)
	d.Set("netmask", p.Netmask)
	d.Set("start_ip", p.Startip)
	d.Set("zone_id", p.Zoneid)
	d.Set("zone_name", p.Zonename)
	return nil
}

func resourceCloudStackPodUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	if d.HasChange("gateway") || d.HasChange("netmask") || d.HasChange("start_ip") || d.HasChange("end_ip") {

		p := cs.Pod.NewUpdatePodParams(d.Id())

		if d.HasChange("gateway") {
			log.Printf("[DEBUG] Gateway changed for %s, starting update", name)
			p.SetGateway(d.Get("gateway").(string))
		}
		if d.HasChange("netmask") {
			log.Printf("[DEBUG] NetMask changed for %s, starting update", name)
			p.SetNetmask(d.Get("netmask").(string))
		}
		if d.HasChange("start_ip") {
			log.Printf("[DEBUG] StartIP changed for %s, starting update", name)
			p.SetStartip(d.Get("start_ip").(string))
		}
		if d.HasChange("end_ip") {
			log.Printf("[DEBUG] endIP changed for %s, starting update", name)
			p.SetEndip(d.Get("end_ip").(string))
		}

		_, err := cs.Pod.UpdatePod(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the pod %s: %s", name, err)
		}
	}
	if d.HasChange("name") {
		log.Printf("[DEBUG] Name changed for %s, starting update", name)

		p := cs.Pod.NewUpdatePodParams(d.Id())

		p.SetName(d.Get("name").(string))

		_, err := cs.Pod.UpdatePod(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the name for pod %s: %s", name, err)
		}

	}

	if d.HasChange("allocation_state") {
		log.Printf("[DEBUG] allocationState changed for %s, starting update", name)

		p := cs.Pod.NewUpdatePodParams(d.Id())

		p.SetAllocationstate(d.Get("allocation_state").(string))

		_, err := cs.Pod.UpdatePod(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the allocation_state for pod %s: %s", name, err)
		}
	}

	return resourceCloudStackPodRead(d, meta)
}

func resourceCloudStackPodDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Pod.NewDeletePodParams(d.Id())
	_, err := cs.Pod.DeletePod(p)

	if err != nil {
		return fmt.Errorf("Error deleting Pod: %s", err)
	}

	return nil
}
