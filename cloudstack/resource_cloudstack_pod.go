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

func resourceCloudStackPod() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackPodCreate,
		Read:   resourceCloudStackPodRead,
		Update: resourceCloudStackPodUpdate,
		Delete: resourceCloudStackPodDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway": {
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
				Required: true,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// VLAN ID is not directly settable in the CreatePodParams
			// It's returned in the response but can't be set during creation
			"vlan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackPodCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	zoneID := d.Get("zone_id").(string)
	gateway := d.Get("gateway").(string)
	netmask := d.Get("netmask").(string)
	startIP := d.Get("start_ip").(string)

	// Create a new parameter struct
	p := cs.Pod.NewCreatePodParams(name, zoneID)

	// Set required parameters
	p.SetGateway(gateway)
	p.SetNetmask(netmask)
	p.SetStartip(startIP)

	// Set optional parameters
	if endIP, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(endIP.(string))
	}

	// Note: VLAN ID is not directly settable in the CreatePodParams

	log.Printf("[DEBUG] Creating Pod %s", name)
	pod, err := cs.Pod.CreatePod(p)
	if err != nil {
		return fmt.Errorf("Error creating Pod %s: %s", name, err)
	}

	d.SetId(pod.Id)

	return resourceCloudStackPodRead(d, meta)
}

func resourceCloudStackPodRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the Pod details
	p := cs.Pod.NewListPodsParams()
	p.SetId(d.Id())

	pods, err := cs.Pod.ListPods(p)
	if err != nil {
		return fmt.Errorf("Error getting Pod %s: %s", d.Id(), err)
	}

	if pods.Count == 0 {
		log.Printf("[DEBUG] Pod %s does no longer exist", d.Id())
		d.SetId("")
		return nil
	}

	pod := pods.Pods[0]

	d.Set("name", pod.Name)
	d.Set("zone_id", pod.Zoneid)
	d.Set("zone_name", pod.Zonename)
	d.Set("gateway", pod.Gateway)
	d.Set("netmask", pod.Netmask)
	d.Set("allocation_state", pod.Allocationstate)

	if len(pod.Startip) > 0 {
		d.Set("start_ip", pod.Startip[0])
	}

	if len(pod.Endip) > 0 {
		d.Set("end_ip", pod.Endip[0])
	}

	if len(pod.Vlanid) > 0 {
		d.Set("vlan_id", pod.Vlanid[0])
	}

	return nil
}

func resourceCloudStackPodUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Pod.NewUpdatePodParams(d.Id())

	if d.HasChange("name") {
		p.SetName(d.Get("name").(string))
	}

	if d.HasChange("gateway") {
		p.SetGateway(d.Get("gateway").(string))
	}

	if d.HasChange("netmask") {
		p.SetNetmask(d.Get("netmask").(string))
	}

	if d.HasChange("start_ip") {
		p.SetStartip(d.Get("start_ip").(string))
	}

	if d.HasChange("end_ip") {
		p.SetEndip(d.Get("end_ip").(string))
	}

	_, err := cs.Pod.UpdatePod(p)
	if err != nil {
		return fmt.Errorf("Error updating Pod %s: %s", d.Get("name").(string), err)
	}

	return resourceCloudStackPodRead(d, meta)
}

func resourceCloudStackPodDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Pod.NewDeletePodParams(d.Id())

	log.Printf("[DEBUG] Deleting Pod %s", d.Get("name").(string))
	_, err := cs.Pod.DeletePod(p)

	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting Pod %s: %s", d.Get("name").(string), err)
	}

	return nil
}
