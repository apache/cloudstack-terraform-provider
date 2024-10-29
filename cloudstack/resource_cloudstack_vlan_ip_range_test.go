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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudstackVlanIpRange_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudstackVlanIpRange_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "gateway", "10.0.0.1"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "start_ip", "10.0.0.2"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "end_ip", "10.0.0.10"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "vlan", "vlan://123"),
				),
			},
			{
				Config: testAccCloudstackVlanIpRange_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "gateway", "10.0.0.1"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "start_ip", "10.0.0.2"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "end_ip", "10.0.0.20"),
					resource.TestCheckResourceAttr("cloudstack_vlan_ip_range.test", "vlan", "vlan://123"),
				),
			},
		},
	})
}

const testAccCloudstackVlanIpRange_basic = `
resource "cloudstack_zone" "test" {
	name          = "acctest"
	dns1          = "8.8.8.8"
	dns2          = "8.8.8.8"
	internal_dns1 = "8.8.4.4"
	internal_dns2 = "8.8.4.4"
	network_type  = "Advanced"
	domain        = "cloudstack.apache.org"
}
resource "cloudstack_physical_network" "test" {
	broadcast_domain_range = "ZONE"
	isolation_methods      = "VLAN"
	name                   = "test01"
	network_speed          = "1G"
	tags                   = "vlan"
	zone_id                = cloudstack_zone.test.id
}
resource "cloudstack_vlan_ip_range" "test" {
    physical_network_id = cloudstack_physical_network.test.id
    for_virtual_network = true
    zone_id = cloudstack_zone.test.id
    gateway = "10.0.0.1"
    netmask = "255.255.255.0"
    start_ip = "10.0.0.2"
    end_ip   = "10.0.0.10"
    vlan     = "vlan://123"
}
`
const testAccCloudstackVlanIpRange_update = `
resource "cloudstack_zone" "test" {
	name          = "acctest"
	dns1          = "8.8.8.8"
	dns2          = "8.8.8.8"
	internal_dns1 = "8.8.4.4"
	internal_dns2 = "8.8.4.4"
	network_type  = "Advanced"
	domain        = "cloudstack.apache.org"
}
resource "cloudstack_physical_network" "test" {
	broadcast_domain_range = "ZONE"
	isolation_methods      = "VLAN"
	name                   = "test01"
	network_speed          = "1G"
	tags                   = "vlan"
	zone_id                = cloudstack_zone.test.id
}
resource "cloudstack_vlan_ip_range" "test" {
    physical_network_id = cloudstack_physical_network.test.id
    for_virtual_network = true
    zone_id = cloudstack_zone.test.id
    gateway = "10.0.0.1"
    netmask = "255.255.255.0"
    start_ip = "10.0.0.2"
    end_ip   = "10.0.0.20"
    vlan     = "vlan://123"
}
`
