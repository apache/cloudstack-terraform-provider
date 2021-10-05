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

func TestAccCloudStackNetworkACL_basic(t *testing.T) {
	var acl cloudstack.NetworkACLList
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACL_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLExists(
						"cloudstack_network_acl.foo", &acl),
					testAccCheckCloudStackNetworkACLBasicAttributes(&acl),
				),
			},
		},
	})
}

func TestAccCloudStackNetworkACL_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACL_basic,
			},

			{
				ResourceName:      "cloudstack_network_acl.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackNetworkACLExists(
	n string, acl *cloudstack.NetworkACLList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		acllist, _, err := cs.NetworkACL.GetNetworkACLListByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if acllist.Id != rs.Primary.ID {
			return fmt.Errorf("Network ACL not found")
		}

		*acl = *acllist

		return nil
	}
}

func testAccCheckCloudStackNetworkACLBasicAttributes(
	acl *cloudstack.NetworkACLList) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if acl.Name != "terraform-acl" {
			return fmt.Errorf("Bad name: %s", acl.Name)
		}

		if acl.Description != "terraform-acl-text" {
			return fmt.Errorf("Bad description: %s", acl.Description)
		}

		return nil
	}
}

func testAccCheckCloudStackNetworkACLDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_network_acl" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ID is set")
		}

		_, _, err := cs.NetworkACL.GetNetworkACLListByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network ACl list %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackNetworkACL_basic = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "foo" {
  name = "terraform-acl"
  description = "terraform-acl-text"
  vpc_id = "${cloudstack_vpc.foo.id}"
}`
