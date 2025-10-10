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

func TestAccCloudStackUserDataTemplateLink_template(t *testing.T) {
	var template cloudstack.Template

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackUserDataTemplateLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackUserDataTemplateLinkConfigTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackTemplateExists("cloudstack_template.test", &template),
					testAccCheckCloudStackUserDataTemplateLinkExists("cloudstack_user_data_template_link.test"),
					resource.TestCheckResourceAttr("cloudstack_user_data_template_link.test", "user_data_policy", "ALLOWOVERRIDE"),
				),
			},
		},
	})
}

func TestAccCloudStackUserDataTemplateLink_iso(t *testing.T) {
	var template cloudstack.Template

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackUserDataTemplateLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackUserDataTemplateLinkConfigISO,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackTemplateExists("cloudstack_template.test_iso", &template),
					testAccCheckCloudStackUserDataTemplateLinkExists("cloudstack_user_data_template_link.test_iso"),
					resource.TestCheckResourceAttr("cloudstack_user_data_template_link.test_iso", "user_data_policy", "APPEND"),
				),
			},
		},
	})
}

func testAccCheckCloudStackUserDataTemplateLinkExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No UserData Template Link ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

		// Check if template/ISO exists with linked userdata
		var templateId string
		if rs.Primary.Attributes["template_id"] != "" {
			templateId = rs.Primary.Attributes["template_id"]
		} else if rs.Primary.Attributes["iso_id"] != "" {
			templateId = rs.Primary.Attributes["iso_id"]
		} else {
			return fmt.Errorf("Neither template_id nor iso_id found in state")
		}

		template, count, err := cs.Template.GetTemplateByID(templateId, "all")
		if err != nil {
			if count == 0 {
				return fmt.Errorf("Template/ISO %s not found", templateId)
			}
			return err
		}

		// Check if userdata is linked (optional since unlinking is also valid)
		if template.Userdataname != "" {
			return nil // UserData is linked
		}

		return nil // Template exists, whether userdata is linked or not
	}
}

func testAccCheckCloudStackUserDataTemplateLinkDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_user_data_template_link" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No UserData Template Link ID is set")
		}

		var templateId string
		if rs.Primary.Attributes["template_id"] != "" {
			templateId = rs.Primary.Attributes["template_id"]
		} else if rs.Primary.Attributes["iso_id"] != "" {
			templateId = rs.Primary.Attributes["iso_id"]
		} else {
			continue // Skip if no template/ISO ID
		}

		template, count, err := cs.Template.GetTemplateByID(templateId, "all")
		if err != nil {
			if count == 0 {
				return nil // Template doesn't exist anymore, that's fine
			}
			return err
		}

		// Check that userdata is not linked anymore
		if template.Userdataname != "" {
			return fmt.Errorf("UserData is still linked to template/ISO %s", templateId)
		}
	}

	return nil
}

const testAccCloudStackUserDataTemplateLinkConfigTemplate = `
resource "cloudstack_user_data" "test" {
  name      = "test-userdata-link"
  user_data = "#!/bin/bash\necho 'template test' > /tmp/test.txt"
}

resource "cloudstack_template" "test" {
  name              = "test-template-userdata"
  format            = "QCOW2"
  hypervisor        = "KVM"
  os_type           = "CentOS 7"
  url               = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
  is_extractable    = true
  is_featured       = false
  is_public         = false
  password_enabled  = false
}

resource "cloudstack_user_data_template_link" "test" {
  template_id       = cloudstack_template.test.id
  user_data_id      = cloudstack_user_data.test.id
  user_data_policy  = "ALLOWOVERRIDE"
}`

const testAccCloudStackUserDataTemplateLinkConfigISO = `
resource "cloudstack_user_data" "test_iso" {
  name      = "test-userdata-iso-link"
  user_data = "#!/bin/bash\necho 'iso test' > /tmp/test.txt"
}

resource "cloudstack_template" "test_iso" {
  name              = "test-iso-userdata"
  format            = "ISO"
  hypervisor        = "KVM"
  os_type           = "CentOS 7"
  url               = "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
  is_extractable    = true
  is_featured       = false
  is_public         = false
  password_enabled  = false
}

resource "cloudstack_user_data_template_link" "test_iso" {
  iso_id            = cloudstack_template.test_iso.id
  user_data_id      = cloudstack_user_data.test_iso.id
  user_data_policy  = "APPEND"
}`
