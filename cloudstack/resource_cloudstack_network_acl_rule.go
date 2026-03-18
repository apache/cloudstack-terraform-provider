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
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNetworkACLRuleCreate,
		Read:   resourceCloudStackNetworkACLRuleRead,
		Update: resourceCloudStackNetworkACLRuleUpdate,
		Delete: resourceCloudStackNetworkACLRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudStackNetworkACLRuleImport,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			// Force replacement for migration from deprecated 'ports' to 'port' field
			if diff.HasChange("rule") {
				oldRules, newRules := diff.GetChange("rule")
				oldRulesList := oldRules.([]interface{})
				newRulesList := newRules.([]interface{})

				log.Printf("[DEBUG] CustomizeDiff: checking %d old rules -> %d new rules", len(oldRulesList), len(newRulesList))

				// Check if ANY old rule uses deprecated 'ports' field
				hasDeprecatedPorts := false
				for i, oldRule := range oldRulesList {
					oldRuleMap := oldRule.(map[string]interface{})
					protocol := oldRuleMap["protocol"].(string)

					if protocol == "tcp" || protocol == "udp" {
						if portsSet, hasPortsSet := oldRuleMap["ports"].(*schema.Set); hasPortsSet && portsSet.Len() > 0 {
							log.Printf("[DEBUG] CustomizeDiff: OLD rule %d has deprecated ports field with %d ports: %v", i, portsSet.Len(), portsSet.List())
							hasDeprecatedPorts = true
							break
						}
					}
				}

				// Check if ANY new rule uses new 'port' field
				hasNewPortFormat := false
				for i, newRule := range newRulesList {
					newRuleMap := newRule.(map[string]interface{})
					protocol := newRuleMap["protocol"].(string)

					if protocol == "tcp" || protocol == "udp" {
						if portStr, hasPort := newRuleMap["port"].(string); hasPort && portStr != "" {
							log.Printf("[DEBUG] CustomizeDiff: NEW rule %d has port field: %s", i, portStr)
							hasNewPortFormat = true
							break
						}
					}
				}

				// Force replacement if migrating from deprecated ports to new port format
				if hasDeprecatedPorts && hasNewPortFormat {
					log.Printf("[DEBUG] CustomizeDiff: MIGRATION DETECTED - old rules use deprecated 'ports', new rules use 'port' - FORCING REPLACEMENT")
					diff.ForceNew("rule")
					return nil
				}

				// Also force replacement if old rules have deprecated ports but new rules don't use ports at all
				if hasDeprecatedPorts && !hasNewPortFormat {
					log.Printf("[DEBUG] CustomizeDiff: POTENTIAL MIGRATION - old rules use deprecated 'ports' but new rules don't - FORCING REPLACEMENT")
					diff.ForceNew("rule")
					return nil
				}

				log.Printf("[DEBUG] CustomizeDiff: No migration detected - hasDeprecatedPorts=%t, hasNewPortFormat=%t", hasDeprecatedPorts, hasNewPortFormat)
			}

			// WORKAROUND: Filter out ghost entries from ruleset
			// The SDK creates ghost entries when rules are removed from a TypeSet that has Computed: true
			// This happens because the SDK tries to preserve Computed fields (like uuid) when elements are removed
			if diff.HasChange("ruleset") {
				_, newRuleset := diff.GetChange("ruleset")
				if newSet, ok := newRuleset.(*schema.Set); ok {
					cleanRules, ghostCount := filterGhostEntries(newSet.List(), "CustomizeDiff")

					if ghostCount > 0 {
						// Create a new Set with the clean rules
						rulesetResource := resourceCloudStackNetworkACLRule().Schema["ruleset"].Elem.(*schema.Resource)
						hashFunc := schema.HashResource(rulesetResource)
						cleanSet := schema.NewSet(hashFunc, cleanRules)
						if err := diff.SetNew("ruleset", cleanSet); err != nil {
							log.Printf("[ERROR] CustomizeDiff: Failed to set clean ruleset: %v", err)
							return err
						}
					}
				}
			}

			return nil
		},

		Schema: map[string]*schema.Schema{
			"acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"managed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "allow",
						},

						"cidr_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},

						"icmp_type": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"icmp_code": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"ports": {
							Type:       schema.TypeSet,
							Optional:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Set:        schema.HashString,
							Deprecated: "Use 'port' instead. The 'ports' field is deprecated and will be removed in a future version.",
						},

						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"traffic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ingress",
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"uuids": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},

			"ruleset": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed is required to allow CustomizeDiff to use SetNew() for filtering ghost entries.
				// Ghost entries are created by the SDK when elements are removed from a TypeSet that
				// contains Computed fields (like uuid). The SDK preserves the Computed fields but zeros
				// out the required fields, creating invalid "ghost" entries in the state.
				// By marking the field as Computed, we can use CustomizeDiff to filter these out before
				// the Update phase, preventing them from being persisted to the state.
				Computed:      true,
				ConflictsWith: []string{"rule"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_number": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "allow",
						},

						"cidr_list": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},

						"icmp_type": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},

						"icmp_code": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},

						"port": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// Treat empty string as equivalent to not set (for "all" protocol)
								return old == "" && new == ""
							},
						},

						"traffic_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ingress",
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"parallelism": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},
	}
}

// Helper functions for UUID handling to abstract differences between
// 'rule' (uses uuids map) and 'ruleset' (uses uuid string)

// getRuleUUID gets the UUID for a rule, handling both formats
// For ruleset: returns the uuid string
// For rule with key: returns the UUID from uuids map for the given key
// For rule without key: returns the first UUID from uuids map
func getRuleUUID(rule map[string]interface{}, key string) (string, bool) {
	// Try uuid string first (ruleset format)
	if uuidVal, ok := rule["uuid"]; ok && uuidVal != nil {
		if uuid, ok := uuidVal.(string); ok && uuid != "" {
			return uuid, true
		}
	}

	// Try uuids map (rule format)
	if uuidsVal, ok := rule["uuids"]; ok && uuidsVal != nil {
		if uuids, ok := uuidsVal.(map[string]interface{}); ok {
			if key != "" {
				// Get specific key
				if idVal, ok := uuids[key]; ok && idVal != nil {
					if id, ok := idVal.(string); ok {
						return id, true
					}
				}
			} else {
				// Get first non-nil UUID
				for _, idVal := range uuids {
					if idVal != nil {
						if id, ok := idVal.(string); ok {
							return id, true
						}
					}
				}
			}
		}
	}

	return "", false
}

// setRuleUUID sets the UUID for a rule, handling both formats
// For ruleset: sets the uuid string
// For rule: sets the UUID in uuids map with the given key
func setRuleUUID(rule map[string]interface{}, key string, uuid string) {
	// Check if this is a ruleset (has uuid field) or rule (has uuids field)
	if _, hasUUID := rule["uuid"]; hasUUID {
		// Ruleset format
		rule["uuid"] = uuid
	} else {
		// Rule format - ensure uuids map exists
		var uuids map[string]interface{}
		if uuidsVal, ok := rule["uuids"]; ok && uuidsVal != nil {
			uuids = uuidsVal.(map[string]interface{})
		} else {
			uuids = make(map[string]interface{})
			rule["uuids"] = uuids
		}
		uuids[key] = uuid
	}
}

// hasRuleUUID checks if a rule has any UUID set
func hasRuleUUID(rule map[string]interface{}) bool {
	// Check uuid string (ruleset format)
	if uuidVal, ok := rule["uuid"]; ok && uuidVal != nil {
		if uuid, ok := uuidVal.(string); ok && uuid != "" {
			return true
		}
	}

	// Check uuids map (rule format)
	if uuidsVal, ok := rule["uuids"]; ok && uuidsVal != nil {
		if uuids, ok := uuidsVal.(map[string]interface{}); ok && len(uuids) > 0 {
			return true
		}
	}

	return false
}

// isRulesetRule checks if a rule is from a ruleset (has uuid field) vs rule (has uuids field)
func isRulesetRule(rule map[string]interface{}) bool {
	_, hasUUID := rule["uuid"]
	return hasUUID
}

// isGhostEntry checks if a rule is a ghost entry created by the SDK
// Ghost entries have empty protocol and rule_number=0 but may have a UUID
func isGhostEntry(rule map[string]interface{}) bool {
	protocol, _ := rule["protocol"].(string)
	ruleNumber, _ := rule["rule_number"].(int)
	return protocol == "" && ruleNumber == 0
}

// filterGhostEntries removes ghost entries from a list of rules
// Returns the cleaned list and the count of ghosts removed
func filterGhostEntries(rules []interface{}, logPrefix string) ([]interface{}, int) {
	var cleanRules []interface{}
	ghostCount := 0

	for i, r := range rules {
		rMap := r.(map[string]interface{})
		if isGhostEntry(rMap) {
			log.Printf("[DEBUG] %s: Filtering out ghost entry at index %d (uuid=%v)", logPrefix, i, rMap["uuid"])
			ghostCount++
			continue
		}
		cleanRules = append(cleanRules, r)
	}

	if ghostCount > 0 {
		log.Printf("[DEBUG] %s: Filtered %d ghost entries (%d -> %d rules)", logPrefix, ghostCount, len(rules), len(cleanRules))
	}

	return cleanRules, ghostCount
}

// assignRuleNumbers assigns rule numbers to rules that don't have them
// Rules are numbered sequentially starting from 1
// If a rule has an explicit rule_number, nextNumber advances to ensure no duplicates
// For rules using the deprecated 'ports' field with multiple ports, reserves enough numbers
func assignRuleNumbers(rules []interface{}) []interface{} {
	result := make([]interface{}, len(rules))
	nextNumber := 1

	for i, rule := range rules {
		ruleMap := make(map[string]interface{})
		// Copy the rule
		for k, v := range rule.(map[string]interface{}) {
			ruleMap[k] = v
		}

		// Check if rule_number is set
		if ruleNum, ok := ruleMap["rule_number"].(int); ok && ruleNum > 0 {
			// Rule has explicit number, ensure nextNumber never decreases
			// to prevent duplicate or decreasing rule numbers
			if ruleNum >= nextNumber {
				nextNumber = ruleNum + 1
			}
			log.Printf("[DEBUG] Rule at index %d has explicit rule_number=%d, nextNumber=%d", i, ruleNum, nextNumber)
		} else {
			// Auto-assign sequential number
			ruleMap["rule_number"] = nextNumber
			log.Printf("[DEBUG] Auto-assigned rule_number=%d to rule at index %d", nextNumber, i)

			// Check if this rule uses the deprecated 'ports' field with multiple ports
			// If so, we need to reserve additional rule numbers for the expanded rules
			if portsSet, ok := ruleMap["ports"].(*schema.Set); ok && portsSet.Len() > 1 {
				// Reserve portsSet.Len() numbers (one for each port)
				// The first port gets nextNumber, subsequent ports get nextNumber+1, nextNumber+2, etc.
				nextNumber += portsSet.Len()
				log.Printf("[DEBUG] Rule uses deprecated ports field with %d ports, reserved numbers up to %d", portsSet.Len(), nextNumber-1)
			} else {
				nextNumber++
			}
		}

		result[i] = ruleMap
	}

	return result
}

func resourceCloudStackNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Entering resourceCloudStackNetworkACLRuleCreate with acl_id=%s", d.Get("acl_id").(string))

	// Make sure all required parameters are there
	if err := verifyNetworkACLParams(d); err != nil {
		log.Printf("[ERROR] Failed parameter verification: %v", err)
		return err
	}

	// Handle 'rule' (TypeList with auto-numbering)
	if nrs := d.Get("rule").([]interface{}); len(nrs) > 0 {
		// Create an empty rule list to hold all newly created rules
		rules := make([]interface{}, 0)

		log.Printf("[DEBUG] Processing %d rules from 'rule' field", len(nrs))

		// Validate rules BEFORE assigning numbers, so we can detect user-provided rule_number
		if err := validateRulesList(d, nrs, "rule"); err != nil {
			return err
		}

		// Assign rule numbers to rules that don't have them
		rulesWithNumbers := assignRuleNumbers(nrs)

		err := createNetworkACLRules(d, meta, &rules, rulesWithNumbers)
		if err != nil {
			log.Printf("[ERROR] Failed to create network ACL rules: %v", err)
			return err
		}

		// Set the resource ID only after successful creation
		log.Printf("[DEBUG] Setting resource ID to acl_id=%s", d.Get("acl_id").(string))
		d.SetId(d.Get("acl_id").(string))

		// Update state with created rules
		if err := d.Set("rule", rules); err != nil {
			log.Printf("[ERROR] Failed to set rule attribute: %v", err)
			return err
		}
	} else if nrs := d.Get("ruleset").(*schema.Set); nrs.Len() > 0 {
		// Handle 'ruleset' (TypeSet with mandatory rule_number)
		rules := make([]interface{}, 0)

		log.Printf("[DEBUG] Processing %d rules from 'ruleset' field", nrs.Len())

		// Convert Set to list (no auto-numbering needed, rule_number is required)
		rulesList := nrs.List()

		// Validate rules BEFORE creating them
		if err := validateRulesList(d, rulesList, "ruleset"); err != nil {
			return err
		}

		err := createNetworkACLRules(d, meta, &rules, rulesList)
		if err != nil {
			log.Printf("[ERROR] Failed to create network ACL rules: %v", err)
			return err
		}

		// Set the resource ID only after successful creation
		log.Printf("[DEBUG] Setting resource ID to acl_id=%s", d.Get("acl_id").(string))
		d.SetId(d.Get("acl_id").(string))

		// Update state with created rules
		if err := d.Set("ruleset", rules); err != nil {
			log.Printf("[ERROR] Failed to set ruleset attribute: %v", err)
			return err
		}
	} else {
		log.Printf("[DEBUG] No rules provided, setting ID to acl_id=%s", d.Get("acl_id").(string))
		d.SetId(d.Get("acl_id").(string))
	}

	log.Printf("[DEBUG] Calling resourceCloudStackNetworkACLRuleRead")
	return resourceCloudStackNetworkACLRuleRead(d, meta)
}

func createNetworkACLRules(d *schema.ResourceData, meta interface{}, rules *[]interface{}, nrs []interface{}) error {
	log.Printf("[DEBUG] Creating %d network ACL rules", len(nrs))
	var errs *multierror.Error

	results := make([]map[string]interface{}, len(nrs))
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(len(nrs))

	sem := make(chan struct{}, d.Get("parallelism").(int))
	for i, rule := range nrs {
		// Put in a tiny sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(rule map[string]interface{}, index int) {
			defer wg.Done()
			sem <- struct{}{}

			log.Printf("[DEBUG] Creating rule #%d: %+v", index+1, rule)

			// Create a single rule
			err := createNetworkACLRule(d, meta, rule)
			if err != nil {
				log.Printf("[ERROR] Failed to create rule #%d: %v", index+1, err)
				mu.Lock()
				errs = multierror.Append(errs, fmt.Errorf("rule #%d: %v", index+1, err))
				mu.Unlock()
			} else {
				// Check if rule was created successfully (has uuid or uuids)
				if hasRuleUUID(rule) {
					log.Printf("[DEBUG] Successfully created rule #%d, storing at index %d", index+1, index)
					results[index] = rule
				} else {
					log.Printf("[WARN] Rule #%d created but has no UUID/UUIDs", index+1)
				}
			}

			<-sem
		}(rule.(map[string]interface{}), i)
	}

	wg.Wait()

	if err := errs.ErrorOrNil(); err != nil {
		log.Printf("[ERROR] Errors occurred while creating rules: %v", err)
		return err
	}

	for i, result := range results {
		if result != nil {
			*rules = append(*rules, result)
			log.Printf("[DEBUG] Added rule #%d to final rules list", i+1)
		}
	}

	log.Printf("[DEBUG] Successfully created all rules")
	return nil
}

func createNetworkACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	protocol := rule["protocol"].(string)
	action := rule["action"].(string)
	trafficType := rule["traffic_type"].(string)

	log.Printf("[DEBUG] Creating network ACL rule with protocol=%s, action=%s, traffic_type=%s", protocol, action, trafficType)

	// Note: Parameter verification is done before assignRuleNumbers in resourceCloudStackNetworkACLRuleCreate

	// Create a new parameter struct
	p := cs.NetworkACL.NewCreateNetworkACLParams(protocol)
	log.Printf("[DEBUG] Initialized CreateNetworkACLParams")

	// If a rule ID is specified, set it
	if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
		p.SetNumber(ruleNum)
		log.Printf("[DEBUG] Set rule_number=%d", ruleNum)
	}

	// Set the acl ID from the configuration
	aclID := d.Get("acl_id").(string)
	p.SetAclid(aclID)
	log.Printf("[DEBUG] Set aclid=%s", aclID)

	// Set the action
	p.SetAction(action)
	log.Printf("[DEBUG] Set action=%s", action)

	// Set the CIDR list
	var cidrList []string
	if cidrSet, ok := rule["cidr_list"].(*schema.Set); ok {
		for _, cidr := range cidrSet.List() {
			cidrList = append(cidrList, cidr.(string))
		}
	} else {
		// Fallback for 'rule' field which uses TypeList
		for _, cidr := range rule["cidr_list"].([]interface{}) {
			cidrList = append(cidrList, cidr.(string))
		}
	}
	p.SetCidrlist(cidrList)
	log.Printf("[DEBUG] Set cidr_list=%v", cidrList)

	// Set the traffic type
	p.SetTraffictype(trafficType)
	log.Printf("[DEBUG] Set traffic_type=%s", trafficType)

	// Set the description
	if desc, ok := rule["description"].(string); ok && desc != "" {
		p.SetReason(desc)
		log.Printf("[DEBUG] Set description=%s", desc)
	}

	// If the protocol is ICMP set the needed ICMP parameters
	if rule["protocol"].(string) == "icmp" {
		p.SetIcmptype(rule["icmp_type"].(int))
		p.SetIcmpcode(rule["icmp_code"].(int))
		log.Printf("[DEBUG] Set icmp_type=%d, icmp_code=%d", rule["icmp_type"].(int), rule["icmp_code"].(int))

		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			log.Printf("[ERROR] Failed to create ICMP rule: %v", err)
			return err
		}
		ruleID := r.(*cloudstack.CreateNetworkACLResponse).Id
		setRuleUUID(rule, "icmp", ruleID)
		log.Printf("[DEBUG] Created ICMP rule with ID=%s", ruleID)
	}

	// If the protocol is ALL set the needed parameters
	if rule["protocol"].(string) == "all" {
		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			log.Printf("[ERROR] Failed to create ALL rule: %v", err)
			return err
		}
		ruleID := r.(*cloudstack.CreateNetworkACLResponse).Id
		setRuleUUID(rule, "all", ruleID)
		log.Printf("[DEBUG] Created ALL rule with ID=%s", ruleID)
	}

	// If protocol is TCP or UDP, create the rule (with or without port)
	if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
		// Check if deprecated ports field is used (for backward compatibility)
		portsSet, hasPortsSet := rule["ports"].(*schema.Set)
		portStr, hasPort := rule["port"].(string)

		if hasPortsSet && portsSet.Len() > 0 {
			// Handle deprecated ports field for backward compatibility
			// Create a separate rule for each port in the set, each with a unique rule number
			log.Printf("[DEBUG] Using deprecated ports field for backward compatibility, creating %d rules", portsSet.Len())

			// Get the base rule number - this should always be set by assignRuleNumbers
			baseRuleNum := 0
			if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
				baseRuleNum = ruleNum
			}

			// Convert TypeSet to sorted list for deterministic rule number assignment
			// This ensures that rule numbers are stable across runs
			portsList := portsSet.List()
			portsStrings := make([]string, len(portsList))
			for i, port := range portsList {
				portsStrings[i] = port.(string)
			}
			sort.Strings(portsStrings)
			log.Printf("[DEBUG] Sorted ports for deterministic numbering: %v", portsStrings)

			portIndex := 0
			for _, portValue := range portsStrings {

				// Check if this port already has a UUID
				if _, hasUUID := getRuleUUID(rule, portValue); !hasUUID {
					m := splitPorts.FindStringSubmatch(portValue)
					if m == nil {
						log.Printf("[ERROR] Invalid port format: %s", portValue)
						return fmt.Errorf("%q is not a valid port value. Valid options are '80' or '80-90'", portValue)
					}

					startPort, err := strconv.Atoi(m[1])
					if err != nil {
						log.Printf("[ERROR] Failed to parse start port %s: %v", m[1], err)
						return err
					}

					endPort := startPort
					if m[2] != "" {
						endPort, err = strconv.Atoi(m[2])
						if err != nil {
							log.Printf("[ERROR] Failed to parse end port %s: %v", m[2], err)
							return err
						}
					}

					// Create a new parameter object for this specific port with a unique rule number
					portP := cs.NetworkACL.NewCreateNetworkACLParams(protocol)
					portP.SetAclid(aclID)
					portP.SetAction(action)
					portP.SetCidrlist(cidrList)
					portP.SetTraffictype(trafficType)
					if desc, ok := rule["description"].(string); ok && desc != "" {
						portP.SetReason(desc)
					}

					// Set a unique rule number for each port by adding the port index
					// This ensures each expanded rule gets a unique number
					uniqueRuleNum := baseRuleNum + portIndex
					portP.SetNumber(uniqueRuleNum)
					log.Printf("[DEBUG] Set unique rule_number=%d for port %s (base=%d, index=%d)", uniqueRuleNum, portValue, baseRuleNum, portIndex)

					portP.SetStartport(startPort)
					portP.SetEndport(endPort)
					log.Printf("[DEBUG] Set port start=%d, end=%d for deprecated ports field", startPort, endPort)

					r, err := Retry(4, retryableACLCreationFunc(cs, portP))
					if err != nil {
						log.Printf("[ERROR] Failed to create TCP/UDP rule for port %s: %v", portValue, err)
						return err
					}

					ruleID := r.(*cloudstack.CreateNetworkACLResponse).Id
					setRuleUUID(rule, portValue, ruleID)
					log.Printf("[DEBUG] Created TCP/UDP rule for port %s with ID=%s (deprecated ports field)", portValue, ruleID)

					portIndex++
				} else {
					log.Printf("[DEBUG] Port %s already has UUID, skipping", portValue)
					portIndex++
				}
			}
		} else if hasPort && portStr != "" {
			// Handle single port
			log.Printf("[DEBUG] Processing single port for TCP/UDP rule: %s", portStr)

			// Check if this port already has a UUID (for 'rule' field with uuids map)
			if _, hasUUID := getRuleUUID(rule, portStr); !hasUUID {
				m := splitPorts.FindStringSubmatch(portStr)
				if m == nil {
					log.Printf("[ERROR] Invalid port format: %s", portStr)
					return fmt.Errorf("%q is not a valid port value. Valid options are '80' or '80-90'", portStr)
				}

				startPort, err := strconv.Atoi(m[1])
				if err != nil {
					log.Printf("[ERROR] Failed to parse start port %s: %v", m[1], err)
					return err
				}

				endPort := startPort
				if m[2] != "" {
					endPort, err = strconv.Atoi(m[2])
					if err != nil {
						log.Printf("[ERROR] Failed to parse end port %s: %v", m[2], err)
						return err
					}
				}

				p.SetStartport(startPort)
				p.SetEndport(endPort)
				log.Printf("[DEBUG] Set port start=%d, end=%d", startPort, endPort)

				r, err := Retry(4, retryableACLCreationFunc(cs, p))
				if err != nil {
					log.Printf("[ERROR] Failed to create TCP/UDP rule for port %s: %v", portStr, err)
					return err
				}

				ruleID := r.(*cloudstack.CreateNetworkACLResponse).Id
				setRuleUUID(rule, portStr, ruleID)
				log.Printf("[DEBUG] Created TCP/UDP rule for port %s with ID=%s", portStr, ruleID)
			} else {
				log.Printf("[DEBUG] Port %s already has UUID, skipping", portStr)
			}
		} else {
			// No port specified - create rule for all ports
			log.Printf("[DEBUG] No port specified for TCP/UDP rule, creating rule for all ports")
			r, err := Retry(4, retryableACLCreationFunc(cs, p))
			if err != nil {
				log.Printf("[ERROR] Failed to create TCP/UDP rule for all ports: %v", err)
				return err
			}
			ruleID := r.(*cloudstack.CreateNetworkACLResponse).Id
			setRuleUUID(rule, "all_ports", ruleID)
			log.Printf("[DEBUG] Created TCP/UDP rule for all ports with ID=%s", ruleID)
		}
	}

	log.Printf("[DEBUG] Successfully created rule")
	return nil
}

func processTCPUDPRule(rule map[string]interface{}, ruleMap map[string]*cloudstack.NetworkACL, rules *[]interface{}) {
	// Check for deprecated ports field first (for reading existing state during migration)
	// This is only applicable to the legacy 'rule' field, not 'ruleset'
	ps, hasPortsSet := rule["ports"].(*schema.Set)
	portStr, hasPort := rule["port"].(string)

	if hasPortsSet && ps.Len() > 0 {
		// Only legacy 'rule' field supports deprecated ports
		log.Printf("[DEBUG] Processing deprecated ports field with %d ports during state read", ps.Len())

		// Create a new rule object to accumulate all ports
		newRule := make(map[string]interface{})
		newRule["uuids"] = make(map[string]interface{})

		// Process each port in the deprecated ports set during state read
		for _, port := range ps.List() {
			portStr := port.(string)

			if portRule, ok := processPortForRuleUnified(portStr, rule, ruleMap); ok {
				// Merge the port rule data into newRule
				for k, v := range portRule {
					if k == "uuids" {
						// Merge uuids maps
						if uuids, ok := v.(map[string]interface{}); ok {
							for uk, uv := range uuids {
								newRule["uuids"].(map[string]interface{})[uk] = uv
							}
						}
					} else {
						newRule[k] = v
					}
				}
				log.Printf("[DEBUG] Processed deprecated port %s during state read", portStr)
			}
		}

		// Only add the rule if we found at least one port
		if uuids, ok := newRule["uuids"].(map[string]interface{}); ok && len(uuids) > 0 {
			// Copy the ports field from the original rule
			newRule["ports"] = ps
			*rules = append(*rules, newRule)
			log.Printf("[DEBUG] Added TCP/UDP rule with deprecated ports to state during read: %+v", newRule)
		}

	} else if hasPort && portStr != "" {
		// Handle single port - works for both 'rule' and 'ruleset'
		log.Printf("[DEBUG] Processing single port for TCP/UDP rule: %s", portStr)

		if newRule, ok := processPortForRuleUnified(portStr, rule, ruleMap); ok {
			newRule["port"] = portStr
			*rules = append(*rules, newRule)
			log.Printf("[DEBUG] Added TCP/UDP rule with single port to state: %+v", newRule)
		}

	} else {
		// No port specified - create rule for all ports
		// Works for both 'rule' and 'ruleset'
		log.Printf("[DEBUG] Processing TCP/UDP rule with no port specified")

		if newRule, ok := processPortForRuleUnified("all_ports", rule, ruleMap); ok {
			*rules = append(*rules, newRule)
			log.Printf("[DEBUG] Added TCP/UDP rule with no port to state: %+v", newRule)
		}
	}
}

func processPortForRuleUnified(portKey string, rule map[string]interface{}, ruleMap map[string]*cloudstack.NetworkACL) (map[string]interface{}, bool) {
	// Get the UUID for this port (handles both 'rule' and 'ruleset' formats)
	id, ok := getRuleUUID(rule, portKey)
	if !ok {
		log.Printf("[DEBUG] No UUID for port %s, skipping", portKey)
		return nil, false
	}

	r, ok := ruleMap[id]
	if !ok {
		log.Printf("[DEBUG] TCP/UDP rule for port %s with ID %s not found", portKey, id)
		return nil, false
	}

	// Delete the known rule so only unknown rules remain in the ruleMap
	delete(ruleMap, id)

	// Create a NEW rule object instead of modifying the existing one
	newRule := make(map[string]interface{})

	// Create a list or set with all CIDR's depending on field type
	// Check if this is a ruleset rule (has uuid field) vs rule (has uuids field)
	_, isRuleset := rule["uuid"]
	if isRuleset {
		cidrs := &schema.Set{F: schema.HashString}
		for _, cidr := range strings.Split(r.Cidrlist, ",") {
			cidrs.Add(cidr)
		}
		newRule["cidr_list"] = cidrs
	} else {
		var cidrs []interface{}
		for _, cidr := range strings.Split(r.Cidrlist, ",") {
			cidrs = append(cidrs, cidr)
		}
		newRule["cidr_list"] = cidrs
	}

	newRule["action"] = strings.ToLower(r.Action)
	newRule["protocol"] = r.Protocol
	newRule["traffic_type"] = strings.ToLower(r.Traffictype)
	newRule["rule_number"] = r.Number
	newRule["description"] = r.Reason
	// Set ICMP fields to 0 for non-ICMP protocols to avoid spurious diffs
	newRule["icmp_type"] = 0
	newRule["icmp_code"] = 0

	// Copy the UUID field if it exists (for ruleset)
	if isRuleset {
		newRule["uuid"] = id
	} else {
		// For legacy 'rule' attribute, set uuids map
		newRule["uuids"] = map[string]interface{}{portKey: id}
	}

	return newRule, true
}

func resourceCloudStackNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// First check if the ACL itself still exists
	_, count, err := cs.NetworkACL.GetNetworkACLListByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project").(string)),
	)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Network ACL list %s does not exist", d.Id())
			d.SetId("")
			return nil
		}
		log.Printf("[ERROR] Failed to get ACL list by ID: %v", err)
		return err
	}

	// Get all the rules from the running environment with retries
	p := cs.NetworkACL.NewListNetworkACLsParams()
	p.SetAclid(d.Id())
	p.SetListall(true)

	var l *cloudstack.ListNetworkACLsResponse
	retryErr := retry.RetryContext(context.Background(), 30*time.Second, func() *retry.RetryError {
		var err error
		l, err = cs.NetworkACL.ListNetworkACLs(p)
		if err != nil {
			log.Printf("[DEBUG] Failed to list network ACL rules, retrying: %v", err)
			return retry.RetryableError(err)
		}
		if l.Count == 0 {
			log.Printf("[DEBUG] No network ACL rules found for ACL %s, retrying", d.Id())
			return retry.RetryableError(fmt.Errorf("no network ACL rules found for ACL %s", d.Id()))
		}
		log.Printf("[DEBUG] Found %d network ACL rules for ACL %s", l.Count, d.Id())
		return nil
	})

	if retryErr != nil {
		log.Printf("[WARN] Network ACL rules for %s not found after retries", d.Id())
		d.SetId("")
		return nil
	}

	// Make a map of all the rules so we can easily find a rule
	ruleMap := make(map[string]*cloudstack.NetworkACL, l.Count)
	for _, r := range l.NetworkACLs {
		ruleMap[r.Id] = r
	}
	log.Printf("[DEBUG] Loaded %d rules into ruleMap", len(ruleMap))

	// Create an empty rule list to hold all rules
	var rules []interface{}

	// Determine which field is being used and get the rules list
	var configuredRules []interface{}
	usingRuleset := false

	if rs := d.Get("ruleset").(*schema.Set); rs != nil && rs.Len() > 0 {
		usingRuleset = true
		configuredRules = rs.List()
	} else if rs := d.Get("rule").([]interface{}); len(rs) > 0 {
		configuredRules = rs
	}

	// Process all configured rules (works for both 'rule' and 'ruleset')
	for _, rule := range configuredRules {
		rule := rule.(map[string]interface{})

		protocol, _ := rule["protocol"].(string)

		if protocol == "" {
			continue
		}

		if protocol == "icmp" {
			id, ok := getRuleUUID(rule, "icmp")
			if !ok {
				log.Printf("[DEBUG] No ICMP UUID found, skipping rule")
				continue
			}

			// Get the rule
			r, ok := ruleMap[id]
			if !ok {
				log.Printf("[DEBUG] ICMP rule with ID %s not found", id)
				continue
			}

			// Delete the known rule so only unknown rules remain in the ruleMap
			delete(ruleMap, id)

			// Create a NEW rule object instead of modifying the existing one
			// This prevents corrupting the configuration data
			newRule := make(map[string]interface{})

			// Create a list or set with all CIDR's depending on field type
			if usingRuleset {
				cidrs := &schema.Set{F: schema.HashString}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs.Add(cidr)
				}
				newRule["cidr_list"] = cidrs
			} else {
				var cidrs []interface{}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs = append(cidrs, cidr)
				}
				newRule["cidr_list"] = cidrs
			}

			// Set the values from CloudStack
			newRule["action"] = strings.ToLower(r.Action)
			newRule["protocol"] = r.Protocol
			newRule["icmp_type"] = r.Icmptype
			newRule["icmp_code"] = r.Icmpcode
			newRule["traffic_type"] = strings.ToLower(r.Traffictype)
			newRule["rule_number"] = r.Number
			newRule["description"] = r.Reason
			if usingRuleset {
				newRule["uuid"] = id
			} else {
				newRule["uuids"] = map[string]interface{}{"icmp": id}
			}
			rules = append(rules, newRule)
			log.Printf("[DEBUG] Added ICMP rule to state: %+v", newRule)
		}

		if rule["protocol"].(string) == "all" {
			id, ok := getRuleUUID(rule, "all")
			if !ok {
				log.Printf("[DEBUG] No ALL UUID found, skipping rule")
				continue
			}

			// Get the rule
			r, ok := ruleMap[id]
			if !ok {
				log.Printf("[DEBUG] ALL rule with ID %s not found", id)
				continue
			}

			// Delete the known rule so only unknown rules remain in the ruleMap
			delete(ruleMap, id)

			// Create a NEW rule object instead of modifying the existing one
			newRule := make(map[string]interface{})

			// Create a list or set with all CIDR's depending on field type
			if usingRuleset {
				cidrs := &schema.Set{F: schema.HashString}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs.Add(cidr)
				}
				newRule["cidr_list"] = cidrs
			} else {
				var cidrs []interface{}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs = append(cidrs, cidr)
				}
				newRule["cidr_list"] = cidrs
			}

			// Set the values from CloudStack
			newRule["action"] = strings.ToLower(r.Action)
			newRule["protocol"] = r.Protocol
			newRule["traffic_type"] = strings.ToLower(r.Traffictype)
			newRule["rule_number"] = r.Number
			newRule["description"] = r.Reason
			// Set ICMP fields to 0 for non-ICMP protocols to avoid spurious diffs
			newRule["icmp_type"] = 0
			newRule["icmp_code"] = 0
			if usingRuleset {
				newRule["uuid"] = id
			} else {
				newRule["uuids"] = map[string]interface{}{"all": id}
			}
			rules = append(rules, newRule)
			log.Printf("[DEBUG] Added ALL rule to state: %+v", newRule)
		}

		if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
			processTCPUDPRule(rule, ruleMap, &rules)
		}
	}

	// If this is a managed ACL, add all unknown rules as out-of-band rule placeholders
	managed := d.Get("managed").(bool)
	if managed && len(ruleMap) > 0 {
		// Find the highest rule_number to avoid conflicts when creating out-of-band rule placeholders
		maxRuleNumber := 0
		for _, rule := range rules {
			if ruleMap, ok := rule.(map[string]interface{}); ok {
				if ruleNum, ok := ruleMap["rule_number"].(int); ok && ruleNum > maxRuleNumber {
					maxRuleNumber = ruleNum
				}
			}
		}

		// Start assigning out-of-band rule numbers after the highest existing rule_number
		outOfBandRuleNumber := maxRuleNumber + 1

		for uuid := range ruleMap {
			// Make a placeholder rule to hold the unknown UUID
			// Format differs between 'rule' and 'ruleset'
			var rule map[string]interface{}
			if usingRuleset {
				// For ruleset: use 'uuid' string and include rule_number
				// cidr_list is a TypeSet for ruleset
				cidrs := &schema.Set{F: schema.HashString}
				cidrs.Add(uuid)

				// Include all fields with defaults to avoid spurious diffs
				rule = map[string]interface{}{
					"cidr_list":    cidrs,
					"protocol":     uuid,
					"uuid":         uuid,
					"rule_number":  outOfBandRuleNumber,
					"action":       "allow",   // default value
					"traffic_type": "ingress", // default value
					"icmp_type":    0,         // default value
					"icmp_code":    0,         // default value
					"description":  "",        // empty string for optional field
					"port":         "",        // empty string for optional field
				}
				outOfBandRuleNumber++
			} else {
				// For rule: use 'uuids' map
				// cidr_list is a TypeList for rule
				cidrs := []interface{}{uuid}
				rule = map[string]interface{}{
					"cidr_list": cidrs,
					"protocol":  uuid,
					"uuids":     map[string]interface{}{uuid: uuid},
				}
			}

			// Add the out-of-band rule placeholder to the rules list
			rules = append(rules, rule)
			log.Printf("[DEBUG] Added out-of-band rule placeholder for UUID %s (usingRuleset=%t)", uuid, usingRuleset)
		}
	}

	// Always set the rules in state, even if empty (for managed=true case)
	if usingRuleset {
		// WORKAROUND: Filter out any ghost entries from the rules we're about to set
		// The SDK can create ghost entries with empty protocol/rule_number
		rules, _ = filterGhostEntries(rules, "Read")

		// For TypeSet, we need to be very careful about state updates
		// The SDK has issues with properly clearing removed elements from TypeSet
		// So we explicitly set to empty first, then set the new value
		// Use schema.HashResource to match the default hash function
		rulesetResource := resourceCloudStackNetworkACLRule().Schema["ruleset"].Elem.(*schema.Resource)
		hashFunc := schema.HashResource(rulesetResource)

		// First, clear the ruleset completely
		emptySet := schema.NewSet(hashFunc, []interface{}{})
		if err := d.Set("ruleset", emptySet); err != nil {
			log.Printf("[ERROR] Failed to clear ruleset attribute: %v", err)
			return err
		}

		// Now set the new rules
		newSet := schema.NewSet(hashFunc, rules)
		if err := d.Set("ruleset", newSet); err != nil {
			return err
		}
	} else {
		if err := d.Set("rule", rules); err != nil {
			log.Printf("[ERROR] Failed to set rule attribute: %v", err)
			return err
		}
	}

	if len(rules) == 0 && !managed {
		log.Printf("[DEBUG] No rules found and not managed, clearing ID")
		d.SetId("")
	}

	log.Printf("[DEBUG] Completed resourceCloudStackNetworkACLRuleRead")
	return nil
}

func resourceCloudStackNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	// Make sure all required parameters are there
	if err := verifyNetworkACLParams(d); err != nil {
		return err
	}

	// Check if the rule list has changed
	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		oldRules := o.([]interface{})
		newRules := n.([]interface{})

		log.Printf("[DEBUG] Rule list changed: %d old rules -> %d new rules", len(oldRules), len(newRules))

		// Check for migration from deprecated 'ports' to 'port' field
		migrationDetected := isPortsMigration(oldRules, newRules)

		if migrationDetected {
			log.Printf("[DEBUG] Migration detected - performing complete rule replacement")

			return performPortsMigration(d, meta, oldRules, newRules)
		}

		log.Printf("[DEBUG] Rule list changed, performing efficient updates")

		// Validate new rules BEFORE assigning numbers
		if err := validateRulesList(d, newRules, "rule"); err != nil {
			return err
		}

		// Assign rule numbers to new rules that don't have them
		newRulesWithNumbers := assignRuleNumbers(newRules)

		err := updateNetworkACLRules(d, meta, oldRules, newRulesWithNumbers)
		if err != nil {
			return err
		}
	}

	// Check if the ruleset has changed
	if d.HasChange("ruleset") {
		o, n := d.GetChange("ruleset")

		// WORKAROUND: The Terraform SDK has a bug where it creates "ghost" entries
		// when rules are removed from a TypeSet. These ghost entries have empty
		// protocol and rule_number=0 but retain the UUID from the deleted rule.
		// We need to filter them out BEFORE doing Set operations.
		cleanNewRules, _ := filterGhostEntries(n.(*schema.Set).List(), "Update")
		cleanOldRules, _ := filterGhostEntries(o.(*schema.Set).List(), "Update old")

		// Use the same sophisticated reconciliation logic as the 'rule' attribute
		// This will match rules by rule_number, update changed rules, and only
		// delete/create rules that truly disappeared/appeared
		cs := meta.(*cloudstack.CloudStackClient)
		err := performNormalRuleUpdates(d, meta, cs, cleanOldRules, cleanNewRules)
		if err != nil {
			return err
		}
	}

	// Call Read to refresh the state from the API
	// Read() already filters ghost entries, so we don't need to do it again here
	return resourceCloudStackNetworkACLRuleRead(d, meta)
}

func resourceCloudStackNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// Delete all rules from 'rule' field
	if ors := d.Get("rule").([]interface{}); len(ors) > 0 {
		for _, rule := range ors {
			ruleMap := rule.(map[string]interface{})
			err := deleteNetworkACLRule(d, meta, ruleMap)
			if err != nil {
				log.Printf("[ERROR] Failed to delete rule: %v", err)
				return err
			}
		}
	}

	// Delete all rules from 'ruleset' field
	if ors := d.Get("ruleset").(*schema.Set); ors != nil && ors.Len() > 0 {
		for _, rule := range ors.List() {
			ruleMap := rule.(map[string]interface{})
			err := deleteNetworkACLRule(d, meta, ruleMap)
			if err != nil {
				log.Printf("[ERROR] Failed to delete ruleset rule: %v", err)
				return err
			}
		}
	}

	return nil
}

func deleteNetworkACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	if isRulesetRule(rule) {
		// For ruleset, delete the single UUID
		if uuid, ok := getRuleUUID(rule, ""); ok {
			if err := deleteSingleACL(cs, uuid); err != nil {
				return err
			}
			// Don't modify the rule object - it's from the old state and modifying it
			// can cause issues with TypeSet state management
		}
	} else {
		// For rule, delete all UUIDs from the map
		if uuidsVal, ok := rule["uuids"]; ok && uuidsVal != nil {
			uuids := uuidsVal.(map[string]interface{})
			for k, id := range uuids {
				// Skip the count field
				if k == "%" {
					continue
				}
				if idStr, ok := id.(string); ok {
					if err := deleteSingleACL(cs, idStr); err != nil {
						return err
					}
					// Don't modify the uuids map - it's from the old state
				}
			}
		}
	}

	return nil
}

func deleteSingleACL(cs *cloudstack.CloudStackClient, id string) error {
	log.Printf("[DEBUG] Deleting ACL rule with UUID=%s", id)

	p := cs.NetworkACL.NewDeleteNetworkACLParams(id)
	if _, err := cs.NetworkACL.DeleteNetworkACL(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", id)) {
			// ID doesn't exist, which is fine for delete
			return nil
		}
		return err
	}
	return nil
}

func verifyNetworkACLParams(d *schema.ResourceData) error {
	managed := d.Get("managed").(bool)
	_, rules := d.GetOk("rule")
	_, ruleset := d.GetOk("ruleset")

	if !rules && !ruleset && !managed {
		return fmt.Errorf(
			"You must supply at least one 'rule' or 'ruleset' when not using the 'managed' firewall feature")
	}

	return nil
}

// validateRulesList validates all rules in a list by calling verifyNetworkACLRuleParams on each
// This helper consolidates the validation logic used in Create and Update paths for both 'rule' and 'ruleset' fields
// Out-of-band rule placeholders (created by managed=true) are skipped as they are markers for deletion
func validateRulesList(d *schema.ResourceData, rules []interface{}, fieldName string) error {
	validatedCount := 0
	for i, rule := range rules {
		ruleMap := rule.(map[string]interface{})

		// Skip validation for out-of-band rule placeholders
		// These are created by managed=true and are just markers for deletion
		if isOutOfBandRulePlaceholder(ruleMap) {
			log.Printf("[DEBUG] Skipping validation for out-of-band rule placeholder at index %d", i)
			continue
		}

		if err := verifyNetworkACLRuleParams(d, ruleMap); err != nil {
			log.Printf("[ERROR] Failed to verify %s rule %d parameters: %v", fieldName, i, err)
			return fmt.Errorf("validation failed for %s rule %d: %w", fieldName, i, err)
		}
		validatedCount++
	}
	log.Printf("[DEBUG] Successfully validated %d %s rules (skipped %d out-of-band placeholders)", validatedCount, fieldName, len(rules)-validatedCount)
	return nil
}

func verifyNetworkACLRuleParams(d *schema.ResourceData, rule map[string]interface{}) error {
	log.Printf("[DEBUG] Verifying parameters for rule: %+v", rule)

	if ruleNum, ok := rule["rule_number"]; ok && ruleNum != nil {
		if number, ok := ruleNum.(int); ok && number != 0 {
			// Validate only if rule_number is explicitly set (non-zero)
			if number < 1 || number > 65535 {
				log.Printf("[ERROR] Invalid rule_number: %d", number)
				return fmt.Errorf(
					"%q must be between %d and %d inclusive, got: %d", "rule_number", 1, 65535, number)
			}
		}
	}

	action := rule["action"].(string)
	if action != "allow" && action != "deny" {
		log.Printf("[ERROR] Invalid action: %s", action)
		return fmt.Errorf("Parameter action only accepts 'allow' or 'deny' as values")
	}

	protocol := rule["protocol"].(string)
	log.Printf("[DEBUG] Validating protocol: %s", protocol)
	switch protocol {
	case "icmp":
		if _, ok := rule["icmp_type"]; !ok {
			log.Printf("[ERROR] Missing icmp_type for ICMP protocol")
			return fmt.Errorf(
				"Parameter icmp_type is a required parameter when using protocol 'icmp'")
		}
		if _, ok := rule["icmp_code"]; !ok {
			log.Printf("[ERROR] Missing icmp_code for ICMP protocol")
			return fmt.Errorf(
				"Parameter icmp_code is a required parameter when using protocol 'icmp'")
		}
	case "all":
		// No additional test are needed
		log.Printf("[DEBUG] Protocol 'all' validated")
	case "tcp", "udp":
		// The deprecated 'ports' field is allowed for backward compatibility
		// but users should migrate to the 'port' field
		portsSet, hasPortsSet := rule["ports"].(*schema.Set)
		portStr, hasPort := rule["port"].(string)

		// Allow deprecated ports field for backward compatibility
		// The schema already marks it as deprecated with a warning
		if hasPortsSet && portsSet.Len() > 0 {
			log.Printf("[DEBUG] Using deprecated ports field for backward compatibility")

			// When using deprecated ports field with multiple values, rule_number cannot be specified
			// because we auto-generate sequential rule numbers for each port
			if portsSet.Len() > 1 {
				if ruleNum, ok := rule["rule_number"]; ok && ruleNum != nil {
					if number, ok := ruleNum.(int); ok && number > 0 {
						log.Printf("[ERROR] Cannot specify rule_number when using deprecated ports field with multiple values")
						return fmt.Errorf(
							"Cannot specify 'rule_number' when using deprecated 'ports' field with multiple values. " +
								"Rule numbers are auto-generated for each port (starting from the auto-assigned base number). " +
								"Either use a single port in 'ports', or omit 'rule_number', or migrate to the 'port' field.")
					}
				}
			}

			// Validate each port in the set
			for _, p := range portsSet.List() {
				portValue := p.(string)
				m := splitPorts.FindStringSubmatch(portValue)
				if m == nil {
					log.Printf("[ERROR] Invalid port format in ports field: %s", portValue)
					return fmt.Errorf(
						"%q is not a valid port value. Valid options are '80' or '80-90'", portValue)
				}
			}
		}

		// Validate the new port field if used
		if hasPort && portStr != "" {
			log.Printf("[DEBUG] Found port for TCP/UDP: %s", portStr)
			m := splitPorts.FindStringSubmatch(portStr)
			if m == nil {
				log.Printf("[ERROR] Invalid port format: %s", portStr)
				return fmt.Errorf(
					"%q is not a valid port value. Valid options are '80' or '80-90'", portStr)
			}
		}

		// If neither port nor ports is specified, that's also valid (allows all ports)
		if (!hasPort || portStr == "") && (!hasPortsSet || portsSet.Len() == 0) {
			log.Printf("[DEBUG] No port specified for TCP/UDP, allowing all ports")
		}
	default:
		_, err := strconv.ParseInt(protocol, 0, 0)
		if err != nil {
			log.Printf("[ERROR] Invalid protocol: %s", protocol)
			return fmt.Errorf(
				"%q is not a valid protocol. Valid options are 'tcp', 'udp', 'icmp', 'all' or a valid protocol number", protocol)
		}
	}

	traffic := rule["traffic_type"].(string)
	if traffic != "ingress" && traffic != "egress" {
		log.Printf("[ERROR] Invalid traffic_type: %s", traffic)
		return fmt.Errorf(
			"Parameter traffic_type only accepts 'ingress' or 'egress' as values")
	}

	log.Printf("[DEBUG] Rule parameters verified successfully")
	return nil
}

func retryableACLCreationFunc(
	cs *cloudstack.CloudStackClient,
	p *cloudstack.CreateNetworkACLParams) func() (interface{}, error) {
	return func() (interface{}, error) {
		r, err := cs.NetworkACL.CreateNetworkACL(p)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}

func resourceCloudStackNetworkACLRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	cs := meta.(*cloudstack.CloudStackClient)

	aclID := d.Id()

	log.Printf("[DEBUG] Attempting to import ACL list with ID: %s", aclID)
	if aclExists, err := checkACLListExists(cs, aclID); err != nil {
		return nil, fmt.Errorf("error checking ACL list existence: %v", err)
	} else if !aclExists {
		return nil, fmt.Errorf("ACL list with ID %s does not exist", aclID)
	}

	log.Printf("[DEBUG] Found ACL list with ID: %s", aclID)
	d.Set("acl_id", aclID)

	log.Printf("[DEBUG] Setting managed=true for ACL list import")
	d.Set("managed", true)

	return []*schema.ResourceData{d}, nil
}

func checkACLListExists(cs *cloudstack.CloudStackClient, aclID string) (bool, error) {
	log.Printf("[DEBUG] Checking if ACL list exists: %s", aclID)
	_, count, err := cs.NetworkACL.GetNetworkACLListByID(aclID)
	if err != nil {
		log.Printf("[DEBUG] Error getting ACL list by ID: %v", err)
		return false, err
	}

	log.Printf("[DEBUG] ACL list check result: count=%d", count)
	return count > 0, nil
}

// isOutOfBandRulePlaceholder checks if a rule is a placeholder for an out-of-band rule
// (created by managed=true for rules that exist in CloudStack but not in config)
// Out-of-band rule placeholders are identified by having protocol == uuid, OR by having
// an empty protocol with rule_number == 0 (TypeSet reconciliation creates these)
func isOutOfBandRulePlaceholder(rule map[string]interface{}) bool {
	protocol, hasProtocol := rule["protocol"].(string)
	uuid, hasUUID := getRuleUUID(rule, "")

	if !hasUUID || uuid == "" {
		return false
	}

	// Case 1: protocol equals uuid (original out-of-band rule placeholder in state)
	if hasProtocol && protocol == uuid {
		return true
	}

	// Case 2: protocol is empty and rule_number is 0
	// This happens when TypeSet reconciles an out-of-band rule placeholder from state but zeros out the fields
	if hasProtocol && protocol == "" {
		if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum == 0 {
			return true
		}
	}

	return false
}

func updateNetworkACLRules(d *schema.ResourceData, meta interface{}, oldRules, newRules []interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Updating ACL rules: %d old rules, %d new rules", len(oldRules), len(newRules))

	log.Printf("[DEBUG] Performing normal rule updates")
	return performNormalRuleUpdates(d, meta, cs, oldRules, newRules)
}

func performNormalRuleUpdates(d *schema.ResourceData, meta interface{}, cs *cloudstack.CloudStackClient, oldRules, newRules []interface{}) error {
	rulesToUpdate := make(map[string]map[string]interface{}) // UUID -> new rule mapping
	rulesToDelete := make([]map[string]interface{}, 0)
	rulesToCreate := make([]map[string]interface{}, 0)

	// Track which new rules match existing old rules
	usedNewRules := make(map[int]bool)

	// For each old rule, try to find a matching new rule
	for _, oldRule := range oldRules {
		oldRuleMap := oldRule.(map[string]interface{})
		foundMatch := false

		for newIdx, newRule := range newRules {
			if usedNewRules[newIdx] {
				continue
			}

			newRuleMap := newRule.(map[string]interface{})
			log.Printf("[DEBUG] Comparing old rule %+v with new rule %+v", oldRuleMap, newRuleMap)

			// For ruleset rules, match by rule_number only
			// For regular rules, use the full rulesMatch function
			var matched bool
			if isRulesetRule(oldRuleMap) && isRulesetRule(newRuleMap) {
				matched = rulesetRulesMatchByNumber(oldRuleMap, newRuleMap)
			} else {
				matched = rulesMatch(oldRuleMap, newRuleMap)
			}

			if matched {
				log.Printf("[DEBUG] Found matching new rule for old rule")

				// Copy UUID from old rule to new rule (following port_forward pattern)
				// This preserves the UUID across updates
				if isRulesetRule(oldRuleMap) {
					// Ruleset format: single uuid string
					if uuid, ok := oldRuleMap["uuid"].(string); ok && uuid != "" {
						newRuleMap["uuid"] = uuid
					}
				} else {
					// Rule format: uuids map
					if uuids, ok := oldRuleMap["uuids"].(map[string]interface{}); ok {
						newRuleMap["uuids"] = uuids
					}
				}

				if ruleNeedsUpdate(oldRuleMap, newRuleMap) {
					log.Printf("[DEBUG] Rule needs updating")
					// Get UUID for update (use empty key to get first UUID)
					if updateUUID, ok := getRuleUUID(oldRuleMap, ""); ok {
						rulesToUpdate[updateUUID] = newRuleMap
					}
				}

				usedNewRules[newIdx] = true
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			log.Printf("[DEBUG] Old rule has no match, will be deleted")
			rulesToDelete = append(rulesToDelete, oldRuleMap)
		}
	}

	for newIdx, newRule := range newRules {
		if !usedNewRules[newIdx] {
			newRuleMap := newRule.(map[string]interface{})

			// Skip out-of-band rule placeholders (created by managed=true for out-of-band rules)
			// These placeholders should not be created - they're just markers for deletion
			if isOutOfBandRulePlaceholder(newRuleMap) {
				log.Printf("[DEBUG] Skipping out-of-band rule placeholder (will not create)")
				continue
			}

			log.Printf("[DEBUG] New rule has no match, will be created")
			rulesToCreate = append(rulesToCreate, newRuleMap)
		}
	}

	for _, ruleToDelete := range rulesToDelete {
		log.Printf("[DEBUG] Deleting unmatched old rule")
		err := deleteNetworkACLRule(d, meta, ruleToDelete)
		if err != nil {
			return fmt.Errorf("failed to delete old rule: %v", err)
		}
	}

	for uuid, newRule := range rulesToUpdate {
		log.Printf("[DEBUG] Updating rule with UUID %s", uuid)

		tempOldRule := make(map[string]interface{})
		tempOldRule["uuids"] = map[string]interface{}{"update": uuid}

		err := updateNetworkACLRule(cs, tempOldRule, newRule)
		if err != nil {
			return fmt.Errorf("failed to update rule UUID %s: %v", uuid, err)
		}
	}

	if len(rulesToCreate) > 0 {
		log.Printf("[DEBUG] Creating %d new rules", len(rulesToCreate))

		var createdRules []interface{}
		var rulesToCreateInterface []interface{}
		for _, rule := range rulesToCreate {
			rulesToCreateInterface = append(rulesToCreateInterface, rule)
		}

		err := createNetworkACLRules(d, meta, &createdRules, rulesToCreateInterface)
		if err != nil {
			return fmt.Errorf("failed to create new rules: %v", err)
		}
	}

	return nil
}

// rulesetRulesMatchByNumber matches ruleset rules by rule_number only
// This allows changes to other fields (CIDR, port, protocol, etc.) to be detected as updates
func rulesetRulesMatchByNumber(oldRule, newRule map[string]interface{}) bool {
	oldRuleNum, oldHasRuleNum := oldRule["rule_number"].(int)
	newRuleNum, newHasRuleNum := newRule["rule_number"].(int)

	// Both must have rule_number and they must match
	if !oldHasRuleNum || !newHasRuleNum {
		return false
	}

	return oldRuleNum == newRuleNum
}

func rulesMatch(oldRule, newRule map[string]interface{}) bool {
	oldProtocol := oldRule["protocol"].(string)
	newProtocol := newRule["protocol"].(string)
	oldTrafficType := oldRule["traffic_type"].(string)
	newTrafficType := newRule["traffic_type"].(string)
	oldAction := oldRule["action"].(string)
	newAction := newRule["action"].(string)

	if oldProtocol != newProtocol ||
		oldTrafficType != newTrafficType ||
		oldAction != newAction {
		return false
	}

	protocol := newProtocol

	if protocol == "tcp" || protocol == "udp" {
		oldPort, oldHasPort := oldRule["port"].(string)
		newPort, newHasPort := newRule["port"].(string)

		if oldHasPort && newHasPort {
			return oldPort == newPort
		}

		if oldHasPort != newHasPort {
			return false
		}

		return true
	}

	switch protocol {
	case "icmp":
		return oldRule["icmp_type"].(int) == newRule["icmp_type"].(int) &&
			oldRule["icmp_code"].(int) == newRule["icmp_code"].(int)

	case "all":
		return true

	default:
		return true
	}
}

func ruleNeedsUpdate(oldRule, newRule map[string]interface{}) bool {
	oldAction := oldRule["action"].(string)
	newAction := newRule["action"].(string)
	if oldAction != newAction {
		log.Printf("[DEBUG] Action changed: %s -> %s", oldAction, newAction)
		return true
	}

	oldProtocol := oldRule["protocol"].(string)
	newProtocol := newRule["protocol"].(string)
	if oldProtocol != newProtocol {
		log.Printf("[DEBUG] Protocol changed: %s -> %s", oldProtocol, newProtocol)
		return true
	}

	oldTrafficType := oldRule["traffic_type"].(string)
	newTrafficType := newRule["traffic_type"].(string)
	if oldTrafficType != newTrafficType {
		log.Printf("[DEBUG] Traffic type changed: %s -> %s", oldTrafficType, newTrafficType)
		return true
	}

	// Check rule_number
	oldRuleNum, oldHasRuleNum := oldRule["rule_number"].(int)
	newRuleNum, newHasRuleNum := newRule["rule_number"].(int)
	if oldHasRuleNum != newHasRuleNum || (oldHasRuleNum && newHasRuleNum && oldRuleNum != newRuleNum) {
		log.Printf("[DEBUG] Rule number changed: %d -> %d", oldRuleNum, newRuleNum)
		return true
	}

	oldDesc, oldHasDesc := oldRule["description"].(string)
	newDesc, newHasDesc := newRule["description"].(string)
	if oldHasDesc != newHasDesc || (oldHasDesc && newHasDesc && oldDesc != newDesc) {
		log.Printf("[DEBUG] Description changed: %s -> %s", oldDesc, newDesc)
		return true
	}

	// Use newProtocol from earlier
	switch newProtocol {
	case "icmp":
		// Helper function to get int value with default
		getInt := func(rule map[string]interface{}, key string, defaultVal int) int {
			if val, ok := rule[key]; ok && val != nil {
				if i, ok := val.(int); ok {
					return i
				}
			}
			return defaultVal
		}

		oldIcmpType := getInt(oldRule, "icmp_type", 0)
		newIcmpType := getInt(newRule, "icmp_type", 0)
		if oldIcmpType != newIcmpType {
			log.Printf("[DEBUG] ICMP type changed: %d -> %d", oldIcmpType, newIcmpType)
			return true
		}

		oldIcmpCode := getInt(oldRule, "icmp_code", 0)
		newIcmpCode := getInt(newRule, "icmp_code", 0)
		if oldIcmpCode != newIcmpCode {
			log.Printf("[DEBUG] ICMP code changed: %d -> %d", oldIcmpCode, newIcmpCode)
			return true
		}
	case "tcp", "udp":
		oldPort, oldHasPort := oldRule["port"].(string)
		newPort, newHasPort := newRule["port"].(string)
		if oldHasPort != newHasPort || (oldHasPort && newHasPort && oldPort != newPort) {
			log.Printf("[DEBUG] Port changed: %s -> %s", oldPort, newPort)
			return true
		}
	}

	// Handle cidr_list comparison - can be TypeSet (ruleset) or TypeList (rule)
	var oldCidrStrs, newCidrStrs []string

	// Extract old CIDRs
	if oldSet, ok := oldRule["cidr_list"].(*schema.Set); ok {
		for _, cidr := range oldSet.List() {
			oldCidrStrs = append(oldCidrStrs, cidr.(string))
		}
	} else if oldList, ok := oldRule["cidr_list"].([]interface{}); ok {
		for _, cidr := range oldList {
			oldCidrStrs = append(oldCidrStrs, cidr.(string))
		}
	}

	// Extract new CIDRs
	if newSet, ok := newRule["cidr_list"].(*schema.Set); ok {
		for _, cidr := range newSet.List() {
			newCidrStrs = append(newCidrStrs, cidr.(string))
		}
	} else if newList, ok := newRule["cidr_list"].([]interface{}); ok {
		for _, cidr := range newList {
			newCidrStrs = append(newCidrStrs, cidr.(string))
		}
	}

	if len(oldCidrStrs) != len(newCidrStrs) {
		log.Printf("[DEBUG] CIDR list length changed: %d -> %d", len(oldCidrStrs), len(newCidrStrs))
		return true
	}

	sort.Strings(oldCidrStrs)
	sort.Strings(newCidrStrs)

	for i, oldCidr := range oldCidrStrs {
		if oldCidr != newCidrStrs[i] {
			log.Printf("[DEBUG] CIDR changed at index %d: %s -> %s", i, oldCidr, newCidrStrs[i])
			return true
		}
	}

	return false
}

func updateNetworkACLRule(cs *cloudstack.CloudStackClient, oldRule, newRule map[string]interface{}) error {
	uuids := oldRule["uuids"].(map[string]interface{})

	for key, uuid := range uuids {
		if key == "%" {
			continue
		}

		log.Printf("[DEBUG] Updating ACL rule with UUID: %s", uuid.(string))
		p := cs.NetworkACL.NewUpdateNetworkACLItemParams(uuid.(string))

		p.SetAction(newRule["action"].(string))

		var cidrList []string
		if cidrSet, ok := newRule["cidr_list"].(*schema.Set); ok {
			for _, cidr := range cidrSet.List() {
				cidrList = append(cidrList, cidr.(string))
			}
		} else {
			for _, cidr := range newRule["cidr_list"].([]interface{}) {
				cidrList = append(cidrList, cidr.(string))
			}
		}
		p.SetCidrlist(cidrList)

		// Set description from the new rule
		if desc, ok := newRule["description"].(string); ok {
			p.SetReason(desc)
		} else {
			p.SetReason("")
		}

		p.SetProtocol(newRule["protocol"].(string))

		p.SetTraffictype(newRule["traffic_type"].(string))

		// Set rule number if provided and non-zero
		if ruleNum, ok := newRule["rule_number"].(int); ok && ruleNum > 0 {
			p.SetNumber(ruleNum)
			log.Printf("[DEBUG] Set rule_number=%d", ruleNum)
		}

		protocol := newRule["protocol"].(string)
		switch protocol {
		case "icmp":
			if icmpType, ok := newRule["icmp_type"].(int); ok {
				p.SetIcmptype(icmpType)
				log.Printf("[DEBUG] Set icmp_type=%d", icmpType)
			}
			if icmpCode, ok := newRule["icmp_code"].(int); ok {
				p.SetIcmpcode(icmpCode)
				log.Printf("[DEBUG] Set icmp_code=%d", icmpCode)
			}
		case "tcp", "udp":
			if portStr, hasPort := newRule["port"].(string); hasPort && portStr != "" {
				m := splitPorts.FindStringSubmatch(portStr)
				if m != nil {
					startPort, err := strconv.Atoi(m[1])
					if err == nil {
						endPort := startPort
						if m[2] != "" {
							if ep, err := strconv.Atoi(m[2]); err == nil {
								endPort = ep
							}
						}
						p.SetStartport(startPort)
						p.SetEndport(endPort)
						log.Printf("[DEBUG] Set port start=%d, end=%d", startPort, endPort)
					}
				}
			}
		}

		_, err := cs.NetworkACL.UpdateNetworkACLItem(p)
		if err != nil {
			log.Printf("[ERROR] Failed to update ACL rule %s: %v", uuid.(string), err)
			return err
		}

		log.Printf("[DEBUG] Successfully updated ACL rule %s", uuid.(string))
	}

	return nil
}

func hasDeprecatedPortsInOldRules(oldRules []interface{}) bool {
	for _, oldRule := range oldRules {
		oldRuleMap := oldRule.(map[string]interface{})
		protocol := oldRuleMap["protocol"].(string)

		if protocol == "tcp" || protocol == "udp" {
			if portsSet, hasPortsSet := oldRuleMap["ports"].(*schema.Set); hasPortsSet && portsSet.Len() > 0 {
				return true
			}
		}
	}
	return false
}

func containsMixedPortFields(oldRules, newRules []interface{}) bool {
	hasDeprecatedInOld := hasDeprecatedPortsInOldRules(oldRules)
	hasNewInNew := hasPortFieldInNewRules(newRules)

	hasDeprecatedInNew := hasDeprecatedPortsInOldRules(newRules)

	// Migration detected if:
	// 1. Old rules have deprecated ports OR
	// 2. We have a mix of deprecated and new port fields anywhere
	return hasDeprecatedInOld || (hasDeprecatedInNew && hasNewInNew)
}

// Checks if any new rule uses the new 'port' field
func hasPortFieldInNewRules(newRules []interface{}) bool {
	for _, newRule := range newRules {
		newRuleMap := newRule.(map[string]interface{})
		protocol := newRuleMap["protocol"].(string)

		if protocol == "tcp" || protocol == "udp" {
			if portStr, hasPort := newRuleMap["port"].(string); hasPort && portStr != "" {
				return true
			}
		}
	}
	return false
}

// Detects if we're migrating from deprecated 'ports' to 'port' field
func isPortsMigration(oldRules, newRules []interface{}) bool {
	log.Printf("[DEBUG] Migration detection: checking %d old rules and %d new rules", len(oldRules), len(newRules))

	hasDeprecatedPorts := false
	hasNewPortFormat := false

	for i, oldRule := range oldRules {
		oldRuleMap := oldRule.(map[string]interface{})
		protocol := oldRuleMap["protocol"].(string)
		log.Printf("[DEBUG] Migration detection: old rule %d has protocol %s", i, protocol)

		if protocol == "tcp" || protocol == "udp" {
			if portsSet, hasPortsSet := oldRuleMap["ports"].(*schema.Set); hasPortsSet && portsSet.Len() > 0 {
				log.Printf("[DEBUG] Migration detection: old rule %d has deprecated ports field with %d ports", i, portsSet.Len())
				hasDeprecatedPorts = true
			}

			oldPort, oldHasPort := oldRuleMap["port"].(string)
			if !oldHasPort || oldPort == "" {
				log.Printf("[DEBUG] Migration detection: old rule %d has no port field, checking if new rules use port field", i)
				for j, newRule := range newRules {
					newRuleMap := newRule.(map[string]interface{})
					newProtocol := newRuleMap["protocol"].(string)
					if newProtocol == protocol {
						if newPortStr, newHasPort := newRuleMap["port"].(string); newHasPort && newPortStr != "" {
							log.Printf("[DEBUG] Migration detection: new rule %d has port field '%s' while old rule had none - potential migration", j, newPortStr)
							hasDeprecatedPorts = true
							break
						}
					}
				}
			}
		}
	}

	for i, newRule := range newRules {
		newRuleMap := newRule.(map[string]interface{})
		protocol := newRuleMap["protocol"].(string)
		log.Printf("[DEBUG] Migration detection: new rule %d has protocol %s", i, protocol)

		if protocol == "tcp" || protocol == "udp" {
			if portStr, hasPort := newRuleMap["port"].(string); hasPort && portStr != "" {
				log.Printf("[DEBUG] Migration detection: new rule %d has port field with value: %s", i, portStr)
				hasNewPortFormat = true
			}

			if portsSet, hasPortsSet := newRuleMap["ports"].(*schema.Set); hasPortsSet && portsSet.Len() > 0 {
				log.Printf("[DEBUG] Migration detection: new rule %d still has deprecated ports, not a migration", i)
				return false
			}
		}
	}

	migrationDetected := hasDeprecatedPorts && hasNewPortFormat
	log.Printf("[DEBUG] Migration detection result: hasDeprecatedPorts=%t, hasNewPortFormat=%t, migrationDetected=%t", hasDeprecatedPorts, hasNewPortFormat, migrationDetected)

	// Migration is detected if:
	// 1. We have old rules with deprecated ports OR no port field AND
	// 2. We have new rules with port format (no deprecated ports)
	return migrationDetected
}

func performPortsMigration(d *schema.ResourceData, meta interface{}, oldRules, newRules []interface{}) error {
	log.Printf("[DEBUG] Starting ports->port migration")
	cs := meta.(*cloudstack.CloudStackClient)

	// Build a map of all UUIDs that need to be deleted
	uuidsToDelete := make([]string, 0)

	for _, oldRule := range oldRules {
		oldRuleMap := oldRule.(map[string]interface{})
		uuids, ok := oldRuleMap["uuids"].(map[string]interface{})
		if !ok {
			continue
		}

		for key, uuid := range uuids {
			if key != "%" && uuid != nil {
				uuidStr := uuid.(string)
				if uuidStr != "" {
					uuidsToDelete = append(uuidsToDelete, uuidStr)
				}
			}
		}
	}

	log.Printf("[DEBUG] Total UUIDs to delete: %d", len(uuidsToDelete))

	// Delete all old rules by UUID and wait for completion
	for _, uuidToDelete := range uuidsToDelete {
		p := cs.NetworkACL.NewDeleteNetworkACLParams(uuidToDelete)
		_, err := cs.NetworkACL.DeleteNetworkACL(p)

		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf(
				"Invalid parameter id value=%s due to incorrect long value format, "+
					"or entity does not exist", uuidToDelete)) {
				continue
			}

			return fmt.Errorf("failed to delete old rule UUID %s during migration: %v", uuidToDelete, err)
		}
	}

	// Wait a moment for CloudStack to process the deletions
	if len(uuidsToDelete) > 0 {
		log.Printf("[DEBUG] Waiting for CloudStack to process %d rule deletions", len(uuidsToDelete))
		time.Sleep(3 * time.Second)

		for _, uuidToCheck := range uuidsToDelete {
			listParams := cs.NetworkACL.NewListNetworkACLsParams()
			listParams.SetId(uuidToCheck)

			listResp, err := cs.NetworkACL.ListNetworkACLs(listParams)
			if err == nil && listResp.Count > 0 {
				time.Sleep(2 * time.Second)
				break
			}
		}
	}

	// Create all new rules with fresh UUIDs
	if len(newRules) > 0 {
		log.Printf("[DEBUG] Creating %d new rules with port field", len(newRules))

		var rulesToCreate []interface{}
		for _, newRule := range newRules {
			newRuleMap := newRule.(map[string]interface{})

			cleanRule := make(map[string]interface{})
			for k, v := range newRuleMap {
				cleanRule[k] = v
			}
			cleanRule["uuids"] = make(map[string]interface{})

			rulesToCreate = append(rulesToCreate, cleanRule)
		}

		// Assign rule numbers to new rules that don't have them
		rulesToCreateWithNumbers := assignRuleNumbers(rulesToCreate)

		var createdRules []interface{}
		err := createNetworkACLRules(d, meta, &createdRules, rulesToCreateWithNumbers)
		if err != nil {
			return fmt.Errorf("failed to create new rules during migration: %v", err)
		}

		log.Printf("[DEBUG] Successfully created %d new rules during migration", len(createdRules))

		if err := d.Set("rule", createdRules); err != nil {
			return fmt.Errorf("failed to update state with migrated rules: %v", err)
		}
		log.Printf("[DEBUG] Updated Terraform state with %d migrated rules", len(createdRules))
	}

	log.Printf("[DEBUG] Ports->port migration completed successfully")
	return nil
}
