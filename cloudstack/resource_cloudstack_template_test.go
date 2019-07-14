package cloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func TestAccCloudStackTemplate_basic(t *testing.T) {
	if CLOUDSTACK_TEMPLATE_URL == "" {
		t.Skip("This test requires an upload URL")
	}

	var template cloudstack.Template

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackTemplate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackTemplateExists("cloudstack_template.foo", &template),
					testAccCheckCloudStackTemplateBasicAttributes(&template),
					resource.TestCheckResourceAttr(
						"cloudstack_template.foo", "display_text", "terraform-test"),
					testAccCheckResourceTags(&template),
				),
			},
		},
	})
}

func TestAccCloudStackTemplate_update(t *testing.T) {
	if CLOUDSTACK_TEMPLATE_URL == "" {
		t.Skip("This test requires an upload URL")
	}

	var template cloudstack.Template

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackTemplate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackTemplateExists("cloudstack_template.foo", &template),
					testAccCheckCloudStackTemplateBasicAttributes(&template),
				),
			},

			{
				Config: testAccCloudStackTemplate_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackTemplateExists(
						"cloudstack_template.foo", &template),
					testAccCheckCloudStackTemplateUpdatedAttributes(&template),
					resource.TestCheckResourceAttr(
						"cloudstack_template.foo", "display_text", "terraform-updated"),
				),
			},
		},
	})
}

func testAccCheckCloudStackTemplateExists(
	n string, template *cloudstack.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No template ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		tmpl, _, err := cs.Template.GetTemplateByID(rs.Primary.ID, "executable")

		if err != nil {
			return err
		}

		if tmpl.Id != rs.Primary.ID {
			return fmt.Errorf("Template not found")
		}

		*template = *tmpl

		return nil
	}
}

func testAccCheckCloudStackTemplateBasicAttributes(
	template *cloudstack.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if template.Name != "terraform-test" {
			return fmt.Errorf("Bad name: %s", template.Name)
		}

		if template.Format != "QCOW2" {
			return fmt.Errorf("Bad format: %s", template.Format)
		}

		if template.Hypervisor != "Simulator" {
			return fmt.Errorf("Bad hypervisor: %s", template.Hypervisor)
		}

		if template.Ostypename != "Centos 5.6 (64-bit)" {
			return fmt.Errorf("Bad os type: %s", template.Ostypename)
		}

		if template.Zonename != "Sandbox-simulator" {
			return fmt.Errorf("Bad zone: %s", template.Zonename)
		}

		return nil
	}
}

func testAccCheckCloudStackTemplateUpdatedAttributes(
	template *cloudstack.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if template.Displaytext != "terraform-updated" {
			return fmt.Errorf("Bad name: %s", template.Displaytext)
		}

		if !template.Isdynamicallyscalable {
			return fmt.Errorf("Bad is_dynamically_scalable: %t", template.Isdynamicallyscalable)
		}

		if !template.Passwordenabled {
			return fmt.Errorf("Bad password_enabled: %t", template.Passwordenabled)
		}

		return nil
	}
}

func testAccCheckCloudStackTemplateDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_template" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No template ID is set")
		}

		_, _, err := cs.Template.GetTemplateByID(rs.Primary.ID, "executable")
		if err == nil {
			return fmt.Errorf("Template %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

var testAccCloudStackTemplate_basic = fmt.Sprintf(`
resource "cloudstack_template" "foo" {
  name = "terraform-test"
  format = "QCOW2"
  hypervisor = "Simulator"
  os_type = "Centos 5.6 (64-bit)"
  url = "%s"
  zone = "Sandbox-simulator"
  tags = {
    terraform-tag = "true"
  }
}`, CLOUDSTACK_TEMPLATE_URL)

var testAccCloudStackTemplate_update = fmt.Sprintf(`
resource "cloudstack_template" "foo" {
  name = "terraform-test"
  display_text = "terraform-updated"
  format = "QCOW2"
  hypervisor = "Simulator"
  os_type = "Centos 5.6 (64-bit)"
  url = "%s"
  is_dynamically_scalable = true
  password_enabled = true
  zone = "Sandbox-simulator"
}`, CLOUDSTACK_TEMPLATE_URL)
