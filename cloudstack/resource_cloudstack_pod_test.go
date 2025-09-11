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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudStackPod_basic(t *testing.T) {
	var pod cloudstack.Pod

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackPodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackPod_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackPodExists(
						"cloudstack_pod.foo", &pod),
					testAccCheckCloudStackPodAttributes(&pod),
					resource.TestCheckResourceAttr(
						"cloudstack_pod.foo", "name", "terraform-pod"),
				),
			},
		},
	})
}

func TestAccCloudStackPod_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackPodDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackPod_basic,
			},

			{
				ResourceName:      "cloudstack_pod.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackPodExists(
	n string, pod *cloudstack.Pod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pod ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Pod.NewListPodsParams()
		p.SetId(rs.Primary.ID)

		list, err := cs.Pod.ListPods(p)
		if err != nil {
			return err
		}

		if list.Count != 1 || list.Pods[0].Id != rs.Primary.ID {
			return fmt.Errorf("Pod not found")
		}

		*pod = *list.Pods[0]

		return nil
	}
}

func testAccCheckCloudStackPodAttributes(
	pod *cloudstack.Pod) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if pod.Name != "terraform-pod" {
			return fmt.Errorf("Bad name: %s", pod.Name)
		}

		if pod.Gateway != "192.168.56.1" {
			return fmt.Errorf("Bad gateway: %s", pod.Gateway)
		}

		if pod.Netmask != "255.255.255.0" {
			return fmt.Errorf("Bad netmask: %s", pod.Netmask)
		}

		return nil
	}
}

func testAccCheckCloudStackPodDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_pod" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pod ID is set")
		}

		p := cs.Pod.NewListPodsParams()
		p.SetId(rs.Primary.ID)

		list, err := cs.Pod.ListPods(p)
		if err != nil {
			return nil
		}

		if list.Count > 0 {
			return fmt.Errorf("Pod %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackPod_basic = `
data "cloudstack_zone" "zone" {
  filter {
    name = "name"
    value = "Sandbox-simulator"
  }
}

# Create a pod in the zone
resource "cloudstack_pod" "foo" {
  name = "terraform-pod"
  zone_id = data.cloudstack_zone.zone.id
  gateway = "192.168.56.1"
  netmask = "255.255.255.0"
  start_ip = "192.168.56.2"
  end_ip = "192.168.56.254"
}`
