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

func TestAccServiceOfferingFixed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccMuxProvider,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceOfferingFixed1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_fixed.fixed1", "name", "fixed1"),
				),
			},
			{
				Config: testAccServiceOfferingFixed2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_fixed.fixed2", "name", "fixed2"),
				),
			},
			{
				Config: testAccServiceOfferingFixed2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_fixed.fixed2", "name", "fixed2update"),
				),
			},
			{
				Config: testAccServiceOfferingFixed_disk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_fixed.disk", "name", "disk"),
				),
			},
			{
				Config: testAccServiceOfferingFixed_disk_hypervisor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_fixed.disk_hypervisor", "name", "disk_hypervisor"),
				),
			},
			{
				Config: testAccServiceOfferingFixed_disk_storage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_fixed.disk_storage", "name", "disk_storage"),
				),
			},
		},
	})
}

const testAccServiceOfferingFixed1 = `
resource "cloudstack_service_offering_fixed" "fixed1" {
	display_text = "fixed1"
	name         = "fixed1"

	// compute
	cpu_number     = 2
	cpu_speed      = 2500
	memory         = 2048

	// other
	host_tags          = "test0101, test0202"
	network_rate       = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false

}
`

const testAccServiceOfferingFixed2 = `
resource "cloudstack_service_offering_fixed" "fixed2" {
	display_text = "fixed2"
	name         = "fixed2"

	cpu_number     = 2
	cpu_speed      = 2500
	memory         = 2048

	host_tags          = "test0101,test0202"
	network_rate       = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true
}
`

const testAccServiceOfferingFixed2_update = `
resource "cloudstack_service_offering_fixed" "fixed2" {
	display_text = "fixed2update"
	name         = "fixed2update"

	cpu_number     = 2
	cpu_speed      = 2500
	memory         = 2048

	host_tags          = "test0101,test0202"
	network_rate       = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true
}
`

const testAccServiceOfferingFixed_disk = `
resource "cloudstack_service_offering_fixed" "disk" {
	display_text = "disk"
	name         = "disk"

	// compute
	cpu_number     = 2
	cpu_speed      = 2500
	memory         = 2048

	// other
	host_tags          = "test0101, test0202"
	network_rate       = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false

	disk_offering = {
		storage_type = "local"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		tags = "FOO"
		disk_offering_strictness = false
	}
}
`

const testAccServiceOfferingFixed_disk_hypervisor = `
resource "cloudstack_service_offering_fixed" "disk_hypervisor" {
	display_text = "disk_hypervisor"
	name         = "disk_hypervisor"

	// compute
	cpu_number     = 2
	cpu_speed      = 2500
	memory         = 2048

	// other
	host_tags          = "test0101, test0202"
	network_rate       = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false

	disk_offering = {
		storage_type = "local"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		tags = "FOO"
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

const testAccServiceOfferingFixed_disk_storage = `
resource "cloudstack_service_offering_fixed" "disk_storage" {
	display_text = "disk_storage"
	name         = "disk_storage"

	// compute
	cpu_number     = 2
	cpu_speed      = 2500
	memory         = 2048

	// other
	host_tags          = "test0101, test0202"
	network_rate       = 1024
	deployment_planner = "UserDispersingPlanner"

	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false

	disk_offering = {
		storage_type = "local"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		tags = "FOO"
		disk_offering_strictness = false
	}
	disk_storage = {
		min_iops = 100
		max_iops = 100
	}
}
`
