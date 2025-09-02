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

func resourceCloudStackPhysicalNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackPhysicalNetworkCreate,
		Read:   resourceCloudStackPhysicalNetworkRead,
		Update: resourceCloudStackPhysicalNetworkUpdate,
		Delete: resourceCloudStackPhysicalNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: importStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"broadcast_domain_range": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ZONE",
				ForceNew: true,
			},

			"isolation_methods": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"network_speed": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vlan": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCloudStackPhysicalNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	// Retrieve the zone ID
	zoneid, e := retrieveID(cs, "zone", d.Get("zone").(string))
	if e != nil {
		return e.Error()
	}

	// Create a new parameter struct
	p := cs.Network.NewCreatePhysicalNetworkParams(name, zoneid)

	// Set optional parameters
	if broadcastDomainRange, ok := d.GetOk("broadcast_domain_range"); ok {
		p.SetBroadcastdomainrange(broadcastDomainRange.(string))
	}

	if isolationMethods, ok := d.GetOk("isolation_methods"); ok {
		methods := make([]string, len(isolationMethods.([]interface{})))
		for i, v := range isolationMethods.([]interface{}) {
			methods[i] = v.(string)
		}
		p.SetIsolationmethods(methods)
	}

	if networkSpeed, ok := d.GetOk("network_speed"); ok {
		p.SetNetworkspeed(networkSpeed.(string))
	}

	if vlan, ok := d.GetOk("vlan"); ok {
		p.SetVlan(vlan.(string))
	}

	// Create the physical network
	r, err := cs.Network.CreatePhysicalNetwork(p)
	if err != nil {
		return fmt.Errorf("Error creating physical network %s: %s", name, err)
	}

	d.SetId(r.Id)

	// Physical networks don't support tags in CloudStack API

	return resourceCloudStackPhysicalNetworkRead(d, meta)
}

func resourceCloudStackPhysicalNetworkRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the physical network details
	p, count, err := cs.Network.GetPhysicalNetworkByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Physical network %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", p.Name)
	d.Set("broadcast_domain_range", p.Broadcastdomainrange)
	d.Set("network_speed", p.Networkspeed)
	d.Set("vlan", p.Vlan)

	// Set isolation methods
	if p.Isolationmethods != "" {
		methods := strings.Split(p.Isolationmethods, ",")
		d.Set("isolation_methods", methods)
	}

	// Set the zone
	setValueOrID(d, "zone", p.Zonename, p.Zoneid)

	// Physical networks don't support tags in CloudStack API

	return nil
}

func resourceCloudStackPhysicalNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Network.NewUpdatePhysicalNetworkParams(d.Id())

	// The UpdatePhysicalNetworkParams struct doesn't have a SetName method
	// so we can't update the name

	if d.HasChange("network_speed") {
		p.SetNetworkspeed(d.Get("network_speed").(string))
	}

	if d.HasChange("vlan") {
		p.SetVlan(d.Get("vlan").(string))
	}

	// Update the physical network
	_, err := cs.Network.UpdatePhysicalNetwork(p)
	if err != nil {
		return fmt.Errorf("Error updating physical network %s: %s", d.Get("name").(string), err)
	}

	// Physical networks don't support tags in CloudStack API

	return resourceCloudStackPhysicalNetworkRead(d, meta)
}

func resourceCloudStackPhysicalNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Network.NewDeletePhysicalNetworkParams(d.Id())

	// Delete the physical network
	_, err := cs.Network.DeletePhysicalNetwork(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting physical network %s: %s", d.Get("name").(string), err)
	}

	return nil
}
