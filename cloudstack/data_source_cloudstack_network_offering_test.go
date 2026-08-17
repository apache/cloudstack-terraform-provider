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

func TestAccNetworkOfferingDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_network_offering.net-off-resource"
	datasourceName := "data.cloudstack_network_offering.net-off-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkOfferingDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "display_text", resourceName, "display_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "guest_ip_type", resourceName, "guest_ip_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "traffic_type", resourceName, "traffic_type"),
				),
			},
		},
	})
}

func TestAccNetworkOfferingDataSource_withAdditionalParams(t *testing.T) {
	resourceName := "cloudstack_network_offering.net-off-resource"
	datasourceName := "data.cloudstack_network_offering.net-off-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkOfferingDataSourceConfig_withAdditionalParams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "display_text", resourceName, "display_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "guest_ip_type", resourceName, "guest_ip_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "traffic_type", resourceName, "traffic_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "network_rate", resourceName, "network_rate"),
					resource.TestCheckResourceAttrPair(datasourceName, "conserve_mode", resourceName, "conserve_mode"),
					resource.TestCheckResourceAttrPair(datasourceName, "for_vpc", resourceName, "for_vpc"),
					resource.TestCheckResourceAttrPair(datasourceName, "specify_vlan", resourceName, "specify_vlan"),
					resource.TestCheckResourceAttrPair(datasourceName, "enable", resourceName, "enable"),
				),
			},
		},
	})
}

func TestAccNetworkOfferingDataSource_withServices(t *testing.T) {
	resourceName := "cloudstack_network_offering.net-off-resource"
	datasourceName := "data.cloudstack_network_offering.net-off-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkOfferingDataSourceConfig_withServices,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "supported_services.#", resourceName, "supported_services.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "service_provider_list.%", resourceName, "service_provider_list.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "service_provider_list.Dhcp", resourceName, "service_provider_list.Dhcp"),
					resource.TestCheckResourceAttrPair(datasourceName, "service_provider_list.Dns", resourceName, "service_provider_list.Dns"),
					resource.TestCheckResourceAttrPair(datasourceName, "enable", resourceName, "enable"),
					resource.TestCheckResourceAttrPair(datasourceName, "max_connections", resourceName, "max_connections"),
				),
			},
		},
	})
}

func TestAccNetworkOfferingDataSource_allOptionalParams(t *testing.T) {
	resourceName := "cloudstack_network_offering.net-off-resource"
	datasourceName := "data.cloudstack_network_offering.net-off-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkOfferingDataSourceConfig_allOptionalParams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "display_text", resourceName, "display_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "for_nsx", resourceName, "for_nsx"),
					resource.TestCheckResourceAttrPair(datasourceName, "specify_as_number", resourceName, "specify_as_number"),
					resource.TestCheckResourceAttrPair(datasourceName, "internet_protocol", resourceName, "internet_protocol"),
					resource.TestCheckResourceAttr(datasourceName, "enable", "true"),
					resource.TestCheckResourceAttr(datasourceName, "for_nsx", "false"),
				),
			},
		},
	})
}

const testNetworkOfferingDataSourceConfig_basic = `
resource "cloudstack_network_offering" "net-off-resource"{
  name       = "TestNetworkDisplay01"
  display_text = "TestNetworkDisplay01"
  guest_ip_type = "Isolated"
  traffic_type = "Guest"
  }

  data "cloudstack_network_offering" "net-off-data-source"{

    filter{
    name = "name"
    value="TestNetworkDisplay01"
    }
	  depends_on = [
	  cloudstack_network_offering.net-off-resource
	]
  }
  `

const testNetworkOfferingDataSourceConfig_withAdditionalParams = `
resource "cloudstack_network_offering" "net-off-resource"{
  name              = "TestNetworkDisplayAdvanced01"
  display_text      = "TestNetworkDisplayAdvanced01"
  guest_ip_type     = "Isolated"
  traffic_type      = "Guest"
  network_rate      = 100
  conserve_mode     = true
  enable            = true
  for_vpc           = false
  specify_vlan      = true
  supported_services = ["Dhcp", "Dns", "Firewall", "Lb", "SourceNat"]
  service_provider_list = {
    Dhcp      = "VirtualRouter"
    Dns       = "VirtualRouter"
    Firewall  = "VirtualRouter"
    Lb        = "VirtualRouter"
    SourceNat = "VirtualRouter"
  }
}

data "cloudstack_network_offering" "net-off-data-source"{
  filter{
    name = "name"
    value = "TestNetworkDisplayAdvanced01"
  }
  depends_on = [
    cloudstack_network_offering.net-off-resource
  ]
}
`

const testNetworkOfferingDataSourceConfig_withServices = `
resource "cloudstack_network_offering" "net-off-resource"{
  name              = "TestNetworkServices01"
  display_text      = "TestNetworkServices01"
  guest_ip_type     = "Isolated"
  traffic_type      = "Guest"
  enable            = true
  supported_services = ["Dhcp", "Dns", "Firewall", "Lb", "SourceNat", "StaticNat", "PortForwarding"]
  service_provider_list = {
    Dhcp           = "VirtualRouter"
    Dns            = "VirtualRouter"
    Firewall       = "VirtualRouter"
    Lb             = "VirtualRouter"
    SourceNat      = "VirtualRouter"
    StaticNat      = "VirtualRouter"
    PortForwarding = "VirtualRouter"
  }
}

data "cloudstack_network_offering" "net-off-data-source"{
  filter{
    name = "name"
    value = "TestNetworkServices01"
  }
  depends_on = [
    cloudstack_network_offering.net-off-resource
  ]
}
`

const testNetworkOfferingDataSourceConfig_allOptionalParams = `
resource "cloudstack_network_offering" "net-off-resource"{
  name              = "TestNetworkDisplayAll01"
  display_text      = "TestNetworkDisplayAll01"
  guest_ip_type     = "Isolated"
  traffic_type      = "Guest"
  network_rate      = 200
  conserve_mode     = true
  enable            = true
  for_vpc           = false
  for_nsx           = false
  specify_vlan      = true
  specify_as_number = false
  internet_protocol = "IPv4"
  max_connections   = 1000
  supported_services = ["Dhcp", "Dns", "Firewall", "Lb", "SourceNat"]
  service_provider_list = {
    Dhcp      = "VirtualRouter"
    Dns       = "VirtualRouter"
    Firewall  = "VirtualRouter"
    Lb        = "VirtualRouter"
    SourceNat = "VirtualRouter"
  }
}

data "cloudstack_network_offering" "net-off-data-source"{
  filter{
    name = "name"
    value = "TestNetworkDisplayAll01"
  }
  depends_on = [
    cloudstack_network_offering.net-off-resource
  ]
}
`
