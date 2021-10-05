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

func TestAccCloudStackSecurityGroup_basic(t *testing.T) {
	var sg cloudstack.SecurityGroup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackSecurityGroup_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackSecurityGroupExists(
						"cloudstack_security_group.foo", &sg),
					testAccCheckCloudStackSecurityGroupBasicAttributes(&sg),
				),
			},
		},
	})
}

func TestAccCloudStackSecurityGroup_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackSecurityGroup_basic,
			},

			{
				ResourceName:      "cloudstack_security_group.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackSecurityGroupExists(
	n string, sg *cloudstack.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No security group ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		resp, _, err := cs.SecurityGroup.GetSecurityGroupByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp.Id != rs.Primary.ID {
			return fmt.Errorf("Network ACL not found")
		}

		*sg = *resp

		return nil
	}
}

func testAccCheckCloudStackSecurityGroupBasicAttributes(
	sg *cloudstack.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if sg.Name != "terraform-security-group" {
			return fmt.Errorf("Bad name: %s", sg.Name)
		}

		if sg.Description != "terraform-security-group-text" {
			return fmt.Errorf("Bad description: %s", sg.Description)
		}

		return nil
	}
}

func testAccCheckCloudStackSecurityGroupDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_security_group" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No security group ID is set")
		}

		_, _, err := cs.SecurityGroup.GetSecurityGroupByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Security group list %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackSecurityGroup_basic = `
resource "cloudstack_security_group" "foo" {
  name = "terraform-security-group"
	description = "terraform-security-group-text"
}`
