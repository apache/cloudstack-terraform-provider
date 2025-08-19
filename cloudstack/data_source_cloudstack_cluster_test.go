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

func TestAccClusterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testClusterDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudstack_cluster.test", "name", "terraform-test-cluster"),
				),
			},
		},
	})
}

const testClusterDataSourceConfig_basic = `
data "cloudstack_zone" "zone" {
  filter {
    name = "name"
    value = "Sandbox-simulator"
  }
}

data "cloudstack_pod" "pod" {
  filter {
    name = "name"
    value = "POD0"
  }
}

# Create a cluster first
resource "cloudstack_cluster" "test_cluster" {
  name = "terraform-test-cluster"
  cluster_type = "CloudManaged"
  hypervisor = "KVM"
  pod_id = data.cloudstack_pod.pod.id
  zone_id = data.cloudstack_zone.zone.id
  arch = "x86_64"
}

# Then query it with the data source
data "cloudstack_cluster" "test" {
  filter {
    name = "name"
    value = "terraform-test-cluster"
  }
  depends_on = [cloudstack_cluster.test_cluster]
}
`
