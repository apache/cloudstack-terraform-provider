//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for Attachitional information
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
	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackAttachVolume() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCloudStackAttachVolumeRead,
		Create: resourceCloudStackAttachVolumeCreate,
		Delete: resourceCloudStackAttachVolumeDelete,
		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the ID of the disk volume",
			},
			"virtual_machine_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the ID of the virtual machine",
			},
			"device_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the device to map the volume to the guest OS. ",
			},
			"attached": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				Description: "the date the volume was attached to a VM instance",
			},
		},
	}
}

func resourceCloudStackAttachVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Volume.NewAttachVolumeParams(d.Get("volume_id").(string), d.Get("virtual_machine_id").(string))
	if v, ok := d.GetOk("device_id"); ok {
		p.SetDeviceid(v.(int64))
	}

	r, err := cs.Volume.AttachVolume(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudStackAttachVolumeRead(d, meta)
}

func resourceCloudStackAttachVolumeRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.Volume.GetVolumeByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("volume_id", r.Id)
	d.Set("virtual_machine_id", r.Virtualmachineid)
	d.Set("device_id", r.Deviceid)
	d.Set("attached", r.Attached)

	return nil
}

func resourceCloudStackAttachVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Volume.NewDetachVolumeParams()
	p.SetId(d.Id())
	_, err := cs.Volume.DetachVolume(p)
	if err != nil {
		return err
	}

	return nil
}
