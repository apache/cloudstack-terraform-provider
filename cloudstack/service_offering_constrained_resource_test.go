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

func TestAccServiceOfferingConstrained(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccMuxProvider,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceOfferingCustomConstrained1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_constrained.constrained1", "name", "constrained1"),
				),
			},
			{
				Config: testAccServiceOfferingCustomConstrained1ZoneAll,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_constrained.constrained1", "name", "constrained1"),
				),
			},
			{
				Config: testAccServiceOfferingCustomConstrained2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_constrained.constrained2", "name", "constrained2"),
				),
			},
			{
				Config: testAccServiceOfferingCustomConstrained2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_constrained.constrained2", "name", "constrained2update"),
				),
			},
			{
				Config: testAccServiceOfferingCustomConstrained_disk,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_constrained.constrained1", "name", "constrained1"),
				),
			},
			{
				Config: testAccServiceOfferingCustomConstrained_disk_hypervisor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_constrained.disk_hypervisor", "name", "disk_hypervisor"),
				),
			},
			{
				Config: testAccServiceOfferingCustomConstrained_disk_storage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering_constrained.disk_storage", "name", "disk_storage"),
				),
			},
		},
	})
}

const testAccServiceOfferingCustomConstrained1 = `
resource "cloudstack_zone" "test" {
	name          = "acctest"
	dns1          = "8.8.8.8"
	internal_dns1 = "8.8.4.4"
	network_type  = "Advanced"
}

resource "cloudstack_service_offering_constrained" "constrained1" {
	display_text = "constrained1"
	name         = "constrained1"

	// compute
	cpu_speed  = 2500
	max_cpu_number = 10
	min_cpu_number = 2
	
	// memory
	max_memory     = 4096
	min_memory     = 1024
	
	// other
	host_tags    = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"
	
	// Feature flags
	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false
	zone_ids = [cloudstack_zone.test.id]

}
`

const testAccServiceOfferingCustomConstrained1ZoneAll = `
resource "cloudstack_service_offering_constrained" "constrained1" {
	display_text = "constrained11"
	name         = "constrained1"

	// compute
	cpu_speed  = 2500
	max_cpu_number = 10
	min_cpu_number = 2

	// memory
	max_memory     = 4096
	min_memory     = 1024

	// other
	host_tags    = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	// Feature flags
	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false
	zone_ids = []
}
`

const testAccServiceOfferingCustomConstrained2 = `
resource "cloudstack_service_offering_constrained" "constrained2" {
	display_text = "constrained2"
	name         = "constrained2"

	// compute
	cpu_speed  = 2500
	max_cpu_number = 10
	min_cpu_number = 2
	
	// memory
	max_memory     = 4096
	min_memory     = 1024
	
	// other
	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"
	
	// Feature flags
	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true
}
`
const testAccServiceOfferingCustomConstrained2_update = `
resource "cloudstack_service_offering_constrained" "constrained2" {
	display_text = "constrained2update"
	name         = "constrained2update"

	// compute
	cpu_speed  = 2500
	max_cpu_number = 10
	min_cpu_number = 2
	
	// memory
	max_memory     = 4096
	min_memory     = 1024
	
	// other
	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"
	
	// Feature flags
	dynamic_scaling_enabled = true
	is_volatile             = true
	limit_cpu_use           = true
	offer_ha                = true
}
`

const testAccServiceOfferingCustomConstrained_disk = `
resource "cloudstack_service_offering_constrained" "constrained1" {
	display_text = "constrained1"
	name         = "constrained1"

	// compute
	cpu_speed  = 2500
	max_cpu_number = 10
	min_cpu_number = 2

	// memory
	max_memory     = 4096
	min_memory     = 1024

	// other
	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	// Feature flags
	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false

	disk_offering = {
		storage_type = "local"
		sdfjklsdf = "sdfjks"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		storage_tags = "FOO"
		disk_offering_strictness = false
	}
}
`

const testAccServiceOfferingCustomConstrained_disk_hypervisor = `
resource "cloudstack_service_offering_constrained" "disk_hypervisor" {
	display_text = "disk_hypervisor"
	name         = "disk_hypervisor"

	// compute
	cpu_speed  = 2500
	max_cpu_number = 10
	min_cpu_number = 2

	// memory
	max_memory     = 4096
	min_memory     = 1024

	// other
	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	// Feature flags
	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false

	disk_offering = {
		storage_type = "local"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		storage_tags = "FOO"
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

const testAccServiceOfferingCustomConstrained_disk_storage = `
resource "cloudstack_service_offering_constrained" "disk_storage" {
	display_text = "disk_storage"
	name         = "disk_storage"

	// compute
	cpu_speed  = 2500
	max_cpu_number = 10
	min_cpu_number = 2

	// memory
	max_memory     = 4096
	min_memory     = 1024

	// other
	host_tags = "test0101,test0202"
	network_rate = 1024
	deployment_planner = "UserDispersingPlanner"

	// Feature flags
	dynamic_scaling_enabled = false
	is_volatile             = false
	limit_cpu_use           = false
	offer_ha                = false

	disk_offering = {
		storage_type = "local"
		provisioning_type = "thin"
		cache_mode = "none"
		root_disk_size = "5"
		storage_tags = "FOO"
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
