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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudstackDomainDataSource_basic(t *testing.T) {
	resourceName := "data.cloudstack_domain.my_domain"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudstackDomainDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudstackDomainDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "ROOT"),
				),
			},
		},
	})
}

func testAccCheckCloudstackDomainDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain ID is set")
		}

		return nil
	}
}

func TestAccCloudstackDomainDataSource_invalidName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudstackDomainDataSource_invalidName(),
				ExpectError: regexp.MustCompile("no domain found with name: badgerbearocto"),
			},
		},
	})
}

func testAccCloudstackDomainDataSource_basic() string {
	return `
data "cloudstack_domain" "my_domain" {
	 filter {
	   name = "name"
	   value = "ROOT"
	 }
}
`
}

func testAccCloudstackDomainDataSource_invalidName() string {
	return `
data "cloudstack_domain" "my_domain" {
	 filter {
	   name = "name"
	   value = "badgerbearocto"
	 }
}
`
}

func TestAccCloudstackDomainDataSource_byID(t *testing.T) {
	domainResourceName := "cloudstack_domain.test_domain"
	dataSourceName := "data.cloudstack_domain.my_domain_by_id"
	testDomainName := "test-domain-" + id.UniqueId()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudstackDomainDataSource_byID(testDomainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudstackDomainDataSourceExists(dataSourceName),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", domainResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "domain_id", domainResourceName, "id"),
				),
			},
		},
	})
}

func testAccCloudstackDomainDataSource_byID(domainName string) string {
	return fmt.Sprintf(`
resource "cloudstack_domain" "test_domain" {
  name = "%s"
}

data "cloudstack_domain" "my_domain_by_id" {
  filter {
    name  = "id"
    value = cloudstack_domain.test_domain.id
  }
}
`, domainName)
}
