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

func TestAccCloudStackLimitsDataSource_basic(t *testing.T) {
	resourceName := "data.cloudstack_limits.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimitsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "limits.#"),
				),
			},
		},
	})
}

func TestAccCloudStackLimitsDataSource_withType(t *testing.T) {
	resourceName := "data.cloudstack_limits.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimitsDataSource_withType(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "limits.#"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
				),
			},
		},
	})
}

func TestAccCloudStackLimitsDataSource_withDomain(t *testing.T) {
	resourceName := "data.cloudstack_limits.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimitsDataSource_domain + testAccCloudStackLimitsDataSource_withDomain(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
				),
			},
		},
	})
}

func TestAccCloudStackLimitsDataSource_withAccount(t *testing.T) {
	resourceName := "data.cloudstack_limits.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimitsDataSource_domain + testAccCloudStackLimitsDataSource_withAccount(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "limits.#"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttrSet(resourceName, "account"),
				),
			},
		},
	})
}

func TestAccCloudStackLimitsDataSource_multipleTypes(t *testing.T) {
	resourceName := "data.cloudstack_limits.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackLimitsDataSource_multipleTypes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackLimitsDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "limits.#"),
					resource.TestCheckResourceAttr(resourceName, "type", "volume"),
				),
			},
		},
	})
}

func testAccCheckCloudStackLimitsDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Limits data source ID is set")
		}

		return nil
	}
}

// Test configurations

const testAccCloudStackLimitsDataSource_domain = `
resource "cloudstack_domain" "test_domain" {
  name = "test-domain-limits-ds"
}

resource "cloudstack_account" "test_account" {
  username     = "test-account-limits-ds"
  password     = "password"
  first_name   = "Test"
  last_name    = "Account"
  email        = "test-account-limits-ds@example.com"
  account_type = 2  # Regular user account type
  role_id      = "4"  # Regular user role
  domain_id    = cloudstack_domain.test_domain.id
}
`

func testAccCloudStackLimitsDataSource_basic() string {
	return `
data "cloudstack_limits" "test" {
}
`
}

func testAccCloudStackLimitsDataSource_withType() string {
	return `
data "cloudstack_limits" "test" {
  type = "instance"
}
`
}

func testAccCloudStackLimitsDataSource_withDomain() string {
	return `
data "cloudstack_limits" "test" {
  type      = "volume"
  domain_id = cloudstack_domain.test_domain.id
}
`
}

func testAccCloudStackLimitsDataSource_withAccount() string {
	return `
data "cloudstack_limits" "test" {
  type      = "snapshot"
  account   = cloudstack_account.test_account.username
  domain_id = cloudstack_domain.test_domain.id
}
`
}

func testAccCloudStackLimitsDataSource_multipleTypes() string {
	return `
data "cloudstack_limits" "test" {
  type = "volume"
}
`
}
