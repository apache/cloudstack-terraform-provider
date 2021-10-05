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
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudStackAutoscaleVMProfile_basic(t *testing.T) {
	var vmProfile cloudstack.AutoScaleVmProfile

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackAutoscaleVMProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackAutoscaleVMProfile_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackAutoscaleVMProfileExists("cloudstack_autoscale_vm_profile.foo", &vmProfile),
					testAccCheckCloudStackAutoscaleVMProfileBasicAttributes(&vmProfile),
					resource.TestCheckResourceAttr(
						"cloudstack_autoscale_vm_profile.foo", "zone", "Sandbox-simulator"),
					testAccCheckResourceMetadata(&vmProfile),
				),
			},
		},
	})
}

func TestAccCloudStackAutoscaleVMProfile_update(t *testing.T) {
	var vmProfile cloudstack.AutoScaleVmProfile

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackAutoscaleVMProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackAutoscaleVMProfile_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackAutoscaleVMProfileExists("cloudstack_autoscale_vm_profile.foo", &vmProfile),
					testAccCheckCloudStackAutoscaleVMProfileBasicAttributes(&vmProfile),
				),
			},

			{
				Config: testAccCloudStackAutoscaleVMProfile_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackAutoscaleVMProfileExists(
						"cloudstack_autoscale_vm_profile.foo", &vmProfile),
					testAccCheckCloudStackAutoscaleVMProfileUpdatedAttributes(&vmProfile),
					resource.TestCheckResourceAttr(
						"cloudstack_autoscale_vm_profile.foo", "zone", "Sandbox-simulator"),
				),
			},
		},
	})
}

func testAccCheckResourceMetadata(vmProfile *cloudstack.AutoScaleVmProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Resourcemetadata.NewListResourceDetailsParams("AutoScaleVmProfile")
		p.SetResourceid(vmProfile.Id)
		response, err := cs.Resourcemetadata.ListResourceDetails(p)
		if err != nil {
			return err
		}
		metadata := make(map[string]string)
		for _, item := range response.ResourceDetails {
			metadata[item.Key] = item.Value
		}
		return testAccCheckTags(metadata, "terraform-meta", "true")
	}
}

func testAccCheckCloudStackAutoscaleVMProfileExists(
	n string, vmProfile *cloudstack.AutoScaleVmProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No vmProfile ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		avp, _, err := cs.AutoScale.GetAutoScaleVmProfileByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if avp.Id != rs.Primary.ID {
			return fmt.Errorf("AutoScaleVMProfile not found")
		}

		*vmProfile = *avp

		return nil
	}
}

func testAccCheckCloudStackAutoscaleVMProfileBasicAttributes(
	vmProfile *cloudstack.AutoScaleVmProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

		serviceofferingid, e := retrieveID(cs, "service_offering", "Small Instance")
		if e != nil {
			return e.Error()
		}

		zoneid, e := retrieveID(cs, "zone", "Sandbox-simulator")
		if e != nil {
			return e.Error()
		}

		templateid, e := retrieveTemplateID(cs, zoneid, "CentOS 5.6 (64-bit) no GUI (Simulator)")
		if e != nil {
			return e.Error()
		}

		if vmProfile.Serviceofferingid != serviceofferingid {
			return fmt.Errorf("Bad offering: %s", vmProfile.Serviceofferingid)
		}

		if vmProfile.Templateid != templateid {
			return fmt.Errorf("Bad template: %s", vmProfile.Templateid)
		}

		if vmProfile.Zoneid != zoneid {
			return fmt.Errorf("Bad zone: %s", vmProfile.Zoneid)
		}

		if vmProfile.Otherdeployparams != "displayname=display1&networkids=net1" {
			return fmt.Errorf("Bad otherdeployparams: %s", vmProfile.Otherdeployparams)
		}

		return nil
	}
}

func testAccCheckCloudStackAutoscaleVMProfileUpdatedAttributes(
	vmProfile *cloudstack.AutoScaleVmProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if vmProfile.Destroyvmgraceperiod != 10 {
			return fmt.Errorf("Bad destroy_vm_grace_period: %d", vmProfile.Destroyvmgraceperiod)
		}

		return nil
	}
}

func testAccCheckCloudStackAutoscaleVMProfileDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_autoscale_vm_profile" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No vmProfile ID is set")
		}

		_, _, err := cs.AutoScale.GetAutoScaleVmProfileByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("AutoScaleVMProfile %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCloudStackAutoscaleVMProfile_basic = `
resource "cloudstack_autoscale_vm_profile" "foo" {
  service_offering = "Small Instance"
  template         = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone             = "Sandbox-simulator"

  other_deploy_params = {
    networkids  = "net1"
    displayname = "display1"
  }

  metadata = {
    terraform-meta = "true"
  }
}`

var testAccCloudStackAutoscaleVMProfile_update = `
resource "cloudstack_autoscale_vm_profile" "foo" {
  service_offering        = "Small Instance"
  template                = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone                    = "Sandbox-simulator"
  destroy_vm_grace_period = "10s"

  other_deploy_params = {
    networkids  = "net1"
    displayname = "display1"
  }
}`
