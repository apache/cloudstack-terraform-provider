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

func TestAccCloudStackSnapshotPolicy_basic(t *testing.T) {
	var snapshotPolicy cloudstack.SnapshotPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackSnapshotPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackSnapshotPolicyExists("cloudstack_snapshot_policy.foo", &snapshotPolicy),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "interval_type", "DAILY"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "max_snaps", "7"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "schedule", "02:30"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "timezone", "UTC"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "tags.%", "2"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "tags.Environment", "test"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "tags.Purpose", "backup"),
				),
			},
		},
	})
}

func TestAccCloudStackSnapshotPolicy_hourly(t *testing.T) {
	var snapshotPolicy cloudstack.SnapshotPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackSnapshotPolicy_hourly,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackSnapshotPolicyExists("cloudstack_snapshot_policy.hourly", &snapshotPolicy),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.hourly", "interval_type", "HOURLY"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.hourly", "max_snaps", "6"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.hourly", "schedule", "0"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.hourly", "timezone", "UTC"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.hourly", "custom_id", "test-hourly"),
				),
			},
		},
	})
}

func TestAccCloudStackSnapshotPolicy_update(t *testing.T) {
	var snapshotPolicy cloudstack.SnapshotPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackSnapshotPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackSnapshotPolicyExists("cloudstack_snapshot_policy.foo", &snapshotPolicy),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "max_snaps", "7"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "tags.Environment", "test"),
				),
			},
			{
				Config: testAccCloudStackSnapshotPolicy_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackSnapshotPolicyExists("cloudstack_snapshot_policy.foo", &snapshotPolicy),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "max_snaps", "8"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "tags.Environment", "production"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.foo", "tags.Updated", "true"),
				),
			},
		},
	})
}

func TestAccCloudStackSnapshotPolicy_weekly(t *testing.T) {
	var snapshotPolicy cloudstack.SnapshotPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackSnapshotPolicy_weekly,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackSnapshotPolicyExists("cloudstack_snapshot_policy.weekly", &snapshotPolicy),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.weekly", "interval_type", "WEEKLY"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.weekly", "schedule", "1:03:00"),
				),
			},
		},
	})
}

func TestAccCloudStackSnapshotPolicy_monthly(t *testing.T) {
	var snapshotPolicy cloudstack.SnapshotPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackSnapshotPolicy_monthly,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackSnapshotPolicyExists("cloudstack_snapshot_policy.monthly", &snapshotPolicy),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.monthly", "interval_type", "MONTHLY"),
					resource.TestCheckResourceAttr("cloudstack_snapshot_policy.monthly", "schedule", "15:01:00"),
				),
			},
		},
	})
}

func testAccCheckCloudStackSnapshotPolicyExists(n string, snapshotPolicy *cloudstack.SnapshotPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No snapshot policy ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		sp, _, err := cs.Snapshot.GetSnapshotPolicyByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if sp.Id != rs.Primary.ID {
			return fmt.Errorf("Snapshot policy not found")
		}

		*snapshotPolicy = *sp
		return nil
	}
}

func testAccCheckCloudStackSnapshotPolicyDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_snapshot_policy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No snapshot policy ID is set")
		}

		_, _, err := cs.Snapshot.GetSnapshotPolicyByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Snapshot policy %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackSnapshotPolicy_basic = `
data "cloudstack_zone" "zone" {
  filter {
    name   = "name"
    value  = "Sandbox-simulator"
  }
}

resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = data.cloudstack_zone.zone.name
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering = "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = data.cloudstack_zone.zone.name
  expunge = true
}

resource "cloudstack_disk_offering" "foo" {
  name               = "terraform-disk-offering"
  display_text       = "terraform-disk-offering"
  disk_size          = 10
}

resource "cloudstack_disk" "foo" {
  name             = "terraform-disk"
  attach           = true
  disk_offering    = cloudstack_disk_offering.foo.name
  virtual_machine_id = cloudstack_instance.foobar.id
  zone             = data.cloudstack_zone.zone.name
}

resource "cloudstack_snapshot_policy" "foo" {
  volume_id     = cloudstack_disk.foo.id
  interval_type = "DAILY"
  max_snaps     = 7
  schedule      = "02:30"
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone.id]

  tags = {
    Environment = "test"
    Purpose     = "backup"
  }
}
`

const testAccCloudStackSnapshotPolicy_update = `
data "cloudstack_zone" "zone" {
  filter {
    name   = "name"
    value  = "Sandbox-simulator"
  }
}

resource "cloudstack_network" "foo" {
  name = "terraform-network"
  display_text = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = data.cloudstack_zone.zone.name
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering = "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = data.cloudstack_zone.zone.name
  expunge = true
}

resource "cloudstack_disk_offering" "foo" {
  name               = "terraform-disk-offering"
  display_text       = "terraform-disk-offering"
  disk_size          = 10
}

resource "cloudstack_disk" "foo" {
  name             = "terraform-disk"
  attach           = true
  disk_offering    = cloudstack_disk_offering.foo.name
  virtual_machine_id = cloudstack_instance.foobar.id
  zone             = data.cloudstack_zone.zone.name
}

resource "cloudstack_snapshot_policy" "foo" {
  volume_id     = cloudstack_disk.foo.id
  interval_type = "DAILY"
  max_snaps     = 8
  schedule      = "02:30"
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone.id]

  tags = {
    Environment = "production"
    Purpose     = "backup"
    Updated     = "true"
  }
}
`

const testAccCloudStackSnapshotPolicy_hourly = `
data "cloudstack_zone" "zone" {
  filter {
    name   = "name"
    value  = "Sandbox-simulator"
  }
}

resource "cloudstack_network" "foo" {
  name = "terraform-network-hourly"
  display_text = "terraform-network-hourly"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = data.cloudstack_zone.zone.name
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test-hourly"
  display_name = "terraform-hourly"
  service_offering = "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = data.cloudstack_zone.zone.name
  expunge = true
}

resource "cloudstack_disk_offering" "foo" {
  name               = "terraform-disk-offering-hourly"
  display_text       = "terraform-disk-offering-hourly"
  disk_size          = 10
}

resource "cloudstack_disk" "foo" {
  name             = "terraform-disk-hourly"
  attach           = true
  disk_offering    = cloudstack_disk_offering.foo.name
  virtual_machine_id = cloudstack_instance.foobar.id
  zone             = data.cloudstack_zone.zone.name
}

resource "cloudstack_snapshot_policy" "hourly" {
  volume_id     = cloudstack_disk.foo.id
  interval_type = "HOURLY"
  max_snaps     = 6
  schedule      = "0"
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone.id]
  custom_id     = "test-hourly"
  
  tags = {}
}
`

const testAccCloudStackSnapshotPolicy_weekly = `
data "cloudstack_zone" "zone" {
  filter {
    name   = "name"
    value  = "Sandbox-simulator"
  }
}

resource "cloudstack_network" "foo" {
  name = "terraform-network-weekly"
  display_text = "terraform-network-weekly"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = data.cloudstack_zone.zone.name
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test-weekly"
  display_name = "terraform-weekly"
  service_offering = "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = data.cloudstack_zone.zone.name
  expunge = true
}

resource "cloudstack_disk_offering" "foo" {
  name               = "terraform-disk-offering-weekly"
  display_text       = "terraform-disk-offering-weekly"
  disk_size          = 10
}

resource "cloudstack_disk" "foo" {
  name             = "terraform-disk-weekly"
  attach           = true
  disk_offering    = cloudstack_disk_offering.foo.name
  virtual_machine_id = cloudstack_instance.foobar.id
  zone             = data.cloudstack_zone.zone.name
}

resource "cloudstack_snapshot_policy" "weekly" {
  volume_id     = cloudstack_disk.foo.id
  interval_type = "WEEKLY"
  max_snaps     = 4
  schedule      = "1:03:00"  # Monday at 3:00 AM
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone.id]
  
  tags = {}
}
`

const testAccCloudStackSnapshotPolicy_monthly = `
data "cloudstack_zone" "zone" {
  filter {
    name   = "name"
    value  = "Sandbox-simulator"
  }
}

resource "cloudstack_network" "foo" {
  name = "terraform-network-monthly"
  display_text = "terraform-network-monthly"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = data.cloudstack_zone.zone.name
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test-monthly"
  display_name = "terraform-monthly"
  service_offering = "Small Instance"
  network_id = cloudstack_network.foo.id
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = data.cloudstack_zone.zone.name
  expunge = true
}

resource "cloudstack_disk_offering" "foo" {
  name               = "terraform-disk-offering-monthly"
  display_text       = "terraform-disk-offering-monthly"
  disk_size          = 10
}

resource "cloudstack_disk" "foo" {
  name             = "terraform-disk-monthly"
  attach           = true
  disk_offering    = cloudstack_disk_offering.foo.name
  virtual_machine_id = cloudstack_instance.foobar.id
  zone             = data.cloudstack_zone.zone.name
}

resource "cloudstack_snapshot_policy" "monthly" {
  volume_id     = cloudstack_disk.foo.id
  interval_type = "MONTHLY"
  max_snaps     = 8
  schedule      = "15:01:00"  # 15th day at 1:00 AM
  timezone      = "UTC"
  zone_ids      = [data.cloudstack_zone.zone.id]
  
  tags = {}
}
`
