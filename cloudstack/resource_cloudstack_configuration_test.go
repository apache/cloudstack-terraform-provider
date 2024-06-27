package cloudstack

import (
	"fmt"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudStackConfiguration_basic(t *testing.T) {
	var configuration cloudstack.ListConfigurationsResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfiguration(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackConfigurationExists("cloudstack_configuration.test", &configuration),
					resource.TestCheckResourceAttr("cloudstack_configuration.test", "value", "test_host"),
				),
			},
		},
	})
}

func TestAccCloudStackConfiguration_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConfiguration(),
			},

			{
				ResourceName:      "cloudstack_configuration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackConfigurationExists(n string, configuration *cloudstack.ListConfigurationsResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("configuration ID not set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		p := cs.Configuration.NewListConfigurationsParams()
		p.SetName(rs.Primary.ID)

		cfg, err := cs.Configuration.ListConfigurations(p)
		if err != nil {
			return err
		}

		*configuration = *cfg

		return nil
	}
}

func testAccCheckCloudStackConfigurationAttributes(configuration *cloudstack.ListConfigurationsResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, v := range configuration.Configurations {
			if v.Name == "host" {
				if v.Value != "test_host" {
					return fmt.Errorf("Bad value: %s", v.Value)
				}
				return nil
			}
		}
		return fmt.Errorf("Bad name: %s", "host")
	}
}

func testAccCheckCloudStackConfigurationDestroy(s *terraform.State) error {
	return nil

}

func testAccResourceConfiguration() string {
	return fmt.Sprintf(`
	resource "cloudstack_configuration" "test" {
		name  = "host"
		value = "test_host"
	}
`)
}
