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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackDiskOffering_basic(t *testing.T) {
	var do cloudstack.DiskOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDiskOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDiskOffering_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskOfferingExists("cloudstack_disk_offering.test1", &do),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "name", "disk_offering_1"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "display_text", "Test"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "disk_size", "10"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "storage_type", "shared"),
				),
			},
			{
				ResourceName:      "cloudstack_disk_offering.test1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCloudStackDiskOffering_basic = `
resource "cloudstack_disk_offering" "test1" {
  name         = "disk_offering_1"
  display_text = "Test"
  disk_size    = 10
}
`

func TestAccCloudStackDiskOffering_customized(t *testing.T) {
	var do cloudstack.DiskOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDiskOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDiskOffering_customized,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskOfferingExists("cloudstack_disk_offering.custom", &do),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.custom", "customized", "true"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.custom", "storage_type", "local"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.custom", "provisioning_type", "thin"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.custom", "tags", "ssd"),
				),
			},
		},
	})
}

const testAccCloudStackDiskOffering_customized = `
resource "cloudstack_disk_offering" "custom" {
  name              = "custom_disk_offering"
  display_text      = "Custom Test"
  customized        = true
  storage_type      = "local"
  provisioning_type = "thin"
  tags              = "ssd"
}
`

func TestAccCloudStackDiskOffering_update(t *testing.T) {
	var do cloudstack.DiskOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackDiskOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackDiskOffering_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskOfferingExists("cloudstack_disk_offering.test1", &do),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "display_text", "Test"),
				),
			},
			{
				Config: testAccCloudStackDiskOffering_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackDiskOfferingExists("cloudstack_disk_offering.test1", &do),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "name", "disk_offering_1_updated"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "display_text", "Test Updated"),
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "tags", "gold"),
				),
			},
		},
	})
}

const testAccCloudStackDiskOffering_update = `
resource "cloudstack_disk_offering" "test1" {
  name         = "disk_offering_1_updated"
  display_text = "Test Updated"
  disk_size    = 10
  tags         = "gold"
}
`

func testAccCheckCloudStackDiskOfferingExists(n string, do *cloudstack.DiskOffering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk offering ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		resp, _, err := cs.DiskOffering.GetDiskOfferingByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp.Id != rs.Primary.ID {
			return fmt.Errorf("Disk offering not found")
		}

		*do = *resp

		return nil
	}
}

func testAccCheckCloudStackDiskOfferingDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_disk_offering" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No disk offering ID is set")
		}

		_, _, err := cs.DiskOffering.GetDiskOfferingByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Disk offering %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
