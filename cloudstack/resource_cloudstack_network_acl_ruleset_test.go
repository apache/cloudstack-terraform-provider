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
	"regexp"
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func testAccCheckCloudStackNetworkACLRulesetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ruleset ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		foundRules := 0

		for k, id := range rs.Primary.Attributes {
			// Check for ruleset format: rule.*.uuid
			if strings.Contains(k, "rule.") && strings.HasSuffix(k, ".uuid") && id != "" {
				_, count, err := cs.NetworkACL.GetNetworkACLByID(id)

				if err != nil {
					// Check if this is a "not found" error
					if strings.Contains(err.Error(), "No match found") {
						continue
					}
					return err
				}

				if count == 0 {
					continue
				}
				foundRules++
			}
		}

		if foundRules == 0 {
			return fmt.Errorf("No network ACL rules found in state for %s", n)
		}

		return nil
	}
}

func testAccCheckCloudStackNetworkACLRulesetDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_network_acl_ruleset" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL ruleset ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			// Check for ruleset format: rule.*.uuid
			if strings.Contains(k, "rule.") && strings.HasSuffix(k, ".uuid") && id != "" {
				_, _, err := cs.NetworkACL.GetNetworkACLByID(id)
				if err == nil {
					return fmt.Errorf("Network ACL rule %s still exists", rs.Primary.ID)
				}
			}
		}
	}

	return nil
}

// testAccCreateOutOfBandACLRule creates an ACL rule outside of Terraform
// to simulate an out-of-band change for testing managed=true behavior
func testAccCreateOutOfBandACLRule(t *testing.T, aclID string) {
	client := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	p := client.NetworkACL.NewCreateNetworkACLParams("tcp")
	p.SetAclid(aclID)
	p.SetCidrlist([]string{"10.0.0.0/8"})
	p.SetStartport(443)
	p.SetEndport(443)
	p.SetTraffictype("ingress")
	p.SetAction("allow")
	p.SetNumber(30)

	_, err := client.NetworkACL.CreateNetworkACL(p)
	if err != nil {
		t.Fatalf("Failed to create out-of-band ACL rule: %v", err)
	}
}

// testAccCheckOutOfBandACLRuleExists verifies that the out-of-band rule still exists
func testAccCheckOutOfBandACLRuleExists(aclID string) error {
	client := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	p := client.NetworkACL.NewListNetworkACLsParams()
	p.SetAclid(aclID)

	resp, err := client.NetworkACL.ListNetworkACLs(p)
	if err != nil {
		return fmt.Errorf("Failed to list ACL rules: %v", err)
	}

	// Check that the out-of-band rule (rule number 30) still exists
	for _, rule := range resp.NetworkACLs {
		if rule.Number == 30 {
			return nil // Found it - success!
		}
	}

	return fmt.Errorf("Out-of-band rule (number 30) was deleted even though managed=false")
}

// testAccCheckOutOfBandACLRuleDeleted verifies that the out-of-band rule was deleted
func testAccCheckOutOfBandACLRuleDeleted(aclID string) error {
	client := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	p := client.NetworkACL.NewListNetworkACLsParams()
	p.SetAclid(aclID)

	resp, err := client.NetworkACL.ListNetworkACLs(p)
	if err != nil {
		return fmt.Errorf("Failed to list ACL rules: %v", err)
	}

	// Check that the out-of-band rule (rule number 30) was deleted
	for _, rule := range resp.NetworkACLs {
		if rule.Number == 30 {
			return fmt.Errorf("Out-of-band rule (number 30) still exists even though managed=true")
		}
	}

	return nil // Not found - success!
}

func TestAccCloudStackNetworkACLRuleset_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.bar", "rule.#", "8"),
					// Check for the expected rules using TypeSet elem matching
					// Test minimum rule number (1)
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "1",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic - min rule number",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "40",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					// Test UDP protocol
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "50",
							"action":       "allow",
							"protocol":     "udp",
							"port":         "53",
							"traffic_type": "ingress",
							"description":  "Allow DNS",
						}),
					// Test optional description field (rule without description)
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "60",
							"action":       "allow",
							"protocol":     "udp",
							"port":         "123",
							"traffic_type": "ingress",
						}),
					// Test maximum port number
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "100",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "65535",
							"traffic_type": "ingress",
							"description":  "Max port number",
						}),
					// Test maximum rule number (65535)
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "65535",
							"action":       "deny",
							"protocol":     "all",
							"traffic_type": "egress",
							"description":  "Max rule number",
						}),
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_basic = `
resource "cloudstack_vpc" "bar" {
  name = "terraform-vpc-ruleset"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "bar" {
  name = "terraform-acl-ruleset"
  description = "terraform-acl-ruleset-text"
  vpc_id = cloudstack_vpc.bar.id
}

resource "cloudstack_network_acl_ruleset" "bar" {
  acl_id = cloudstack_network_acl.bar.id

  rule {
  	rule_number = 1
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic - min rule number"
  }

  rule {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = "-1"
    icmp_code = "-1"
    traffic_type = "ingress"
	description = "Allow ICMP traffic"
  }

  rule {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  rule {
  	rule_number = 40
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  rule {
  	rule_number = 50
  	action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "udp"
    port = "53"
    traffic_type = "ingress"
	description = "Allow DNS"
  }

  rule {
  	rule_number = 60
  	action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "udp"
    port = "123"
    traffic_type = "ingress"
	# No description - testing optional field
  }

  rule {
  	rule_number = 100
  	action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "65535"
    traffic_type = "ingress"
	description = "Max port number"
  }

  rule {
  	rule_number = 65535
  	action = "deny"
    cidr_list = ["0.0.0.0/0"]
    protocol = "all"
    traffic_type = "egress"
	description = "Max rule number"
  }
}`

func TestAccCloudStackNetworkACLRuleset_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.bar", "rule.#", "8"),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRuleset_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.bar", "rule.#", "6"),
					// Check for the expected rules using TypeSet elem matching
					// Rule 10: Changed action from allow to deny AND changed CIDR list from single to multiple
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "deny",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					// Rule 20: Changed action and added CIDR
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "deny",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					// Rule 30: Unchanged
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					// Rule 40: Unchanged
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "40",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					// Rule 50: New egress rule
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "50",
							"action":       "deny",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "egress",
							"description":  "Deny specific TCP ports",
						}),
					// Rule 60: New egress rule with port range
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.bar", "rule.*", map[string]string{
							"rule_number":  "60",
							"action":       "deny",
							"protocol":     "tcp",
							"port":         "1000-2000",
							"traffic_type": "egress",
							"description":  "Deny specific TCP ports",
						}),
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_update = `
resource "cloudstack_vpc" "bar" {
  name = "terraform-vpc-ruleset"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "bar" {
  name = "terraform-acl-ruleset"
  description = "terraform-acl-ruleset-text"
  vpc_id = cloudstack_vpc.bar.id
}

resource "cloudstack_network_acl_ruleset" "bar" {
  acl_id = cloudstack_network_acl.bar.id

  rule {
  	rule_number = 10
  	action = "deny"
    cidr_list = ["172.18.100.0/24", "192.168.1.0/24", "10.0.0.0/8"]  # Changed from single to multiple CIDRs
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic"
  }

  rule {
  	rule_number = 20
  	action = "deny"
	cidr_list = ["172.18.100.0/24", "172.18.101.0/24"]
    protocol = "icmp"
    icmp_type = "-1"
    icmp_code = "-1"
    traffic_type = "ingress"
	description = "Allow ICMP traffic"
  }

  rule {
  	rule_number = 30
	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  rule {
  	rule_number = 40
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  rule {
  	rule_number = 50
	action = "deny"
    cidr_list = ["10.0.0.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "egress"
	description = "Deny specific TCP ports"
  }

  rule {
  	rule_number = 60
	action = "deny"
    cidr_list = ["10.0.0.0/24"]
    protocol = "tcp"
    port = "1000-2000"
    traffic_type = "egress"
	description = "Deny specific TCP ports"
  }
}`

func TestAccCloudStackNetworkACLRuleset_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.managed"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.managed", "managed", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.managed", "rule.#", "2"),
					// Capture the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_ruleset.managed"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_ruleset.managed")
						}
						aclID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				// Add an out-of-band rule via the API
				PreConfig: func() {
					// Create a rule outside of Terraform
					testAccCreateOutOfBandACLRule(t, aclID)
				},
				Config: testAccCloudStackNetworkACLRuleset_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.managed"),
					// With managed=true, the out-of-band rule should be DELETED
					// Verify only the 2 configured rules exist
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.managed", "rule.#", "2"),
					// Verify the out-of-band rule was deleted
					func(s *terraform.State) error {
						return testAccCheckOutOfBandACLRuleDeleted(aclID)
					},
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_managed = `
resource "cloudstack_vpc" "managed" {
  name = "terraform-vpc-managed"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "managed" {
  name = "terraform-acl-managed"
  description = "terraform-acl-managed-text"
  vpc_id = cloudstack_vpc.managed.id
}

resource "cloudstack_network_acl_ruleset" "managed" {
  acl_id = cloudstack_network_acl.managed.id
  managed = true

  rule {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH"
  }

  rule {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }
}`

func TestAccCloudStackNetworkACLRuleset_insert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_insert_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.baz"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.baz", "rule.#", "3"),
					// Initial rules: 10, 30, 50
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.baz", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.baz", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.baz", "rule.*", map[string]string{
							"rule_number":  "50",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "3306",
							"traffic_type": "ingress",
							"description":  "Allow MySQL",
						}),
				),
			},

			{
				Config: testAccCloudStackNetworkACLRuleset_insert_middle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.baz"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.baz", "rule.#", "4"),
					// After inserting rule 20 in the middle, all original rules should still exist
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.baz", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					// NEW RULE inserted in the middle
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.baz", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.baz", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.baz", "rule.*", map[string]string{
							"rule_number":  "50",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "3306",
							"traffic_type": "ingress",
							"description":  "Allow MySQL",
						}),
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_insert_initial = `
resource "cloudstack_vpc" "baz" {
  name = "terraform-vpc-ruleset-insert"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "baz" {
  name = "terraform-acl-ruleset-insert"
  description = "terraform-acl-ruleset-insert-text"
  vpc_id = cloudstack_vpc.baz.id
}

resource "cloudstack_network_acl_ruleset" "baz" {
  acl_id = cloudstack_network_acl.baz.id

  rule {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  rule {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  rule {
  	rule_number = 50
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "3306"
    traffic_type = "ingress"
	description = "Allow MySQL"
  }
}`

const testAccCloudStackNetworkACLRuleset_insert_middle = `
resource "cloudstack_vpc" "baz" {
  name = "terraform-vpc-ruleset-insert"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "baz" {
  name = "terraform-acl-ruleset-insert"
  description = "terraform-acl-ruleset-insert-text"
  vpc_id = cloudstack_vpc.baz.id
}

resource "cloudstack_network_acl_ruleset" "baz" {
  acl_id = cloudstack_network_acl.baz.id

  rule {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  # NEW RULE INSERTED IN THE MIDDLE
  rule {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  rule {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  rule {
  	rule_number = 50
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "3306"
    traffic_type = "ingress"
	description = "Allow MySQL"
  }
}`

func TestAccCloudStackNetworkACLRuleset_not_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.not_managed"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.not_managed", "rule.#", "2"),
					// Capture the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_ruleset.not_managed"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_ruleset.not_managed")
						}
						aclID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				// Add an out-of-band rule via the API
				PreConfig: func() {
					// Create a rule outside of Terraform
					testAccCreateOutOfBandACLRule(t, aclID)
				},
				Config: testAccCloudStackNetworkACLRuleset_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.not_managed"),
					// With managed=false (default), the out-of-band rule should be PRESERVED
					// Verify the out-of-band rule still exists
					func(s *terraform.State) error {
						return testAccCheckOutOfBandACLRuleExists(aclID)
					},
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_not_managed = `
resource "cloudstack_vpc" "not_managed" {
  name = "terraform-vpc-not-managed"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "not_managed" {
  name = "terraform-acl-not-managed"
  description = "terraform-acl-not-managed-text"
  vpc_id = cloudstack_vpc.not_managed.id
}

resource "cloudstack_network_acl_ruleset" "not_managed" {
  acl_id = cloudstack_network_acl.not_managed.id
  # managed = false is the default, so we don't set it explicitly

  rule {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH"
  }

  rule {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }
}`

func TestAccCloudStackNetworkACLRuleset_insert_plan_check(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_plan_check_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.plan_check"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.plan_check", "rule.#", "3"),
					// Initial rules: 10, 30, 50
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.plan_check", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.plan_check", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.plan_check", "rule.*", map[string]string{
							"rule_number":  "50",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "3306",
							"traffic_type": "ingress",
							"description":  "Allow MySQL",
						}),
				),
			},

			{
				Config: testAccCloudStackNetworkACLRuleset_plan_check_insert,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify that only 1 rule is being added (the new rule 20)
						// and the existing rules (10, 30, 50) are not being modified
						plancheck.ExpectResourceAction("cloudstack_network_acl_ruleset.plan_check", plancheck.ResourceActionUpdate),
						// Verify that rule.# is changing from 3 to 4 (exactly one block added)
						plancheck.ExpectKnownValue(
							"cloudstack_network_acl_ruleset.plan_check",
							tfjsonpath.New("rule"),
							knownvalue.SetSizeExact(4),
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.plan_check"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.plan_check", "rule.#", "4"),
					// After inserting rule 20 in the middle, all original rules should still exist
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.plan_check", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					// NEW RULE inserted in the middle
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.plan_check", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.plan_check", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.plan_check", "rule.*", map[string]string{
							"rule_number":  "50",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "3306",
							"traffic_type": "ingress",
							"description":  "Allow MySQL",
						}),
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_plan_check_initial = `
resource "cloudstack_vpc" "plan_check" {
  name = "terraform-vpc-ruleset-plan-check"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "plan_check" {
  name = "terraform-acl-ruleset-plan-check"
  description = "terraform-acl-ruleset-plan-check-text"
  vpc_id = cloudstack_vpc.plan_check.id
}

resource "cloudstack_network_acl_ruleset" "plan_check" {
  acl_id = cloudstack_network_acl.plan_check.id

  rule {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  rule {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  rule {
  	rule_number = 50
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "3306"
    traffic_type = "ingress"
	description = "Allow MySQL"
  }
}
`

const testAccCloudStackNetworkACLRuleset_plan_check_insert = `
resource "cloudstack_vpc" "plan_check" {
  name = "terraform-vpc-ruleset-plan-check"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "plan_check" {
  name = "terraform-acl-ruleset-plan-check"
  description = "terraform-acl-ruleset-plan-check-text"
  vpc_id = cloudstack_vpc.plan_check.id
}

resource "cloudstack_network_acl_ruleset" "plan_check" {
  acl_id = cloudstack_network_acl.plan_check.id

  rule {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  # NEW RULE INSERTED IN THE MIDDLE
  rule {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  rule {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  rule {
  	rule_number = 50
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "3306"
    traffic_type = "ingress"
	description = "Allow MySQL"
  }
}
`

func TestAccCloudStackNetworkACLRuleset_field_changes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_field_changes_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.field_changes"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.field_changes", "rule.#", "4"),
					// Initial rules with specific values
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "8",
							"icmp_code":    "0",
							"traffic_type": "ingress",
							"description":  "Allow ping",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "40",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "egress",
							"description":  "Allow all egress",
						}),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRuleset_field_changes_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.field_changes"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.field_changes", "rule.#", "4"),
					// Same rule numbers but with changed fields
					// Rule 10: Changed port and CIDR list
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "2222", // Changed port
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					// Rule 20: Changed action
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "deny", // Changed action
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					// Rule 30: Changed ICMP type
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "0", // Changed ICMP type
							"icmp_code":    "0",
							"traffic_type": "ingress",
							"description":  "Allow ping",
						}),
					// Rule 40: Changed action
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.field_changes", "rule.*", map[string]string{
							"rule_number":  "40",
							"action":       "deny", // Changed action
							"protocol":     "all",
							"traffic_type": "egress",
							"description":  "Allow all egress",
						}),
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_field_changes_initial = `
resource "cloudstack_vpc" "field_changes" {
  name = "terraform-vpc-field-changes"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "field_changes" {
  name = "terraform-acl-field-changes"
  description = "terraform-acl-field-changes-text"
  vpc_id = cloudstack_vpc.field_changes.id
}

resource "cloudstack_network_acl_ruleset" "field_changes" {
  acl_id = cloudstack_network_acl.field_changes.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
    description = "Allow SSH"
  }

  rule {
    rule_number = 20
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
    description = "Allow HTTP"
  }

  rule {
    rule_number = 30
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = 8
    icmp_code = 0
    traffic_type = "ingress"
    description = "Allow ping"
  }

  rule {
    rule_number = 40
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "egress"
    description = "Allow all egress"
  }
}
`

const testAccCloudStackNetworkACLRuleset_field_changes_updated = `
resource "cloudstack_vpc" "field_changes" {
  name = "terraform-vpc-field-changes"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "field_changes" {
  name = "terraform-acl-field-changes"
  description = "terraform-acl-field-changes-text"
  vpc_id = cloudstack_vpc.field_changes.id
}

resource "cloudstack_network_acl_ruleset" "field_changes" {
  acl_id = cloudstack_network_acl.field_changes.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["192.168.1.0/24", "10.0.0.0/8"]  # Changed CIDR list
    protocol = "tcp"
    port = "2222"  # Changed from 22
    traffic_type = "ingress"
    description = "Allow SSH"
  }

  rule {
    rule_number = 20
    action = "deny"  # Changed from allow
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
    description = "Allow HTTP"
  }

  rule {
    rule_number = 30
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = 0  # Changed from 8
    icmp_code = 0
    traffic_type = "ingress"
    description = "Allow ping"
  }

  rule {
    rule_number = 40
    action = "deny"  # Changed from allow
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "egress"
    description = "Allow all egress"
  }
}
`

func TestAccCloudStackNetworkACLRuleset_protocol_transitions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_protocol_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.protocol_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.*", map[string]string{
							"rule_number":  "10",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
						}),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRuleset_protocol_icmp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.protocol_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.#", "1"),
					// Verify protocol changed to ICMP
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.*", map[string]string{
							"rule_number":  "10",
							"protocol":     "icmp",
							"icmp_type":    "8",
							"icmp_code":    "0",
							"traffic_type": "ingress",
						}),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRuleset_protocol_all,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.protocol_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.#", "1"),
					// Verify protocol changed to all
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.*", map[string]string{
							"rule_number":  "10",
							"protocol":     "all",
							"traffic_type": "ingress",
						}),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRuleset_protocol_udp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.protocol_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.#", "1"),
					// Verify protocol changed back to UDP with port
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.protocol_test", "rule.*", map[string]string{
							"rule_number":  "10",
							"protocol":     "udp",
							"port":         "53",
							"traffic_type": "ingress",
						}),
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_protocol_tcp = `
resource "cloudstack_vpc" "protocol_test" {
  name = "terraform-vpc-protocol-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "protocol_test" {
  name = "terraform-acl-protocol-test"
  description = "terraform-acl-protocol-test-text"
  vpc_id = cloudstack_vpc.protocol_test.id
}

resource "cloudstack_network_acl_ruleset" "protocol_test" {
  acl_id = cloudstack_network_acl.protocol_test.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
    description = "TCP with port"
  }
}
`

const testAccCloudStackNetworkACLRuleset_protocol_icmp = `
resource "cloudstack_vpc" "protocol_test" {
  name = "terraform-vpc-protocol-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "protocol_test" {
  name = "terraform-acl-protocol-test"
  description = "terraform-acl-protocol-test-text"
  vpc_id = cloudstack_vpc.protocol_test.id
}

resource "cloudstack_network_acl_ruleset" "protocol_test" {
  acl_id = cloudstack_network_acl.protocol_test.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "icmp"
    icmp_type = 8
    icmp_code = 0
    traffic_type = "ingress"
    description = "ICMP ping"
  }
}
`

const testAccCloudStackNetworkACLRuleset_protocol_all = `
resource "cloudstack_vpc" "protocol_test" {
  name = "terraform-vpc-protocol-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "protocol_test" {
  name = "terraform-acl-protocol-test"
  description = "terraform-acl-protocol-test-text"
  vpc_id = cloudstack_vpc.protocol_test.id
}

resource "cloudstack_network_acl_ruleset" "protocol_test" {
  acl_id = cloudstack_network_acl.protocol_test.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "all"
    traffic_type = "ingress"
    description = "All protocols"
  }
}
`

const testAccCloudStackNetworkACLRuleset_protocol_udp = `
resource "cloudstack_vpc" "protocol_test" {
  name = "terraform-vpc-protocol-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "protocol_test" {
  name = "terraform-acl-protocol-test"
  description = "terraform-acl-protocol-test-text"
  vpc_id = cloudstack_vpc.protocol_test.id
}

resource "cloudstack_network_acl_ruleset" "protocol_test" {
  acl_id = cloudstack_network_acl.protocol_test.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "udp"
    port = "53"
    traffic_type = "ingress"
    description = "UDP DNS"
  }
}
`

func TestAccCloudStackNetworkACLRuleset_no_spurious_diff(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_no_spurious_diff_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.no_spurious_diff"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.no_spurious_diff", "rule.#", "3"),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRuleset_no_spurious_diff_change_one_rule,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify that only an update occurs (not a replacement)
						plancheck.ExpectResourceAction("cloudstack_network_acl_ruleset.no_spurious_diff", plancheck.ResourceActionUpdate),
						// Verify that rule.# stays at 3 (no rules added or removed)
						// This proves that rules 10 and 30 are not being deleted and recreated
						plancheck.ExpectKnownValue(
							"cloudstack_network_acl_ruleset.no_spurious_diff",
							tfjsonpath.New("rule"),
							knownvalue.SetSizeExact(3),
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.no_spurious_diff"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.no_spurious_diff", "rule.#", "3"),
					// Verify rule 20 was updated
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.no_spurious_diff", "rule.*", map[string]string{
							"rule_number": "20",
							"port":        "8080", // Changed from 80
						}),
					// Verify rules 10 and 30 still exist with their original values
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.no_spurious_diff", "rule.*", map[string]string{
							"rule_number": "10",
							"port":        "22",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.no_spurious_diff", "rule.*", map[string]string{
							"rule_number": "30",
							"port":        "443",
						}),
				),
			},
		},
	})
}

// Test that changing the action field on one rule doesn't cause spurious diffs on other rules
func TestAccCloudStackNetworkACLRuleset_no_spurious_diff_action_change(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_action_change_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.action_change"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.action_change", "rule.#", "2"),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRuleset_action_change_deny,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify that only an update occurs (not a replacement)
						plancheck.ExpectResourceAction("cloudstack_network_acl_ruleset.action_change", plancheck.ResourceActionUpdate),
						// Verify that rule.# stays at 2 (no rules added or removed)
						// This proves that rule 42002 is not being deleted and recreated
						plancheck.ExpectKnownValue(
							"cloudstack_network_acl_ruleset.action_change",
							tfjsonpath.New("rule"),
							knownvalue.SetSizeExact(2),
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.action_change"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.action_change", "rule.#", "2"),
					// Verify rule 42001 was updated to deny
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.action_change", "rule.*", map[string]string{
							"rule_number":  "42001",
							"action":       "deny", // Changed from allow
							"traffic_type": "egress",
						}),
					// Verify rule 42002 still exists with its original values
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.action_change", "rule.*", map[string]string{
							"rule_number":  "42002",
							"action":       "allow", // Unchanged
							"traffic_type": "ingress",
						}),
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_no_spurious_diff_initial = `
resource "cloudstack_vpc" "no_spurious_diff" {
  name = "terraform-vpc-no-spurious-diff"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "no_spurious_diff" {
  name = "terraform-acl-no-spurious-diff"
  description = "terraform-acl-no-spurious-diff-text"
  vpc_id = cloudstack_vpc.no_spurious_diff.id
}

resource "cloudstack_network_acl_ruleset" "no_spurious_diff" {
  acl_id = cloudstack_network_acl.no_spurious_diff.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
    description = "SSH"
  }

  rule {
    rule_number = 20
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
    description = "HTTP"
  }

  rule {
    rule_number = 30
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
    description = "HTTPS"
  }
}
`

const testAccCloudStackNetworkACLRuleset_action_change_initial = `
resource "cloudstack_vpc" "action_change" {
  name = "terraform-vpc-action-change"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "action_change" {
  name = "terraform-acl-action-change"
  description = "terraform-acl-action-change-text"
  vpc_id = cloudstack_vpc.action_change.id
}

resource "cloudstack_network_acl_ruleset" "action_change" {
  acl_id = cloudstack_network_acl.action_change.id

  rule {
    rule_number = 42001
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "all"
    traffic_type = "egress"
    description = "to any vpc: allow egress"
  }

  rule {
    rule_number = 42002
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "all"
    traffic_type = "ingress"
    description = "from anywhere: allow ingress"
  }
}
`

const testAccCloudStackNetworkACLRuleset_action_change_deny = `
resource "cloudstack_vpc" "action_change" {
  name = "terraform-vpc-action-change"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "action_change" {
  name = "terraform-acl-action-change"
  description = "terraform-acl-action-change-text"
  vpc_id = cloudstack_vpc.action_change.id
}

resource "cloudstack_network_acl_ruleset" "action_change" {
  acl_id = cloudstack_network_acl.action_change.id

  rule {
    rule_number = 42001
    action = "deny"
    cidr_list = ["0.0.0.0/0"]
    protocol = "all"
    traffic_type = "egress"
    description = "to any vpc: deny egress"
  }

  rule {
    rule_number = 42002
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "all"
    traffic_type = "ingress"
    description = "from anywhere: allow ingress"
  }
}
`

const testAccCloudStackNetworkACLRuleset_no_spurious_diff_change_one_rule = `
resource "cloudstack_vpc" "no_spurious_diff" {
  name = "terraform-vpc-no-spurious-diff"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "no_spurious_diff" {
  name = "terraform-acl-no-spurious-diff"
  description = "terraform-acl-no-spurious-diff-text"
  vpc_id = cloudstack_vpc.no_spurious_diff.id
}

resource "cloudstack_network_acl_ruleset" "no_spurious_diff" {
  acl_id = cloudstack_network_acl.no_spurious_diff.id

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
    description = "SSH"
  }

  rule {
    rule_number = 20
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "8080"  # Changed from 80
    traffic_type = "ingress"
    description = "HTTP"
  }

  rule {
    rule_number = 30
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
    description = "HTTPS"
  }
}
`

func TestAccCloudStackNetworkACLRuleset_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_import_config,
			},
			{
				ResourceName:            "cloudstack_network_acl_ruleset.import_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"managed"},
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_import_config = `
resource "cloudstack_vpc" "import_test" {
  name = "terraform-vpc-import-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "import_test" {
  name = "terraform-acl-import-test"
  description = "terraform-acl-import-test-text"
  vpc_id = cloudstack_vpc.import_test.id
}

resource "cloudstack_network_acl_ruleset" "import_test" {
  acl_id = cloudstack_network_acl.import_test.id
  managed = false  # Don't delete rules on destroy, so they can be imported

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
    description = "Allow HTTP"
  }

  rule {
    rule_number = 20
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "icmp"
    icmp_type = 8
    icmp_code = 0
    traffic_type = "ingress"
    description = "Allow ping"
  }
}
`

func TestAccCloudStackNetworkACLRuleset_numeric_protocol_error(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// Don't check destroy since the resource was never created
		Steps: []resource.TestStep{
			{
				Config:      testAccCloudStackNetworkACLRuleset_numeric_protocol,
				ExpectError: regexp.MustCompile("numeric protocols are not supported"),
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_numeric_protocol = `
resource "cloudstack_network_acl_ruleset" "numeric_test" {
  acl_id = "test-acl-id"

  rule {
    rule_number = 10
    action = "allow"
    cidr_list = ["0.0.0.0/0"]
    protocol = "6"  # Numeric protocol (6 = TCP)
    port = "80"
    traffic_type = "ingress"
    description = "This should fail"
  }
}
`

func TestAccCloudStackNetworkACLRuleset_remove(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_remove_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.remove_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.remove_test", "rule.#", "4"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.remove_test", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.remove_test", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.remove_test", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.remove_test", "rule.*", map[string]string{
							"rule_number":  "40",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
				),
			},

			{
				Config: testAccCloudStackNetworkACLRuleset_remove_after,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify that we're only removing rules, not adding ghost entries
						plancheck.ExpectResourceAction("cloudstack_network_acl_ruleset.remove_test", plancheck.ResourceActionUpdate),
						// The plan should show exactly 2 rules in the ruleset after removal
						// No ghost entries with empty cidr_list should appear
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.remove_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_ruleset.remove_test", "rule.#", "2"),
					// Only rules 10 and 30 should remain
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.remove_test", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_ruleset.remove_test", "rule.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
				),
			},
			{
				// Re-apply the same config to verify no permadiff
				// This ensures that Computed: true doesn't cause unexpected diffs
				Config:   testAccCloudStackNetworkACLRuleset_remove_after,
				PlanOnly: true, // Should show no changes
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_remove_initial = `
resource "cloudstack_vpc" "remove_test" {
  name = "terraform-vpc-remove-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "remove_test" {
  name = "terraform-acl-remove-test"
  description = "terraform-acl-remove-test-text"
  vpc_id = cloudstack_vpc.remove_test.id
}

resource "cloudstack_network_acl_ruleset" "remove_test" {
  acl_id = cloudstack_network_acl.remove_test.id

  rule {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic"
  }

  rule {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = "-1"
    icmp_code = "-1"
    traffic_type = "ingress"
	description = "Allow ICMP traffic"
  }

  rule {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  rule {
  	rule_number = 40
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }
}`

const testAccCloudStackNetworkACLRuleset_remove_after = `
resource "cloudstack_vpc" "remove_test" {
  name = "terraform-vpc-remove-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "remove_test" {
  name = "terraform-acl-remove-test"
  description = "terraform-acl-remove-test-text"
  vpc_id = cloudstack_vpc.remove_test.id
}

resource "cloudstack_network_acl_ruleset" "remove_test" {
  acl_id = cloudstack_network_acl.remove_test.id

  rule {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic"
  }

  rule {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }
}
`

func TestAccCloudStackNetworkACLRuleset_icmp_defaults(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRulesetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRuleset_icmp_defaults,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.icmp_test"),
					resource.TestCheckResourceAttr("cloudstack_network_acl_ruleset.icmp_test", "rule.#", "1"),
					// Verify that icmp_type and icmp_code default to -1
					resource.TestCheckTypeSetElemNestedAttrs("cloudstack_network_acl_ruleset.icmp_test", "rule.*", map[string]string{
						"rule_number":  "100",
						"action":       "allow",
						"protocol":     "icmp",
						"icmp_type":    "-1",
						"icmp_code":    "-1",
						"traffic_type": "ingress",
					}),
				),
			},
			{
				// Re-apply the same config to ensure no spurious diff
				Config: testAccCloudStackNetworkACLRuleset_icmp_defaults,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				// Transition from ICMP to TCP
				Config: testAccCloudStackNetworkACLRuleset_icmp_to_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.icmp_test"),
					resource.TestCheckResourceAttr("cloudstack_network_acl_ruleset.icmp_test", "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("cloudstack_network_acl_ruleset.icmp_test", "rule.*", map[string]string{
						"rule_number":  "100",
						"action":       "allow",
						"protocol":     "tcp",
						"port":         "80",
						"traffic_type": "ingress",
					}),
				),
			},
			{
				// Re-apply to ensure no spurious diff after protocol transition
				Config: testAccCloudStackNetworkACLRuleset_icmp_to_tcp,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				// Transition back from TCP to ICMP
				Config: testAccCloudStackNetworkACLRuleset_icmp_defaults,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesetExists("cloudstack_network_acl_ruleset.icmp_test"),
					resource.TestCheckResourceAttr("cloudstack_network_acl_ruleset.icmp_test", "rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("cloudstack_network_acl_ruleset.icmp_test", "rule.*", map[string]string{
						"rule_number":  "100",
						"action":       "allow",
						"protocol":     "icmp",
						"icmp_type":    "-1",
						"icmp_code":    "-1",
						"traffic_type": "ingress",
					}),
				),
			},
			{
				// Final check - no spurious diff
				Config: testAccCloudStackNetworkACLRuleset_icmp_defaults,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

const testAccCloudStackNetworkACLRuleset_icmp_defaults = `
resource "cloudstack_vpc" "icmp_test" {
  name = "terraform-vpc-icmp-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "icmp_test" {
  name = "terraform-acl-icmp-test"
  description = "Testing ICMP defaults"
  vpc_id = cloudstack_vpc.icmp_test.id
}

resource "cloudstack_network_acl_ruleset" "icmp_test" {
  acl_id = cloudstack_network_acl.icmp_test.id

  rule {
    rule_number  = 100
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "icmp"
    # icmp_type and icmp_code not specified - should default to -1
    traffic_type = "ingress"
    description  = "Allow all ICMP"
  }
}
`

const testAccCloudStackNetworkACLRuleset_icmp_to_tcp = `
resource "cloudstack_vpc" "icmp_test" {
  name = "terraform-vpc-icmp-test"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "icmp_test" {
  name = "terraform-acl-icmp-test"
  description = "Testing ICMP defaults"
  vpc_id = cloudstack_vpc.icmp_test.id
}

resource "cloudstack_network_acl_ruleset" "icmp_test" {
  acl_id = cloudstack_network_acl.icmp_test.id

  rule {
    rule_number  = 100
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }
}
`
