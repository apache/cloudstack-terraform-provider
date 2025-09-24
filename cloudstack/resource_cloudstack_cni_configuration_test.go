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

func TestAccCloudStackCniConfiguration_basic(t *testing.T) {
	var cniConfig cloudstack.CniConfiguration

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackCniConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackCniConfiguration_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackCniConfigurationExists("cloudstack_cni_configuration.foo", &cniConfig),
					resource.TestCheckResourceAttr("cloudstack_cni_configuration.foo", "name", "test-cni-config"),
					resource.TestCheckResourceAttr("cloudstack_cni_configuration.foo", "params.#", "2"),
				),
			},
		},
	})
}

func testAccCheckCloudStackCniConfigurationExists(n string, cniConfig *cloudstack.CniConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CNI configuration ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		config, _, err := cs.Configuration.GetCniConfigurationByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if config.Id != rs.Primary.ID {
			return fmt.Errorf("CNI configuration not found")
		}

		*cniConfig = *config
		return nil
	}
}

func testAccCheckCloudStackCniConfigurationDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_cni_configuration" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CNI configuration ID is set")
		}

		_, _, err := cs.Configuration.GetCniConfigurationByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CNI configuration %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackCniConfiguration_basic = `
resource "cloudstack_cni_configuration" "foo" {
  name       = "test-cni-config"
  cni_config = <<EOF
{
  "cniVersion": "0.4.0",
  "name": "test-network",
  "type": "bridge",
  "bridge": "cni0",
  "isGateway": true,
  "ipMasq": true,
  "ipam": {
    "type": "host-local",
    "subnet": "10.244.0.0/16",
    "routes": [
      { "dst": "0.0.0.0/0" }
    ]
  }
}
EOF
  
  params = ["subnet", "gateway"]
}
`
