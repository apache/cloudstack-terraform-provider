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

// func TestAccCloudStackStoragePool_basic(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCloudStackStoragePoolConfig_basic,
// 			},
// 			{
// 				Config: testAccCloudStackStoragePoolConfig_update,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("cloudstack_storage_pool.test", "name", "accprimarystorage1"),
// 				),
// 			},
// 		},
// 	})
// }

const testAccCloudStackStoragePoolConfig_basic = `
resource "cloudstack_storage_pool" "test" {
	name         = "accprimarystorage"
	url          = "nfs://10.147.28.6/export/home/sandbox/primary11"
	zone_id      = "0ed38eb3-f279-4951-ac20-fef39ebab20c"
	cluster_id   = "9daeeb36-d8b7-497a-9b53-bbebba88c817"
	pod_id       = "2ff52b73-139e-4c40-a0a3-5b7d87d8e3c4"
	scope        = "CLUSTER"
	hypervisor   = "Simulator"
	tags         = "XYZ,123"
}
`

const testAccCloudStackStoragePoolConfig_update = `
resource "cloudstack_storage_pool" "test" {
	name         = "accprimarystorage1"
	url          = "nfs://10.147.28.6/export/home/sandbox/primary11"
	zone_id      = "0ed38eb3-f279-4951-ac20-fef39ebab20c"
	cluster_id   = "9daeeb36-d8b7-497a-9b53-bbebba88c817"
	pod_id       = "2ff52b73-139e-4c40-a0a3-5b7d87d8e3c4"
	scope        = "CLUSTER"
	hypervisor   = "Simulator"
	state        = "Maintenance"
	tags         = "XYZ,123,456"
}
`
