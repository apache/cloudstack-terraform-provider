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
			"allocation_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudStackPodCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Pod.NewCreatePodParams(d.Get("name").(string), d.Get("zone_id").(string))
	if v, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(v.(string))
	}
	if v, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(v.(string))
	}
	if v, ok := d.GetOk("gateway"); ok {
		p.SetGateway(v.(string))
	}
	if v, ok := d.GetOk("netmask"); ok {
		p.SetNetmask(v.(string))
	}
	if v, ok := d.GetOk("start_ip"); ok {
		p.SetStartip(v.(string))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		p.SetZoneid(v.(string))
	}

	r, err := cs.Pod.CreatePod(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudStackPodRead(d, meta)
}

func resourceCloudStackPodRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.Pod.GetPodByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("allocation_state", r.Allocationstate)
	d.Set("end_ip", strings.Join(r.Endip, " "))
	d.Set("gateway", r.Gateway)
	d.Set("name", r.Name)
	d.Set("netmask", r.Netmask)
	d.Set("start_ip", strings.Join(r.Startip, " "))
	d.Set("zone_id", r.Zoneid)

	return nil
}

func resourceCloudStackPodUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Pod.NewUpdatePodParams(d.Id())
	if v, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(v.(string))
	}
	if v, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(v.(string))
	}
	if v, ok := d.GetOk("gateway"); ok {
		p.SetGateway(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}
	if v, ok := d.GetOk("netmask"); ok {
		p.SetNetmask(v.(string))
	}
	if v, ok := d.GetOk("start_ip"); ok {
		p.SetStartip(v.(string))
	}

	_, err := cs.Pod.UpdatePod(p)
	if err != nil {
		return err
	}

	return resourceCloudStackPodRead(d, meta)
}

func resourceCloudStackPodDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.Pod.DeletePod(cs.Pod.NewDeletePodParams(d.Id()))
	if err != nil {
		return err
	}

	return nil
}
