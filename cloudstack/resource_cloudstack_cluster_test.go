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

func TestAccCloudStackCluster_basic(t *testing.T) {
	var cluster cloudstack.Cluster

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackCluster_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackClusterExists(
						"cloudstack_cluster.foo", &cluster),
					testAccCheckCloudStackClusterAttributes(&cluster),
					resource.TestCheckResourceAttr(
						"cloudstack_cluster.foo", "name", "terraform-cluster"),
				),
			},
		},
	})
}

func TestAccCloudStackCluster_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackCluster_basic,
			},

			{
				ResourceName:      "cloudstack_cluster.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "vsm_password",
				},
			},
		},
	})
}

func testAccCheckCloudStackClusterExists(
	n string, cluster *cloudstack.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		c, count, err := cs.Cluster.GetClusterByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *c

		return nil
	}
}

func testAccCheckCloudStackClusterAttributes(
	cluster *cloudstack.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if cluster.Name != "terraform-cluster" {
			return fmt.Errorf("Bad name: %s", cluster.Name)
		}

		if cluster.Clustertype != "CloudManaged" {
			return fmt.Errorf("Bad cluster type: %s", cluster.Clustertype)
		}

		if cluster.Hypervisortype != "KVM" {
			return fmt.Errorf("Bad hypervisor: %s", cluster.Hypervisortype)
		}

		return nil
	}
}

func testAccCheckCloudStackClusterDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_cluster" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster ID is set")
		}

		_, count, err := cs.Cluster.GetClusterByID(rs.Primary.ID)
		if err != nil {
			return nil
		}

		if count > 0 {
			return fmt.Errorf("Cluster %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackCluster_basic = `
data "cloudstack_zone" "zone" {
  filter {
    name = "name"
    value = "Sandbox-simulator"
  }
}

resource "cloudstack_pod" "foopod" {
  name = "terraform-pod"
  zone_id = data.cloudstack_zone.zone.id
  gateway = "192.168.56.1"
  netmask = "255.255.255.0"
  start_ip = "192.168.56.2"
  end_ip = "192.168.56.254"
}

resource "cloudstack_cluster" "foo" {
  name = "terraform-cluster"
  cluster_type = "CloudManaged"
  hypervisor = "KVM"
  pod_id = cloudstack_pod.foopod.id
  zone_id = data.cloudstack_zone.zone.id
  arch = "x86_64"
}`
