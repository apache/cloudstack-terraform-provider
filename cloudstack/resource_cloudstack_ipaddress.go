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
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackIPAddressCreate,
		Read:   resourceCloudStackIPAddressRead,
		Delete: resourceCloudStackIPAddressDelete,

		Schema: map[string]*schema.Schema{
			"is_portable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_source_nat": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceCloudStackIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if err := verifyIPAddressParams(d); err != nil {
		return err
	}

	// Create a new parameter struct
	p := cs.Address.NewAssociateIpAddressParams()

	if d.Get("is_portable").(bool) {
		p.SetIsportable(true)
	}

	if networkid, ok := d.GetOk("network_id"); ok {
		// Set the networkid
		p.SetNetworkid(networkid.(string))
	}

	if vpcid, ok := d.GetOk("vpc_id"); ok {
		// Set the vpcid
		p.SetVpcid(vpcid.(string))
	}

	if zone, ok := d.GetOk("zone"); ok {
		// Retrieve the zone ID
		zoneid, e := retrieveID(cs, "zone", zone.(string))
		if e != nil {
			return e.Error()
		}

		// Set the zoneid
		p.SetZoneid(zoneid)
	}

	// If there is a project supplied, we retrieve and set the project id
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	// Associate a new IP address
	r, err := cs.Address.AssociateIpAddress(p)
	if err != nil {
		return fmt.Errorf("Error associating a new IP address: %s", err)
	}

	d.SetId(r.Id)

	// Set tags if necessary
	err = setTags(cs, d, "PublicIpAddress")
	if err != nil {
		return fmt.Errorf("Error setting tags on the IP address: %s", err)
	}

	return resourceCloudStackIPAddressRead(d, meta)
}

func resourceCloudStackIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the IP address details
	ip, count, err := cs.Address.GetPublicIpAddressByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project").(string)),
	)
	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] IP address with ID %s is no longer associated", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("is_portable", ip.Isportable)
	d.Set("is_source_nat", ip.Issourcenat)

	// Updated the IP address
	d.Set("ip_address", ip.Ipaddress)

	if _, ok := d.GetOk("network_id"); ok {
		d.Set("network_id", ip.Associatednetworkid)
	}

	if _, ok := d.GetOk("vpc_id"); ok {
		d.Set("vpc_id", ip.Vpcid)
	}

	if _, ok := d.GetOk("zone"); ok {
		setValueOrID(d, "zone", ip.Zonename, ip.Zoneid)
	}

	tags := make(map[string]interface{})
	for _, tag := range ip.Tags {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	setValueOrID(d, "project", ip.Project, ip.Projectid)

	return nil
}

func resourceCloudStackIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	if !d.Get("is_source_nat").(bool) {
		cs := meta.(*cloudstack.CloudStackClient)

		// Create a new parameter struct
		p := cs.Address.NewDisassociateIpAddressParams(d.Id())

		// Disassociate the IP address
		if _, err := cs.Address.DisassociateIpAddress(p); err != nil {
			// This is a very poor way to be told the ID does no longer exist :(
			if strings.Contains(err.Error(), fmt.Sprintf(
				"Invalid parameter id value=%s due to incorrect long value format, "+
					"or entity does not exist", d.Id())) {
				return nil
			}

			return fmt.Errorf("Error disassociating IP address %s: %s", d.Id(), err)
		}
	}

	return nil
}

func verifyIPAddressParams(d *schema.ResourceData) error {
	_, portable := d.GetOk("is_portable")
	_, network := d.GetOk("network_id")
	_, vpc := d.GetOk("vpc_id")
	_, zone := d.GetOk("zone")

	if portable && ((network && vpc) || (!network && !vpc)) {
		return fmt.Errorf(
			"You must supply a value for either (so not both) the 'network_id' or 'vpc_id' parameter for a portable IP")
	}

	if !portable && !zone && !network {
		return fmt.Errorf("You must supply a value for the 'network_id' and/or 'zone' parameters for a non portable IP")
	}

	return nil
}
