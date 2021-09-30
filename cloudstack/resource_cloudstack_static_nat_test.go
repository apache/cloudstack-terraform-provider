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
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudStackStaticNAT_basic(t *testing.T) {
	var ipaddr cloudstack.PublicIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackStaticNATDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackStaticNAT_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackStaticNATExists(
						"cloudstack_static_nat.foo", &ipaddr),
					testAccCheckCloudStackStaticNATAttributes(&ipaddr),
				),
			},
		},
	})
}

func testAccCheckCloudStackStaticNATExists(
	n string, ipaddr *cloudstack.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static NAT ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		ip, _, err := cs.Address.GetPublicIpAddressByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ip.Id != rs.Primary.ID {
			return fmt.Errorf("Static NAT not found")
		}

		if !ip.Isstaticnat {
			return fmt.Errorf("Static NAT not enabled")
		}

		*ipaddr = *ip

		return nil
	}
}

func testAccCheckCloudStackStaticNATAttributes(
	ipaddr *cloudstack.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if ipaddr.Associatednetworkname != "terraform-network" {
			return fmt.Errorf("Bad network name: %s", ipaddr.Associatednetworkname)
		}

		return nil
	}
}

func testAccCheckCloudStackStaticNATDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_static_nat" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static NAT ID is set")
		}

		ip, _, err := cs.Address.GetPublicIpAddressByID(rs.Primary.ID)
		if err == nil && ip.Isstaticnat {
			return fmt.Errorf("Static NAT %s still enabled", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackStaticNAT_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
	source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_ipaddress" "foo" {
  network_id = "${cloudstack_network.foo.id}"
}

resource "cloudstack_static_nat" "foo" {
	ip_address_id = "${cloudstack_ipaddress.foo.id}"
  virtual_machine_id = "${cloudstack_instance.foobar.id}"
}`
