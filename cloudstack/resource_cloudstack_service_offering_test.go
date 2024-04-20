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

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudStackServiceOffering_basic(t *testing.T) {
	// ctx := context.Background()
	rName := acctest.RandomWithPrefix("service_offering")

	// var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProvidersV6,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "cloudstack_service_offering" "test" {
					name 			= "%s"
					display_text 	= "Test"
					cpu_number		= 2
					cpu_speed		= 2200
					memory          = 8096
				}`, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_service_offering.test", "cpu_number", "2"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test", "cpu_speed", "2200"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test", "memory", "8096"),
				),
			},
			{
				ResourceName:      "cloudstack_service_offering.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"custom_iops",
					"customized",
					"disk_offering_strictness",
					"dynamic_scaling",
					"encrypt",
					"host_tags",
					"provisioning_type",
					"qos_type",
					"volatile",
					"write_cache_type",
				},
			},
			{
				Config: fmt.Sprintf(`
				resource "cloudstack_service_offering" "test2" {
					name 			= "%s"
					display_text 	= "Test"
					customized      = true
					cpu_number_min  = 2
					cpu_number_max  = 4
					memory_min      = 8096
					memory_max      = 16384
					cpu_speed		= 2200
				}`, rName),
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckCloudStackServiceOfferingExists("cloudstack_service_offering.test2", &so),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test2", "cpu_number_min", "2"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test2", "cpu_number_max", "4"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test2", "memory_min", "8096"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test2", "memory_max", "16384"),
					resource.TestCheckResourceAttr("cloudstack_service_offering.test2", "cpu_speed", "2200"),
				),
			},
		},
	})
}
