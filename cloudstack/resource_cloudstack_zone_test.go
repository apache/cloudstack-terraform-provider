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

func TestAccCloudStackZone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackZone_basic,
			},
			{
				Config: testAccCloudStackZone_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudstack_zone.test", "name", "acctestupdated"),
				),
			},
		},
	})
}

const testAccCloudStackZone_basic = `
resource "cloudstack_zone" "test" {
	name          = "acctest"
	dns1          = "8.8.8.8"
	dns2          = "8.8.8.8"
	internal_dns1 = "8.8.4.4"
	internal_dns2 = "8.8.4.4"
	network_type  = "Advanced"
	domain        = "cloudstack.apache.org"
}
`

const testAccCloudStackZone_update = `
resource "cloudstack_zone" "test" {
	name             = "acctestupdated"
	dns1             = "8.8.4.4"
	dns2             = "8.8.4.4"
	internal_dns1    = "8.8.8.8"
	internal_dns2    = "8.8.8.8"
	network_type     = "Advanced"
	domain           = "cloudstack.apache.org"
	guest_cidr_address = "172.29.2.0/20"
}
`
