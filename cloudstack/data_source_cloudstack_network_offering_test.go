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

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkOfferingDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_network_offering.net-off-resource"
	datasourceName := "data.cloudstack_network_offering.net-off-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNetworkOfferingDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testNetworkOfferingDataSourceConfig_basic = `
resource "cloudstack_network_offering" "net-off-resource"{
  name       = "TestNetworkDisplay01"
  display_text = "TestNetworkDisplay01"
  guest_ip_type = "Isolated"
  traffic_type = "Guest"
  }

  data "cloudstack_network_offering" "net-off-data-source"{

    filter{
    name = "name"
    value="TestNetworkDisplay01"
    }
	  depends_on = [
	  cloudstack_network_offering.net-off-resource
	]
  }
  `
