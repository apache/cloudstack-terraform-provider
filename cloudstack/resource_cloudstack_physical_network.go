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
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"broadcast_domain_range": {
				Description: "the broadcast domain range for the physical network[Pod or Zone]. In Acton release it can be Zone only in Advance zone, and Pod in Basic",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"domain_id": {
				Description: "domain ID of the account owning a physical network",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"isolation_methods": {
				Description: "the isolation method for the physical network[VLAN/L3/GRE]",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "the name of the physical network",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"network_speed": {
				Description: "the speed for the physical network[1G/10G]",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": {
				Description: "Tag the physical network",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vlan": {
				Description: "the VLAN for the physical network",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"zone_id": {
				Description: "zone id of the physical network",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudStackPhysicalNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Network.NewCreatePhysicalNetworkParams(d.Get("name").(string), d.Get("zone_id").(string))
	if v, ok := d.GetOk("broadcast_domain_range"); ok {
		p.SetBroadcastdomainrange(strings.ToUpper(v.(string)))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("isolation_methods"); ok {
		p.SetIsolationmethods([]string{v.(string)})
	}
	if v, ok := d.GetOk("network_speed"); ok {
		p.SetNetworkspeed(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags([]string{v.(string)})
	}
	if v, ok := d.GetOk("vlan"); ok {
		p.SetVlan(v.(string))
	}

	r, err := cs.Network.CreatePhysicalNetwork(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudStackPhysicalNetworkRead(d, meta)
}

func resourceCloudStackPhysicalNetworkRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	pn, _, err := cs.Network.GetPhysicalNetworkByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("broadcast_domain_range", pn.Broadcastdomainrange)
	d.Set("domain_id", pn.Domainid)
	d.Set("isolation_methods", pn.Isolationmethods)
	d.Set("name", pn.Name)
	d.Set("network_speed", pn.Networkspeed)
	d.Set("tags", pn.Tags)
	d.Set("vlan", pn.Vlan)
	d.Set("zone_id", pn.Zoneid)

	return nil
}

func resourceCloudStackPhysicalNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Network.NewUpdatePhysicalNetworkParams(d.Id())
	if v, ok := d.GetOk("network_speed"); ok {
		p.SetNetworkspeed(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags([]string{v.(string)})
	}
	if v, ok := d.GetOk("vlan"); ok {
		p.SetVlan(v.(string))
	}

	_, err := cs.Network.UpdatePhysicalNetwork(p)
	if err != nil {
		return fmt.Errorf("Error deleting physical network: %s", err)
	}

	return resourceCloudStackPhysicalNetworkRead(d, meta)
}

func resourceCloudStackPhysicalNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Network.NewDeletePhysicalNetworkParams(d.Id())
	_, err := cs.Network.DeletePhysicalNetwork(p)
	if err != nil {
		return fmt.Errorf("Error deleting phsyical network: %s", err)
	}

	return nil
}
