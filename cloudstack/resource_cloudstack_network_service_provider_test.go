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

func TestAccCloudStackNetworkServiceProvider_basic(t *testing.T) {
	var provider cloudstack.NetworkServiceProvider

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkServiceProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkServiceProvider_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkServiceProviderExists(
						"cloudstack_network_service_provider.foo", &provider),
					testAccCheckCloudStackNetworkServiceProviderBasicAttributes(&provider),
					resource.TestCheckResourceAttr(
						"cloudstack_network_service_provider.foo", "name", "VirtualRouter"),
				),
			},
		},
	})
}

func TestAccCloudStackNetworkServiceProvider_securityGroup(t *testing.T) {
	var provider cloudstack.NetworkServiceProvider

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkServiceProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkServiceProvider_securityGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkServiceProviderExists(
						"cloudstack_network_service_provider.security_group", &provider),
					testAccCheckCloudStackNetworkServiceProviderSecurityGroupAttributes(&provider),
					resource.TestCheckResourceAttr(
						"cloudstack_network_service_provider.security_group", "name", "SecurityGroupProvider"),
				),
			},
		},
	})
}

func TestAccCloudStackNetworkServiceProvider_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkServiceProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkServiceProvider_basic,
			},
			{
				ResourceName:      "cloudstack_network_service_provider.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackNetworkServiceProviderExists(
	n string, provider *cloudstack.NetworkServiceProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network service provider ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Network.NewListNetworkServiceProvidersParams()
		p.SetPhysicalnetworkid(rs.Primary.Attributes["physical_network_id"])

		l, err := cs.Network.ListNetworkServiceProviders(p)
		if err != nil {
			return err
		}

		// Find the network service provider with the matching ID
		var found bool
		for _, p := range l.NetworkServiceProviders {
			if p.Id == rs.Primary.ID {
				*provider = *p
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Network service provider not found")
		}

		return nil
	}
}

func testAccCheckCloudStackNetworkServiceProviderBasicAttributes(
	provider *cloudstack.NetworkServiceProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if provider.Name != "VirtualRouter" {
			return fmt.Errorf("Bad name: %s", provider.Name)
		}

		// We don't check the state for VirtualRouter as it requires configuration first

		return nil
	}
}

func testAccCheckCloudStackNetworkServiceProviderSecurityGroupAttributes(
	provider *cloudstack.NetworkServiceProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if provider.Name != "SecurityGroupProvider" {
			return fmt.Errorf("Bad name: %s", provider.Name)
		}

		// We don't check the service list for SecurityGroupProvider as it's predefined
		// and can't be modified

		return nil
	}
}

func testAccCheckCloudStackNetworkServiceProviderDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_network_service_provider" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network service provider ID is set")
		}

		// Get the physical network ID from the state
		physicalNetworkID := rs.Primary.Attributes["physical_network_id"]
		if physicalNetworkID == "" {
			continue // If the resource is gone, that's okay
		}

		p := cs.Network.NewListNetworkServiceProvidersParams()
		p.SetPhysicalnetworkid(physicalNetworkID)
		l, err := cs.Network.ListNetworkServiceProviders(p)
		if err != nil {
			return nil
		}

		// Check if the network service provider still exists
		for _, p := range l.NetworkServiceProviders {
			if p.Id == rs.Primary.ID {
				return fmt.Errorf("Network service provider %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccCloudStackNetworkServiceProvider_basic = `
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

resource "cloudstack_network_service_provider" "foo" {
  name = "VirtualRouter"
  physical_network_id = cloudstack_physicalnetwork.foo.id
  service_list = ["Dhcp", "Dns"]
  # Note: We don't set state for VirtualRouter as it requires configuration first
}`

const testAccCloudStackNetworkServiceProvider_securityGroup = `
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

resource "cloudstack_network_service_provider" "security_group" {
  name = "SecurityGroupProvider"
  physical_network_id = cloudstack_physicalnetwork.foo.id
  # Note: We don't set service_list for SecurityGroupProvider as it doesn't support updating
}`
