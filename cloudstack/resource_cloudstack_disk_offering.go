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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackDiskOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackDiskOfferingCreate,
		Read:   resourceCloudStackDiskOfferingRead,
		Update: resourceCloudStackDiskOfferingUpdate,
		Delete: resourceCloudStackDiskOfferingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the disk offering",
				Type:        schema.TypeString,
				Required:    true,
			},
			"display_text": {
				Description: "The display text of the disk offering",
				Type:        schema.TypeString,
				Required:    true,
			},
			"disk_size": {
				Description:   "The size of the disk offering in GB",
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"customized"},
			},
			"customized": {
				Description:   "Whether the disk offering allows a custom disk size at deployment time",
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Default:       false,
				ConflictsWith: []string{"disk_size"},
			},
			"storage_type": {
				Description: "The storage type of the disk offering. Values are local and shared",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "shared",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "local" || v == "shared" {
						return
					}
					errs = append(errs, fmt.Errorf("storage type should be either local or shared, got %s", v))
					return
				},
			},
			"provisioning_type": {
				Description: "Provisioning type used to create volumes. Values are thin, sparse and fat",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "thin",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "thin" || v == "sparse" || v == "fat" {
						return
					}
					errs = append(errs, fmt.Errorf("provisioning type should be one of thin, sparse or fat, got %s", v))
					return
				},
			},
			"tags": {
				Description: "The storage tags for the disk offering",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"display_offering": {
				Description: "Whether the disk offering is displayed to the end user",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceCloudStackDiskOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	displayText := d.Get("display_text").(string)

	// Create a new parameter struct
	p := cs.DiskOffering.NewCreateDiskOfferingParams(displayText, name)

	if v, ok := d.GetOk("disk_size"); ok {
		p.SetDisksize(int64(v.(int)))
	}

	customized := false
	if v, ok := d.GetOk("customized"); ok {
		customized = v.(bool)
	}
	if _, ok := d.GetOk("disk_size"); !ok {
		customized = true
	}
	p.SetCustomized(customized)

	if v, ok := d.GetOk("storage_type"); ok {
		p.SetStoragetype(v.(string))
	}

	if v, ok := d.GetOk("provisioning_type"); ok {
		p.SetProvisioningtype(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(v.(string))
	}

	// display_offering defaults to true, so read it directly rather than via GetOk
	p.SetDisplayoffering(d.Get("display_offering").(bool))

	log.Printf("[DEBUG] Creating Disk Offering %s", name)
	diskOff, err := cs.DiskOffering.CreateDiskOffering(p)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Disk Offering %s successfully created", name)
	d.SetId(diskOff.Id)

	return resourceCloudStackDiskOfferingRead(d, meta)
}

func resourceCloudStackDiskOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Retrieving Disk Offering %s", d.Id())

	// Get the Disk Offering details
	diskOff, count, err := cs.DiskOffering.GetDiskOfferingByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Disk Offering %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", diskOff.Name)
	d.Set("display_text", diskOff.Displaytext)
	d.Set("disk_size", int(diskOff.Disksize))
	d.Set("customized", diskOff.Iscustomized)
	d.Set("storage_type", diskOff.Storagetype)
	d.Set("provisioning_type", diskOff.Provisioningtype)
	d.Set("tags", diskOff.Tags)
	d.Set("display_offering", diskOff.Displayoffering)

	return nil
}

func resourceCloudStackDiskOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)

	if d.HasChange("name") || d.HasChange("display_text") ||
		d.HasChange("tags") || d.HasChange("display_offering") {

		// Create a new parameter struct
		p := cs.DiskOffering.NewUpdateDiskOfferingParams(d.Id())

		p.SetName(d.Get("name").(string))
		p.SetDisplaytext(d.Get("display_text").(string))
		p.SetTags(d.Get("tags").(string))
		p.SetDisplayoffering(d.Get("display_offering").(bool))

		log.Printf("[DEBUG] Updating Disk Offering %s", name)
		_, err := cs.DiskOffering.UpdateDiskOffering(p)
		if err != nil {
			return fmt.Errorf("Error updating Disk Offering %s: %s", name, err)
		}
	}

	return resourceCloudStackDiskOfferingRead(d, meta)
}

func resourceCloudStackDiskOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.DiskOffering.NewDeleteDiskOfferingParams(d.Id())

	log.Printf("[DEBUG] Deleting Disk Offering %s", d.Get("name").(string))
	_, err := cs.DiskOffering.DeleteDiskOffering(p)
	if err != nil {
		return fmt.Errorf("Error deleting Disk Offering: %s", err)
	}

	return nil
}
