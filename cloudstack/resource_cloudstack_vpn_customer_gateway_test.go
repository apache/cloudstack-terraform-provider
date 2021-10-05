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
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudStackVPNCustomerGateway_basic(t *testing.T) {
	var vpnCustomerGateway cloudstack.VpnCustomerGateway

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNCustomerGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNCustomerGatewayExists(
						"cloudstack_vpn_customer_gateway.foo", &vpnCustomerGateway),
					testAccCheckCloudStackVPNCustomerGatewayAttributes(&vpnCustomerGateway),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.foo", "name", "terraform-foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.bar", "name", "terraform-bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.bar", "esp_policy", "aes256-sha1"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.foo", "ike_policy", "aes256-sha1;modp1536"),
				),
			},
		},
	})
}

func TestAccCloudStackVPNCustomerGateway_update(t *testing.T) {
	var vpnCustomerGateway cloudstack.VpnCustomerGateway

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNCustomerGateway_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNCustomerGatewayExists(
						"cloudstack_vpn_customer_gateway.foo", &vpnCustomerGateway),
					testAccCheckCloudStackVPNCustomerGatewayAttributes(&vpnCustomerGateway),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.foo", "name", "terraform-foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.bar", "name", "terraform-bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.bar", "esp_policy", "aes256-sha1"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.foo", "ike_policy", "aes256-sha1;modp1536"),
				),
			},

			{
				Config: testAccCloudStackVPNCustomerGateway_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackVPNCustomerGatewayExists(
						"cloudstack_vpn_customer_gateway.foo", &vpnCustomerGateway),
					testAccCheckCloudStackVPNCustomerGatewayUpdatedAttributes(&vpnCustomerGateway),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.foo", "name", "terraform-foo-bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.bar", "name", "terraform-bar-foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.bar", "esp_policy", "3des-md5"),
					resource.TestCheckResourceAttr(
						"cloudstack_vpn_customer_gateway.foo", "ike_policy", "3des-md5;modp1536"),
				),
			},
		},
	})
}

func TestAccCloudStackVPNCustomerGateway_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackVPNCustomerGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackVPNCustomerGateway_basic,
			},

			{
				ResourceName:      "cloudstack_vpn_customer_gateway.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackVPNCustomerGatewayExists(
	n string, vpnCustomerGateway *cloudstack.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN CustomerGateway ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		v, _, err := cs.VPN.GetVpnCustomerGatewayByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if v.Id != rs.Primary.ID {
			return fmt.Errorf("VPN CustomerGateway not found")
		}

		*vpnCustomerGateway = *v

		return nil
	}
}

func testAccCheckCloudStackVPNCustomerGatewayAttributes(
	vpnCustomerGateway *cloudstack.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpnCustomerGateway.Esppolicy != "aes256-sha1" {
			return fmt.Errorf("Bad ESP policy: %s", vpnCustomerGateway.Esppolicy)
		}

		if vpnCustomerGateway.Ikepolicy != "aes256-sha1;modp1536" {
			return fmt.Errorf("Bad IKE policy: %s", vpnCustomerGateway.Ikepolicy)
		}

		if vpnCustomerGateway.Ipsecpsk != "terraform" {
			return fmt.Errorf("Bad IPSEC pre-shared key: %s", vpnCustomerGateway.Ipsecpsk)
		}

		return nil
	}
}

func testAccCheckCloudStackVPNCustomerGatewayUpdatedAttributes(
	vpnCustomerGateway *cloudstack.VpnCustomerGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vpnCustomerGateway.Esppolicy != "3des-md5" {
			return fmt.Errorf("Bad ESP policy: %s", vpnCustomerGateway.Esppolicy)
		}

		if vpnCustomerGateway.Ikepolicy != "3des-md5;modp1536" {
			return fmt.Errorf("Bad IKE policy: %s", vpnCustomerGateway.Ikepolicy)
		}

		if vpnCustomerGateway.Ipsecpsk != "terraform" {
			return fmt.Errorf("Bad IPSEC pre-shared key: %s", vpnCustomerGateway.Ipsecpsk)
		}

		return nil
	}
}

func testAccCheckCloudStackVPNCustomerGatewayDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_vpn_customer_gateway" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN Customer Gateway ID is set")
		}

		_, _, err := cs.VPN.GetVpnCustomerGatewayByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("VPN Customer Gateway %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackVPNCustomerGateway_basic = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.1.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_vpc" "bar" {
  name = "terraform-vpc"
  cidr = "10.2.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_vpn_gateway" "foo" {
	vpc_id = "${cloudstack_vpc.foo.id}"
}

resource "cloudstack_vpn_gateway" "bar" {
	vpc_id = "${cloudstack_vpc.bar.id}"
}

resource "cloudstack_vpn_customer_gateway" "foo" {
	name = "terraform-foo"
	cidr = "${cloudstack_vpc.foo.cidr}"
	esp_policy = "aes256-sha1"
	gateway = "${cloudstack_vpn_gateway.foo.public_ip}"
	ike_policy = "aes256-sha1;modp1536"
	ipsec_psk = "terraform"
}

resource "cloudstack_vpn_customer_gateway" "bar" {
  name = "terraform-bar"
  cidr = "${cloudstack_vpc.bar.cidr}"
  esp_policy = "aes256-sha1"
  gateway = "${cloudstack_vpn_gateway.bar.public_ip}"
  ike_policy = "aes256-sha1;modp1536"
	ipsec_psk = "terraform"
}`

const testAccCloudStackVPNCustomerGateway_update = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.1.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_vpc" "bar" {
  name = "terraform-vpc"
  cidr = "10.2.0.0/16"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_vpn_gateway" "foo" {
  vpc_id = "${cloudstack_vpc.foo.id}"
}

resource "cloudstack_vpn_gateway" "bar" {
  vpc_id = "${cloudstack_vpc.bar.id}"
}

resource "cloudstack_vpn_customer_gateway" "foo" {
  name = "terraform-foo-bar"
  cidr = "${cloudstack_vpc.foo.cidr}"
  esp_policy = "3des-md5"
  gateway = "${cloudstack_vpn_gateway.foo.public_ip}"
  ike_policy = "3des-md5;modp1536"
  ipsec_psk = "terraform"
}

resource "cloudstack_vpn_customer_gateway" "bar" {
  name = "terraform-bar-foo"
  cidr = "${cloudstack_vpc.bar.cidr}"
  esp_policy = "3des-md5"
  gateway = "${cloudstack_vpn_gateway.bar.public_ip}"
  ike_policy = "3des-md5;modp1536"
  ipsec_psk = "terraform"
}`
