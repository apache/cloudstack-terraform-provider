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

func TestAccCloudStackRole_basic(t *testing.T) {
	var role cloudstack.Role

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackRole_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRoleExists("cloudstack_role.foo", &role),
					resource.TestCheckResourceAttr(
						"cloudstack_role.foo", "name", "terraform-role"),
					resource.TestCheckResourceAttr(
						"cloudstack_role.foo", "description", "terraform test role"),
					resource.TestCheckResourceAttr(
						"cloudstack_role.foo", "is_public", "true"),
				),
			},
		},
	})
}

func testAccCheckCloudStackRoleExists(n string, role *cloudstack.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		r, _, err := cs.Role.GetRoleByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if r.Id != rs.Primary.ID {
			return fmt.Errorf("Role not found")
		}

		*role = *r

		return nil
	}
}

func testAccCheckCloudStackRoleDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_role" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role ID is set")
		}

		// Use a defer/recover to catch the panic that might occur when trying to access l.Roles[0]
		var err error
		func() {
			defer func() {
				if r := recover(); r != nil {
					// If a panic occurs, it means the role doesn't exist, which is what we want
					err = nil
				}
			}()
			r, _, e := cs.Role.GetRoleByID(rs.Primary.ID)
			if e == nil && r != nil && r.Id == rs.Primary.ID {
				err = fmt.Errorf("Role %s still exists", rs.Primary.ID)
			}
		}()

		if err != nil {
			return err
		}
	}

	return nil
}

const testAccCloudStackRole_basic = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  description = "terraform test role"
  is_public = true
  type = "User"
}
`
