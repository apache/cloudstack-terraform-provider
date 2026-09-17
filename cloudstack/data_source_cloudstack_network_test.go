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

func TestAccNetworkDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_network.network-resource"
	datasourceName := "data.cloudstack_network.network-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "display_text", resourceName, "display_text"),
					resource.TestCheckResourceAttrPair(datasourceName, "cidr", resourceName, "cidr"),
					resource.TestCheckResourceAttrPair(datasourceName, "network_offering_name", resourceName, "network_offering"),
				),
			},
		},
	})
}

const testNetworkDataSourceConfig_basic = `
resource "cloudstack_network" "network-resource" {
  name             = "terraform-network"
  display_text     = "terraform-network"
  cidr             = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone             = "Sandbox-simulator"
}

data "cloudstack_network" "network-data-source" {
  filter {
    name  = "name"
    value = "terraform-network"
  }

  filter {
    name  = "cidr"
    value = "10.1.1.0/24"
  }

  depends_on = [
    cloudstack_network.network-resource
  ]
}

output "network-output" {
  value = data.cloudstack_network.network-data-source
}
`
