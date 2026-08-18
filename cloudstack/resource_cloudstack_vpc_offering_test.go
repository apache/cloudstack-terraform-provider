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

func TestAccCloudStackVPCOffering_basic(t *testing.T) {
	var vpcOffering cloudstack.VPCOffering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPCOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPCOffering_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPCOfferingExists(
						"cloudstack_vpc_offering.foo", &vpcOffering),
					testAccCheckCloudStackVPCOfferingAttributes(&vpcOffering),
					resource.TestCheckResourceAttr(
						"cloudstack_vpc_offering.foo", "enable", "true"),
				),
			},
		},
	})
}

func TestAccCloudStackVPCOffering_serviceCapabilityList(t *testing.T) {
	var vpcOffering cloudstack.VPCOffering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPCOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPCOffering_serviceCapabilityList,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPCOfferingExists(
						"cloudstack_vpc_offering.capability", &vpcOffering),
					resource.TestCheckResourceAttr(
						"cloudstack_vpc_offering.capability", "service_capability_list.0.service", "SourceNat"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpc_offering.capability", "service_capability_list.0.capability_type", "RedundantRouter"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpc_offering.capability", "service_capability_list.0.capability_value", "true"),
				),
			},
		},
	})
}

func TestAccCloudStackVPCOffering_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPCOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPCOffering_basic,
			},

			{
				ResourceName:      "cloudstack_vpc_offering.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackVPCOfferingExists(
	n string, vpcOffering *cloudstack.VPCOffering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC offering ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		o, _, err := cs.VPC.GetVPCOfferingByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if o.Id != rs.Primary.ID {
			return fmt.Errorf("VPC offering not found")
		}

		*vpcOffering = *o

		return nil
	}
}

func testAccCheckCloudStackVPCOfferingAttributes(
	vpcOffering *cloudstack.VPCOffering) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpcOffering.Name != "terraform-vpc-offering" {
			return fmt.Errorf("Bad name: %s", vpcOffering.Name)
		}

		if vpcOffering.Displaytext != "terraform-vpc-offering-text" {
			return fmt.Errorf("Bad display text: %s", vpcOffering.Displaytext)
		}

		return nil
	}
}

func testAccCheckCloudStackVPCOfferingDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_vpc_offering" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC offering ID is set")
		}

		_, _, err := cs.VPC.GetVPCOfferingByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPC offering %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackVPCOffering_basic = `
resource "cloudstack_vpc_offering" "foo" {
  name         = "terraform-vpc-offering"
  display_text = "terraform-vpc-offering-text"
  enable       = true
  supported_services = ["Dhcp", "Dns", "SourceNat", "PortForwarding", "Lb", "UserData", "StaticNat", "NetworkACL"]
  service_provider_list = {
    Dhcp           = "VpcVirtualRouter"
    Dns            = "VpcVirtualRouter"
    SourceNat      = "VpcVirtualRouter"
    PortForwarding = "VpcVirtualRouter"
    Lb             = "VpcVirtualRouter"
    UserData       = "VpcVirtualRouter"
    StaticNat      = "VpcVirtualRouter"
    NetworkACL     = "VpcVirtualRouter"
  }
}`

const testAccCloudStackVPCOffering_serviceCapabilityList = `
resource "cloudstack_vpc_offering" "capability" {
  name         = "terraform-vpc-offering-capability"
  display_text = "terraform-vpc-offering-capability-text"
  // CloudStack always injects SourceNat and NetworkACL support into a VPC
  // offering even if omitted here, so they must be declared explicitly to
  // avoid permanent post-apply drift on these ForceNew attributes.
  supported_services = ["SourceNat", "NetworkACL"]
  service_provider_list = {
    SourceNat  = "VpcVirtualRouter"
    NetworkACL = "VpcVirtualRouter"
  }
  service_capability_list {
    service          = "SourceNat"
    capability_type  = "RedundantRouter"
    capability_value = "true"
  }
}`
