package cloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func TestAccCloudStackNetwork_basic(t *testing.T) {
	var network cloudstack.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetwork_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkExists(
						"cloudstack_network.foo", &network),
					testAccCheckCloudStackNetworkBasicAttributes(&network),
					// testAccCheckResourceTags(&network),
				),
			},
		},
	})
}

func TestAccCloudStackNetwork_vpc(t *testing.T) {
	var network cloudstack.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetwork_vpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkExists(
						"cloudstack_network.foo", &network),
					testAccCheckCloudStackNetworkVPCAttributes(&network),
				),
			},
		},
	})
}

func TestAccCloudStackNetwork_updateACL(t *testing.T) {
	var network cloudstack.Network

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetwork_acl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkExists(
						"cloudstack_network.foo", &network),
					testAccCheckCloudStackNetworkVPCAttributes(&network),
				),
			},

			{
				Config: testAccCloudStackNetwork_updateACL,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkExists(
						"cloudstack_network.foo", &network),
					testAccCheckCloudStackNetworkVPCAttributes(&network),
				),
			},
		},
	})
}

func TestAccCloudStackNetwork_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetwork_basic,
			},

			{
				ResourceName:      "cloudstack_network.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudStackNetworkExists(
	n string, network *cloudstack.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		ntwrk, _, err := cs.Network.GetNetworkByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if ntwrk.Id != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = *ntwrk

		return nil
	}
}

func testAccCheckCloudStackNetworkBasicAttributes(
	network *cloudstack.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if network.Name != "terraform-network" {
			return fmt.Errorf("Bad name: %s", network.Name)
		}

		if network.Displaytext != "terraform-network" {
			return fmt.Errorf("Bad display name: %s", network.Displaytext)
		}

		if network.Cidr != "10.1.1.0/24" {
			return fmt.Errorf("Bad CIDR: %s", network.Cidr)
		}

		if network.Networkofferingname != "DefaultIsolatedNetworkOfferingWithSourceNatService" {
			return fmt.Errorf("Bad network offering: %s", network.Networkofferingname)
		}

		return nil
	}
}

func testAccCheckCloudStackNetworkVPCAttributes(
	network *cloudstack.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if network.Name != "terraform-network" {
			return fmt.Errorf("Bad name: %s", network.Name)
		}

		if network.Displaytext != "terraform-network" {
			return fmt.Errorf("Bad display name: %s", network.Displaytext)
		}

		if network.Cidr != "10.1.1.0/24" {
			return fmt.Errorf("Bad CIDR: %s", network.Cidr)
		}

		if network.Networkofferingname != "DefaultIsolatedNetworkOfferingForVpcNetworks" {
			return fmt.Errorf("Bad network offering: %s", network.Networkofferingname)
		}

		return nil
	}
}

func testAccCheckCloudStackNetworkDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_network" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ID is set")
		}

		_, _, err := cs.Network.GetNetworkByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Network %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCloudStackNetwork_basic = `
resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingWithSourceNatService"
  zone = "Sandbox-simulator"
  #tags = {
  #  terraform-tag = "true"
  #}
}`

const testAccCloudStackNetwork_vpc = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id = "${cloudstack_vpc.foo.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}`

const testAccCloudStackNetwork_acl = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "foo" {
  name = "foo"
  vpc_id = "${cloudstack_vpc.foo.id}"
}

resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id = "${cloudstack_vpc.foo.id}"
  acl_id = "${cloudstack_network_acl.foo.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}`

const testAccCloudStackNetwork_updateACL = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "bar" {
  name = "bar"
  vpc_id = "${cloudstack_vpc.foo.id}"
}

resource "cloudstack_network" "foo" {
  name = "terraform-network"
  cidr = "10.1.1.0/24"
  network_offering = "DefaultIsolatedNetworkOfferingForVpcNetworks"
  vpc_id = "${cloudstack_vpc.foo.id}"
  acl_id = "${cloudstack_network_acl.bar.id}"
  zone = "${cloudstack_vpc.foo.zone}"
}`
