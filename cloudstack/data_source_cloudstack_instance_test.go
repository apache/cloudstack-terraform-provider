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

// basic acceptance to check if the display_name attribute has same value in
// the created instance and its data source respectively.
func TestAccInstanceDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_instance.my_instance"
	datasourceName := "data.cloudstack_instance.my_instance_test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "display_name", resourceName, "display_name"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccInstanceDataSourceConfig_basic = `
	resource "cloudstack_instance" "my_instance" {
		name             = "server-a"
		service_offering = "Small Instance"
		template         = "CentOS 5.6 (64-bit) no GUI (Simulator)"
		zone             = "Sandbox-simulator"
	  }
	  data "cloudstack_instance" "my_instance_test" {
		filter {
		name = "display_name" 
		value = "server-a"
	  }
		depends_on = [
		cloudstack_instance.my_instance
	  ]
	}
`
