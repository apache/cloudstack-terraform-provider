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
	p := cs.NetworkOffering.NewCreateNetworkOfferingParams(display_text, guest_ip_type, name, []string{}, traffic_type)

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
	return nil
}

func resourceCloudStackNetworkOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudStackNetworkOfferingRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
