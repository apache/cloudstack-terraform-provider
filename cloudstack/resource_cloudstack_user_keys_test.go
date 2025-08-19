package cloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudStackUserKeys_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackUserKeysDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackUserKeysConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackUserKeysExists("cloudstack_user_keys.test"),
					resource.TestCheckResourceAttrSet("cloudstack_user_keys.test", "user_id"),
					resource.TestCheckResourceAttrSet("cloudstack_user_keys.test", "api_key"),
					resource.TestCheckResourceAttrSet("cloudstack_user_keys.test", "secret_key"),
				),
			},
		},
	})
}

func testAccCheckCloudStackUserKeysExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No user keys ID is set")
		}

		return nil
	}
}

func testAccCheckCloudStackUserKeysDestroy(s *terraform.State) error {
	// Since keys cannot be deleted through CloudStack API,
	// we just verify the resource is removed from state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_user_keys" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("User keys still exist")
		}
	}

	return nil
}

const testAccCloudStackUserKeysConfig_basic = `
data "cloudstack_user" "test" {
  username = "admin"
}

resource "cloudstack_user_keys" "test" {
  user_id = data.cloudstack_user.test.id
}
`
