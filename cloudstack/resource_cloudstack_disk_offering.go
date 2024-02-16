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
	"log"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackDiskOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackDiskOfferingCreate,
		Read:   resourceCloudStackDiskOfferingRead,
		Update: resourceCloudStackDiskOfferingUpdate,
		Delete: resourceCloudStackDiskOfferingDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"encrypt": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackDiskOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	display_text := d.Get("display_text").(string)
	disk_size := d.Get("disk_size").(int)

	// Create a new parameter struct
	p := cs.DiskOffering.NewCreateDiskOfferingParams(name, display_text)
	p.SetDisksize(int64(disk_size))

	p.SetDisksize(int64(disk_size))
	if encrypt, ok := d.GetOk("encrypt"); ok {
		p.SetEncrypt(encrypt.(bool))
	}
	if tags, ok := d.GetOk("tags"); ok {
		p.SetTags(tags.(string))
	}

	log.Printf("[DEBUG] Creating Disk Offering %s", name)
	diskOff, err := cs.DiskOffering.CreateDiskOffering(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Disk Offering %s successfully created", name)
	d.SetId(diskOff.Id)

	return resourceCloudStackDiskOfferingRead(d, meta)
}

func resourceCloudStackDiskOfferingRead(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceCloudStackDiskOfferingUpdate(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceCloudStackDiskOfferingDelete(d *schema.ResourceData, meta interface{}) error { 
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.DiskOffering.NewDeleteDiskOfferingParams(d.Id())
	_, err := cs.DiskOffering.DeleteDiskOffering(p)
	if err != nil {
		return fmt.Errorf("Error deleting disk offering: %s", err)
	}
	
	
	return nil

}
