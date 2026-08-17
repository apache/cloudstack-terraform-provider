// //
// // Licensed to the Apache Software Foundation (ASF) under one
// // or more contributor license agreements.  See the NOTICE file
// // distributed with this work for additional information
// // regarding copyright ownership.  The ASF licenses this file
// // to you under the Apache License, Version 2.0 (the
// // "License"); you may not use this file except in compliance
// // with the License.  You may obtain a copy of the License at
// //
// //   http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing,
// // software distributed under the License is distributed on an
// // "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// // KIND, either express or implied.  See the License for the
// // specific language governing permissions and limitations
// // under the License.
// //

package cloudstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceOfferingUnconstrained(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccMuxProvider,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceOfferingUnconstrained1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_unconstrained.unconstrained1", "name", "unconstrained1"),
				),
			},
			{
				Config: testAccServiceOfferingUnconstrained2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_unconstrained.unconstrained2", "name", "unconstrained2"),
				),
			},
			{
				Config: testAccServiceOfferingUnconstrained2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_unconstrained.unconstrained2", "name", "unconstrained2update"),
				),
			},
			{
				Config: testAccServiceOfferingUnconstrained_disk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_unconstrained.disk", "name", "disk"),
				),
			},
			{
				Config: testAccServiceOfferingUnconstrained_disk_hypervisor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_unconstrained.disk_hypervisor", "name", "disk_hypervisor"),
				),
			},
			{
				Config: testAccServiceOfferingUnconstrained_disk_storage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_unconstrained.disk_storage", "name", "disk_storage"),
				),
			},
		},
	})
}

const testAccServiceOfferingUnconstrained1 = `
resource "cloudstack_service_offering_unconstrained" "unconstrained1" {
	display_text = "unconstrained1"
	name         = "unconstrained1"

	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false
}
`

const testAccServiceOfferingUnconstrained2 = `
resource "cloudstack_service_offering_unconstrained" "unconstrained2" {
	display_text = "unconstrained2"
	name         = "unconstrained2"

	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true
}
`

const testAccServiceOfferingUnconstrained2_update = `
resource "cloudstack_service_offering_unconstrained" "unconstrained2" {
	display_text = "unconstrained2update"
	name         = "unconstrained2update"

	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true
}
`

const testAccServiceOfferingUnconstrained_disk = `
resource "cloudstack_service_offering_unconstrained" "disk" {
	display_text = "disk"
	name         = "disk"

	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true

	disk_offering = {
		storage_type = "shared"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		storage_tags = "test0101,test0202"
		disk_offering_strictness = false
	}
}
`

const testAccServiceOfferingUnconstrained_disk_hypervisor = `
resource "cloudstack_service_offering_unconstrained" "disk_hypervisor" {
	display_text = "disk_hypervisor"
	name         = "disk_hypervisor"

	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true

	disk_offering = {
		storage_type = "shared"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		storage_tags = "test0101,test0202"
		disk_offering_strictness = false
	}
	disk_hypervisor = {
		bytes_read_rate             = 1024
		bytes_read_rate_max         = 1024
		bytes_read_rate_max_length  = 1024
		bytes_write_rate            = 1024
		bytes_write_rate_max        = 1024
		bytes_write_rate_max_length = 1024
	}

}
`

const testAccServiceOfferingUnconstrained_disk_storage = `
resource "cloudstack_service_offering_unconstrained" "disk_storage" {
	display_text = "disk_storage"
	name         = "disk_storage"

	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true

	disk_offering = {
		storage_type = "shared"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		storage_tags = "test0101,test0202"
		disk_offering_strictness = false
	}
	disk_storage = {
		min_iops = 100
		max_iops = 100
	}
}
`
