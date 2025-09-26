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

func TestAccCloudStackEgressFirewall_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackEgressFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.cidr_list.0", "10.1.1.10/32"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.ports.0", "8080"),
				),
			},
		},
	})
}

func TestAccCloudStackEgressFirewall_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackEgressFirewall_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.cidr_list.0", "10.1.1.10/32"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.ports.0", "8080"),
				),
			},

			{
				Config: testAccCloudStackEgressFirewall_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.#", "2"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.cidr_list.0", "10.1.1.10/32"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.cidr_list.1", "10.1.1.11/32"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.ports.0", "8080"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1.cidr_list.0", "10.1.1.10/32"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1.protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.1.ports.1", "80"),
				),
			},
		},
	})
}

func testAccCheckCloudStackEgressFirewallRulesExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No firewall ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.%") {
				continue
			}

			cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
			_, count, err := cs.Firewall.GetEgressFirewallRuleByID(id)

			if err != nil {
				return err
			}

			if count == 0 {
				return fmt.Errorf("Firewall rule for %s not found", k)
			}
		}

		return nil
	}
}

func testAccCheckCloudStackEgressFirewallDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_egress_firewall" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, ".uuids.") || strings.HasSuffix(k, ".uuids.%") {
				continue
			}

			_, _, err := cs.Firewall.GetEgressFirewallRuleByID(id)
			if err == nil {
				return fmt.Errorf("Egress rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccCloudStackEgressFirewall_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "cloudstack_egress_firewall" "foo" {
  network_id = cloudstack_network.foo.id

  rule {
    cidr_list = ["10.1.1.10/32"]
    protocol = "tcp"
    ports = ["8080"]
  }
}`

const testAccCloudStackEgressFirewall_update = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "cloudstack_egress_firewall" "foo" {
  network_id = cloudstack_network.foo.id

  rule {
    cidr_list = ["10.1.1.10/32", "10.1.1.11/32"]
    protocol = "tcp"
    ports = ["8080"]
  }

  rule {
    cidr_list = ["10.1.1.10/32"]
    protocol = "tcp"
    ports = ["80", "1000-2000"]
  }
}`

func TestAccCloudStackEgressFirewall_allPortsTCP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackEgressFirewall_allPorts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.cidr_list.0", "10.1.1.10/32"),
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.protocol", "tcp"),
					// No ports should be set when omitting the ports parameter
					resource.TestCheckResourceAttr(
						"cloudstack_egress_firewall.foo", "rule.0.ports.#", "0"),
				),
			},
		},
	})
}

const testAccCloudStackEgressFirewall_allPorts = `
resource "cloudstack_network" "foo" {
  name = "terraform-network-tcp"
  display_text = "terraform-network-tcp"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
}

resource "cloudstack_egress_firewall" "foo" {
  network_id = cloudstack_network.foo.id

  rule {
    cidr_list = ["10.1.1.10/32"]
    protocol  = "tcp"
    # No ports specified - should create a rule for all ports
  }
}`

func TestAccCloudStackEgressFirewall_allPortsUDP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackEgressFirewall_allPortsUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.foo", "rule.0.protocol", "udp"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.foo", "rule.0.ports.#", "0"),
				),
			},
		},
	})
}

const testAccCloudStackEgressFirewall_allPortsUDP = `
resource "cloudstack_network" "foo" {
  name             = "tf-egress-udp-all"
  display_text     = "tf-egress-udp-all"
  cidr             = "10.8.0.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone             = "Sandbox-simulator"
}

resource "cloudstack_egress_firewall" "foo" {
  network_id = cloudstack_network.foo.id

  rule {
    cidr_list = ["10.8.0.10/32"]
    protocol  = "udp"
    # no ports => all ports
  }
}`

func TestAccCloudStackEgressFirewall_allPortsCombined(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackEgressFirewall_allPortsCombined,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.mixed"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.mixed", "rule.#", "2"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.mixed", "rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.mixed", "rule.0.ports.#", "2"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.mixed", "rule.1.protocol", "udp"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.mixed", "rule.1.ports.#", "0"),
				),
			},
		},
	})
}

const testAccCloudStackEgressFirewall_allPortsCombined = `
resource "cloudstack_network" "foo" {
  name             = "terraform-network-mixed"
  display_text     = "terraform-network-mixed"
  cidr             = "10.1.3.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone             = "Sandbox-simulator"
}

resource "cloudstack_egress_firewall" "mixed" {
  network_id = cloudstack_network.foo.id

  rule {
    cidr_list = [cloudstack_network.foo.cidr]
    protocol  = "tcp"
    ports     = ["80", "443"]
  }

  rule {
    cidr_list = ["${cidrhost(cloudstack_network.foo.cidr, 10)}"]
    protocol  = "udp"
    # no ports => all ports
  }
}`

func TestAccCloudStackEgressFirewall_portsToAllPorts(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackEgressFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackEgressFirewall_specificPorts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.foo", "rule.0.ports.#", "2"),
				),
			},
			{
				Config: testAccCloudStackEgressFirewall_allPortsTransition,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackEgressFirewallRulesExist("cloudstack_egress_firewall.foo"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr("cloudstack_egress_firewall.foo", "rule.0.ports.#", "0"),
				),
			},
		},
	})
}

const testAccCloudStackEgressFirewall_specificPorts = `
resource "cloudstack_network" "foo" {
  name             = "terraform-network-transition"
  display_text     = "terraform-network-transition"
  cidr             = "10.1.4.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone             = "Sandbox-simulator"
}

resource "cloudstack_egress_firewall" "foo" {
  network_id = cloudstack_network.foo.id

  rule {
    cidr_list = ["${cidrhost(cloudstack_network.foo.cidr, 10)}"]
    protocol  = "tcp"
    ports     = ["80", "1000-2000"]
  }
}
`

const testAccCloudStackEgressFirewall_allPortsTransition = `
resource "cloudstack_network" "foo" {
  name             = "terraform-network-transition"
  display_text     = "terraform-network-transition"
  cidr             = "10.1.4.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone             = "Sandbox-simulator"
}

resource "cloudstack_egress_firewall" "foo" {
  network_id = cloudstack_network.foo.id

  rule {
    cidr_list = ["${cidrhost(cloudstack_network.foo.cidr, 10)}"]
    protocol  = "tcp"
    # no ports => all ports
  }
}
`
