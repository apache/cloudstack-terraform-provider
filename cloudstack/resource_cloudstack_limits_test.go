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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackLimits_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.foo", "type", "instance"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.foo", "max", "10"),
				),
			},
		},
	})
}

func TestAccCloudStackLimits_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.foo", "max", "10"),
				),
			},
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.foo", "max", "20"),
				),
			},
		},
	})
}

func testAccCheckCloudStackLimitsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Limits ID is set")
		}

		return nil
	}
}

func testAccCheckCloudStackLimitsDestroy(s *terraform.State) error {
	return nil
}

const testAccCloudStackLimits_basic = `
resource "cloudstack_limits" "foo" {
  type         = "instance"
  max          = 10
  domainid     = cloudstack_domain.test_domain.id
}
`

const testAccCloudStackLimits_update = `
resource "cloudstack_limits" "foo" {
  type         = "instance"
  max          = 20
  domainid     = cloudstack_domain.test_domain.id
}
`

func TestAccCloudStackLimits_domain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_domain_limit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.domain_limit"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.domain_limit", "type", "volume"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.domain_limit", "max", "50"),
				),
			},
		},
	})
}

func TestAccCloudStackLimits_account(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_account,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.account_limit"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.account_limit", "type", "snapshot"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.account_limit", "max", "100"),
				),
			},
		},
	})
}

func TestAccCloudStackLimits_project(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_project,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.project_limit"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.project_limit", "type", "primarystorage"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.project_limit", "max", "1000"),
				),
			},
		},
	})
}

func TestAccCloudStackLimits_unlimited(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_unlimited,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.unlimited"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.unlimited", "type", "cpu"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.unlimited", "max", "-1"),
				),
			},
		},
	})
}

func TestAccCloudStackLimits_stringType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_stringType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.string_type"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.string_type", "type", "network"),
					resource.TestCheckResourceAttr(
						"cloudstack_limits.string_type", "max", "30"),
				),
			},
		},
	})
}

func TestAccCloudStackLimits_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackLimitsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimits_domain + testAccCloudStackLimits_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsExists("cloudstack_limits.foo"),
				),
			},
			{
				ResourceName:            "cloudstack_limits.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domainid", "type", "type", "max"},
			},
		},
	})
}

// Test configurations for different resource types
const testAccCloudStackLimits_domain = `
resource "cloudstack_domain" "test_domain" {
  name = "test-domain-limits"
}
`

const testAccCloudStackLimits_domain_limit = `
resource "cloudstack_limits" "domain_limit" {
  type 		   = "volume"
  max          = 50
  domainid     = cloudstack_domain.test_domain.id
}
`

const testAccCloudStackLimits_account = `
resource "cloudstack_account" "test_account" {
  username     = "test-account-limits"
  password     = "password"
  first_name   = "Test"
  last_name    = "Account"
  email        = "test-account-limits@example.com"
  account_type = 2  # Regular user account type
  role_id      = 4  # Regular user role
  domainid     = cloudstack_domain.test_domain.id
}

resource "cloudstack_limits" "account_limit" {
  type         = "snapshot"
  max          = 100
  account      = cloudstack_account.test_account.username
  domainid     = cloudstack_domain.test_domain.id
}
`

const testAccCloudStackLimits_project = `
resource "cloudstack_limits" "project_limit" {
  type         = "primarystorage"
  max          = 1000
  domainid     = cloudstack_domain.test_domain.id
}
`

const testAccCloudStackLimits_unlimited = `
resource "cloudstack_limits" "unlimited" {
  type         = "cpu"
  max          = -1
  domainid     = cloudstack_domain.test_domain.id
}
`

const testAccCloudStackLimits_stringType = `
resource "cloudstack_limits" "string_type" {
  type         = "network"
  max          = 30
  domainid     = cloudstack_domain.test_domain.id
}
`
