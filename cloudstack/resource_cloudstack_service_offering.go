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

func resourceCloudStackServiceOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackServiceOfferingCreate,
		Read:   resourceCloudStackServiceOfferingRead,
		Update: resourceCloudStackServiceOfferingUpdate,
		Delete: resourceCloudStackServiceOfferingDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudStackServiceOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	display_name := d.Get("display_name").(string)

	// Create a new parameter struct
	p := cs.ServiceOffering.NewCreateServiceOfferingParams(display_name, name)

	log.Printf("[DEBUG] Creating Service Offering %s", name)
	s, err := cs.ServiceOffering.CreateServiceOffering(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Service Offering %s successfully created", name)
	d.SetId(s.Id)

	return resourceCloudStackServiceOfferingRead(d, meta)
}

func resourceCloudStackServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Service Offering %s", d.Get("name").(string))

	// Get the Service Offering details
	s, count, err := cs.ServiceOffering.GetServiceOfferingByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Service Offering %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(s.Id)
	d.Set("name", s.Name)
	d.Set("display_name", s.Displaytext)

	return nil
}

func resourceCloudStackServiceOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	d.Partial(true)

	name := d.Get("name").(string)

	// Check if the name is changed and if so, update the service offering
	if d.HasChange("name") {
		log.Printf("[DEBUG] Name changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())

		// Set the new name
		p.SetName(d.Get("name").(string))

		// Update the name
		_, err := cs.ServiceOffering.UpdateServiceOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the name for service offering %s: %s", name, err)
		}

		d.SetPartial("name")
	}

	// Check if the display text is changed and if so, update seervice offering
	if d.HasChange("display_text") {
		log.Printf("[DEBUG] Display text changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())

		// Set the new display text
		p.SetName(d.Get("display_text").(string))

		// Update the display text
		_, err := cs.ServiceOffering.UpdateServiceOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the display text for service offering %s: %s", name, err)
		}

		d.SetPartial("display_text")
	}

	d.Partial(false)

	return resourceCloudStackServiceOfferingRead(d, meta)
}

func resourceCloudStackServiceOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.ServiceOffering.NewDeleteServiceOfferingParams(d.Id())
	_, err := cs.ServiceOffering.DeleteServiceOffering(p)

	if err != nil {
		return fmt.Errorf("Error deleting Service Offering: %s", err)
	}

	return nil
}
