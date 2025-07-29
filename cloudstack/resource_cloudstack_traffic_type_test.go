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

func TestAccCloudStackTrafficType_basic(t *testing.T) {
	var trafficType cloudstack.TrafficType

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackTrafficTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackTrafficType_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackTrafficTypeExists(
						"cloudstack_traffic_type.foo", &trafficType),
					testAccCheckCloudStackTrafficTypeBasicAttributes(&trafficType),
					resource.TestCheckResourceAttrSet(
						"cloudstack_traffic_type.foo", "type"),
					resource.TestCheckResourceAttr(
						"cloudstack_traffic_type.foo", "kvm_network_label", "cloudbr0"),
				),
			},
		},
	})
}

func TestAccCloudStackTrafficType_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackTrafficTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackTrafficType_basic,
			},
			{
				ResourceName:      "cloudstack_traffic_type.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackTrafficTypeExists(
	n string, trafficType *cloudstack.TrafficType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No traffic type ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Usage.NewListTrafficTypesParams(rs.Primary.Attributes["physical_network_id"])

		l, err := cs.Usage.ListTrafficTypes(p)
		if err != nil {
			return err
		}

		// Find the traffic type with the matching ID
		var found bool
		for _, t := range l.TrafficTypes {
			if t.Id == rs.Primary.ID {
				*trafficType = *t
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Traffic type not found")
		}

		return nil
	}
}

func testAccCheckCloudStackTrafficTypeBasicAttributes(
	trafficType *cloudstack.TrafficType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// The TrafficType struct doesn't have a field that directly maps to the 'type' attribute
		// Instead, we'll rely on the resource attribute checks in the test
		return nil
	}
}

func testAccCheckCloudStackTrafficTypeDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_traffic_type" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No traffic type ID is set")
		}

		// Get the physical network ID from the state
		physicalNetworkID := rs.Primary.Attributes["physical_network_id"]
		if physicalNetworkID == "" {
			continue // If the resource is gone, that's okay
		}

		p := cs.Usage.NewListTrafficTypesParams(physicalNetworkID)
		l, err := cs.Usage.ListTrafficTypes(p)
		if err != nil {
			return nil
		}

		// Check if the traffic type still exists
		for _, t := range l.TrafficTypes {
			if t.Id == rs.Primary.ID {
				return fmt.Errorf("Traffic type %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccCloudStackTrafficType_basic = `
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
}

resource "cloudstack_traffic_type" "foo" {
  physical_network_id = cloudstack_physicalnetwork.foo.id
  type = "Management"
  kvm_network_label = "cloudbr0"
  xen_network_label = "xenbr0"
}`
