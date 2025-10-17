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
	"regexp"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackQuotaTariff_basic(t *testing.T) {
	var quotaTariff cloudstack.QuotaTariffList

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckQuotaSupport(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackQuotaTariffDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaTariff_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaTariffExists("cloudstack_quota_tariff.test", &quotaTariff),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "name", "Test CPU Tariff"),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "usage_type", "1"),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "value", "0.05"),
				),
			},
		},
	})
}

func TestAccCloudStackQuotaTariff_update(t *testing.T) {
	var quotaTariff cloudstack.QuotaTariffList

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckQuotaSupport(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackQuotaTariffDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaTariff_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaTariffExists("cloudstack_quota_tariff.test", &quotaTariff),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "name", "Test CPU Tariff"),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "value", "0.05"),
				),
			},
			{
				Config: testAccCloudStackQuotaTariff_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaTariffExists("cloudstack_quota_tariff.test", &quotaTariff),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "name", "Updated CPU Tariff"),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "value", "0.10"),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccCloudStackQuotaTariff_import(t *testing.T) {
	resourceName := "cloudstack_quota_tariff.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckQuotaSupport(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackQuotaTariffDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaTariff_basic(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackQuotaTariffExists(n string, quotaTariff *cloudstack.QuotaTariffList) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No quota tariff ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Quota.NewQuotaTariffListParams()
		p.SetId(rs.Primary.ID)

		r, err := cs.Quota.QuotaTariffList(p)
		if err != nil {
			return err
		}

		if len(r.QuotaTariffList) == 0 {
			return fmt.Errorf("Quota tariff not found")
		}

		*quotaTariff = *r.QuotaTariffList[0]
		return nil
	}
}

func testAccCheckCloudStackQuotaTariffDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_quota_tariff" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No quota tariff ID is set")
		}

		p := cs.Quota.NewQuotaTariffListParams()
		p.SetId(rs.Primary.ID)

		r, err := cs.Quota.QuotaTariffList(p)
		if err == nil && len(r.QuotaTariffList) > 0 && r.QuotaTariffList[0].Removed == "" {
			return fmt.Errorf("Quota tariff %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCloudStackQuotaTariff_basic() string {
	return `
resource "cloudstack_quota_tariff" "test" {
  name        = "Test CPU Tariff"
  usage_type  = 1
  value       = 0.05
  description = "Test tariff for CPU usage"
}
`
}

func testAccCloudStackQuotaTariff_update() string {
	return `
resource "cloudstack_quota_tariff" "test" {
  name        = "Updated CPU Tariff"
  usage_type  = 1
  value       = 0.10
  description = "Updated description"
  start_date  = "2026-01-01"
  end_date    = "2026-12-31"
}
`
}

// Test validation errors
func TestAccCloudStackQuotaTariff_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckQuotaSupport(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudStackQuotaTariff_invalidUsageType(),
				ExpectError: regexp.MustCompile(`"usage_type" must be between 1 and 25`),
			},
			{
				Config:      testAccCloudStackQuotaTariff_negativeValue(),
				ExpectError: regexp.MustCompile(`"value" cannot be negative`),
			},
		},
	})
}

// Test activation rules
func TestAccCloudStackQuotaTariff_activationRules(t *testing.T) {
	var quotaTariff cloudstack.QuotaTariffList

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckQuotaSupport(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackQuotaTariffDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaTariff_withActivationRule(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaTariffExists("cloudstack_quota_tariff.test", &quotaTariff),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "activation_rule", "serviceOffering.id == 'test-so-123'"),
				),
			},
		},
	})
}

// Test complex activation rules
func TestAccCloudStackQuotaTariff_complexActivationRules(t *testing.T) {
	var quotaTariff cloudstack.QuotaTariffList

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckQuotaSupport(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackQuotaTariffDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackQuotaTariff_complexActivationRule(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackQuotaTariffExists("cloudstack_quota_tariff.test", &quotaTariff),
					resource.TestCheckResourceAttr("cloudstack_quota_tariff.test", "activation_rule", "serviceOffering.id == 'test-so-123' && zone.name == 'premium-zone'"),
				),
			},
		},
	})
}

func testAccCloudStackQuotaTariff_invalidUsageType() string {
	return `
resource "cloudstack_quota_tariff" "test" {
  name        = "Invalid Usage Type Test"
  usage_type  = 999
  value       = 0.05
  description = "This should fail"
}
`
}

func testAccCloudStackQuotaTariff_negativeValue() string {
	return `
resource "cloudstack_quota_tariff" "test" {
  name        = "Negative Value Test"
  usage_type  = 1
  value       = -0.05
  description = "This should fail"
}
`
}

func testAccCloudStackQuotaTariff_withActivationRule() string {
	return `
resource "cloudstack_quota_tariff" "test" {
  name            = "Test Activation Rule Tariff"
  usage_type      = 1
  value           = 0.15
  description     = "Test tariff with activation rule"
  activation_rule = "serviceOffering.id == 'test-so-123'"
}
`
}

func testAccCloudStackQuotaTariff_complexActivationRule() string {
	return `
resource "cloudstack_quota_tariff" "test" {
  name            = "Complex Activation Rule Tariff"
  usage_type      = 1
  value           = 0.25
  description     = "Test tariff with complex activation rule"
  activation_rule = "serviceOffering.id == 'test-so-123' && zone.name == 'premium-zone'"
}
`
}
