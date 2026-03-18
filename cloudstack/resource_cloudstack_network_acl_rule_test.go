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
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCloudStackNetworkACLRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_basic,
				Check: resource.ComposeTestCheckFunc(

					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.foo", "rule.#", "4"),
					// Don't rely on specific rule ordering as TypeSet doesn't guarantee order
					// Just check that we have the expected rules with their attributes
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
				),
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.foo", "rule.#", "4"),
					// Don't rely on specific rule ordering as TypeSet doesn't guarantee order
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
				),
			},

			{
				Config: testAccCloudStackNetworkACLRule_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.foo", "rule.#", "6"),
					// Check for the expected rules using TypeSet elem matching
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "deny",
							"protocol":     "all",
							"traffic_type": "ingress",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "deny",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Deny ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
							"action":       "deny",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "egress",
							"description":  "Deny specific TCP ports",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.foo", "rule.*", map[string]string{
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

func testAccCheckCloudStackNetworkACLRulesExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL rule ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		foundRules := 0

		for k, id := range rs.Primary.Attributes {
			// Check for legacy 'rule' format: rule.*.uuids.<key>
			if strings.Contains(k, ".uuids.") && !strings.HasSuffix(k, ".uuids.%") {
				_, count, err := cs.NetworkACL.GetNetworkACLByID(id)

				if err != nil {
					return err
				}

				if count == 0 {
					return fmt.Errorf("Network ACL rule %s not found", k)
				}
				foundRules++
			}

			// Check for new 'ruleset' format: ruleset.*.uuid
			if strings.Contains(k, "ruleset.") && strings.HasSuffix(k, ".uuid") && id != "" {
				_, count, err := cs.NetworkACL.GetNetworkACLByID(id)

				if err != nil {
					// Check if this is a "not found" error
					// This can happen if an out-of-band rule placeholder was deleted but the state hasn't been fully refreshed yet
					if strings.Contains(err.Error(), "No match found") {
						continue
					}
					return err
				}

				if count == 0 {
					// Don't fail - just skip this UUID
					// This can happen if an out-of-band rule placeholder was deleted but the state hasn't been fully refreshed yet
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

func testAccCheckCloudStackNetworkACLRuleDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_network_acl_rule" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No network ACL rule ID is set")
		}

		for k, id := range rs.Primary.Attributes {
			// Check for legacy 'rule' format: rule.*.uuids.<key>
			if strings.Contains(k, ".uuids.") && !strings.HasSuffix(k, ".uuids.%") {
				_, _, err := cs.NetworkACL.GetNetworkACLByID(id)
				if err == nil {
					return fmt.Errorf("Network ACL rule %s still exists", rs.Primary.ID)
				}
			}

			// Check for new 'ruleset' format: ruleset.*.uuid
			if strings.Contains(k, "ruleset.") && strings.HasSuffix(k, ".uuid") && id != "" {
				_, _, err := cs.NetworkACL.GetNetworkACLByID(id)
				if err == nil {
					return fmt.Errorf("Network ACL rule %s still exists", rs.Primary.ID)
				}
			}
		}
	}

	return nil
}

func TestAccCloudStackNetworkACLRule_ruleset_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.bar", "ruleset.#", "4"),
					// Check for the expected rules using TypeSet elem matching
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "40",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
				),
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_ruleset_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.bar", "ruleset.#", "4"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
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
				Config: testAccCloudStackNetworkACLRule_ruleset_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.bar"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.bar", "ruleset.#", "6"),
					// Check for the expected rules using TypeSet elem matching
					// Rule 10: Changed action from allow to deny
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "deny",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					// Rule 20: Changed action from allow to deny, added CIDR
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "deny",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					// Rule 30: No changes
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					// Rule 40: No changes
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "40",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					// Rule 50: New rule
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
							"rule_number":  "50",
							"action":       "deny",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "egress",
							"description":  "Deny specific TCP ports",
						}),
					// Rule 60: New rule
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.bar", "ruleset.*", map[string]string{
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

func TestAccCloudStackNetworkACLRule_ruleset_insert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_insert_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.baz"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.baz", "ruleset.#", "3"),
					// Initial rules: 10, 30, 50
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.baz", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.baz", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.baz", "ruleset.*", map[string]string{
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
				Config: testAccCloudStackNetworkACLRule_ruleset_insert_middle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.baz"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.baz", "ruleset.#", "4"),
					// After inserting rule 20 in the middle, all original rules should still exist
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.baz", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					// NEW RULE inserted in the middle
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.baz", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.baz", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.baz", "ruleset.*", map[string]string{
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

func TestAccCloudStackNetworkACLRule_ruleset_insert_plan_check(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_plan_check_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.plan_check"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.plan_check", "ruleset.#", "3"),
					// Initial rules: 10, 30, 50
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.plan_check", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.plan_check", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.plan_check", "ruleset.*", map[string]string{
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
				Config: testAccCloudStackNetworkACLRule_ruleset_plan_check_insert,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify that only 1 rule is being added (the new rule 20)
						// and the existing rules (10, 30, 50) are not being modified
						plancheck.ExpectResourceAction("cloudstack_network_acl_rule.plan_check", plancheck.ResourceActionUpdate),
						// Verify that ruleset.# is changing from 3 to 4 (exactly one block added)
						plancheck.ExpectKnownValue(
							"cloudstack_network_acl_rule.plan_check",
							tfjsonpath.New("ruleset"),
							knownvalue.SetSizeExact(4),
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.plan_check"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.plan_check", "ruleset.#", "4"),
					// After inserting rule 20 in the middle, all original rules should still exist
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.plan_check", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					// NEW RULE inserted in the middle
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.plan_check", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.plan_check", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "443",
							"traffic_type": "ingress",
							"description":  "Allow HTTPS",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.plan_check", "ruleset.*", map[string]string{
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

func TestAccCloudStackNetworkACLRule_ruleset_field_changes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_field_changes_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.field_changes"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.field_changes", "ruleset.#", "4"),
					// Initial rules with specific values
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "22",
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "8",
							"icmp_code":    "0",
							"traffic_type": "ingress",
							"description":  "Allow ping",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
							"rule_number":  "40",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "egress",
							"description":  "Allow all egress",
						}),
				),
			},
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_field_changes_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.field_changes"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.field_changes", "ruleset.#", "4"),
					// Same rule numbers but with changed fields
					// Rule 10: Changed port and CIDR list
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "2222", // Changed port
							"traffic_type": "ingress",
							"description":  "Allow SSH",
						}),
					// Rule 20: Changed action
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "deny", // Changed action
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					// Rule 30: Changed ICMP type
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
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
						"cloudstack_network_acl_rule.field_changes", "ruleset.*", map[string]string{
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

func TestAccCloudStackNetworkACLRule_ruleset_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.managed"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.managed", "managed", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.managed", "ruleset.#", "2"),
					// Store the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_rule.managed"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_rule.managed")
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
				Config: testAccCloudStackNetworkACLRule_ruleset_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.managed"),
					// With managed=true, the out-of-band rule should be DELETED from CloudStack
					// Verify the out-of-band rule was actually deleted
					func(s *terraform.State) error {
						return testAccCheckOutOfBandACLRuleDeleted(aclID)
					},
				),
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_ruleset_not_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.not_managed"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.not_managed", "ruleset.#", "2"),
					// Capture the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_rule.not_managed"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_rule.not_managed")
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
				Config: testAccCloudStackNetworkACLRule_ruleset_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.not_managed"),
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

func TestAccCloudStackNetworkACLRule_ruleset_remove(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_ruleset_remove_initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.remove_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.remove_test", "ruleset.#", "4"),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.remove_test", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.remove_test", "ruleset.*", map[string]string{
							"rule_number":  "20",
							"action":       "allow",
							"protocol":     "icmp",
							"icmp_type":    "-1",
							"icmp_code":    "-1",
							"traffic_type": "ingress",
							"description":  "Allow ICMP traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.remove_test", "ruleset.*", map[string]string{
							"rule_number":  "30",
							"action":       "allow",
							"protocol":     "tcp",
							"port":         "80",
							"traffic_type": "ingress",
							"description":  "Allow HTTP",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.remove_test", "ruleset.*", map[string]string{
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
				Config: testAccCloudStackNetworkACLRule_ruleset_remove_after,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Verify that we're only removing rules, not adding ghost entries
						plancheck.ExpectResourceAction("cloudstack_network_acl_rule.remove_test", plancheck.ResourceActionUpdate),
						// The plan should show exactly 2 rules in the ruleset after removal
						// No ghost entries with empty cidr_list should appear
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.remove_test"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.remove_test", "ruleset.#", "2"),
					// Only rules 10 and 30 should remain
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.remove_test", "ruleset.*", map[string]string{
							"rule_number":  "10",
							"action":       "allow",
							"protocol":     "all",
							"traffic_type": "ingress",
							"description":  "Allow all traffic",
						}),
					resource.TestCheckTypeSetElemNestedAttrs(
						"cloudstack_network_acl_rule.remove_test", "ruleset.*", map[string]string{
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
				Config:   testAccCloudStackNetworkACLRule_ruleset_remove_after,
				PlanOnly: true, // Should show no changes
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_rule_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_rule_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.managed_legacy"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.managed_legacy", "managed", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.managed_legacy", "rule.#", "2"),
					// Capture the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_rule.managed_legacy"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_rule.managed_legacy")
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
				Config: testAccCloudStackNetworkACLRule_rule_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.managed_legacy"),
					// With managed=true, the out-of-band rule should be DELETED from CloudStack
					// Verify the out-of-band rule was actually deleted
					func(s *terraform.State) error {
						return testAccCheckOutOfBandACLRuleDeleted(aclID)
					},
				),
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_rule_not_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_rule_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.not_managed"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.not_managed", "rule.#", "2"),
					// Capture the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_rule.not_managed"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_rule.not_managed")
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
				Config: testAccCloudStackNetworkACLRule_rule_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.not_managed"),
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

const testAccCloudStackNetworkACLRule_basic = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "foo" {
  name = "terraform-acl"
  description = "terraform-acl-text"
  vpc_id = cloudstack_vpc.foo.id
}

resource "cloudstack_network_acl_rule" "foo" {
  acl_id = cloudstack_network_acl.foo.id

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
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  rule {
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }
}`

const testAccCloudStackNetworkACLRule_update = `
resource "cloudstack_vpc" "foo" {
  name = "terraform-vpc"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "foo" {
  name = "terraform-acl"
  description = "terraform-acl-text"
  vpc_id = cloudstack_vpc.foo.id
}

resource "cloudstack_network_acl_rule" "foo" {
  acl_id = cloudstack_network_acl.foo.id

  rule {
  	action = "deny"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
  }

  rule {
  	action = "deny"
	cidr_list = ["172.18.100.0/24", "172.18.101.0/24"]
    protocol = "icmp"
    icmp_type = "-1"
    icmp_code = "-1"
    traffic_type = "ingress"
	description = "Deny ICMP traffic"
  }

  rule {
	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  rule {
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  rule {
	action = "deny"
    cidr_list = ["10.0.0.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "egress"
	description = "Deny specific TCP ports"
  }

  rule {
	action = "deny"
    cidr_list = ["10.0.0.0/24"]
    protocol = "tcp"
    port = "1000-2000"
    traffic_type = "egress"
	description = "Deny specific TCP ports"
  }
}`

const testAccCloudStackNetworkACLRule_ruleset_basic = `
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

resource "cloudstack_network_acl_rule" "bar" {
  acl_id = cloudstack_network_acl.bar.id

  ruleset {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic"
  }

  ruleset {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = "-1"
    icmp_code = "-1"
    traffic_type = "ingress"
	description = "Allow ICMP traffic"
  }

  ruleset {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  ruleset {
  	rule_number = 40
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }
}`

const testAccCloudStackNetworkACLRule_ruleset_update = `
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

resource "cloudstack_network_acl_rule" "bar" {
  acl_id = cloudstack_network_acl.bar.id

  ruleset {
  	rule_number = 10
  	action = "deny"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic"
  }

  ruleset {
  	rule_number = 20
  	action = "deny"
	cidr_list = ["172.18.100.0/24", "172.18.101.0/24"]
    protocol = "icmp"
    icmp_type = "-1"
    icmp_code = "-1"
    traffic_type = "ingress"
	description = "Allow ICMP traffic"
  }

  ruleset {
  	rule_number = 30
	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  ruleset {
  	rule_number = 40
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  ruleset {
  	rule_number = 50
	action = "deny"
    cidr_list = ["10.0.0.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "egress"
	description = "Deny specific TCP ports"
  }

  ruleset {
  	rule_number = 60
	action = "deny"
    cidr_list = ["10.0.0.0/24"]
    protocol = "tcp"
    port = "1000-2000"
    traffic_type = "egress"
	description = "Deny specific TCP ports"
  }
}`

const testAccCloudStackNetworkACLRule_ruleset_insert_initial = `
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

resource "cloudstack_network_acl_rule" "baz" {
  acl_id = cloudstack_network_acl.baz.id

  ruleset {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  ruleset {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  ruleset {
  	rule_number = 50
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "3306"
    traffic_type = "ingress"
	description = "Allow MySQL"
  }
}`

const testAccCloudStackNetworkACLRule_ruleset_insert_middle = `
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

resource "cloudstack_network_acl_rule" "baz" {
  acl_id = cloudstack_network_acl.baz.id

  ruleset {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  # NEW RULE INSERTED IN THE MIDDLE
  ruleset {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  ruleset {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  ruleset {
  	rule_number = 50
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "3306"
    traffic_type = "ingress"
	description = "Allow MySQL"
  }
}`

const testAccCloudStackNetworkACLRule_ruleset_plan_check_initial = `
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

resource "cloudstack_network_acl_rule" "plan_check" {
  acl_id = cloudstack_network_acl.plan_check.id

  ruleset {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  ruleset {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  ruleset {
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

const testAccCloudStackNetworkACLRule_ruleset_plan_check_insert = `
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

resource "cloudstack_network_acl_rule" "plan_check" {
  acl_id = cloudstack_network_acl.plan_check.id

  ruleset {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
	description = "Allow SSH"
  }

  # NEW RULE INSERTED IN THE MIDDLE
  ruleset {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  ruleset {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }

  ruleset {
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

const testAccCloudStackNetworkACLRule_ruleset_field_changes_initial = `
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

resource "cloudstack_network_acl_rule" "field_changes" {
  acl_id = cloudstack_network_acl.field_changes.id

  ruleset {
    rule_number = 10
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "22"
    traffic_type = "ingress"
    description = "Allow SSH"
  }

  ruleset {
    rule_number = 20
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
    description = "Allow HTTP"
  }

  ruleset {
    rule_number = 30
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = 8
    icmp_code = 0
    traffic_type = "ingress"
    description = "Allow ping"
  }

  ruleset {
    rule_number = 40
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "egress"
    description = "Allow all egress"
  }
}
`

const testAccCloudStackNetworkACLRule_ruleset_field_changes_updated = `
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

resource "cloudstack_network_acl_rule" "field_changes" {
  acl_id = cloudstack_network_acl.field_changes.id

  ruleset {
    rule_number = 10
    action = "allow"
    cidr_list = ["192.168.1.0/24", "10.0.0.0/8"]  # Changed CIDR list
    protocol = "tcp"
    port = "2222"  # Changed from 22
    traffic_type = "ingress"
    description = "Allow SSH"
  }

  ruleset {
    rule_number = 20
    action = "deny"  # Changed from allow
    cidr_list = ["172.18.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
    description = "Allow HTTP"
  }

  ruleset {
    rule_number = 30
    action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = 0  # Changed from 8
    icmp_code = 0
    traffic_type = "ingress"
    description = "Allow ping"
  }

  ruleset {
    rule_number = 40
    action = "deny"  # Changed from allow
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "egress"
    description = "Allow all egress"
  }
}
`

func TestAccCloudStackNetworkACLRule_icmp_fields_no_spurious_diff(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_icmp_fields_config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.foo", "ruleset.#", "3"),
				),
			},
			{
				// Second apply with same config should show no changes
				Config: testAccCloudStackNetworkACLRule_icmp_fields_config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_icmp_fields_add_remove_rule(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with 2 rules
				Config: testAccCloudStackNetworkACLRule_icmp_fields_two_rules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.foo", "ruleset.#", "2"),
				),
			},
			{
				// Step 2: Add a third rule
				Config: testAccCloudStackNetworkACLRule_icmp_fields_three_rules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.foo", "ruleset.#", "3"),
				),
			},
			{
				// Step 3: Remove the third rule - should not cause spurious diff on remaining rules
				Config: testAccCloudStackNetworkACLRule_icmp_fields_two_rules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.foo"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.foo", "ruleset.#", "2"),
				),
			},
			{
				// Step 4: Plan should be empty after removing the rule
				Config: testAccCloudStackNetworkACLRule_icmp_fields_two_rules,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

const testAccCloudStackNetworkACLRule_icmp_fields_config = `
resource "cloudstack_vpc" "foo" {
  name         = "terraform-vpc"
  display_text = "terraform-vpc"
  cidr         = "10.0.0.0/16"
  zone         = "Sandbox-simulator"
  vpc_offering = "Default VPC offering"
}

resource "cloudstack_network_acl" "foo" {
  name   = "terraform-acl"
  vpc_id = cloudstack_vpc.foo.id
}

resource "cloudstack_network_acl_rule" "foo" {
  acl_id = cloudstack_network_acl.foo.id

  ruleset {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["0.0.0.0/0"]
    protocol     = "all"
    traffic_type = "ingress"
    description  = "Allow all ingress - protocol all with icmp_type=0, icmp_code=0 in config"
  }

  ruleset {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH - protocol tcp with icmp_type=0, icmp_code=0 in config"
  }

  ruleset {
    rule_number  = 30
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "icmp"
    icmp_type    = 8
    icmp_code    = 0
    traffic_type = "ingress"
    description  = "Allow ICMP echo - protocol icmp with explicit icmp_type and icmp_code"
  }
}
`

const testAccCloudStackNetworkACLRule_icmp_fields_two_rules = `
resource "cloudstack_vpc" "foo" {
  name         = "terraform-vpc-add-remove"
  display_text = "terraform-vpc-add-remove"
  cidr         = "10.0.0.0/16"
  zone         = "Sandbox-simulator"
  vpc_offering = "Default VPC offering"
}

resource "cloudstack_network_acl" "foo" {
  name   = "terraform-acl-add-remove"
  vpc_id = cloudstack_vpc.foo.id
}

resource "cloudstack_network_acl_rule" "foo" {
  acl_id = cloudstack_network_acl.foo.id

  ruleset {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH ingress"
  }

  ruleset {
    rule_number  = 100
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "443"
    traffic_type = "egress"
    description  = "Allow HTTPS egress"
  }
}
`

const testAccCloudStackNetworkACLRule_icmp_fields_three_rules = `
resource "cloudstack_vpc" "foo" {
  name         = "terraform-vpc-add-remove"
  display_text = "terraform-vpc-add-remove"
  cidr         = "10.0.0.0/16"
  zone         = "Sandbox-simulator"
  vpc_offering = "Default VPC offering"
}

resource "cloudstack_network_acl" "foo" {
  name   = "terraform-acl-add-remove"
  vpc_id = cloudstack_vpc.foo.id
}

resource "cloudstack_network_acl_rule" "foo" {
  acl_id = cloudstack_network_acl.foo.id

  ruleset {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH ingress"
  }

  ruleset {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP ingress"
  }

  ruleset {
    rule_number  = 100
    action       = "allow"
    cidr_list    = ["10.0.0.0/8"]
    protocol     = "tcp"
    port         = "443"
    traffic_type = "egress"
    description  = "Allow HTTPS egress"
  }
}
`

const testAccCloudStackNetworkACLRule_ruleset_managed = `
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

resource "cloudstack_network_acl_rule" "managed" {
  acl_id = cloudstack_network_acl.managed.id
  managed = true

  ruleset {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH"
  }

  ruleset {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }
}
`

const testAccCloudStackNetworkACLRule_ruleset_not_managed = `
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

resource "cloudstack_network_acl_rule" "not_managed" {
  acl_id = cloudstack_network_acl.not_managed.id
  # managed = false is the default, so we don't set it explicitly

  ruleset {
    rule_number  = 10
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "22"
    traffic_type = "ingress"
    description  = "Allow SSH"
  }

  ruleset {
    rule_number  = 20
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    port         = "80"
    traffic_type = "ingress"
    description  = "Allow HTTP"
  }
}
`

const testAccCloudStackNetworkACLRule_ruleset_remove_initial = `
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

resource "cloudstack_network_acl_rule" "remove_test" {
  acl_id = cloudstack_network_acl.remove_test.id

  ruleset {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic"
  }

  ruleset {
  	rule_number = 20
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "icmp"
    icmp_type = "-1"
    icmp_code = "-1"
    traffic_type = "ingress"
	description = "Allow ICMP traffic"
  }

  ruleset {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }

  ruleset {
  	rule_number = 40
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "443"
    traffic_type = "ingress"
	description = "Allow HTTPS"
  }
}`

const testAccCloudStackNetworkACLRule_ruleset_remove_after = `
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

resource "cloudstack_network_acl_rule" "remove_test" {
  acl_id = cloudstack_network_acl.remove_test.id

  ruleset {
  	rule_number = 10
  	action = "allow"
    cidr_list = ["172.18.100.0/24"]
    protocol = "all"
    traffic_type = "ingress"
	description = "Allow all traffic"
  }

  ruleset {
  	rule_number = 30
  	action = "allow"
    cidr_list = ["172.16.100.0/24"]
    protocol = "tcp"
    port = "80"
    traffic_type = "ingress"
	description = "Allow HTTP"
  }
}`

const testAccCloudStackNetworkACLRule_rule_managed = `
resource "cloudstack_vpc" "managed_legacy" {
  name = "terraform-vpc-managed-legacy"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "managed_legacy" {
  name = "terraform-acl-managed-legacy"
  description = "terraform-acl-managed-legacy-text"
  vpc_id = cloudstack_vpc.managed_legacy.id
}

resource "cloudstack_network_acl_rule" "managed_legacy" {
  acl_id = cloudstack_network_acl.managed_legacy.id
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
}
`

const testAccCloudStackNetworkACLRule_rule_not_managed = `
resource "cloudstack_vpc" "not_managed_legacy" {
  name = "terraform-vpc-not-managed-legacy"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "not_managed_legacy" {
  name = "terraform-acl-not-managed-legacy"
  description = "terraform-acl-not-managed-legacy-text"
  vpc_id = cloudstack_vpc.not_managed_legacy.id
}

resource "cloudstack_network_acl_rule" "not_managed" {
  acl_id = cloudstack_network_acl.not_managed_legacy.id
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
}
`

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

// testAccCheckOutOfBandACLRuleDeleted verifies that the out-of-band rule was deleted
func testAccCheckOutOfBandACLRuleDeleted(aclID string) error {
	client := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	p := client.NetworkACL.NewListNetworkACLsParams()
	p.SetAclid(aclID)

	resp, err := client.NetworkACL.ListNetworkACLs(p)
	if err != nil {
		return fmt.Errorf("Failed to list ACL rules: %v", err)
	}

	// Check that only the 2 configured rules exist (rule numbers 10 and 20)
	// The out-of-band rule (rule number 30) should have been deleted
	for _, rule := range resp.NetworkACLs {
		if rule.Number == 30 {
			return fmt.Errorf("Out-of-band rule (number 30) was not deleted by managed=true")
		}
	}

	return nil
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

func TestAccCloudStackNetworkACLRule_deprecated_ports(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_deprecated_ports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.deprecated"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.deprecated", "rule.#", "2"),
				),
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_deprecated_ports_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_deprecated_ports_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.deprecated_managed"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.deprecated_managed", "managed", "true"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.deprecated_managed", "rule.#", "2"),
					// Store the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_rule.deprecated_managed"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_rule.deprecated_managed")
						}
						aclID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					// Create an out-of-band ACL rule
					testAccCreateOutOfBandACLRule(t, aclID)
				},
				Config: testAccCloudStackNetworkACLRule_deprecated_ports_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.deprecated_managed"),
					// Verify that the out-of-band rule was deleted by managed=true
					func(s *terraform.State) error {
						return testAccCheckOutOfBandACLRuleDeleted(aclID)
					},
				),
			},
		},
	})
}

func TestAccCloudStackNetworkACLRule_deprecated_ports_not_managed(t *testing.T) {
	var aclID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudStackNetworkACLRule_deprecated_ports_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.deprecated_not_managed"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.deprecated_not_managed", "managed", "false"),
					resource.TestCheckResourceAttr(
						"cloudstack_network_acl_rule.deprecated_not_managed", "rule.#", "2"),
					// Store the ACL ID for later use
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["cloudstack_network_acl_rule.deprecated_not_managed"]
						if !ok {
							return fmt.Errorf("Not found: cloudstack_network_acl_rule.deprecated_not_managed")
						}
						aclID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					// Create an out-of-band ACL rule
					testAccCreateOutOfBandACLRule(t, aclID)
				},
				Config: testAccCloudStackNetworkACLRule_deprecated_ports_not_managed,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackNetworkACLRulesExist("cloudstack_network_acl_rule.deprecated_not_managed"),
					// Verify that the out-of-band rule still exists with managed=false
					func(s *terraform.State) error {
						return testAccCheckOutOfBandACLRuleExists(aclID)
					},
				),
			},
		},
	})
}

const testAccCloudStackNetworkACLRule_deprecated_ports = `
resource "cloudstack_vpc" "deprecated" {
  name = "terraform-vpc-deprecated-ports"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "deprecated" {
  name = "terraform-acl-deprecated-ports"
  description = "terraform-acl-deprecated-ports-text"
  vpc_id = cloudstack_vpc.deprecated.id
}

resource "cloudstack_network_acl_rule" "deprecated" {
  acl_id = cloudstack_network_acl.deprecated.id

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    ports        = ["80", "443"]
    traffic_type = "ingress"
    description  = "Allow HTTP and HTTPS using deprecated ports field"
  }

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    ports        = ["8000-8100"]
    traffic_type = "ingress"
    description  = "Allow port range using deprecated ports field"
  }
}
`

const testAccCloudStackNetworkACLRule_deprecated_ports_managed = `
resource "cloudstack_vpc" "deprecated_managed" {
  name = "terraform-vpc-deprecated-ports-managed"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "deprecated_managed" {
  name = "terraform-acl-deprecated-ports-managed"
  description = "terraform-acl-deprecated-ports-managed-text"
  vpc_id = cloudstack_vpc.deprecated_managed.id
}

resource "cloudstack_network_acl_rule" "deprecated_managed" {
  acl_id  = cloudstack_network_acl.deprecated_managed.id
  managed = true

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    ports        = ["80", "443"]
    traffic_type = "ingress"
    description  = "Allow HTTP and HTTPS using deprecated ports field"
  }

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    ports        = ["22"]
    traffic_type = "ingress"
    description  = "Allow SSH using deprecated ports field"
  }
}
`

const testAccCloudStackNetworkACLRule_deprecated_ports_not_managed = `
resource "cloudstack_vpc" "deprecated_not_managed" {
  name = "terraform-vpc-deprecated-ports-not-managed"
  cidr = "10.0.0.0/8"
  vpc_offering = "Default VPC offering"
  zone = "Sandbox-simulator"
}

resource "cloudstack_network_acl" "deprecated_not_managed" {
  name = "terraform-acl-deprecated-ports-not-managed"
  description = "terraform-acl-deprecated-ports-not-managed-text"
  vpc_id = cloudstack_vpc.deprecated_not_managed.id
}

resource "cloudstack_network_acl_rule" "deprecated_not_managed" {
  acl_id  = cloudstack_network_acl.deprecated_not_managed.id
  managed = false

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    ports        = ["80", "443"]
    traffic_type = "ingress"
    description  = "Allow HTTP and HTTPS using deprecated ports field"
  }

  rule {
    action       = "allow"
    cidr_list    = ["172.18.100.0/24"]
    protocol     = "tcp"
    ports        = ["22"]
    traffic_type = "ingress"
    description  = "Allow SSH using deprecated ports field"
  }
}
`
