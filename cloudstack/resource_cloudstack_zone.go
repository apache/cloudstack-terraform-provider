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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackZoneCreate,
		Read:   resourceCloudStackZoneRead,
		Update: resourceCloudStackZoneUpdate,
		Delete: resourceCloudStackZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"allocation_state": {
				Description: "Allocation state of this Zone for allocation of new resources",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"dhcp_provider": {
				Description: "the dhcp Provider for the Zone",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"dns1": {
				Description: "the first DNS for the Zone",
				Type:        schema.TypeString,
				Required:    true,
			},
			"dns2": {
				Description: "the second DNS for the Zone",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain": {
				Description: "Network domain name for the networks in the zone",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"domain_id": {
				Description: "the ID of the containing domain, null for public zones",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"guest_cidr_address": {
				Description: "the guest CIDR address for the Zone",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"internal_dns1": {
				Description: "the first internal DNS for the Zone",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_dns2": {
				Description: "the second internal DNS for the Zone",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ip6_dns1": {
				Description: "the first DNS for IPv6 network in the Zone",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ip6_dns2": {
				Description: "the second DNS for IPv6 network in the Zone",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"local_storage_enabled": {
				Description: "true if local storage offering enabled, false otherwise",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Description: "the name of the Zone",
				Type:        schema.TypeString,
				Required:    true,
			},
			"network_type": {
				Description: "network type of the zone, can be Basic or Advanced",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"security_group_enabled": {
				Description: "true if network is security group enabled, false otherwise",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudStackZoneCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Zone.NewCreateZoneParams(d.Get("dns1").(string), d.Get("internal_dns1").(string), d.Get("name").(string), d.Get("network_type").(string))
	if v, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(v.(string))
	}
	if v, ok := d.GetOk("dns2"); ok {
		p.SetDns2(v.(string))
	}
	if v, ok := d.GetOk("domain"); ok {
		p.SetDomain(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("guest_cidr_address"); ok {
		p.SetGuestcidraddress(v.(string))
	}
	if v, ok := d.GetOk("internal_dns2"); ok {
		p.SetInternaldns2(v.(string))
	}
	if v, ok := d.GetOk("ip6_dns1"); ok {
		p.SetIp6dns1(v.(string))
	}
	if v, ok := d.GetOk("ip6_dns2"); ok {
		p.SetIp6dns2(v.(string))
	}
	if v, ok := d.GetOk("local_storage_enabled"); ok {
		p.SetLocalstorageenabled(v.(bool))
	}
	if v, ok := d.GetOk("security_group_enabled"); ok {
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

	z, _, err := cs.Zone.GetZoneByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("allocation_state", z.Allocationstate)
	d.Set("dhcp_provider", z.Dhcpprovider)
	d.Set("dns1", z.Dns1)
	d.Set("dns2", z.Dns2)
	d.Set("domain", z.Domain)
	d.Set("domain_id", z.Domainid)
	d.Set("guest_cidr_address", z.Guestcidraddress)
	d.Set("internal_dns1", z.Internaldns1)
	d.Set("internal_dns2", z.Internaldns2)
	d.Set("ip6_dns1", z.Ip6dns1)
	d.Set("ip6_dns2", z.Ip6dns2)
	d.Set("local_storage_enabled", z.Localstorageenabled)
	d.Set("name", z.Name)
	d.Set("network_type", z.Networktype)
	d.Set("security_group_enabled", z.Securitygroupsenabled)

	return nil
}

func resourceCloudStackZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Zone.NewUpdateZoneParams(d.Id())
	if v, ok := d.GetOk("allocation_state"); ok {
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
	if v, ok := d.GetOk("guest_cidr_address"); ok {
		p.SetGuestcidraddress(v.(string))
	}
	if v, ok := d.GetOk("internal_dns1"); ok {
		p.SetInternaldns1(v.(string))
	}
	if v, ok := d.GetOk("internal_dns2"); ok {
		p.SetInternaldns2(v.(string))
	}
	if v, ok := d.GetOk("ip6_dns1"); ok {
		p.SetIp6dns1(v.(string))
	}
	if v, ok := d.GetOk("ip6_dns2"); ok {
		p.SetIp6dns2(v.(string))
	}
	if v, ok := d.GetOk("local_storage_enabled"); ok {
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

	return nil
}
