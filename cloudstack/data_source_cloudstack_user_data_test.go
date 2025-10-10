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

func TestAccDataSourceCloudStackUserData_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackUserDataDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudStackUserDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudstack_user_data.test", "name", "terraform-test-userdata"),
					resource.TestCheckResourceAttr("data.cloudstack_user_data.test", "user_data", "#!/bin/bash\\necho 'Hello World'\\n"),
				),
			},
		},
	})
}

const testAccDataSourceCloudStackUserDataConfig = `
resource "cloudstack_user_data" "test" {
  name      = "terraform-test-userdata"
  user_data = <<-EOF
    #!/bin/bash
    echo 'Hello World'
  EOF
}

data "cloudstack_user_data" "test" {
  name = cloudstack_user_data.test.name
}
`
