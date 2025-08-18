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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackPhysicalNetwork_basic(t *testing.T) {
	var physicalNetwork cloudstack.PhysicalNetwork

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackPhysicalNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackPhysicalNetwork_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackPhysicalNetworkExists(
						"cloudstack_physicalnetwork.foo", &physicalNetwork),
					testAccCheckCloudStackPhysicalNetworkBasicAttributes(&physicalNetwork),
				),
			},
		},
	})
}

func TestAccCloudStackPhysicalNetwork_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackPhysicalNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackPhysicalNetwork_basic,
			},
			{
				ResourceName:      "cloudstack_physicalnetwork.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackPhysicalNetworkExists(
	n string, physicalNetwork *cloudstack.PhysicalNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No physical network ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p, _, err := cs.Network.GetPhysicalNetworkByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if p.Id != rs.Primary.ID {
			return fmt.Errorf("Physical network not found")
		}

		*physicalNetwork = *p

		return nil
	}
}

func testAccCheckCloudStackPhysicalNetworkBasicAttributes(
	physicalNetwork *cloudstack.PhysicalNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if physicalNetwork.Name != "terraform-physical-network" {
			return fmt.Errorf("Bad name: %s", physicalNetwork.Name)
		}

		if physicalNetwork.Broadcastdomainrange != "ZONE" {
			return fmt.Errorf("Bad broadcast domain range: %s", physicalNetwork.Broadcastdomainrange)
		}

		return nil
	}
}

func testAccCheckCloudStackPhysicalNetworkDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_physicalnetwork" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No physical network ID is set")
		}

		_, _, err := cs.Network.GetPhysicalNetworkByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Physical network %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackPhysicalNetwork_basic = `
resource "cloudstack_zone" "foo" {
  name = "terraform-zone"
  dns1 = "8.8.8.8"
  internal_dns1 = "8.8.4.4"
  network_type = "Advanced"
}

resource "cloudstack_physicalnetwork" "foo" {
  name = "terraform-physical-network"
  zone = cloudstack_zone.foo.name
  broadcast_domain_range = "ZONE"
  isolation_methods = ["VLAN"]
}`
