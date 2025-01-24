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

func TestAccCloudStackCluster_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackCluster_basic,
			},
		},
	})
}

const testAccCloudStackCluster_basic = `
resource "cloudstack_zone" "test" {
	name          = "acc_zone"
	dns1          = "8.8.8.8"
	dns2          = "8.8.8.8"
	internal_dns1 = "8.8.4.4"
	internal_dns2 = "8.8.4.4"
	network_type  = "Advanced"
	domain        = "cloudstack.apache.org"
}
resource "cloudstack_pod" "test" {
	allocation_state = "Disabled"
	gateway          = "172.30.0.1"
	name             = "acc_pod"
	netmask          = "255.255.240.0"
	start_ip         = "172.30.0.2"
	zone_id          = cloudstack_zone.test.id
}
resource "cloudstack_cluster" "test" {
	cluster_name = "acc_cluster"
	cluster_type = "CloudManaged"
	hypervisor   = "KVM"
	pod_id       = cloudstack_pod.test.id
	zone_id      = cloudstack_zone.test.id
}
`
