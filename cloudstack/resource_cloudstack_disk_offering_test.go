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

func TestAccDiskOffering(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskOffering1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test1", "name", "test1"),
				),
			},
			{
				Config: testAccDiskOffering2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test2", "name", "test2"),
				),
			},
			{
				Config: testAccDiskOffering3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test3", "name", "test3"),
				),
			},
			{
				Config: testAccDiskOffering3Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test3", "name", "test3update"),
				),
			},
			{
				Config: testAccDiskOffering4,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test4", "name", "test4"),
				),
			},
			{
				Config: testAccDiskOffering5,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test5", "name", "test5"),
				),
			},
			{
				Config: testAccDiskOffering5Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_disk_offering.test5", "name", "test5update"),
				),
			},
		},
	})
}

const testAccDiskOffering1 = `
resource "cloudstack_disk_offering" "test1" {
	display_text = "test1"
	name         = "test1"

	storage_type = "shared"
	provisioning_type = "fat"
	cache_mode        = "writeback"
	tags              = "test1,test2"

	disk_size = 7
}
`

const testAccDiskOffering2 = `
resource "cloudstack_disk_offering" "test2" {
	display_text = "test2"
	name         = "test2"

	storage_type = "local"
	provisioning_type = "thin"
	cache_mode        = "writethrough"
	tags              = "test1,test2"

	disk_size = 7
}
`

const testAccDiskOffering3 = `
resource "cloudstack_disk_offering" "test3" {
	display_text = "test3"
	name         = "test3"

	storage_type = "local"
	provisioning_type = "thin"
	cache_mode        = "writethrough"
	tags              = "test1,test2"

	disk_size = 7

	storage = {
		min_iops = 100
		max_iops = 100
		customized_iops = false
		disk_offering_strictness = true
	}
}
`

const testAccDiskOffering3Update = `
resource "cloudstack_disk_offering" "test3" {
	display_text = "test3update"
	name         = "test3update"

	storage_type = "shared"
	provisioning_type = "fat"
	cache_mode        = "writeback"
	tags              = "test1,test2,test3"

	disk_size = 7

	storage = {
		min_iops = 200
		max_iops = 300
		customized_iops = false
		disk_offering_strictness = true
	}
}
`

const testAccDiskOffering4 = `
resource "cloudstack_disk_offering" "test4" {
	display_text = "test4"
	name         = "test4"

	storage_type = "local"
	provisioning_type = "thin"
	cache_mode        = "writethrough"
	tags              = "test1,test2"

	disk_size = 7

	storage = {
		customized_iops = true
	}
}
`
const testAccDiskOffering5 = `
resource "cloudstack_disk_offering" "test5" {
	display_text = "test5"
	name         = "test5"

	storage_type = "local"
	provisioning_type = "thin"
	cache_mode        = "writethrough"
	tags              = "test1,test2"

	disk_size = 7

	hypervisor = {
		bytes_read_rate             = 1024
		bytes_read_rate_max         = 1024
		bytes_read_rate_max_length  = 1024
		bytes_write_rate            = 1024
		bytes_write_rate_max        = 1024
		bytes_write_rate_max_length = 1024
	}
}
`

const testAccDiskOffering5Update = `
resource "cloudstack_disk_offering" "test5" {
	display_text = "test5update"
	name         = "test5update"

	storage_type = "shared"
	provisioning_type = "fat"
	cache_mode        = "writeback"
	tags              = "test1,test2,test3"

	disk_size = 7

	hypervisor = {
		bytes_read_rate             = 2048
		bytes_read_rate_max         = 2048
		bytes_read_rate_max_length  = 2048
		bytes_write_rate            = 2048
		bytes_write_rate_max        = 2048
		bytes_write_rate_max_length = 2048
	}
}
`
