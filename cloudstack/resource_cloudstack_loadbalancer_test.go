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

// TestAccCloudStackLoadBalancer_basic exercises cloudstack_loadbalancer, the
// resource backed by CloudStack's dedicated internal LB API (createLoadBalancer
// with scheme=Internal). Unlike cloudstack_loadbalancer_rule (a public,
// IP-bound LB rule), this is CloudStack's real no-public-IP internal LB
// mechanism, routed through the InternalLbVm provider.
//
// This deliberately does NOT manage the InternalLbVm network_service_provider_state
// itself: that provider is zone-wide, shared state, and every other VPC test
// in this suite already relies on it (and on VPCVirtualRouter) being enabled --
// CloudStack's own "Default VPC offering" maps its Lb service to both
// VPCVirtualRouter and InternalLbVm, so any VPC creation validates InternalLbVm
// is enabled regardless of which offering the *network* underneath uses. A
// per-test resource here would disable it again on teardown (see
// resourceCloudStackNetworkServiceProviderStateDelete) and break every VPC
// test that runs afterwards in the same zone.
func TestAccCloudStackLoadBalancer_basic(t *testing.T) {
	var lb cloudstack.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLoadBalancer_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLoadBalancerExists(
						"cloudstack_loadbalancer.foo", &lb),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer.foo", "name", "terraform-ilb"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer.foo", "algorithm", "roundrobin"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer.foo", "scheme", "Internal"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer.foo", "instanceport", "8080"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer.foo", "sourceport", "8080"),
					resource.TestCheckResourceAttr(
						"cloudstack_loadbalancer.foo", "virtualmachineids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckCloudStackLoadBalancerExists(
	n string, lb *cloudstack.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No load balancer ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		found, _, err := cs.LoadBalancer.GetLoadBalancerByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Load balancer not found")
		}

		*lb = *found

		return nil
	}
}

func testAccCheckCloudStackLoadBalancerDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_loadbalancer" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No load balancer ID is set")
		}

		_, _, err := cs.LoadBalancer.GetLoadBalancerByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Load balancer %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackLoadBalancer_basic = `
resource "cloudstack_vpc" "foo" {
  name         = "terraform-vpc"
  cidr         = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone         = "Sandbox-simulator"
}

resource "cloudstack_network" "foo" {
  name              = "terraform-network"
  display_text      = "terraform-network"
  cidr              = "10.1.1.0/24"
  network_offering  = "DefaultIsolatedNetworkOfferingForVpcNetworksWithInternalLB"
  vpc_id            = cloudstack_vpc.foo.id
  zone              = cloudstack_vpc.foo.zone
}

resource "cloudstack_instance" "foobar1" {
  name             = "terraform-server1"
  display_name     = "terraform"
  service_offering = "Small Instance"
  network_id       = cloudstack_network.foo.id
  template         = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone             = cloudstack_network.foo.zone
  expunge          = true
}

resource "cloudstack_loadbalancer" "foo" {
  name                     = "terraform-ilb"
  algorithm                = "roundrobin"
  instanceport             = 8080
  networkid                = cloudstack_network.foo.id
  scheme                   = "Internal"
  sourceipaddressnetworkid = cloudstack_network.foo.id
  sourceport               = 8080
  virtualmachineids        = [cloudstack_instance.foobar1.id]
}`
