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

func TestAccCloudStackVPNGateway_basic(t *testing.T) {
	var vpnGateway cloudstack.VpnGateway

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNGatewayExists(
						"cloudstack_vpn_gateway.foo", &vpnGateway),
				),
			},
		},
	})
}

func TestAccCloudStackVPNGateway_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNGateway_basic,
			},

			{
				ResourceName:      "cloudstack_vpn_gateway.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackVPNGatewayExists(
	n string, vpnGateway *cloudstack.VpnGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Gateway ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		v, _, err := cs.VPN.GetVpnGatewayByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPN Gateway not found")
		}

		*vpnGateway = *v

		return nil
	}
}

func testAccCheckCloudStackVPNGatewayDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_vpn_gateway" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Gateway ID is set")
		}

		_, _, err := cs.VPN.GetVpnGatewayByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN Gateway %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackVPNGateway_basic = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_vpn_gateway" "foo" {
  vpc_id = "${cloudstack_vpc.foo.id}"
}`
