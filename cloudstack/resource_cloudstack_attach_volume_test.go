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

func TestAccCloudstackAttachVolume_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudstackAttachVolume_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_attach_volume.foo", "device_id", "1"),
				),
			},
		},
	})
}

const testAccCloudstackAttachVolume_basic = `
resource "cloudstack_network" "foo" {
	name = "terraform-network"
	cidr = "10.1.1.0/24"
	network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
	zone = "Sandbox-simulator"
}
  
  resource "cloudstack_instance" "foobar" {
	name = "terraform-test"
	display_name = "terraform"
	service_offering= "Small Instance"
	network_id = cloudstack_network.foo.id
	template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
	zone = cloudstack_network.foo.zone
	expunge = true
}
  
  resource "cloudstack_disk" "foo" {
	name = "terraform-disk"
	disk_offering = "Small"
	zone = cloudstack_instance.foobar.zone
}

resource "cloudstack_attach_volume" "foo" {
	volume_id          = cloudstack_disk.foo.id
	virtual_machine_id = cloudstack_instance.foobar.id
}
`
