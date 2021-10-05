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

func TestAccCloudStackLoadBalancerRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_update(t *testing.T) {
	var id string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			{
				Config: testAccCloudStackLoadBalancerRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", &id),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_forceNew(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			{
				Config: testAccCloudStackLoadBalancerRule_forcenew,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "443"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "443"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "protocol", "tcp-proxy"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_vpc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "80"),
				),
			},
		},
	})
}

func TestAccCloudStackLoadBalancerRule_vpcUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancerRule_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "80"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "80"),
				),
			},

			{
				Config: testAccCloudStackLoadBalancerRule_vpc_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerRuleExist("cloudstack_loadbalancer_rule.foo", nil),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "name", "terraform-lb-update"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "algorithm", "leastconn"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "public_port", "443"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer_rule.foo", "private_port", "443"),
				),
			},
		},
	})
}

func testAccCheckCloudStackLoadBalancerRuleExist(n string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No loadbalancer rule ID is set")
		}

		if id != nil {
			if *id != "" && *id != rs.Primary.ID {
				return fmt.Errorf("Resource ID has changed!")
			}

			*id = rs.Primary.ID
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		_, count, err := cs.LoadBalancer.GetLoadBalancerRuleByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Loadbalancer rule %s not found", n)
		}

		return nil
	}
}

func testAccCheckCloudStackLoadBalancerRuleDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_loadbalancer_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Loadbalancer rule ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			if !strings.Contains(k, "uuid") {
				continue
			}

			_, _, err := cs.LoadBalancer.GetLoadBalancerRuleByID(id)
			if err == nil {
				return fmt.Errorf("Loadbalancer rule %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccCloudStackLoadBalancerRule_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_ipaddress" "foo" {
  network_id = "${cloudstack_network.foo.id}"
}

resource "cloudstack_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_loadbalancer_rule" "foo" {
  name = "terraform-lb"
  ip_address_id = "${cloudstack_ipaddress.foo.id}"
  algorithm = "roundrobin"
  public_port = 80
  private_port = 80
  member_ids = ["${cloudstack_instance.foobar1.id}"]
}`

const testAccCloudStackLoadBalancerRule_update = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_ipaddress" "foo" {
  network_id = "${cloudstack_network.foo.id}"
}

resource "cloudstack_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = "${cloudstack_ipaddress.foo.id}"
  algorithm = "leastconn"
  public_port = 80
  private_port = 80
  member_ids = ["${cloudstack_instance.foobar1.id}"]
}`

const testAccCloudStackLoadBalancerRule_forcenew = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_ipaddress" "foo" {
  network_id = "${cloudstack_network.foo.id}"
}

resource "cloudstack_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = "${cloudstack_ipaddress.foo.id}"
  algorithm = "leastconn"
  public_port = 443
  private_port = 443
	protocol = "tcp-proxy"
  member_ids = ["${cloudstack_instance.foobar1.id}"]
}`

const testAccCloudStackLoadBalancerRule_vpc = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id = "${cloudstack_vpc.foo.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}

resource "cloudstack_ipaddress" "foo" {
  vpc_id = "${cloudstack_vpc.foo.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}

resource "cloudstack_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_loadbalancer_rule" "foo" {
  name = "terraform-lb"
  ip_address_id = "${cloudstack_ipaddress.foo.id}"
  algorithm = "roundrobin"
  network_id = "${cloudstack_network.foo.id}"
  public_port = 80
  private_port = 80
  member_ids = ["${cloudstack_instance.foobar1.id}"]
}`

const testAccCloudStackLoadBalancerRule_vpc_update = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id = "${cloudstack_vpc.foo.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}

resource "cloudstack_ipaddress" "foo" {
  vpc_id = "${cloudstack_vpc.foo.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}

resource "cloudstack_instance" "foobar1" {
  name = "terraform-server1"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_instance" "foobar2" {
  name = "terraform-server2"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_loadbalancer_rule" "foo" {
  name = "terraform-lb-update"
  ip_address_id = "${cloudstack_ipaddress.foo.id}"
  algorithm = "leastconn"
  network_id = "${cloudstack_network.foo.id}"
  public_port = 443
  private_port = 443
  member_ids = ["${cloudstack_instance.foobar1.id}", "${cloudstack_instance.foobar2.id}"]
}`
