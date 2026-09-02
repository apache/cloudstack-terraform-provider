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

func resourceCloudStackIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackIPAddressCreate,
		Read:   resourceCloudStackIPAddressRead,
		Delete: resourceCloudStackIPAddressDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudStackIPAddressImport,
		},

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
				Optional: true,
				Computed: true,
				ForceNew: true,
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
		if vpcid, ok := d.GetOk("vpc_id"); ok && vpcid.(string) != "" {
			return fmt.Errorf("set only network_id or vpc_id")
		}

		// If no project is explicitly set, try to inherit it from the network
		if _, ok := d.GetOk("project"); !ok {
			// Get the network to retrieve its project
			// Use projectid=-1 to search across all projects
			network, count, err := cs.Network.GetNetworkByID(networkid.(string), cloudstack.WithProject("-1"))
			if err == nil && count > 0 && network.Projectid != "" {
				log.Printf("[DEBUG] Inheriting project %s from network %s", network.Projectid, networkid.(string))
				p.SetProjectid(network.Projectid)
			}
		}
	}

	if vpcid, ok := d.GetOk("vpc_id"); ok {
		// Set the vpcid
		p.SetVpcid(vpcid.(string))

		// If no project is explicitly set, try to inherit it from the VPC
		if _, ok := d.GetOk("project"); !ok {
			// Get the VPC to retrieve its project
			// Use projectid=-1 to search across all projects
			vpc, count, err := cs.VPC.GetVPCByID(vpcid.(string), cloudstack.WithProject("-1"))
			if err == nil && count > 0 && vpc.Projectid != "" {
				log.Printf("[DEBUG] Inheriting project %s from VPC %s", vpc.Projectid, vpcid.(string))
				p.SetProjectid(vpc.Projectid)
			}
		}
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
	// This will override the inherited project from VPC or network if explicitly set
	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	if ipaddress, ok := d.GetOk("ip_address"); ok {
		p.SetIpaddress(ipaddress.(string))
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

func resourceCloudStackIPAddressImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	cs := meta.(*cloudstack.CloudStackClient)

	// Try to split the ID to extract the optional project name.
	s := strings.SplitN(d.Id(), "/", 2)
	if len(s) == 2 {
		d.Set("project", s[0])
	}

	ipAddressID := s[len(s)-1]
	d.SetId(ipAddressID)

	ip, count, err := cs.Address.GetPublicIpAddressByID(
		ipAddressID,
		cloudstack.WithProject(d.Get("project").(string)),
	)
	if err != nil {
		if count == 0 {
			return nil, fmt.Errorf("IP address with ID %s does not exist", ipAddressID)
		}
		return nil, err
	}

	// Seed whichever of network_id/vpc_id actually applies before Read runs:
	// Read only refreshes these when already present in state, so import
	// needs to set the right one here first.
	if ip.Vpcid != "" {
		d.Set("vpc_id", ip.Vpcid)
	} else if ip.Associatednetworkid != "" {
		d.Set("network_id", ip.Associatednetworkid)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceCloudStackIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the IP address details
	// First try with the project from state (if any)
	project := d.Get("project").(string)
	ip, count, err := cs.Address.GetPublicIpAddressByID(
		d.Id(),
		cloudstack.WithProject(project),
	)

	// If not found and no explicit project was set, try with projectid=-1
	// This handles the case where the project was inherited from the VPC or network
	if count == 0 && project == "" {
		ip, count, err = cs.Address.GetPublicIpAddressByID(
			d.Id(),
			cloudstack.WithProject("-1"),
		)
	}

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

	// Only refresh network_id/vpc_id if already present in state: a plain
	// zone-scoped IP (neither set) can still come back from the API with an
	// associated network under the hood, and syncing that into an Optional,
	// non-Computed, ForceNew field would create a permanent diff. The
	// importer is responsible for seeding whichever one applies before this
	// Read runs.
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

		// A source NAT IP can't be disassociated while its network/VPC still exists;
		// deleting that network/VPC releases it instead, so treat this as a no-op.
		if strings.Contains(err.Error(), "used for source nat purposes and can not be disassociated") {
			return nil
		}

		return fmt.Errorf("Error disassociating IP address %s: %s", d.Id(), err)
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
