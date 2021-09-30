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

func TestAccCloudStackStaticRoute_basic(t *testing.T) {
	var staticroute cloudstack.StaticRoute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackStaticRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackStaticRoute_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackStaticRouteExists(
						"cloudstack_static_route.foo", &staticroute),
					testAccCheckCloudStackStaticRouteAttributes(&staticroute),
				),
			},
		},
	})
}

func testAccCheckCloudStackStaticRouteExists(
	n string, staticroute *cloudstack.StaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Static Route ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		route, _, err := cs.VPC.GetStaticRouteByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if route.Id != rs.Primary.ID {
			return fmt.Errorf("Static Route not found")
		}

		*staticroute = *route

		return nil
	}
}

func testAccCheckCloudStackStaticRouteAttributes(
	staticroute *cloudstack.StaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if staticroute.Cidr != "172.16.0.0/16" {
			return fmt.Errorf("Bad CIDR: %s", staticroute.Cidr)
		}

		return nil
	}
}

func testAccCheckCloudStackStaticRouteDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_static_route" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static route ID is set")
		}

		staticroute, _, err := cs.VPC.GetStaticRouteByID(rs.Primary.ID)
		if err == nil && staticroute.Id != "" {
			return fmt.Errorf("Static route %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackStaticRoute_basic = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "foo" {
  name = "terraform-acl"
  vpc_id = "${cloudstack_vpc.foo.id}"
}

resource "cloudstack_private_gateway" "foo" {
  gateway = "10.1.1.254"
  ip_address = "192.168.0.1"
  netmask = "255.255.255.0"
  vlan = "1"
  vpc_id = "${cloudstack_vpc.foo.id}"
  acl_id = "${cloudstack_network_acl.foo.id}"
}

resource "cloudstack_static_route" "foo" {
  cidr = "172.16.0.0/16"
  gateway_id = "${cloudstack_private_gateway.foo.id}"
}`
