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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackRolePermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackRolePermission_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"listVirtualMachines"}),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "permission.0.rule", "listVirtualMachines"),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "permission.0.permission", "allow"),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "permission.0.description", "terraform test role permission"),
				),
			},
			{
				Config: testAccCloudStackRolePermission_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"listVirtualMachines"}),
					resource.TestCheckResourceAttr(
						"cloudstack_role_permission.foo", "permission.0.permission", "deny"),
				),
			},
		},
	})
}

func TestAccCloudStackRolePermission_orderAfterRecreate(t *testing.T) {
	var wildcardRuleID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackRolePermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackRolePermission_orderWithSpecificRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"listZones", "*"}),
					testAccCaptureCloudStackRolePermissionID("cloudstack_role_permission.foo", "*", &wildcardRuleID),
				),
			},
			{
				Config: testAccCloudStackRolePermission_orderWithoutSpecificRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"*"}),
					testAccCheckCloudStackRolePermissionID("cloudstack_role_permission.foo", "*", &wildcardRuleID),
				),
			},
			{
				Config:   testAccCloudStackRolePermission_orderWithoutSpecificRule,
				PlanOnly: true,
			},
			{
				Config: testAccCloudStackRolePermission_orderWithSpecificRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"listZones", "*"}),
					testAccCheckCloudStackRolePermissionID("cloudstack_role_permission.foo", "*", &wildcardRuleID),
				),
			},
			{
				Config:   testAccCloudStackRolePermission_orderWithSpecificRule,
				PlanOnly: true,
			},
		},
	})
}

func TestMatchCloudStackRolePermissions_shiftedIDsAfterRemoval(t *testing.T) {
	rolePermissions := []*cloudstack.RolePermission{
		{Id: "detach-id", Rule: "detachIso", Permission: "allow"},
		{Id: "start-id", Rule: "startSystemVm", Permission: "allow"},
		{Id: "wildcard-id", Rule: "*", Permission: "deny"},
	}
	desiredPermissions := []rolePermissionSpec{
		{ID: "removed-attach-id", Rule: "detachIso", Permission: "allow"},
		{ID: "detach-id", Rule: "startSystemVm", Permission: "allow"},
		{ID: "start-id", Rule: "*", Permission: "deny"},
	}

	matchedPermissions := matchCloudStackRolePermissions(rolePermissions, desiredPermissions)
	assertRolePermissionIDs(t, matchedPermissions, []string{"detach-id", "start-id", "wildcard-id"})
}

func TestMatchCloudStackRolePermissions_shiftedIDsAfterInsertion(t *testing.T) {
	rolePermissions := []*cloudstack.RolePermission{
		{Id: "detach-id", Rule: "detachIso", Permission: "allow"},
		{Id: "start-id", Rule: "startSystemVm", Permission: "allow"},
	}
	desiredPermissions := []rolePermissionSpec{
		{ID: "detach-id", Rule: "attachIso", Permission: "allow"},
		{ID: "start-id", Rule: "detachIso", Permission: "allow"},
		{Rule: "startSystemVm", Permission: "allow"},
	}

	matchedPermissions := matchCloudStackRolePermissions(rolePermissions, desiredPermissions)
	assertRolePermissionIDs(t, matchedPermissions, []string{"", "detach-id", "start-id"})
}

func assertRolePermissionIDs(t *testing.T, rolePermissions []*cloudstack.RolePermission, expectedIDs []string) {
	t.Helper()

	if len(rolePermissions) != len(expectedIDs) {
		t.Fatalf("Expected %d matched Role Permissions, got %d", len(expectedIDs), len(rolePermissions))
	}

	for i, expectedID := range expectedIDs {
		var actualID string
		if rolePermissions[i] != nil {
			actualID = rolePermissions[i].Id
		}
		if actualID != expectedID {
			t.Errorf("Expected matched Role Permission %d to have ID %q, got %q", i, expectedID, actualID)
		}
	}
}

func TestAccCloudStackRolePermission_orderChange(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackRolePermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackRolePermission_orderWithSpecificRule,
				// The check deliberately changes the API-side order to verify drift detection.
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"listZones", "*"}),
					testAccReorderCloudStackRolePermissions("cloudstack_role_permission.foo", []string{"*", "listZones"}),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"*", "listZones"}),
				),
			},
			{
				Config: testAccCloudStackRolePermission_orderWithSpecificRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"listZones", "*"}),
				),
			},
			{
				Config: testAccCloudStackRolePermission_orderReversed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"*", "listZones"}),
				),
			},
		},
	})
}

func TestAccCloudStackRolePermission_orderChangeThreeRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackRolePermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackRolePermission_orderThreeRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"listZones", "listVirtualMachines", "*"}),
				),
			},
			{
				Config: testAccCloudStackRolePermission_orderThreeRulesReordered,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionOrder("cloudstack_role_permission.foo", []string{"*", "listZones", "listVirtualMachines"}),
				),
			},
		},
	})
}

func TestAccCloudStackRolePermission_authoritative(t *testing.T) {
	var externalRuleID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackRolePermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackRolePermission_subset,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCreateCloudStackRolePermission("cloudstack_role_permission.foo", "listVirtualMachines", "allow", "external role permission", &externalRuleID),
				),
			},
			{
				Config: testAccCloudStackRolePermission_subset,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionRuleExists("cloudstack_role_permission.foo", &externalRuleID),
				),
			},
			{
				Config: testAccCloudStackRolePermission_authoritative,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackRolePermissionExists("cloudstack_role_permission.foo"),
					testAccCheckCloudStackRolePermissionRuleMissing("cloudstack_role_permission.foo", &externalRuleID),
				),
			},
		},
	})
}

func testAccCheckCloudStackRolePermissionExists(n string) resource.TestCheckFunc {
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

		if _, err := cs.Role.ListRolePermissions(p); err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckCloudStackRolePermissionOrder(n string, rules []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rolePermissions, err := testAccListCloudStackRolePermissions(s, n)
		if err != nil {
			return err
		}
		if len(rolePermissions) < len(rules) {
			return fmt.Errorf("Expected at least %d Role Permissions, got %d", len(rules), len(rolePermissions))
		}

		for i, rule := range rules {
			if rolePermissions[i].Rule != rule {
				return fmt.Errorf("Expected Role Permission rule %d to be %q, got %q", i, rule, rolePermissions[i].Rule)
			}
		}

		return nil
	}
}

func testAccReorderCloudStackRolePermissions(n string, rules []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rolePermissions, err := testAccListCloudStackRolePermissions(s, n)
		if err != nil {
			return err
		}

		ruleIDs := make([]string, 0, len(rules))
		for _, rule := range rules {
			var ruleID string
			for _, rp := range rolePermissions {
				if rp.Rule == rule {
					ruleID = rp.Id
					break
				}
			}
			if ruleID == "" {
				return fmt.Errorf("Role Permission rule %q not found", rule)
			}
			ruleIDs = append(ruleIDs, ruleID)
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Role.NewUpdateRolePermissionParams(s.RootModule().Resources[n].Primary.Attributes["role_id"])
		p.SetRuleorder(ruleIDs)
		if _, err := cs.Role.UpdateRolePermission(p); err != nil {
			return fmt.Errorf("Error ordering Role Permissions: %s", err)
		}

		return nil
	}
}

func testAccCaptureCloudStackRolePermissionID(n, rule string, ruleID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rolePermissions, err := testAccListCloudStackRolePermissions(s, n)
		if err != nil {
			return err
		}

		for _, rp := range rolePermissions {
			if rp.Rule == rule {
				*ruleID = rp.Id
				return nil
			}
		}

		return fmt.Errorf("Role Permission rule %q not found", rule)
	}
}

func testAccCheckCloudStackRolePermissionID(n, rule string, expectedRuleID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *expectedRuleID == "" {
			return fmt.Errorf("No expected Role Permission ID is set for rule %q", rule)
		}

		rolePermissions, err := testAccListCloudStackRolePermissions(s, n)
		if err != nil {
			return err
		}

		for _, rp := range rolePermissions {
			if rp.Rule == rule {
				if rp.Id != *expectedRuleID {
					return fmt.Errorf("Expected Role Permission rule %q to keep ID %s, got %s", rule, *expectedRuleID, rp.Id)
				}
				return nil
			}
		}

		return fmt.Errorf("Role Permission rule %q not found", rule)
	}
}

func testAccCreateCloudStackRolePermission(n, rule, permission, description string, ruleID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Role.NewCreateRolePermissionParams(permission, rs.Primary.Attributes["role_id"], rule)
		p.SetDescription(description)

		r, err := cs.Role.CreateRolePermission(p)
		if err != nil {
			return err
		}

		*ruleID = r.Id
		return nil
	}
}

func testAccCheckCloudStackRolePermissionRuleExists(n string, ruleID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *ruleID == "" {
			return fmt.Errorf("No external Role Permission ID is set")
		}

		rolePermissions, err := testAccListCloudStackRolePermissions(s, n)
		if err != nil {
			return err
		}

		for _, rp := range rolePermissions {
			if rp.Id == *ruleID {
				return nil
			}
		}

		return fmt.Errorf("Role Permission %s not found", *ruleID)
	}
}

func testAccCheckCloudStackRolePermissionRuleMissing(n string, ruleID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *ruleID == "" {
			return fmt.Errorf("No external Role Permission ID is set")
		}

		rolePermissions, err := testAccListCloudStackRolePermissions(s, n)
		if err != nil {
			return err
		}

		for _, rp := range rolePermissions {
			if rp.Id == *ruleID {
				return fmt.Errorf("Role Permission %s still exists", *ruleID)
			}
		}

		return nil
	}
}

func testAccListCloudStackRolePermissions(s *terraform.State, n string) ([]*cloudstack.RolePermission, error) {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", n)
	}

	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
	p := cs.Role.NewListRolePermissionsParams()
	p.SetRoleid(rs.Primary.Attributes["role_id"])

	l, err := cs.Role.ListRolePermissions(p)
	if err != nil {
		return nil, err
	}

	return l.RolePermissions, nil
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

		if _, err := cs.Role.ListRolePermissions(p); err != nil {
			// If the parent role is already gone, the permissions are too.
			continue
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
  role_id = cloudstack_role.foo.id

  permission {
    rule        = "listVirtualMachines"
    permission  = "allow"
    description = "terraform test role permission"
  }
}
`

const testAccCloudStackRolePermission_update = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id = cloudstack_role.foo.id

  permission {
    rule        = "listVirtualMachines"
    permission  = "deny"
    description = "terraform test role permission"
  }
}
`

const testAccCloudStackRolePermission_orderWithSpecificRule = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id = cloudstack_role.foo.id

  permission {
    rule       = "listZones"
    permission = "allow"
  }

  permission {
    rule       = "*"
    permission = "deny"
  }
}
`

const testAccCloudStackRolePermission_orderWithoutSpecificRule = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id = cloudstack_role.foo.id

  permission {
    rule       = "*"
    permission = "deny"
  }
}
`

const testAccCloudStackRolePermission_orderReversed = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id = cloudstack_role.foo.id

  permission {
    rule       = "*"
    permission = "deny"
  }

  permission {
    rule       = "listZones"
    permission = "allow"
  }
}
`

const testAccCloudStackRolePermission_orderThreeRules = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id = cloudstack_role.foo.id

  permission {
    rule       = "listZones"
    permission = "allow"
  }

  permission {
    rule       = "listVirtualMachines"
    permission = "allow"
  }

  permission {
    rule       = "*"
    permission = "deny"
  }
}
`

const testAccCloudStackRolePermission_orderThreeRulesReordered = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id = cloudstack_role.foo.id

  permission {
    rule       = "*"
    permission = "deny"
  }

  permission {
    rule       = "listZones"
    permission = "allow"
  }

  permission {
    rule       = "listVirtualMachines"
    permission = "allow"
  }
}
`

const testAccCloudStackRolePermission_subset = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id = cloudstack_role.foo.id

  permission {
    rule       = "listZones"
    permission = "allow"
  }
}
`

const testAccCloudStackRolePermission_authoritative = `
resource "cloudstack_role" "foo" {
  name = "terraform-role"
  type = "User"
}

resource "cloudstack_role_permission" "foo" {
  role_id       = cloudstack_role.foo.id
  authoritative = true

  permission {
    rule       = "listZones"
    permission = "allow"
  }
}
`
