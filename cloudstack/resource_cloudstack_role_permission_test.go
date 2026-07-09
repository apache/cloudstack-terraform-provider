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

func TestAccCloudStackRolePermission_basic(t *testing.T) {
	var rolePermission cloudstack.RolePermission

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackRolePermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackRolePermission_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo", &rolePermission),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "rule", "listVirtualMachines"),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "permission", "allow"),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "description", "terraform test role permission"),
				),
			},
			{
				Config: testAccCloudStackRolePermission_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo", &rolePermission),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "permission", "deny"),
				),
			},
		},
	})
}

func testAccCheckCloudStackRolePermissionExists(n string, rolePermission *cloudstack.RolePermission) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role Permission ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

		p := cs.Role.NewListRolePermissionsParams()
		p.SetRoleid(rs.Primary.Attributes["role_id"])

		l, err := cs.Role.ListRolePermissions(p)
		if err != nil {
			return err
		}

		for _, rp := range l.RolePermissions {
			if rp.Id == rs.Primary.ID {
				*rolePermission = *rp
				return nil
			}
		}

		return fmt.Errorf("Role Permission not found")
	}
}

func testAccCheckCloudStackRolePermissionDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_role_permission" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role Permission ID is set")
		}

		p := cs.Role.NewListRolePermissionsParams()
		p.SetRoleid(rs.Primary.Attributes["role_id"])

		l, err := cs.Role.ListRolePermissions(p)
		if err != nil {
			// If the parent role is already gone, the permission is too.
			continue
		}

		for _, rp := range l.RolePermissions {
			if rp.Id == rs.Primary.ID {
				return fmt.Errorf("Role Permission %s still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

const testAccCloudStackRolePermission_basic = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id     = cloudstack_role.foo.id
  rule        = "listVirtualMachines"
  permission  = "allow"
  description = "terraform test role permission"
}
`

const testAccCloudStackRolePermission_update = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id     = cloudstack_role.foo.id
  rule        = "listVirtualMachines"
  permission  = "deny"
  description = "terraform test role permission"
}
`
