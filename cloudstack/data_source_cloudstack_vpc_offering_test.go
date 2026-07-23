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

func TestAccVPCOfferingDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_vpc_offering.vpc-off-resource"
	datasourceName := "data.cloudstack_vpc_offering.vpc-off-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVPCOfferingDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "display_text", resourceName, "display_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "enable", resourceName, "enable"),
				),
			},
		},
	})
}

func TestAccVPCOfferingDataSource_withServices(t *testing.T) {
	resourceName := "cloudstack_vpc_offering.vpc-off-resource"
	datasourceName := "data.cloudstack_vpc_offering.vpc-off-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVPCOfferingDataSourceConfig_withServices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "supported_services.#", resourceName, "supported_services.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "service_provider_list.%", resourceName, "service_provider_list.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "service_provider_list.Dhcp", resourceName, "service_provider_list.Dhcp"),
					resource.TestCheckResourceAttrPair(datasourceName, "enable", resourceName, "enable"),
				),
			},
		},
	})
}

const testVPCOfferingDataSourceConfig_basic = `
resource "cloudstack_vpc_offering" "vpc-off-resource"{
  name         = "TestVPCOfferingDisplay01"
  display_text = "TestVPCOfferingDisplay01"
  enable       = true
  // CloudStack always injects SourceNat and NetworkACL support into a VPC
  // offering even if omitted here, so they must be declared explicitly to
  // avoid permanent post-apply drift on these ForceNew attributes.
  supported_services = ["SourceNat", "NetworkACL"]
  service_provider_list = {
    SourceNat  = "VpcVirtualRouter"
    NetworkACL = "VpcVirtualRouter"
  }
}

data "cloudstack_vpc_offering" "vpc-off-data-source"{
  filter{
    name = "name"
    value = "TestVPCOfferingDisplay01"
  }
  depends_on = [
    cloudstack_vpc_offering.vpc-off-resource
  ]
}
`

const testVPCOfferingDataSourceConfig_withServices = `
resource "cloudstack_vpc_offering" "vpc-off-resource"{
  name         = "TestVPCOfferingServices01"
  display_text = "TestVPCOfferingServices01"
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
}

data "cloudstack_vpc_offering" "vpc-off-data-source"{
  filter{
    name = "name"
    value = "TestVPCOfferingServices01"
  }
  depends_on = [
    cloudstack_vpc_offering.vpc-off-resource
  ]
}
`
