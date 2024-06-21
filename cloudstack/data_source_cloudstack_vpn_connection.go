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
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackVPNConnection() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackVPNConnectionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"s2s_customer_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"s2s_vpn_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceCloudStackVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.VPN.NewListVpnConnectionsParams()
	csVPNConnections, err := cs.VPN.ListVpnConnections(p)

	if err != nil {
		return fmt.Errorf("Failed to list VPNs: %s", err)
	}

	filters := d.Get("filter")
	var vpnConnections []*cloudstack.VpnConnection

	for _, v := range csVPNConnections.VpnConnections {
		match, err := applyVPNConnectionFilters(v, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			vpnConnections = append(vpnConnections, v)
		}
	}

	if len(vpnConnections) == 0 {
		return fmt.Errorf("No VPN Connection is matching with the specified regex")
	}
	//return the latest VPN Connection from the list of filtered VPN Connections according
	//to its creation date
	vpnConnection, err := latestVPNConnection(vpnConnections)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected VPN Connections: %s\n", vpnConnection.Id)

	return vpnConnectionDescriptionAttributes(d, vpnConnection)
}

func vpnConnectionDescriptionAttributes(d *schema.ResourceData, vpnConnection *cloudstack.VpnConnection) error {
	d.SetId(vpnConnection.Id)
	d.Set("s2s_customer_gateway_id", vpnConnection.S2scustomergatewayid)
	d.Set("s2s_vpn_gateway_id", vpnConnection.S2svpngatewayid)

	return nil
}

func latestVPNConnection(vpnConnections []*cloudstack.VpnConnection) (*cloudstack.VpnConnection, error) {
	var latest time.Time
	var vpnConnection *cloudstack.VpnConnection

	for _, v := range vpnConnections {
		created, err := time.Parse("2006-01-02T15:04:05-0700", v.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a VPN Connection: %s", err)
		}

		if created.After(latest) {
			latest = created
			vpnConnection = v
		}
	}

	return vpnConnection, nil
}

func applyVPNConnectionFilters(vpnConnection *cloudstack.VpnConnection, filters *schema.Set) (bool, error) {
	var vpnConnectionJSON map[string]interface{}
	k, _ := json.Marshal(vpnConnection)
	err := json.Unmarshal(k, &vpnConnectionJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		log.Print(m)
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		log.Print(updatedName)
		vpnConnectionField := vpnConnectionJSON[updatedName].(string)
		if !r.MatchString(vpnConnectionField) {
			return false, nil
		}
	}
	return true, nil
}
