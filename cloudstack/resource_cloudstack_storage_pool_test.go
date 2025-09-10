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

// Simulator does not support storage pool creation
// Error: CloudStack API error 530 (CSExceptionErrorCode: 9999): Failed to add data store: No host up to associate a storage pool with in cluster 1
func TestAccCloudStackStoragePool_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps:     []resource.TestStep{
			// {
			// 	Config: testAccCloudStackStoragePoolConfig_basic,
			// },
			// {
			// 	Config: testAccCloudStackStoragePoolConfig_update,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr("cloudstack_storage_pool.test", "name", "accprimarystorage1"),
			// 	),
			// },
		},
	})
}

const testAccCloudStackStoragePoolConfig_basic = `
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
	gateway          = "172.31.0.1"
	name             = "acc_pod"
	netmask          = "255.255.240.0"
	start_ip         =  "172.31.0.2"
	zone_id          =  cloudstack_zone.test.id
}
resource "cloudstack_cluster" "test" {
	cluster_name = "acc_cluster"
	cluster_type = "CloudManaged"
	hypervisor   = "KVM"
	pod_id       = cloudstack_pod.test.id
	zone_id      = cloudstack_zone.test.id
}

resource "cloudstack_storage_pool" "test" {
	name         = "acc_primarystorage"
	url          = "nfs://10.147.28.6/export/home/sandbox/primary11"
	zone_id      = cloudstack_zone.test.id
	cluster_id   = cloudstack_cluster.test.id
	pod_id       = cloudstack_pod.test.id
	scope        = "CLUSTER"
	hypervisor   = "Simulator"
	tags         = "XYZ,123"
}
`

const testAccCloudStackStoragePoolConfig_update = `
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
	gateway          = "172.31.0.1"
	name             = "acc_pod"
	netmask          = "255.255.240.0"
	start_ip         =  "172.31.0.2"
	zone_id          =  cloudstack_zone.test.id
}
resource "cloudstack_cluster" "test" {
	cluster_name = "acc_cluster"
	cluster_type = "CloudManaged"
	hypervisor   = "KVM"
	pod_id       = cloudstack_pod.test.id
	zone_id      = cloudstack_zone.test.id
}

resource "cloudstack_storage_pool" "test" {
	name         = "acc_primarystorage1"
	url          = "nfs://10.147.28.6/export/home/sandbox/primary11"
	zone_id      = cloudstack_zone.test.id
	cluster_id   = cloudstack_cluster.test.id
	pod_id       = cloudstack_pod.test.id
	scope        = "CLUSTER"
	hypervisor   = "Simulator"
	state        = "Maintenance"
	tags         = "XYZ,123,456"
}
`
