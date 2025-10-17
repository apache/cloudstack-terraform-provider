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

const (
	testServiceOfferingResourceName = "cloudstack_service_offering"
	testServiceOfferingBasic        = testServiceOfferingResourceName + ".test1"
	testServiceOfferingGPU          = testServiceOfferingResourceName + ".gpu"
	testServiceOfferingCustom       = testServiceOfferingResourceName + ".custom"
	testServiceOfferingDiskOpt      = testServiceOfferingResourceName + ".disk_optimized"
	testServiceOfferingHighPriority = testServiceOfferingResourceName + ".high_priority"
)

func TestAccServiceOfferingBasic(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingBasicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingBasic, &so),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "cpu_number", "2"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "cpu_speed", "2200"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "memory", "8096"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "storage_type", "shared"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "offer_ha", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "limit_cpu_use", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "disk_iops_read_rate", "10000"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "disk_iops_write_rate", "10000"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "min_iops", "5000"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "max_iops", "15000"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "dynamic_scaling_enabled", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "is_volatile", "false"),
					resource.TestCheckResourceAttr(testServiceOfferingBasic, "root_disk_size", "50"),
				),
			},
		},
	})
}

func TestAccServiceOfferingWithGPU(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingGPUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingGPU, &so),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "name", "gpu_offering_1"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "cpu_number", "4"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "cpu_speed", "2400"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "memory", "16384"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "gpu_count", "1"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "gpu_display", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "storage_type", "shared"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "offer_ha", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "disk_iops_read_rate", "20000"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "disk_iops_write_rate", "20000"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "min_iops", "10000"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "max_iops", "30000"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "dynamic_scaling_enabled", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "root_disk_size", "100"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingBasicConfig = `
resource "cloudstack_service_offering" "test1" {
  name                     = "service_offering_1"
  display_text            = "Test Basic Offering"
  cpu_number             = 2
  cpu_speed              = 2200
  memory                 = 8096
  storage_type           = "shared"
  offer_ha               = true
  limit_cpu_use         = true
  disk_iops_read_rate   = 10000
  disk_iops_write_rate  = 10000
  min_iops              = 5000
  max_iops              = 15000
  dynamic_scaling_enabled = true
  is_volatile           = false
  root_disk_size        = 50
}
`

const testAccCloudStackServiceOfferingGPUConfig = `
resource "cloudstack_service_offering" "gpu" {
  name                     = "gpu_offering_1"
  display_text            = "Test GPU Offering"
  cpu_number             = 4
  cpu_speed              = 2400
  memory                 = 16384
  gpu_count              = 1
  gpu_display            = true
  storage_type           = "shared"
  offer_ha               = true
  disk_iops_read_rate   = 20000
  disk_iops_write_rate  = 20000
  min_iops              = 10000
  max_iops              = 30000
  dynamic_scaling_enabled = true
  is_volatile           = false
  root_disk_size        = 100
}
`

func testAccCheckCloudStackServiceOfferingExists(n string, so *cloudstack.ServiceOffering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No service offering ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		resp, _, err := cs.ServiceOffering.GetServiceOfferingByID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting service offering: %s", err)
		}

		if resp == nil {
			return fmt.Errorf("Service offering (%s) not found", rs.Primary.ID)
		}

		if resp.Id != rs.Primary.ID {
			return fmt.Errorf("Service offering not found: expected ID %s, got %s", rs.Primary.ID, resp.Id)
		}

		*so = *resp

		return nil
	}
}

func testAccCheckCloudStackServiceOfferingDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testServiceOfferingResourceName {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		resp, _, err := cs.ServiceOffering.GetServiceOfferingByID(rs.Primary.ID)
		if err != nil {
			// CloudStack returns 431 error code when the resource doesn't exist
			// Just return nil in this case as the resource is gone
			return nil
		}

		if resp != nil {
			return fmt.Errorf("Service offering %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func TestAccServiceOfferingCustomized(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingCustomConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingCustom, &so),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "customized", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "min_cpu_number", "1"),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "max_cpu_number", "8"),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "min_memory", "1024"),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "max_memory", "16384"),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "cpu_speed", "1000"),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "encrypt_root", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingCustom, "storage_tags", "production,ssd"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingCustomConfig = `
resource "cloudstack_service_offering" "custom" {
  name             = "custom_service_offering"
  display_text     = "Custom Test"
  customized       = true
  min_cpu_number   = 1
  max_cpu_number   = 8
  min_memory       = 1024
  max_memory       = 16384
  cpu_speed        = 1000
  encrypt_root     = true
  storage_tags     = "production,ssd"
}
`

func TestAccServiceOfferingWithVGPU(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingVGPUConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingGPU, &so),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "name", "gpu_service_offering"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "display_text", "GPU Test"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "cpu_number", "4"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "memory", "16384"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "service_offering_details.pciDevice", "Group of NVIDIA A6000 GPUs"),
					resource.TestCheckResourceAttr(testServiceOfferingGPU, "service_offering_details.vgpuType", "A6000-8A"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingVGPUConfig = `
resource "cloudstack_service_offering" "gpu" {
  name         = "gpu_service_offering"
  display_text = "GPU Test"
  cpu_number   = 4
  memory       = 16384
  cpu_speed    = 1000
  
  service_offering_details = {
    pciDevice = "Group of NVIDIA A6000 GPUs"
    vgpuType  = "A6000-8A"
  }
}
`

func TestAccServiceOfferingDiskOptimized(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingDiskOptimizedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingDiskOpt, &so),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "name", "disk_optimized_offering"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "cpu_number", "4"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "cpu_speed", "2000"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "memory", "4096"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "root_disk_size", "100"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "provisioning_type", "thin"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "encrypt_root", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "min_iops", "1000"),
					resource.TestCheckResourceAttr(testServiceOfferingDiskOpt, "max_iops", "5000"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingDiskOptimizedConfig = `
resource "cloudstack_service_offering" "disk_optimized" {
  name              = "disk_optimized_offering"
  display_text      = "Disk Optimized Test"
  cpu_number        = 4
  cpu_speed         = 2000
  memory            = 4096
  storage_type      = "shared"
  root_disk_size    = 100
  provisioning_type = "thin"
  encrypt_root      = true
  min_iops          = 1000
  max_iops          = 5000
}
`

func TestAccServiceOfferingHighPriority(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingHighPriorityConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingHighPriority, &so),
					resource.TestCheckResourceAttr(testServiceOfferingHighPriority, "name", "high_priority_offering"),
					resource.TestCheckResourceAttr(testServiceOfferingHighPriority, "limit_cpu_use", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingHighPriority, "is_volatile", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingHighPriority, "customized_iops", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingHighPriority, "tags", "production,tier1"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingHighPriorityConfig = `
resource "cloudstack_service_offering" "high_priority" {
  name            = "high_priority_offering"
  display_text    = "High Priority Parameters Test"
  cpu_number      = 4
  cpu_speed       = 3000
  memory          = 8192
  storage_type    = "shared"
  limit_cpu_use   = true
  is_volatile     = true
  customized_iops = true
  tags            = "production,tier1"
}
`

// Phase 2: IOPS/Bandwidth Parameters Tests
func TestAccServiceOfferingPerformance(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingPerformanceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingResourceName+".performance", &so),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "name", "performance_offering"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_iops_read_rate", "2000"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_iops_write_rate", "1000"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_iops_read_rate_max", "5000"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_iops_write_rate_max", "2500"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_iops_read_rate_max_length", "60"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_iops_write_rate_max_length", "60"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_bytes_read_rate", "209715200"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_bytes_write_rate", "104857600"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_bytes_read_rate_max", "524288000"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "disk_bytes_write_rate_max", "262144000"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "bytes_read_rate_max_length", "120"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".performance", "bytes_write_rate_max_length", "120"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingPerformanceConfig = `
resource "cloudstack_service_offering" "performance" {
  name         = "performance_offering"
  display_text = "Performance with IOPS and Bandwidth Limits"
  cpu_number   = 8
  cpu_speed    = 3000
  memory       = 16384
  storage_type = "shared"

  # IOPS limits
  disk_iops_read_rate           = 2000
  disk_iops_write_rate          = 1000
  disk_iops_read_rate_max       = 5000
  disk_iops_write_rate_max      = 2500
  disk_iops_read_rate_max_length = 60
  disk_iops_write_rate_max_length = 60

  # Bandwidth limits (bytes/sec)
  disk_bytes_read_rate          = 209715200  # 200 MB/s
  disk_bytes_write_rate         = 104857600  # 100 MB/s
  disk_bytes_read_rate_max      = 524288000  # 500 MB/s burst
  disk_bytes_write_rate_max     = 262144000  # 250 MB/s burst
  bytes_read_rate_max_length    = 120
  bytes_write_rate_max_length   = 120
}
`

// Phase 3: Hypervisor Parameters Tests
func TestAccServiceOfferingHypervisor(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingHypervisorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingResourceName+".hypervisor", &so),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".hypervisor", "name", "hypervisor_offering"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".hypervisor", "hypervisor_snapshot_reserve", "25"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".hypervisor", "cache_mode", "writeback"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".hypervisor", "deployment_planner", "FirstFitPlanner"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".hypervisor", "storage_policy", "VMware-VSAN-Gold"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingHypervisorConfig = `
resource "cloudstack_service_offering" "hypervisor" {
  name         = "hypervisor_offering"
  display_text = "Hypervisor-Optimized Offering"
  cpu_number   = 16
  cpu_speed    = 3200
  memory       = 32768
  storage_type = "shared"

  # Hypervisor-specific parameters
  hypervisor_snapshot_reserve = 25
  cache_mode                  = "writeback"
  deployment_planner          = "FirstFitPlanner"
  storage_policy              = "VMware-VSAN-Gold"
}
`

// Phase 4: Low Priority Parameters Tests
func TestAccServiceOfferingLowPriority(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingLowPriorityConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingResourceName+".low_priority", &so),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".low_priority", "name", "low_priority_offering"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".low_priority", "network_rate", "1000"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".low_priority", "purge_resources", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".low_priority", "system_vm_type", "domainrouter"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingLowPriorityConfig = `
resource "cloudstack_service_offering" "low_priority" {
  name         = "low_priority_offering"
  display_text = "System VM and Network Rate Test"
  cpu_number   = 1
  cpu_speed    = 1000
  memory       = 512
  storage_type = "local"

  # Low priority parameters
  network_rate    = 1000
  purge_resources = true
  system_vm_type  = "domainrouter"
}
`

// Phase 5: Final SDK Parameters Tests
func TestAccServiceOfferingFinalSDK(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingFinalSDKConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingResourceName+".final_sdk", &so),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".final_sdk", "name", "final_sdk_offering"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".final_sdk", "lease_duration", "86400"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".final_sdk", "lease_expiry_action", "destroy"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".final_sdk", "is_system", "false"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".final_sdk", "disk_offering_strictness", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".final_sdk", "external_details.cmdb_id", "CMDB-12345"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".final_sdk", "external_details.billing_code", "TIER2"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingFinalSDKConfig = `
resource "cloudstack_service_offering" "final_sdk" {
  name         = "final_sdk_offering"
  display_text = "Final SDK Parameters Test"
  cpu_number   = 4
  cpu_speed    = 2400
  memory       = 8192
  storage_type = "shared"

  # Final SDK parameters
  lease_duration           = 86400
  lease_expiry_action      = "destroy"
  is_system                = false
  disk_offering_strictness = true

  external_details = {
    cmdb_id      = "CMDB-12345"
    billing_code = "TIER2"
  }
}
`

// Complete Test: All Parameters
func TestAccServiceOfferingComplete(t *testing.T) {
	var so cloudstack.ServiceOffering
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackServiceOfferingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackServiceOfferingCompleteConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackServiceOfferingExists(testServiceOfferingResourceName+".complete", &so),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "name", "complete_offering"),
					// Core parameters
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "cpu_number", "16"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "memory", "65536"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "storage_type", "shared"),
					// HA
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "offer_ha", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "dynamic_scaling_enabled", "true"),
					// Phase 1
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "limit_cpu_use", "true"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "is_volatile", "false"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "tags", "production,tier1"),
					// Phase 2 - IOPS
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "disk_iops_read_rate", "3000"),
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "disk_bytes_read_rate", "314572800"),
					// Phase 3 - Hypervisor
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "cache_mode", "writeback"),
					// Phase 4 - Network
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "network_rate", "10000"),
					// Phase 5 - Lease
					resource.TestCheckResourceAttr(testServiceOfferingResourceName+".complete", "lease_duration", "604800"),
				),
			},
		},
	})
}

const testAccCloudStackServiceOfferingCompleteConfig = `
resource "cloudstack_service_offering" "complete" {
  name         = "complete_offering"
  display_text = "Complete Test - All 54 Parameters"
  
  # Core compute
  cpu_number   = 16
  cpu_speed    = 3200
  memory       = 65536
  storage_type = "shared"
  root_disk_size = 100

  # HA and scaling
  offer_ha                     = true
  dynamic_scaling_enabled      = true
  encrypt_root                 = true

  # Phase 1: HIGH Priority
  limit_cpu_use   = true
  is_volatile     = false
  customized_iops = true
  tags            = "production,tier1"

  # Phase 2: IOPS/Bandwidth
  disk_iops_read_rate           = 3000
  disk_iops_write_rate          = 1500
  disk_iops_read_rate_max       = 6000
  disk_iops_write_rate_max      = 3000
  disk_iops_read_rate_max_length = 90
  disk_iops_write_rate_max_length = 90
  disk_bytes_read_rate          = 314572800  # 300 MB/s
  disk_bytes_write_rate         = 157286400  # 150 MB/s
  disk_bytes_read_rate_max      = 629145600  # 600 MB/s
  disk_bytes_write_rate_max     = 314572800  # 300 MB/s
  bytes_read_rate_max_length    = 180
  bytes_write_rate_max_length   = 180

  # Phase 3: Hypervisor
  hypervisor_snapshot_reserve = 30
  cache_mode                  = "writeback"
  deployment_planner          = "UserConcentratedPodPlanner"
  storage_policy              = "Premium-Storage"

  # Phase 4: Low Priority
  network_rate    = 10000
  purge_resources = false

  # Phase 5: Final SDK
  lease_duration      = 604800  # 7 days
  lease_expiry_action = "stop"
  is_system           = false

  # Metadata
  service_offering_details = {
    tier = "platinum"
    sla  = "99.99"
  }

  external_details = {
    cmdb_id        = "ASSET-PROD-001"
    cost_center    = "ENGINEERING"
    compliance     = "SOC2"
  }
}
`
