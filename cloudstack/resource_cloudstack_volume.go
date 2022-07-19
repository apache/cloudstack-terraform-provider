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
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackVolumeCreate,
		Read:   resourceCloudStackVolumeRead,
		Delete: resourceCloudStackVolumeDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_offering_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	disk_offering_id := d.Get("disk_offering_id").(string)
	zone_id := d.Get("zone_id").(string)

	//Create a new parameter struct
	p := cs.Volume.NewCreateVolumeParams()
	p.SetDiskofferingid(disk_offering_id)
	p.SetZoneid(zone_id)
	p.SetName(name)

	log.Printf("[DEBUG] Creating Volume %s", name)
	v, err := cs.Volume.CreateVolume(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Volume %s successfully created", name)
	d.SetId(v.Id)

	return resourceCloudStackVolumeRead(d, meta)
}
func resourceCloudStackVolumeRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Volume %s", d.Get("name").(string))

	// Get the Volume details
	v, count, err := cs.Volume.GetVolumeByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Volume %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(v.Id)
	d.Set("name", v.Name)
	d.Set("disk_offering_id", v.Diskofferingid)
	d.Set("zone_id", v.Zoneid)

	return nil
}

func resourceCloudStackVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Volume.NewDeleteVolumeParams(d.Id())
	_, err := cs.Volume.DeleteVolume(p)

	if err != nil {
		return fmt.Errorf("Error deleting Volume: %s", err)
	}

	return nil
}
