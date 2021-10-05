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

func TestAccCloudStackVPNConnection_basic(t *testing.T) {
	var vpnConnection cloudstack.VpnConnection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNConnection_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNConnectionExists(
						"cloudstack_vpn_connection.foo-bar", &vpnConnection),
					testAccCheckCloudStackVPNConnectionExists(
						"cloudstack_vpn_connection.bar-foo", &vpnConnection),
				),
			},
		},
	})
}

func testAccCheckCloudStackVPNConnectionExists(
	n string, vpnConnection *cloudstack.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Connection ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		v, _, err := cs.VPN.GetVpnConnectionByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPN Connection not found")
		}

		*vpnConnection = *v

		return nil
	}
}

func testAccCheckCloudStackVPNConnectionDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_vpn_connection" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Connection ID is set")
		}

		_, _, err := cs.VPN.GetVpnConnectionByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN Connection %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackVPNConnection_basic = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.1.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_vpc" "bar" {
  name = "terraform-vpc"
  cidr = "10.2.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_vpn_gateway" "foo" {
  vpc_id = "${cloudstack_vpc.foo.id}"
}

resource "cloudstack_vpn_gateway" "bar" {
  vpc_id = "${cloudstack_vpc.bar.id}"
}

resource "cloudstack_vpn_customer_gateway" "foo" {
  name = "terraform-foo"
  cidr = "${cloudstack_vpc.foo.cidr}"
  esp_policy = "aes256-sha1"
  gateway = "${cloudstack_vpn_gateway.foo.public_ip}"
  ike_policy = "aes256-sha1;modp1536"
  ipsec_psk = "terraform"
}

resource "cloudstack_vpn_customer_gateway" "bar" {
  name = "terraform-bar"
  cidr = "${cloudstack_vpc.bar.cidr}"
  esp_policy = "aes256-sha1"
  gateway = "${cloudstack_vpn_gateway.bar.public_ip}"
  ike_policy = "aes256-sha1;modp1536"
  ipsec_psk = "terraform"
}

resource "cloudstack_vpn_connection" "foo-bar" {
  customer_gateway_id = "${cloudstack_vpn_customer_gateway.foo.id}"
  vpn_gateway_id = "${cloudstack_vpn_gateway.bar.id}"
}

resource "cloudstack_vpn_connection" "bar-foo" {
  customer_gateway_id = "${cloudstack_vpn_customer_gateway.bar.id}"
  vpn_gateway_id = "${cloudstack_vpn_gateway.foo.id}"
}`
