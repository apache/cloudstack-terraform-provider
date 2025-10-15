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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackQuotaDataSource_basic(t *testing.T) {
	resourceName := "data.cloudstack_quota.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.#"),
				),
			},
		},
	})
}

func TestAccCloudStackQuotaDataSource_withFilters(t *testing.T) {
	resourceName := "data.cloudstack_quota.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaDataSource_withFilters(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "quotas.#"),
				),
			},
		},
	})
}

func TestAccCloudStackQuotaEnabledDataSource_basic(t *testing.T) {
	resourceName := "data.cloudstack_quota_enabled.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaEnabledDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaEnabledDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
				),
			},
		},
	})
}

func TestAccCloudStackQuotaTariffDataSource_basic(t *testing.T) {
	resourceName := "data.cloudstack_quota_tariff.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaTariffDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaTariffDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "tariffs.#"),
				),
			},
		},
	})
}

func TestAccCloudStackQuotaTariffDataSource_withFilters(t *testing.T) {
	resourceName := "data.cloudstack_quota_tariff.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaTariffDataSource_withFilters(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaTariffDataSourceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "tariffs.#"),
				),
			},
		},
	})
}

// Test check functions
func testAccCheckCloudStackQuotaDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No quota data source ID is set")
		}

		return nil
	}
}

func testAccCheckCloudStackQuotaEnabledDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No quota enabled data source ID is set")
		}

		return nil
	}
}

func testAccCheckCloudStackQuotaTariffDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No quota tariff data source ID is set")
		}

		return nil
	}
}

// Test configuration functions
func testAccCloudStackQuotaDataSource_basic() string {
	return `
data "cloudstack_quota" "test" {
}
`
}

func testAccCloudStackQuotaDataSource_withFilters() string {
	return `
data "cloudstack_quota" "test" {
  account   = "admin"
  domain_id = "ROOT"
}
`
}

func testAccCloudStackQuotaEnabledDataSource_basic() string {
	return `
data "cloudstack_quota_enabled" "test" {
}
`
}

func testAccCloudStackQuotaTariffDataSource_basic() string {
	return `
data "cloudstack_quota_tariff" "test" {
}
`
}

func testAccCloudStackQuotaTariffDataSource_withFilters() string {
	return `
data "cloudstack_quota_tariff" "test" {
  usage_type = 1
}
`
}
