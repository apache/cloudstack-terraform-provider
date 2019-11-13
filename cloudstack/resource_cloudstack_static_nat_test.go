package cloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-cloudstack/v2/cloudstack"
)

func TestAccCloudStackStaticNAT_basic(t *testing.T) {
	var ipaddr cloudstack.PublicIpAddress

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackStaticNATDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackStaticNAT_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackStaticNATExists(
						"cloudstack_static_nat.foo", &ipaddr),
					testAccCheckCloudStackStaticNATAttributes(&ipaddr),
				),
			},
		},
	})
}

func testAccCheckCloudStackStaticNATExists(
	n string, ipaddr *cloudstack.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static NAT ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		ip, _, err := cs.Address.GetPublicIpAddressByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ip.Id != rs.Primary.ID {
			return fmt.Errorf("Static NAT not found")
		}

		if !ip.Isstaticnat {
			return fmt.Errorf("Static NAT not enabled")
		}

		*ipaddr = *ip

		return nil
	}
}

func testAccCheckCloudStackStaticNATAttributes(
	ipaddr *cloudstack.PublicIpAddress) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if ipaddr.Associatednetworkname != "terraform-network" {
			return fmt.Errorf("Bad network name: %s", ipaddr.Associatednetworkname)
		}

		return nil
	}
}

func testAccCheckCloudStackStaticNATDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_static_nat" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No static NAT ID is set")
		}

		ip, _, err := cs.Address.GetPublicIpAddressByID(rs.Primary.ID)
		if err == nil && ip.Isstaticnat {
			return fmt.Errorf("Static NAT %s still enabled", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackStaticNAT_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
	source_nat_ip = true
  zone = "Sandbox-simulator"
}

resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "Small Instance"
  network_id = "${cloudstack_network.foo.id}"
  template = "CentOS 5.6 (64-bit) no GUI (Simulator)"
  zone = "${cloudstack_network.foo.zone}"
  expunge = true
}

resource "cloudstack_ipaddress" "foo" {
  network_id = "${cloudstack_network.foo.id}"
}

resource "cloudstack_static_nat" "foo" {
	ip_address_id = "${cloudstack_ipaddress.foo.id}"
  virtual_machine_id = "${cloudstack_instance.foobar.id}"
}`
