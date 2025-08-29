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
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackPortForward_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackPortForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackPortForward_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackPortForwardsExist("cloudstack_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.#", "1"),
					resource.TestCheckResourceAttrSet(
						"cloudstack_port_forward.foo", "forward.0.uuid"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.0.private_port", "443"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.0.public_port", "8443"),
					resource.TestCheckResourceAttrSet(
						"cloudstack_port_forward.foo", "forward.0.virtual_machine_id"),
				),
			},
		},
	})
}

func TestAccCloudStackPortForward_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackPortForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackPortForward_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackPortForwardsExist("cloudstack_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.#", "1"),
				),
			},

			{
				Config: testAccCloudStackPortForward_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackPortForwardsExist("cloudstack_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.#", "2"),
					// Validate first forward rule
					resource.TestCheckResourceAttrSet(
						"cloudstack_port_forward.foo", "forward.0.uuid"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.0.protocol", "tcp"),
					resource.TestCheckResourceAttrSet(
						"cloudstack_port_forward.foo", "forward.0.virtual_machine_id"),
					// Validate second forward rule
					resource.TestCheckResourceAttrSet(
						"cloudstack_port_forward.foo", "forward.1.uuid"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.1.protocol", "tcp"),
					resource.TestCheckResourceAttrSet(
						"cloudstack_port_forward.foo", "forward.1.virtual_machine_id"),
				),
			},
		},
	})
}

func TestAccCloudStackPortForward_portRange(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackPortForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackPortForward_portRange,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackPortForwardsExist("cloudstack_port_forward.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_port_forward.foo", "forward.#", "2"),
					testAccCheckCloudStackPortForwardAttributes("cloudstack_port_forward.foo"),
					// Note: We don't check specific indices since sets are unordered
					// The testAccCheckCloudStackPortForwardAttributes function handles validation
				),
			},
		},
	})
}

func testAccCheckCloudStackPortForwardsExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No port forward ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
			_, count, err := cs.Firewall.GetPortForwardingRuleByID(id)

			if err != nil {
				return err
			}

			if count == 0 {
				return fmt.Errorf("Port forward for %s not found", k)
			}
		}

		return nil
	}
}

func testAccCheckCloudStackPortForwardAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No port forward ID is set")
		}

		// Verify we have 2 forward rules
		if rs.Primary.Attributes["forward.#"] != "2" {
			return fmt.Errorf("Expected 2 forward rules, got %s", rs.Primary.Attributes["forward.#"])
		}

		var foundTCPRange, foundUDPSingle bool

		// Check both forward rules to find the expected configurations
		for i := 0; i < 2; i++ {
			protocolKey := fmt.Sprintf("forward.%d.protocol", i)
			privatePortKey := fmt.Sprintf("forward.%d.private_port", i)
			privateEndPortKey := fmt.Sprintf("forward.%d.private_end_port", i)
			publicPortKey := fmt.Sprintf("forward.%d.public_port", i)
			publicEndPortKey := fmt.Sprintf("forward.%d.public_end_port", i)
			uuidKey := fmt.Sprintf("forward.%d.uuid", i)

			protocol := rs.Primary.Attributes[protocolKey]
			privatePort := rs.Primary.Attributes[privatePortKey]
			privateEndPort := rs.Primary.Attributes[privateEndPortKey]
			publicPort := rs.Primary.Attributes[publicPortKey]
			publicEndPort := rs.Primary.Attributes[publicEndPortKey]
			uuid := rs.Primary.Attributes[uuidKey]

			// Verify basic required fields exist
			if protocol == "" {
				return fmt.Errorf("Missing protocol for forward rule %d", i)
			}
			if privatePort == "" {
				return fmt.Errorf("Missing private_port for forward rule %d", i)
			}
			if publicPort == "" {
				return fmt.Errorf("Missing public_port for forward rule %d", i)
			}
			if uuid == "" {
				return fmt.Errorf("Missing uuid for forward rule %d", i)
			}

			// Check for TCP rule with port range (8080-8090)
			if protocol == "tcp" && privatePort == "8080" && publicPort == "8080" {
				if privateEndPort != "8090" {
					return fmt.Errorf("Expected TCP rule to have private_end_port=8090, got %s", privateEndPort)
				}
				if publicEndPort != "8090" {
					return fmt.Errorf("Expected TCP rule to have public_end_port=8090, got %s", publicEndPort)
				}
				foundTCPRange = true
			}

			// Check for UDP rule with single port (9000)
			if protocol == "udp" && privatePort == "9000" && publicPort == "9000" {
				// For single port rules, end ports should be empty, "0", or equal to start ports
				if privateEndPort != "" && privateEndPort != "0" && privateEndPort != "9000" {
					return fmt.Errorf("Expected UDP rule to have empty, '0', or matching private_end_port, got %s", privateEndPort)
				}
				if publicEndPort != "" && publicEndPort != "0" && publicEndPort != "9000" {
					return fmt.Errorf("Expected UDP rule to have empty, '0', or matching public_end_port, got %s", publicEndPort)
				}
				foundUDPSingle = true
			}
		}

		if !foundTCPRange {
			return fmt.Errorf("Expected to find TCP rule with port range 8080-8090")
		}
		if !foundUDPSingle {
			return fmt.Errorf("Expected to find UDP rule with single port 9000")
		}

		return nil
	}
}

func testAccCheckCloudStackPortForwardDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_port_forward" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No port forward ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			_, _, err := cs.Firewall.GetPortForwardingRuleByID(id)
			if err == nil {
				return fmt.Errorf("Port forward %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccCloudStackPortForward_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-updated"
  service_offering= "Medium Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_port_forward" "foo" {
  ip_address_id = cloudstack_network.foo.source_nat_ip_id

  forward {
    protocol = "tcp"
    private_port = 443
    public_port = 8443
    virtual_machine_id = cloudstack_instance.foobar.id
  }
}`

const testAccCloudStackPortForward_update = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-updated"
  service_offering= "Medium Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_port_forward" "foo" {
  ip_address_id = cloudstack_network.foo.source_nat_ip_id

  forward {
    protocol = "tcp"
    private_port = 443
    public_port = 8443
    virtual_machine_id = cloudstack_instance.foobar.id
  }

  forward {
    protocol = "tcp"
    private_port = 80
    public_port = 8080
    virtual_machine_id = cloudstack_instance.foobar.id
  }
}`

const testAccCloudStackPortForward_portRange = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-updated"
  service_offering= "Medium Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_port_forward" "foo" {
  ip_address_id = cloudstack_network.foo.source_nat_ip_id

  forward {
    protocol = "tcp"
    private_port = 8080
    private_end_port = 8090
    public_port = 8080
    public_end_port = 8090
    virtual_machine_id = cloudstack_instance.foobar.id
  }

  forward {
    protocol = "udp"
    private_port = 9000
    public_port = 9000
    virtual_machine_id = cloudstack_instance.foobar.id
  }
}`
