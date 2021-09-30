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
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
	source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-updated"
  service_offering= "Medium Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_port_forward" "foo" {
  ip_address_id = "${cloudstack_network.foo.source_nat_ip_id}"

  forward {
    protocol = "tcp"
    private_port = 443
    public_port = 8443
    virtual_machine_id = "${cloudstack_instance.foobar.id}"
  }
}`

const testAccCloudStackPortForward_update = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
	source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-updated"
  service_offering= "Medium Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "Sandbox-simulator"
  expunge = true
}

resource "cloudstack_port_forward" "foo" {
  ip_address_id = "${cloudstack_network.foo.source_nat_ip_id}"

  forward {
    protocol = "tcp"
    private_port = 443
    public_port = 8443
    virtual_machine_id = "${cloudstack_instance.foobar.id}"
  }

  forward {
    protocol = "tcp"
    private_port = 80
    public_port = 8080
    virtual_machine_id = "${cloudstack_instance.foobar.id}"
  }
}`
