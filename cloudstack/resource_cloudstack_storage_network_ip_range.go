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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackStorageNetworkIpRange() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackStorageNetworkIpRangeCreate,
		Read:   resourceCloudStackStorageNetworkIpRangeRead,
		Update: resourceCloudStackStorageNetworkIpRangeUpdate,
		Delete: resourceCloudStackStorageNetworkIpRangeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway": {
				Description: "the gateway for the storage network IP range",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"netmask": {
				Description: "the netmask for the storage network IP range",
				Type:        schema.TypeString,
				Required:    true,
			},
			"pod_id": {
				Description: "the Pod ID for the storage network IP range",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"start_ip": {
				Description: "the beginning IP address in the storage network IP range",
				Type:        schema.TypeString,
				Required:    true,
			},
			"end_ip": {
				Description: "the ending IP address in the storage network IP range",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"vlan": {
				Description: "the optional VLAN of the storage network IP range",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"network_id": {
				Description: "the network id of the storage network IP range",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceCloudStackStorageNetworkIpRangeCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Network.NewCreateStorageNetworkIpRangeParams(
		d.Get("gateway").(string),
		d.Get("netmask").(string),
		d.Get("pod_id").(string),
		d.Get("start_ip").(string),
	)

	if v, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(v.(string))
	}
	if v, ok := d.GetOk("vlan"); ok {
		p.SetVlan(v.(int))
	}

	r, err := cs.Network.CreateStorageNetworkIpRange(p)
	if err != nil {
		return fmt.Errorf("Error creating storage network IP range: %s", err)
	}

	d.SetId(r.Id)

	return resourceCloudStackStorageNetworkIpRangeRead(d, meta)
}

func resourceCloudStackStorageNetworkIpRangeRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, count, err := cs.Network.GetStorageNetworkIpRangeByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Storage network IP range %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("gateway", r.Gateway)
	d.Set("netmask", r.Netmask)
	d.Set("pod_id", r.Podid)
	d.Set("start_ip", r.Startip)
	d.Set("end_ip", r.Endip)
	d.Set("vlan", r.Vlan)
	d.Set("network_id", r.Networkid)

	return nil
}

func resourceCloudStackStorageNetworkIpRangeUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Network.NewUpdateStorageNetworkIpRangeParams(d.Id())

	if v, ok := d.GetOk("netmask"); ok {
		p.SetNetmask(v.(string))
	}
	if v, ok := d.GetOk("start_ip"); ok {
		p.SetStartip(v.(string))
	}
	if v, ok := d.GetOk("end_ip"); ok {
		p.SetEndip(v.(string))
	}
	if v, ok := d.GetOk("vlan"); ok {
		p.SetVlan(v.(int))
	}

	_, err := cs.Network.UpdateStorageNetworkIpRange(p)
	if err != nil {
		return fmt.Errorf("Error updating storage network IP range: %s", err)
	}

	return resourceCloudStackStorageNetworkIpRangeRead(d, meta)
}

func resourceCloudStackStorageNetworkIpRangeDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.Network.DeleteStorageNetworkIpRange(cs.Network.NewDeleteStorageNetworkIpRangeParams(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting storage network IP range: %s", err)
	}

	return nil
}
