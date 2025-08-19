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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					testAccCheckCloudStackProjectExists("cloudstack_project.test", &project),
					testAccCheckCloudStackProjectBasicAttributes(&project),
					resource.TestCheckResourceAttr("cloudstack_project.test", "name", "terraform-test"),
					resource.TestCheckResourceAttr("cloudstack_project.test", "display_text", "terraform-test"),
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
					testAccCheckCloudStackProjectExists("cloudstack_project.test", &project),
					testAccCheckCloudStackProjectBasicAttributes(&project),
					resource.TestCheckResourceAttr("cloudstack_project.test", "name", "terraform-test"),
					resource.TestCheckResourceAttr("cloudstack_project.test", "display_text", "terraform-test"),
				),
			},

			{
				Config: testAccCloudStackProject_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectExists("cloudstack_project.test", &project),
					testAccCheckCloudStackProjectUpdatedAttributes(&project),
					resource.TestCheckResourceAttr("cloudstack_project.test", "name", "terraform-test-updated"),
					resource.TestCheckResourceAttr("cloudstack_project.test", "display_text", "terraform-test-updated"),
				),
			},
		},
	})
}

func testAccCheckCloudStackProjectExists(n string, project *cloudstack.Project) resource.TestCheckFunc {
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

		if list.Count == 0 {
			return fmt.Errorf("Project not found")
		}

		if list.Count > 1 {
			return fmt.Errorf("Found more than one project with ID: %s", rs.Primary.ID)
		}

		*project = *list.Projects[0]

		return nil
	}
}

func testAccCheckCloudStackProjectBasicAttributes(project *cloudstack.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if project.Name != "terraform-test" {
			return fmt.Errorf("Bad name: %s", project.Name)
		}

		if project.Displaytext != "terraform-test" {
			return fmt.Errorf("Bad display text: %s", project.Displaytext)
		}

		return nil
	}
}

func testAccCheckCloudStackProjectUpdatedAttributes(project *cloudstack.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if project.Name != "terraform-test-updated" {
			return fmt.Errorf("Bad name: %s", project.Name)
		}

		if project.Displaytext != "terraform-test-updated" {
			return fmt.Errorf("Bad display text: %s", project.Displaytext)
		}

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

		if list.Count > 0 {
			return fmt.Errorf("Project %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackProject_basic = `
resource "cloudstack_project" "test" {
  name         = "terraform-test"
  display_text = "terraform-test"
}`

const testAccCloudStackProject_update = `
resource "cloudstack_project" "test" {
  name         = "terraform-test-updated"
  display_text = "terraform-test-updated"
}`
