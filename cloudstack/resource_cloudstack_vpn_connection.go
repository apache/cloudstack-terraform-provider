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

func resourceCloudStackVPNConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackVPNConnectionCreate,
		Read:   resourceCloudStackVPNConnectionRead,
		Delete: resourceCloudStackVPNConnectionDelete,

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackVPNConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.VPN.NewCreateVpnConnectionParams(
		d.Get("customer_gateway_id").(string),
		d.Get("vpn_gateway_id").(string),
	)

	// Create the new VPN Connection
	v, err := cs.VPN.CreateVpnConnection(p)
	if err != nil {
		return fmt.Errorf("Error creating VPN Connection: %s", err)
	}

	d.SetId(v.Id)

	return resourceCloudStackVPNConnectionRead(d, meta)
}

func resourceCloudStackVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the VPN Connection details
	v, count, err := cs.VPN.GetVpnConnectionByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] VPN Connection does no longer exist")
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("customer_gateway_id", v.S2scustomergatewayid)
	d.Set("vpn_gateway_id", v.S2svpngatewayid)

	return nil
}

func resourceCloudStackVPNConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.VPN.NewDeleteVpnConnectionParams(d.Id())

	// Delete the VPN Connection
	_, err := cs.VPN.DeleteVpnConnection(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting VPN Connection: %s", err)
	}

	return nil
}
