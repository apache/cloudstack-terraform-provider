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
			},
		},
	})
}

const testVolumeDataSourceConfig_basic = `
resource "cloudstack_disk_offering" "disk-offering" {
	name			=	"TestDiskOffering"
	display_text	=	"TestDiskOffering"
	disk_size		=	1
}

data "cloudstack_zone" "zone-data-source" {
	filter {
		name	=	"name"
		value	=	"Sandbox-simulator"
	}
}

resource "cloudstack_volume" "volume-resource"{
	name				=	"TestVolume"
	disk_offering_id	=	cloudstack_disk_offering.disk-offering.id
	zone_id				=	data.cloudstack_zone.zone-data-source.id
  }

data "cloudstack_volume" "volume-data-source"{
	filter {
		name	=	"name"
		value	=	"TestVolume"
    }
	depends_on	=	[cloudstack_volume.volume-resource]
}
`
