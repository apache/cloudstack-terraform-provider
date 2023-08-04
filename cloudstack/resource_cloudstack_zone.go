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
	"errors"
	"fmt"

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
			"allocationstate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domainid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"guestcidraddress": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal_dns1": {
				Type:     schema.TypeString,
				Required: true,
			},
			"internal_dns2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip6dns1": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip6dns2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"localstorageenabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"securitygroupenabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceCloudStackZoneCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameters
	p := cs.Zone.NewCreateZoneParams(d.Get("dns1").(string), d.Get("internal_dns1").(string), d.Get("name").(string), d.Get("network_type").(string))
	if v, ok := d.GetOk("allocationstate"); ok {
		p.SetAllocationstate(v.(string))
	}
	if v, ok := d.GetOk("dns2"); ok {
		p.SetDns2(v.(string))
	}
	if v, ok := d.GetOk("domain"); ok {
		p.SetDomain(v.(string))
	}
	if v, ok := d.GetOk("domainid"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("guestcidraddress"); ok {
		p.SetGuestcidraddress(v.(string))
	}
	if v, ok := d.GetOk("internal_dns2"); ok {
		p.SetInternaldns2(v.(string))
	}
	if v, ok := d.GetOk("ip6dns1"); ok {
		p.SetIp6dns1(v.(string))
	}
	if v, ok := d.GetOk("ip6dns2"); ok {
		p.SetIp6dns2(v.(string))
	}
	if v, ok := d.GetOk("localstorageenabled"); ok {
		p.SetLocalstorageenabled(v.(bool))
	}
	if v, ok := d.GetOk("securitygroupenabled"); ok {
		p.SetSecuritygroupenabled(v.(bool))
	}

	// Create zone
	r, err := cs.Zone.CreateZone(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudStackZoneRead(d, meta)
}

func resourceCloudStackZoneRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	z, count, err := cs.Zone.GetZoneByID(d.Id())
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New(fmt.Sprintf("Multiple zones.  Invalid zone id: %s", d.Id()))
	}

	d.Set("allocationstate", z.Allocationstate)
	d.Set("description", z.Description)
	d.Set("dhcp_provider", z.Dhcpprovider)
	d.Set("dns1", z.Dns1)
	d.Set("dns2", z.Dns2)
	d.Set("domain", z.Domain)
	d.Set("domainid", z.Domainid)
	d.Set("guestcidraddress", z.Guestcidraddress)
	d.Set("internal_dns1", z.Internaldns1)
	d.Set("internal_dns2", z.Internaldns2)
	d.Set("ip6dns1", z.Ip6dns1)
	d.Set("ip6dns2", z.Ip6dns2)
	d.Set("localstorageenabled", z.Localstorageenabled)
	d.Set("name", z.Name)
	d.Set("network_type", z.Networktype)
	d.Set("securitygroupenabled", z.Securitygroupsenabled)

	return nil
}

func resourceCloudStackZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Zone.NewUpdateZoneParams(d.Id())

	if v, ok := d.GetOk("allocationstate"); ok {
		p.SetAllocationstate(v.(string))
	}
	if v, ok := d.GetOk("dhcp_provider"); ok {
		p.SetDhcpprovider(v.(string))
	}
	if v, ok := d.GetOk("dns1"); ok {
		p.SetDns1(v.(string))
	}
	if v, ok := d.GetOk("dns2"); ok {
		p.SetDns2(v.(string))
	}
	if v, ok := d.GetOk("domain"); ok {
		p.SetDomain(v.(string))
	}
	if v, ok := d.GetOk("guestcidraddress"); ok {
		p.SetGuestcidraddress(v.(string))
	}

	if v, ok := d.GetOk("internal_dns1"); ok {
		p.SetInternaldns1(v.(string))
	}
	if v, ok := d.GetOk("internal_dns2"); ok {
		p.SetInternaldns2(v.(string))
	}
	if v, ok := d.GetOk("ip6dns1"); ok {
		p.SetIp6dns1(v.(string))
	}
	if v, ok := d.GetOk("ip6dns2"); ok {
		p.SetIp6dns2(v.(string))
	}
	if v, ok := d.GetOk("localstorageenabled"); ok {
		p.SetLocalstorageenabled(v.(bool))
	}
	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}

	_, err := cs.Zone.UpdateZone(p)
	if err != nil {
		return err
	}

	return resourceCloudStackZoneRead(d, meta)
}

func resourceCloudStackZoneDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Delete zone
	_, err := cs.Zone.DeleteZone(cs.Zone.NewDeleteZoneParams(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting Zone: %s", err)
	}

	d.SetId("")

	return nil
}
