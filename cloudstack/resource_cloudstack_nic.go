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

func resourceCloudStackNIC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNICCreate,
		Read:   resourceCloudStackNICRead,
		Delete: resourceCloudStackNICDelete,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackNICCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.VirtualMachine.NewAddNicToVirtualMachineParams(
		d.Get("network_id").(string),
		d.Get("virtual_machine_id").(string),
	)

	// If there is an ipaddress supplied, add it to the parameter struct
	if ipaddress, ok := d.GetOk("ip_address"); ok {
		p.SetIpaddress(ipaddress.(string))
	}

	// If there is a macaddress supplied, add it to the parameter struct
	if macaddress, ok := d.GetOk("mac_address"); ok {
		p.SetMacaddress(macaddress.(string))
	}

	// Create and attach the new NIC
	r, err := Retry(10, retryableAddNicFunc(cs, p))
	if err != nil {
		return fmt.Errorf("Error creating the new NIC: %s", err)
	}

	return setCloudStackNICStateFromCreateResponse(
		d,
		r.(*cloudstack.AddNicToVirtualMachineResponse),
	)
}

func resourceCloudStackNICRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Nic.NewListNicsParams(d.Get("virtual_machine_id").(string))
	p.SetNicid(d.Id())

	l, err := cs.Nic.ListNics(p)
	if err != nil {
		return err
	}

	switch len(l.Nics) {
	case 0:
		log.Printf("[DEBUG] NIC %s no longer exists", d.Id())
		d.SetId("")
		return nil
	case 1:
		return setCloudStackNICState(d, l.Nics[0])
	default:
		return fmt.Errorf("found more than one NIC for ID %s: %v", d.Id(), l.Nics)
	}
}

func setCloudStackNICStateFromCreateResponse(
	d *schema.ResourceData,
	r *cloudstack.AddNicToVirtualMachineResponse,
) error {
	networkID := d.Get("network_id").(string)
	for i := range r.Nic {
		if r.Nic[i].Networkid == networkID {
			return setCloudStackNICState(d, &r.Nic[i])
		}
	}

	return fmt.Errorf("could not find NIC ID for network ID: %s", networkID)
}

func setCloudStackNICState(d *schema.ResourceData, nic *cloudstack.Nic) error {
	if nic == nil || nic.Id == "" {
		return fmt.Errorf("NIC response did not contain an ID")
	}

	d.Set("ip_address", nic.Ipaddress)
	d.Set("network_id", nic.Networkid)
	// The NIC embedded in an addNicToVirtualMachine response may omit the
	// virtual machine ID. Preserve the configured value in that case so a
	// subsequent refresh does not call listNics with an empty required ID.
	if nic.Virtualmachineid != "" {
		d.Set("virtual_machine_id", nic.Virtualmachineid)
	}
	d.Set("mac_address", nic.Macaddress)
	d.SetId(nic.Id)

	return nil
}

func resourceCloudStackNICDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.VirtualMachine.NewRemoveNicFromVirtualMachineParams(
		d.Id(),
		d.Get("virtual_machine_id").(string),
	)

	// Remove the NIC
	_, err := cs.VirtualMachine.RemoveNicFromVirtualMachine(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting NIC: %s", err)
	}

	return nil
}

func retryableAddNicFunc(cs *cloudstack.CloudStackClient, p *cloudstack.AddNicToVirtualMachineParams) func() (interface{}, error) {
	return func() (interface{}, error) {
		r, err := cs.VirtualMachine.AddNicToVirtualMachine(p)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}
