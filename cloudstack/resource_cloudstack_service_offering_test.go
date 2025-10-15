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

func TestAccCloudStackServiceOffering_basic(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOffering_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists("cloudstack_service_offering.test1", &so),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test1", "cpu_number", "2"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test1", "cpu_speed", "2200"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test1", "memory", "8096"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOffering_basic = `
resource "cloudstack_service_offering" "test1" {
  name 			= "service_offering_1"
  display_text 	= "Test"
  cpu_number	= 2
  cpu_speed		= 2200
  memory        = 8096
}
`

func testAccCheckCloudStackServiceOfferingExists(n string, so *cloudstack.ServiceOffering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No service offering ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		resp, _, err := cs.ServiceOffering.GetServiceOfferingByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp.Id != rs.Primary.ID {
			return fmt.Errorf("Service offering not found")
		}

		*so = *resp

		return nil
	}
}

func TestAccCloudStackServiceOffering_customized(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOffering_customized,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists("cloudstack_service_offering.custom", &so),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "customized", "true"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "min_cpu_number", "1"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "max_cpu_number", "8"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "min_memory", "1024"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "max_memory", "16384"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "cpu_speed", "1000"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "encrypt_root", "true"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.custom", "storage_tags", "production,ssd"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOffering_customized = `
resource "cloudstack_service_offering" "custom" {
  name             = "custom_service_offering"
  display_text     = "Custom Test"
  customized       = true
  min_cpu_number   = 1
  max_cpu_number   = 8
  min_memory       = 1024
  max_memory       = 16384
  cpu_speed        = 1000
  encrypt_root     = true
  storage_tags     = "production,ssd"
}
`

func TestAccCloudStackServiceOffering_gpu(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOffering_gpu,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists("cloudstack_service_offering.gpu", &so),
					resource.TestCheckResourceAttr("cloudstack_service_offering.gpu", "name", "gpu_service_offering"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.gpu", "display_text", "GPU Test"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.gpu", "cpu_number", "4"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.gpu", "memory", "16384"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.gpu", "service_offering_details.pciDevice", "Group of NVIDIA A6000 GPUs"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.gpu", "service_offering_details.vgpuType", "A6000-8A"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOffering_gpu = `
resource "cloudstack_service_offering" "gpu" {
  name         = "gpu_service_offering"
  display_text = "GPU Test"
  cpu_number   = 4
  memory       = 16384
  
  service_offering_details = {
    pciDevice = "Group of NVIDIA A6000 GPUs"
    vgpuType  = "A6000-8A"
  }
}
`
