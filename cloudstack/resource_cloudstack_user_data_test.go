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

func TestAccCloudStackUserData_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackUserDataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackUserData_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackUserDataExists("cloudstack_user_data.foobar"),
					resource.TestCheckResourceAttr("cloudstack_user_data.foobar", "name", "terraform-test-userdata"),
					resource.TestCheckResourceAttr("cloudstack_user_data.foobar", "user_data", "#!/bin/bash\necho 'Hello World'\n"),
				),
			},
		},
	})
}

func TestAccCloudStackUserData_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackUserDataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackUserData_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackUserDataExists("cloudstack_user_data.foobar"),
					resource.TestCheckResourceAttr("cloudstack_user_data.foobar", "name", "terraform-test-userdata"),
					resource.TestCheckResourceAttr("cloudstack_user_data.foobar", "user_data", "#!/bin/bash\necho 'Hello World'\n"),
				),
			},
			{
				Config: testAccCloudStackUserData_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackUserDataExists("cloudstack_user_data.foobar"),
					resource.TestCheckResourceAttr("cloudstack_user_data.foobar", "name", "terraform-test-userdata"),
					resource.TestCheckResourceAttr("cloudstack_user_data.foobar", "user_data", "#!/bin/bash\necho 'Updated Hello World'\n"),
				),
			},
		},
	})
}

func testAccCheckCloudStackUserDataExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No UserData ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		_, count, err := cs.User.GetUserDataByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if count == 0 {
			return fmt.Errorf("UserData not found")
		}

		return nil
	}
}

func testAccCheckCloudStackUserDataDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_user_data" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No UserData ID is set")
		}

		_, count, err := cs.User.GetUserDataByID(rs.Primary.ID)

		if err == nil && count != 0 {
			return fmt.Errorf("UserData %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackUserData_basic = `
resource "cloudstack_user_data" "foobar" {
  name      = "terraform-test-userdata"
  user_data = <<-EOF
    #!/bin/bash
    echo 'Hello World'
  EOF
}
`

const testAccCloudStackUserData_update = `
resource "cloudstack_user_data" "foobar" {
  name      = "terraform-test-userdata"
  user_data = <<-EOF
    #!/bin/bash
    echo 'Updated Hello World'
  EOF
}
`
