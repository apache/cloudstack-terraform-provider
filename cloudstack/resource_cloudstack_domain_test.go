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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudStackDomain_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDomain_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDomainExists("cloudstack_domain.test", nil),
					resource.TestCheckResourceAttr("cloudstack_domain.test", "name", "terraform-test-domain"),
					resource.TestCheckResourceAttr("cloudstack_domain.test", "network_domain", "terraform.test"),
					resource.TestCheckResourceAttrSet("cloudstack_domain.test", "level"),
					resource.TestCheckResourceAttrSet("cloudstack_domain.test", "path"),
				),
			},
		},
	})
}

func TestAccCloudStackDomain_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDomain_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDomainExists("cloudstack_domain.test", nil),
					resource.TestCheckResourceAttr("cloudstack_domain.test", "name", "terraform-test-domain"),
					resource.TestCheckResourceAttr("cloudstack_domain.test", "network_domain", "terraform.test"),
				),
			},
			{
				Config: testAccCloudStackDomain_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDomainExists("cloudstack_domain.test", nil),
					resource.TestCheckResourceAttr("cloudstack_domain.test", "name", "terraform-test-domain-updated"),
					resource.TestCheckResourceAttr("cloudstack_domain.test", "network_domain", "terraform-updated.test"),
				),
			},
		},
	})
}

func testAccCheckCloudStackDomainExists(n string, domain *cloudstack.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No domain ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		d, _, err := cs.Domain.GetDomainByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if d.Id != rs.Primary.ID {
			return fmt.Errorf("Domain not found")
		}

		if domain != nil {
			*domain = *d
		}

		return nil
	}
}

func testAccCheckCloudStackDomainDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_domain" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No domain ID is set")
		}

		_, _, err := cs.Domain.GetDomainByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Domain %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackDomain_basic = `
resource "cloudstack_domain" "test" {
  name           = "terraform-test-domain"
  network_domain = "terraform.test"
}`

const testAccCloudStackDomain_update = `
resource "cloudstack_domain" "test" {
  name           = "terraform-test-domain-updated"
  network_domain = "terraform-updated.test"
}`
