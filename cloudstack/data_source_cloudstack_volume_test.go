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

func TestAccVolumeDataSource_basic(t *testing.T) {
	resourceName := "cloudstack_volume.volume-resource"
	datasourceName := "data.cloudstack_volume.volume-data-source"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVolumeDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testVolumeDataSourceConfig_basic = `
resource "cloudstack_volume" "volume-resource"{
	name       = "TestVolume"
  disk_offering_id = "0038adec-5e3e-47df-b4b4-77b5dc8e3338"
  zone_id = "9a7002b2-09a2-44dc-a332-f2e4e7f01539"
  }

  data "cloudstack_volume" "volume-data-source"{
    filter{
    name = "name"
    value="TestVolume"
    }
	  depends_on = [
	  cloudstack_volume.volume-resource
	]
  }
  `
