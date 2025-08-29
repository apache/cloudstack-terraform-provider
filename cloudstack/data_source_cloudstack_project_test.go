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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_project.project-resource"
	datasourceName := "data.cloudstack_project.project-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testProjectDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "display_text", resourceName, "display_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain", resourceName, "domain"),
				),
			},
		},
	})
}

func TestAccProjectDataSource_withAccount(t *testing.T) {
	resourceName := "cloudstack_project.project-account-resource"
	datasourceName := "data.cloudstack_project.project-account-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testProjectDataSourceConfig_withAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "display_text", resourceName, "display_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain", resourceName, "domain"),
					resource.TestCheckResourceAttrPair(datasourceName, "account", resourceName, "account"),
				),
			},
		},
	})
}

const testProjectDataSourceConfig_basic = `
resource "cloudstack_project" "project-resource" {
  name = "test-project-datasource"
  display_text = "Test Project for Data Source"
}

data "cloudstack_project" "project-data-source" {
  filter {
    name = "name"
    value = "test-project-datasource"
  }
  depends_on = [
    cloudstack_project.project-resource
  ]
}

output "project-output" {
  value = data.cloudstack_project.project-data-source
}
`

const testProjectDataSourceConfig_withAccount = `
resource "cloudstack_project" "project-account-resource" {
  name = "test-project-account-datasource"
  display_text = "Test Project with Account for Data Source"
  account = "admin"
  domain = "ROOT"
}

data "cloudstack_project" "project-account-data-source" {
  filter {
    name = "name"
    value = "test-project-account-datasource"
  }
  filter {
    name = "account"
    value = "admin"
  }
  depends_on = [
    cloudstack_project.project-account-resource
  ]
}

output "project-account-output" {
  value = data.cloudstack_project.project-account-data-source
}
`
