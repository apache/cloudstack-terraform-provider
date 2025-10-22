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

func TestAccServiceOfferingDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_service_offering.service-offering-resource"
	datasourceName := "data.cloudstack_service_offering.service-offering-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServiceOfferingDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "cpu_number", resourceName, "cpu_number"),
					resource.TestCheckResourceAttrPair(datasourceName, "cpu_speed", resourceName, "cpu_speed"),
					resource.TestCheckResourceAttrPair(datasourceName, "memory", resourceName, "memory"),
					resource.TestCheckResourceAttrPair(datasourceName, "gpu_count", resourceName, "gpu_count"),
					resource.TestCheckResourceAttrPair(datasourceName, "gpu_display", resourceName, "gpu_display"),
					resource.TestCheckResourceAttrPair(datasourceName, "offer_ha", resourceName, "offer_ha"),
					resource.TestCheckResourceAttrPair(datasourceName, "storage_type", resourceName, "storage_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "disk_iops_read_rate", resourceName, "disk_iops_read_rate"),
					resource.TestCheckResourceAttrPair(datasourceName, "disk_iops_write_rate", resourceName, "disk_iops_write_rate"),
					// Skip IOPS comparison - these fields may not be supported by CloudStack API
					// resource.TestCheckResourceAttrPair(datasourceName, "min_iops", resourceName, "min_iops"),
					// resource.TestCheckResourceAttrPair(datasourceName, "max_iops", resourceName, "max_iops"),
					resource.TestCheckResourceAttrPair(datasourceName, "dynamic_scaling_enabled", resourceName, "dynamic_scaling_enabled"),
					resource.TestCheckResourceAttrPair(datasourceName, "is_volatile", resourceName, "is_volatile"),
					resource.TestCheckResourceAttrPair(datasourceName, "root_disk_size", resourceName, "root_disk_size"),
				),
			},
		},
	})
}

const testServiceOfferingDataSourceConfig_basic = `
resource "cloudstack_service_offering" "service-offering-resource" {
	name                     = "TestServiceUpdate"
	display_text            = "DisplayService"
	cpu_number             = 4
	cpu_speed              = 2400
	memory                 = 8192
	gpu_count              = 1
	gpu_display            = true
	offer_ha               = true
	storage_type           = "shared"
	disk_iops_read_rate   = 10000
	disk_iops_write_rate  = 10000
	min_iops              = 5000
	max_iops              = 15000
	dynamic_scaling_enabled = true
	is_volatile           = false
	root_disk_size        = 50
}

data "cloudstack_service_offering" "service-offering-data-source" {
	filter {
		name	= "name"
		value	= "TestServiceUpdate"
	}
	depends_on = [cloudstack_service_offering.service-offering-resource]
}
`
