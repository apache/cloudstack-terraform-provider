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

func TestAccCloudStackProject_basic(t *testing.T) {
	var project cloudstack.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.foo", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.foo", "name", "terraform-test-project"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.foo", "display_text", "Terraform Test Project"),
				),
			},
		},
	})
}

func TestAccCloudStackProject_update(t *testing.T) {
	var project cloudstack.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.foo", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.foo", "name", "terraform-test-project"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.foo", "display_text", "Terraform Test Project"),
				),
			},
			{
				Config: testAccCloudStackProject_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.foo", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.foo", "name", "terraform-test-project-updated"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.foo", "display_text", "Terraform Test Project Updated"),
				),
			},
		},
	})
}

func TestAccCloudStackProject_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_basic,
			},
			{
				ResourceName:      "cloudstack_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudStackProject_account(t *testing.T) {
	var project cloudstack.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_account,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.bar", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "name", "terraform-test-project-account"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "display_text", "Terraform Test Project with Account"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "account", "admin"),
				),
			},
		},
	})
}

func TestAccCloudStackProject_updateAccount(t *testing.T) {
	var project cloudstack.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_account,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.bar", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "name", "terraform-test-project-account"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "display_text", "Terraform Test Project with Account"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "account", "admin"),
				),
			},
			{
				Config: testAccCloudStackProject_updateAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.bar", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "name", "terraform-test-project-account"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "display_text", "Terraform Test Project with Account"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.bar", "account", "admin"),
				),
			},
		},
	})
}

func TestAccCloudStackProject_emptyDisplayText(t *testing.T) {
	var project cloudstack.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_emptyDisplayText,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.empty", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.empty", "name", "terraform-test-project-empty-display"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.empty", "display_text", "terraform-test-project-empty-display"),
				),
			},
		},
	})
}

func TestAccCloudStackProject_updateUserid(t *testing.T) {
	var project cloudstack.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_userid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.baz", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.baz", "name", "terraform-test-project-userid"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.baz", "display_text", "Terraform Test Project with Userid"),
				),
			},
			{
				Config: testAccCloudStackProject_updateUserid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists(
						"cloudstack_project.baz", &project),
					resource.TestCheckResourceAttr(
						"cloudstack_project.baz", "name", "terraform-test-project-userid-updated"),
					resource.TestCheckResourceAttr(
						"cloudstack_project.baz", "display_text", "Terraform Test Project with Userid Updated"),
				),
			},
		},
	})
}

func testAccCheckCloudStackProjectExists(
	n string, project *cloudstack.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Project.NewListProjectsParams()
		p.SetId(rs.Primary.ID)

		list, err := cs.Project.ListProjects(p)
		if err != nil {
			return err
		}

		if list.Count != 1 || list.Projects[0].Id != rs.Primary.ID {
			return fmt.Errorf("Project not found")
		}

		*project = *list.Projects[0]

		return nil
	}
}

func testAccCheckCloudStackProjectDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_project" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No project ID is set")
		}

		p := cs.Project.NewListProjectsParams()
		p.SetId(rs.Primary.ID)

		list, err := cs.Project.ListProjects(p)
		if err != nil {
			return err
		}

		if list.Count != 0 {
			return fmt.Errorf("Project %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackProject_basic = `
resource "cloudstack_project" "foo" {
  name = "terraform-test-project"
  display_text = "Terraform Test Project"
}`

const testAccCloudStackProject_update = `
resource "cloudstack_project" "foo" {
  name = "terraform-test-project-updated"
  display_text = "Terraform Test Project Updated"
}`

const testAccCloudStackProject_account = `
resource "cloudstack_project" "bar" {
  name = "terraform-test-project-account"
  display_text = "Terraform Test Project with Account"
  account = "admin"
  domain = "ROOT"
}`

const testAccCloudStackProject_updateAccount = `
resource "cloudstack_project" "bar" {
  name = "terraform-test-project-account"
  display_text = "Terraform Test Project with Account"
  account = "admin"
  domain = "ROOT"
}`

const testAccCloudStackProject_userid = `
resource "cloudstack_project" "baz" {
  name = "terraform-test-project-userid"
  display_text = "Terraform Test Project with Userid"
  domain = "ROOT"
}`

const testAccCloudStackProject_updateUserid = `
resource "cloudstack_project" "baz" {
  name = "terraform-test-project-userid-updated"
  display_text = "Terraform Test Project with Userid Updated"
  domain = "ROOT"
}`

const testAccCloudStackProject_emptyDisplayText = `
resource "cloudstack_project" "empty" {
  name = "terraform-test-project-empty-display"
  display_text = "terraform-test-project-empty-display"
}`
