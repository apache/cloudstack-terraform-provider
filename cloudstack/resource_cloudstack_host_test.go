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
	"regexp"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackHost_basic(t *testing.T) {
	var h cloudstack.Host
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackHost_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackHostExists("cloudstack_host.test", &h),
					resource.TestCheckResourceAttr("cloudstack_host.test", "hypervisor", "Simulator"),
					resource.TestCheckResourceAttr("cloudstack_host.test", "cluster_name", "C1"),
					resource.TestCheckResourceAttrSet("cloudstack_host.test", "state"),
					resource.TestCheckResourceAttrSet("cloudstack_host.test", "name"),
				),
			},
		},
	})
}

const testAccCloudStackHost_basic = `
data "cloudstack_zone" "zone" {
	filter {
		name  = "name"
		value = "Sandbox-simulator"
	}
}

data "cloudstack_pod" "pod" {
	filter {
		name  = "name"
		value = "POD0"
	}
}

resource "cloudstack_host" "test" {
  hypervisor 	= "Simulator"
  pod_id     	= data.cloudstack_pod.pod.id
  url        	= "http://sim/c1/h"
  zone_id    	= data.cloudstack_zone.zone.id
  cluster_name 	= "C1"
  username   	= "root"
  password   	= "password"
}
`

func testAccCheckCloudStackHostExists(n string, h *cloudstack.Host) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No host ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		resp, _, err := cs.Host.GetHostByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp.Id != rs.Primary.ID {
			return fmt.Errorf("Host not found")
		}

		*h = *resp

		return nil
	}
}

func TestAccCloudStackHost_fail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudStackHost_fail,
				ExpectError: regexp.MustCompile("timeout waiting for Host to be created, with error: .*Username and Password need to be provided.*"),
			},
		},
	})
}

const testAccCloudStackHost_fail = `
data "cloudstack_zone" "zone_fail" {
	filter {
		name  = "name"
		value = "Sandbox-simulator"
	}
}

data "cloudstack_pod" "pod_fail" {
	filter {
		name  = "name"
		value = "POD0"
	}
}

resource "cloudstack_host" "test_fail" {
  hypervisor 	 = "Simulator"
  pod_id     	 = data.cloudstack_pod.pod_fail.id
  url        	 = "http://sim/c1/h"
  zone_id    	 = data.cloudstack_zone.zone_fail.id
  cluster_name 	 = "C1"
  username   	 = "root"
  password   	 = ""
  create_timeout = 10
}
`
