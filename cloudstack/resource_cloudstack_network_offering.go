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

func resourceCloudStackNetworkOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNetworkOfferingCreate,
		Read:   resourceCloudStackNetworkOfferingRead,
		Update: resourceCloudStackNetworkOfferingUpdate,
		Delete: resourceCloudStackNetworkOfferingDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"guest_ip_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudStackNetworkOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	display_text := d.Get("display_text").(string)
	guest_ip_type := d.Get("guest_ip_type").(string)
	traffic_type := d.Get("traffic_type").(string)

	// Create a new parameter struct
	p := cs.NetworkOffering.NewCreateNetworkOfferingParams(display_text, guest_ip_type, name, traffic_type)

	if guest_ip_type == "Shared" {
		p.SetSpecifyvlan(true)
		p.SetSpecifyipranges(true)
	}

	log.Printf("[DEBUG] Creating Network Offering %s", name)
	n, err := cs.NetworkOffering.CreateNetworkOffering(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Network Offering %s successfully created", name)
	d.SetId(n.Id)

	return resourceCloudStackNetworkOfferingRead(d, meta)
}

func resourceCloudStackNetworkOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	d.Partial(true)

	name := d.Get("name").(string)

	// Check if the name is changed and if so, update the network offering
	if d.HasChange("name") {
		log.Printf("[DEBUG] Name changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new name
		p.SetName(d.Get("name").(string))

		// Update the name
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the name for network offering %s: %s", name, err)
		}

		d.SetPartial("name")
	}

	// Check if the display text is changed and if so, update the virtual machine
	if d.HasChange("display_text") {
		log.Printf("[DEBUG] Display text changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new display text
		p.SetName(d.Get("display_text").(string))

		// Update the display text
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the display text for network offering %s: %s", name, err)
		}

		d.SetPartial("display_text")
	}

	// Check if the guest ip type is changed and if so, update the virtual machine
	if d.HasChange("guest_ip_type") {
		log.Printf("[DEBUG] Guest ip type changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new guest ip type
		p.SetName(d.Get("guest_ip_type").(string))

		// Update the guest ip type
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the guest ip type for network offering %s: %s", name, err)
		}

		d.SetPartial("guest_ip_type")
	}

	// Check if the traffic type is changed and if so, update the virtual machine
	if d.HasChange("traffic_type") {
		log.Printf("[DEBUG] Traffic type changed for %s, starting update", name)

		// Create a new parameter struct
		p := cs.NetworkOffering.NewUpdateNetworkOfferingParams()

		// Set the new traffic type
		p.SetName(d.Get("traffic_type").(string))

		// Update the traffic type
		_, err := cs.NetworkOffering.UpdateNetworkOffering(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating the traffic type for network offering %s: %s", name, err)
		}

		d.SetPartial("traffic_type")
	}

	d.Partial(false)

	return resourceCloudStackInstanceRead(d, meta)
}

func resourceCloudStackNetworkOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.NetworkOffering.NewDeleteNetworkOfferingParams(d.Id())
	_, err := cs.NetworkOffering.DeleteNetworkOffering(p)

	if err != nil {
		return fmt.Errorf("Error deleting Network Offering: %s", err)
	}

	return nil
}

func resourceCloudStackNetworkOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Network Offering %s", d.Get("name").(string))

	// Get the Network Offering details
	n, count, err := cs.NetworkOffering.GetNetworkOfferingByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Network Offering %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId(n.Id)
	d.Set("name", n.Name)
	d.Set("display_text", n.Displaytext)
	d.Set("guest_ip_type", n.Guestiptype)
	d.Set("traffic_type", n.Traffictype)

	return nil
}
