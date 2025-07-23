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
	"strings"
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

		// Get domain if available
		var domain string
		if domainAttr, ok := rs.Primary.Attributes["domain"]; ok && domainAttr != "" {
			domain = domainAttr
		}

		// First try to find the project by ID with domain if available
		p := cs.Project.NewListProjectsParams()
		p.SetId(rs.Primary.ID)

		// Add domain if available
		if domain != "" {
			domainID, err := retrieveID(cs, "domain", domain)
			if err == nil {
				p.SetDomainid(domainID)
			}
		}

		list, err := cs.Project.ListProjects(p)
		if err == nil && list.Count > 0 && list.Projects[0].Id == rs.Primary.ID {
			// Found by ID, set the project and return
			*project = *list.Projects[0]
			return nil
		}

		// If not found by ID or there was an error, try by name
		if err != nil || list.Count == 0 || list.Projects[0].Id != rs.Primary.ID {
			name := rs.Primary.Attributes["name"]
			if name == "" {
				return fmt.Errorf("Project not found by ID and name attribute is empty")
			}

			// Try to find by name
			p := cs.Project.NewListProjectsParams()
			p.SetName(name)

			// Add domain if available
			if domain, ok := rs.Primary.Attributes["domain"]; ok && domain != "" {
				domainID, err := retrieveID(cs, "domain", domain)
				if err != nil {
					return fmt.Errorf("Error retrieving domain ID: %v", err)
				}
				p.SetDomainid(domainID)
			}

			list, err := cs.Project.ListProjects(p)
			if err != nil {
				return fmt.Errorf("Error retrieving project by name: %s", err)
			}

			if list.Count == 0 {
				return fmt.Errorf("Project with name %s not found", name)
			}

			// Find the project with the matching ID if possible
			found := false
			for _, proj := range list.Projects {
				if proj.Id == rs.Primary.ID {
					*project = *proj
					found = true
					break
				}
			}

			// If we didn't find a project with matching ID, use the first one
			if !found {
				*project = *list.Projects[0]
				// Update the resource ID to match the found project
				rs.Primary.ID = list.Projects[0].Id
			}

			return nil
		}

		return fmt.Errorf("Project not found by ID or name")
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

		// Get domain if available
		var domain string
		if domainAttr, ok := rs.Primary.Attributes["domain"]; ok && domainAttr != "" {
			domain = domainAttr
		}

		// Try to find the project by ID
		p := cs.Project.NewListProjectsParams()
		p.SetId(rs.Primary.ID)

		// Add domain if available
		if domain != "" {
			domainID, err := retrieveID(cs, "domain", domain)
			if err == nil {
				p.SetDomainid(domainID)
			}
		}

		list, err := cs.Project.ListProjects(p)

		// If we get an error, check if it's a "not found" error
		if err != nil {
			if strings.Contains(err.Error(), "not found") ||
				strings.Contains(err.Error(), "does not exist") ||
				strings.Contains(err.Error(), "could not be found") ||
				strings.Contains(err.Error(), fmt.Sprintf(
					"Invalid parameter id value=%s due to incorrect long value format, "+
						"or entity does not exist", rs.Primary.ID)) {
				// Project doesn't exist, which is what we want
				continue
			}
			// For other errors, return them
			return fmt.Errorf("error checking if project %s exists: %s", rs.Primary.ID, err)
		}

		// If we found the project, it still exists
		if list.Count != 0 {
			return fmt.Errorf("project %s still exists (found by ID)", rs.Primary.ID)
		}

		// Also check by name to be thorough
		name := rs.Primary.Attributes["name"]
		if name != "" {
			// Try to find the project by name
			p := cs.Project.NewListProjectsParams()
			p.SetName(name)

			// Add domain if available
			if domain, ok := rs.Primary.Attributes["domain"]; ok && domain != "" {
				domainID, err := retrieveID(cs, "domain", domain)
				if err == nil {
					p.SetDomainid(domainID)
				}
			}

			list, err := cs.Project.ListProjects(p)
			if err != nil {
				// Ignore errors for name lookup
				continue
			}

			// Check if any of the returned projects match our ID
			for _, proj := range list.Projects {
				if proj.Id == rs.Primary.ID {
					return fmt.Errorf("project %s still exists (found by name %s)", rs.Primary.ID, name)
				}
			}
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

func TestAccCloudStackProject_list(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackProject_list,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackProjectsExist("cloudstack_project.project1", "cloudstack_project.project2"),
				),
			},
		},
	})
}

func testAccCheckCloudStackProjectsExist(projectNames ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Get CloudStack client
		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

		// Create a map to track which projects we've found
		foundProjects := make(map[string]bool)
		for _, name := range projectNames {
			// Get the project resource from state
			rs, ok := s.RootModule().Resources[name]
			if !ok {
				return fmt.Errorf("Not found: %s", name)
			}

			if rs.Primary.ID == "" {
				return fmt.Errorf("No project ID is set for %s", name)
			}

			// Add the project ID to our tracking map
			foundProjects[rs.Primary.ID] = false
		}

		// List all projects
		p := cs.Project.NewListProjectsParams()
		list, err := cs.Project.ListProjects(p)
		if err != nil {
			return err
		}

		// Check if all our projects are in the list
		for _, project := range list.Projects {
			if _, exists := foundProjects[project.Id]; exists {
				foundProjects[project.Id] = true
			}
		}

		// Verify all projects were found
		for id, found := range foundProjects {
			if !found {
				return fmt.Errorf("Project with ID %s was not found in the list", id)
			}
		}

		return nil
	}
}

const testAccCloudStackProject_list = `
resource "cloudstack_project" "project1" {
  name = "terraform-test-project-list-1"
  display_text = "Terraform Test Project List 1"
}

resource "cloudstack_project" "project2" {
  name = "terraform-test-project-list-2"
  display_text = "Terraform Test Project List 2"
}`
