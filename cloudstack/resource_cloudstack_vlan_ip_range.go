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

func resourceCloudstackVlanIpRange() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudstackVlanIpRangeCreate,
		Read:   resourceCloudstackVlanIpRangeRead,
		Update: resourceCloudstackVlanIpRangeUpdate,
		Delete: resourceCloudstackVlanIpRangeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account": {
				Description: "account who will own the VLAN. If VLAN is Zone wide, this parameter should be omitted",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"domain_id": {
				Description: "domain ID of the account owning a VLAN",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"end_ip": {
				Description: "the ending IP address in the VLAN IP range",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"end_ipv6": {
				Description: "the ending IPv6 address in the IPv6 network range",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"for_system_vms": {
				Description: "true if IP range is set to system vms, false if not",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"for_virtual_network": {
				Description: "true if VLAN is of Virtual type, false if Direct",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"gateway": {
				Description: "he gateway of the VLAN IP range",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ip6_cidr": {
				Description: "the CIDR of IPv6 network, must be at least /64",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ip6_gateway": {
				Description: "the gateway of the IPv6 network. Required for Shared networks and Isolated networks when it belongs to VPC",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"netmask": {
				Description: "the netmask of the VLAN IP range",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"network_id": {
				Description: "the network id",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"physical_network_id": {
				Description: "the physical network id",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"pod_id": {
				Description: "optional parameter. Have to be specified for Direct Untagged vlan only.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"project_id": {
				Description: "project who will own the VLAN. If VLAN is Zone wide, this parameter should be omitted",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"start_ip": {
				Description: "the beginning IP address in the VLAN IP range",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"start_ipv6": {
				Description: "the beginning IPv6 address in the IPv6 network range",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"vlan": {
				Description: "the ID or VID of the VLAN. If not specified, will be defaulted to the vlan of the network or if vlan of the network is null - to Untagged",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"zone_id": {
				Description: "the Zone ID of the VLAN IP range",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudstackVlanIpRangeCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.VLAN.NewCreateVlanIpRangeParams()
	if v, ok := d.GetOk("account"); ok {
		p.SetAccount(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(v.(string))
	}
	if v, ok := d.GetOk("end_ipv6"); ok {
		p.SetEndipv6(v.(string))
	}
	if v, ok := d.GetOk("for_system_vms"); ok {
		p.SetForsystemvms(v.(bool))
	}
	if v, ok := d.GetOk("for_virtual_network"); ok {
		p.SetForvirtualnetwork(v.(bool))
	}
	if v, ok := d.GetOk("gateway"); ok {
		p.SetGateway(v.(string))
	}
	if v, ok := d.GetOk("ip6_cidr"); ok {
		p.SetIp6cidr(v.(string))
	}
	if v, ok := d.GetOk("ip6_gateway"); ok {
		p.SetIp6gateway(v.(string))
	}
	if v, ok := d.GetOk("netmask"); ok {
		p.SetNetmask(v.(string))
	}
	if v, ok := d.GetOk("network_id"); ok {
		p.SetNetworkid(v.(string))
	}
	if v, ok := d.GetOk("physical_network_id"); ok {
		p.SetPhysicalnetworkid(v.(string))
	}
	if v, ok := d.GetOk("pod_id"); ok {
		p.SetPodid(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(v.(string))
	}
	if v, ok := d.GetOk("start_ip"); ok {
		p.SetStartip(v.(string))
	}
	if v, ok := d.GetOk("start_ipv6"); ok {
		p.SetStartipv6(v.(string))
	}
	if v, ok := d.GetOk("vlan"); ok {
		p.SetVlan(v.(string))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		p.SetZoneid(v.(string))
	}

	// create vlan ip range
	r, err := cs.VLAN.CreateVlanIpRange(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudstackVlanIpRangeRead(d, meta)
}

func resourceCloudstackVlanIpRangeRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.VLAN.GetVlanIpRangeByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("account", r.Account)
	d.Set("domain_id", r.Domainid)
	d.Set("end_ip", r.Endip)
	d.Set("end_ipv6", r.Endipv6)
	d.Set("for_system_vms", r.Forsystemvms)
	d.Set("for_virtual_network", r.Forvirtualnetwork)
	d.Set("gateway", r.Gateway)
	d.Set("ip6_cidr", r.Ip6cidr)
	d.Set("ip6_gateway", r.Ip6gateway)
	d.Set("netmask", r.Netmask)
	d.Set("network_id", r.Networkid)
	d.Set("physical_network_id", r.Physicalnetworkid)
	d.Set("pod_id", r.Podid)
	d.Set("project_id", r.Projectid)
	d.Set("start_ip", r.Startip)
	d.Set("start_ipv6", r.Startipv6)
	d.Set("vlan", r.Vlan)
	d.Set("zone_id", r.Zoneid)

	return nil
}

func resourceCloudstackVlanIpRangeUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.VLAN.NewUpdateVlanIpRangeParams(d.Id())

	if v, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(v.(string))
	}
	if v, ok := d.GetOk("end_ipv6"); ok {
		p.SetEndipv6(v.(string))
	}
	if v, ok := d.GetOk("for_system_vms"); ok {
		p.SetForsystemvms(v.(bool))
	}
	if v, ok := d.GetOk("gateway"); ok {
		p.SetGateway(v.(string))
	}
	if v, ok := d.GetOk("ip6_cidr"); ok {
		p.SetIp6cidr(v.(string))
	}
	if v, ok := d.GetOk("ip6_gateway"); ok {
		p.SetIp6gateway(v.(string))
	}
	if v, ok := d.GetOk("netmask"); ok {
		p.SetNetmask(v.(string))
	}
	if v, ok := d.GetOk("start_ip"); ok {
		p.SetStartip(v.(string))
	}
	if v, ok := d.GetOk("start_ipv6"); ok {
		p.SetStartipv6(v.(string))
	}

	// update vlaniprange
	_, err := cs.VLAN.UpdateVlanIpRange(p)
	if err != nil {
		return err
	}

	return resourceCloudstackVlanIpRangeRead(d, meta)
}

func resourceCloudstackVlanIpRangeDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.VLAN.DeleteVlanIpRange(cs.VLAN.NewDeleteVlanIpRangeParams(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting Vlan Ip Range: %s", err)
	}

	return nil
}
