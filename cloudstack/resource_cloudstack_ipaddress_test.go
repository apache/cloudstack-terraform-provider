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

func TestAccCloudStackIPAddress_basic(t *testing.T) {
	var ipaddr cloudstack.PublicIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackIPAddress_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackIPAddressExists(
						"cloudstack_ipaddress.foo", &ipaddr),
					testAccCheckResourceTags(&ipaddr),
				),
			},
		},
	})
}

func TestAccCloudStackIPAddress_vpc(t *testing.T) {
	var ipaddr cloudstack.PublicIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackIPAddress_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackIPAddressExists(
						"cloudstack_ipaddress.foo", &ipaddr),
				),
			},
		},
	})
}

func testAccCheckCloudStackIPAddressExists(
	n string, ipaddr *cloudstack.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		pip, _, err := cs.Address.GetPublicIpAddressByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if pip.Id != rs.Primary.ID {
			return fmt.Errorf("IP address not found")
		}

		*ipaddr = *pip

		return nil
	}
}

func testAccCheckCloudStackIPAddressDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_ipaddress" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		ip, _, err := cs.Address.GetPublicIpAddressByID(rs.Primary.ID)
		if err == nil && ip.Associatednetworkid != "" {
			return fmt.Errorf("Public IP %s still associated", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackIPAddress_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_ipaddress" "foo" {
  network_id = "${cloudstack_network.foo.id}"
  tags = {
    terraform-tag = "true"
  }
}`

const testAccCloudStackIPAddress_vpc = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_ipaddress" "foo" {
  vpc_id = "${cloudstack_vpc.foo.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}`
