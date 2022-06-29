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

func resourceCloudStackZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackZoneCreate,
		Read:   resourceCloudStackZoneRead,
		Update: resourceCloudStackZoneUpdate,
		Delete: resourceCloudStackZoneDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_dns1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceCloudStackZoneCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	dns1 := d.Get("dns1").(string)
	internal_dns1 := d.Get("internal_dns1").(string)
	network_type := d.Get("network_type").(string)

	// Create a new parameter struct
	p := cs.Zone.NewCreateZoneParams(dns1, internal_dns1, name, network_type)

	log.Printf("[DEBUG] Creating Zone %s", name)
	n, err := cs.Zone.CreateZone(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Zone %s successfully created", name)
	d.SetId(n.Id)

	return resourceCloudStackZoneRead(d, meta)
}

func resourceCloudStackZoneRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Retrieving Zone %s", d.Get("name").(string))

	// Get the Zone details
	z, count, err := cs.Zone.GetZoneByName(d.Get("name").(string))

	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Zone %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(z.Id)
	d.Set("name", z.Name)
	d.Set("dns1", z.Dns1)
	d.Set("internal_dns1", z.Internaldns1)
	d.Set("network_type", z.Networktype)

	return nil
}

func resourceCloudStackZoneUpdate(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceCloudStackZoneDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Zone.NewDeleteZoneParams(d.Id())
	_, err := cs.Zone.DeleteZone(p)

	if err != nil {
		return fmt.Errorf("Error deleting Zone: %s", err)
	}

	return nil
}
