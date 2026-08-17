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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackNetworkACLRuleset() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNetworkACLRulesetCreate,
		Read:   resourceCloudStackNetworkACLRulesetRead,
		Update: resourceCloudStackNetworkACLRulesetUpdate,
		Delete: resourceCloudStackNetworkACLRulesetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudStackNetworkACLRulesetImport,
		},
		// CustomizeDiff is used to eliminate spurious diffs when modifying individual rules.
		// Without this, changing a single rule (e.g., port 80->8080) would show ALL rules
		// as being replaced in the plan because TypeSet uses hashing and any field change
		// changes the hash. This function matches rules by their natural key (rule_number)
		// and uses SetNew to suppress diffs for unchanged rules.
		CustomizeDiff: resourceCloudStackNetworkACLRulesetCustomizeDiff,

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

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed: true is required to allow CustomizeDiff to use SetNew().
				// Without this, we get "Error: SetNew only operates on computed keys".
				// This enables CustomizeDiff to suppress spurious diffs by preserving
				// the old state for unchanged rules.
				Computed: true,
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
							Default:  -1,
						},

						"icmp_code": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  -1,
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

						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudStackNetworkACLRulesetCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// Only apply this logic during updates, not creates
	if d.Id() == "" {
		return nil
	}

	old, new := d.GetChange("rule")
	oldSet := old.(*schema.Set)
	newSet := new.(*schema.Set)

	// If the sets are empty, nothing to do
	if oldSet.Len() == 0 || newSet.Len() == 0 {
		return nil
	}

	// Create maps indexed by rule_number (the natural key)
	oldMap := make(map[int]map[string]interface{})
	for _, v := range oldSet.List() {
		m := v.(map[string]interface{})
		ruleNum := m["rule_number"].(int)
		oldMap[ruleNum] = m
	}

	newMap := make(map[int]map[string]interface{})
	for _, v := range newSet.List() {
		m := v.(map[string]interface{})
		ruleNum := m["rule_number"].(int)
		newMap[ruleNum] = m
	}

	// Build a new set that preserves UUIDs for unchanged rules
	// and uses new values for changed/added rules
	preservedSet := schema.NewSet(newSet.F, []interface{}{})

	for ruleNum, newRule := range newMap {
		oldRule, exists := oldMap[ruleNum]

		if exists && compareACLRules(oldRule, newRule) {
			// Rule exists and is functionally identical - preserve the old rule
			// (including its UUID) to prevent spurious diff
			preservedSet.Add(oldRule)
		} else {
			// Rule is new or changed - use the new rule
			preservedSet.Add(newRule)
		}
	}

	// Set the preserved set as the new value
	// This maintains UUIDs for unchanged rules while allowing changes to show correctly
	return d.SetNew("rule", preservedSet)
}

func compareACLRules(old, new map[string]interface{}) bool {
	// Compare all fields except uuid (which is computed and may differ)
	fields := []string{"rule_number", "action", "protocol", "icmp_type", "icmp_code", "port", "traffic_type", "description"}

	for _, field := range fields {
		oldVal := old[field]
		newVal := new[field]

		// Handle nil/empty string equivalence
		if oldVal == nil && newVal == "" {
			continue
		}
		if oldVal == "" && newVal == nil {
			continue
		}

		if oldVal != newVal {
			return false
		}
	}

	// Compare cidr_list (TypeSet)
	oldCIDR := old["cidr_list"].(*schema.Set)
	newCIDR := new["cidr_list"].(*schema.Set)

	if !oldCIDR.Equal(newCIDR) {
		return false
	}

	return true
}

func resourceCloudStackNetworkACLRulesetCreate(d *schema.ResourceData, meta interface{}) error {
	// We need to set this upfront in order to be able to save a partial state
	d.SetId(d.Get("acl_id").(string))

	// Create all rules that are configured
	if nrs := d.Get("rule").(*schema.Set); nrs.Len() > 0 {
		// Create an empty schema.Set to hold all rules
		rules := resourceCloudStackNetworkACLRuleset().Schema["rule"].ZeroValue().(*schema.Set)

		err := createACLRules(d, meta, rules, nrs)
		if err != nil {
			return err
		}

		// We need to update this first to preserve the correct state
		d.Set("rule", rules)
	}

	return resourceCloudStackNetworkACLRulesetRead(d, meta)
}

func createACLRules(d *schema.ResourceData, meta interface{}, rules *schema.Set, nrs *schema.Set) error {
	var errs *multierror.Error
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(nrs.Len())

	sem := make(chan struct{}, 10)
	for _, rule := range nrs.List() {
		// Put in a tiny sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(rule map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Create a single rule
			err := createACLRule(d, meta, rule)

			// If we have a UUID, we need to save the rule
			if uuid, ok := rule["uuid"].(string); ok && uuid != "" {
				mu.Lock()
				rules.Add(rule)
				mu.Unlock()
			}

			if err != nil {
				mu.Lock()
				errs = multierror.Append(errs, err)
				mu.Unlock()
			}

			<-sem
		}(rule.(map[string]interface{}))
	}

	wg.Wait()

	return errs.ErrorOrNil()
}

func createACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Make sure all required parameters are there
	if err := verifyACLRuleParams(d, rule); err != nil {
		return err
	}

	protocol := rule["protocol"].(string)
	action := rule["action"].(string)
	trafficType := rule["traffic_type"].(string)

	// Create a new parameter struct
	p := cs.NetworkACL.NewCreateNetworkACLParams(protocol)

	// Set the rule number
	p.SetNumber(rule["rule_number"].(int))

	// Set the acl ID
	p.SetAclid(d.Get("acl_id").(string))

	// Set the action
	p.SetAction(action)

	// Set the CIDR list
	var cidrList []string
	for _, cidr := range rule["cidr_list"].(*schema.Set).List() {
		cidrList = append(cidrList, cidr.(string))
	}
	p.SetCidrlist(cidrList)

	// Set the traffic type
	p.SetTraffictype(trafficType)

	// Set the description
	if desc, ok := rule["description"].(string); ok && desc != "" {
		p.SetReason(desc)
	}

	// If the protocol is ICMP set the needed ICMP parameters
	if protocol == "icmp" {
		// icmp_type and icmp_code default to -1 (all) in the schema
		icmpType := rule["icmp_type"].(int)
		icmpCode := rule["icmp_code"].(int)

		p.SetIcmptype(icmpType)
		p.SetIcmpcode(icmpCode)

		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			return err
		}

		rule["uuid"] = r.(*cloudstack.CreateNetworkACLResponse).Id
		return nil
	}

	// If the protocol is ALL set the needed parameters
	if protocol == "all" {
		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			return err
		}

		rule["uuid"] = r.(*cloudstack.CreateNetworkACLResponse).Id
		return nil
	}

	// If protocol is TCP or UDP, create the rule (with or without port)
	if protocol == "tcp" || protocol == "udp" {
		portStr, hasPort := rule["port"].(string)

		if hasPort && portStr != "" {
			// Handle single port
			m := splitPorts.FindStringSubmatch(portStr)
			if m == nil {
				return fmt.Errorf("%q is not a valid port value. Valid options are '80' or '80-90'", portStr)
			}

			startPort, err := strconv.Atoi(m[1])
			if err != nil {
				return err
			}

			endPort := startPort
			if m[2] != "" {
				endPort, err = strconv.Atoi(m[2])
				if err != nil {
					return err
				}
			}

			p.SetStartport(startPort)
			p.SetEndport(endPort)
		}

		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			return err
		}

		rule["uuid"] = r.(*cloudstack.CreateNetworkACLResponse).Id
		return nil
	}

	// If we reach here, it's an unsupported protocol (should have been caught by validation)
	return fmt.Errorf("unsupported protocol %q. Valid protocols are: tcp, udp, icmp, all", protocol)
}

// buildRuleFromAPI converts a CloudStack NetworkACL API response to a rule map
func buildRuleFromAPI(r *cloudstack.NetworkACL) map[string]interface{} {
	cidrs := &schema.Set{F: schema.HashString}
	for _, cidr := range strings.Split(r.Cidrlist, ",") {
		cidrs.Add(cidr)
	}

	rule := map[string]interface{}{
		"cidr_list":    cidrs,
		"action":       strings.ToLower(r.Action),
		"protocol":     r.Protocol,
		"traffic_type": strings.ToLower(r.Traffictype),
		"rule_number":  r.Number,
		"description":  r.Reason,
		"uuid":         r.Id,
	}

	// Set ICMP fields
	if r.Protocol == "icmp" {
		rule["icmp_type"] = r.Icmptype
		rule["icmp_code"] = r.Icmpcode
	} else {
		// For non-ICMP protocols, set to -1 (matches schema default)
		rule["icmp_type"] = -1
		rule["icmp_code"] = -1
	}

	// Set port if applicable
	if r.Protocol == "tcp" || r.Protocol == "udp" {
		if r.Startport != "" && r.Endport != "" {
			if r.Startport == r.Endport {
				rule["port"] = r.Startport
			} else {
				rule["port"] = fmt.Sprintf("%s-%s", r.Startport, r.Endport)
			}
		} else {
			// Explicitly clear port when no ports are set (all ports)
			rule["port"] = ""
		}
	} else {
		// Explicitly clear port when protocol is not tcp/udp
		rule["port"] = ""
	}

	return rule
}

func resourceCloudStackNetworkACLRulesetRead(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}

	// Get all the rules from the running environment
	p := cs.NetworkACL.NewListNetworkACLsParams()
	p.SetAclid(d.Id())
	p.SetListall(true)

	if err := setProjectid(p, cs, d); err != nil {
		return err
	}

	l, err := cs.NetworkACL.ListNetworkACLs(p)
	if err != nil {
		return err
	}

	// Make a map of all the rules so we can easily find a rule
	ruleMap := make(map[string]*cloudstack.NetworkACL, l.Count)
	for _, r := range l.NetworkACLs {
		ruleMap[r.Id] = r
	}

	// Create an empty schema.Set to hold all rules
	rules := resourceCloudStackNetworkACLRuleset().Schema["rule"].ZeroValue().(*schema.Set)

	// Read all rules that are configured
	rs := d.Get("rule").(*schema.Set)
	if rs.Len() > 0 {
		for _, oldRule := range rs.List() {
			oldRule := oldRule.(map[string]interface{})

			id, ok := oldRule["uuid"]
			if !ok || id.(string) == "" {
				continue
			}

			// Get the rule
			r, ok := ruleMap[id.(string)]
			if !ok {
				// Rule no longer exists in the API, skip it
				continue
			}

			// Delete the known rule so only unknown rules remain in the ruleMap
			delete(ruleMap, id.(string))

			// Create a NEW map with the updated values (don't mutate the old one)
			rule := buildRuleFromAPI(r)
			rules.Add(rule)
		}
	} else {
		// If no rules in state (e.g., during import), read all remote rules
		for _, r := range ruleMap {
			rule := buildRuleFromAPI(r)
			rules.Add(rule)
			// Remove from ruleMap so we don't add it again as a dummy rule
			delete(ruleMap, r.Id)
		}
	}

	// If this is a managed resource, add all unknown rules to dummy rules
	// This allows Terraform to detect them and trigger an update to delete them
	managed := d.Get("managed").(bool)
	if managed && len(ruleMap) > 0 {
		log.Printf("[DEBUG] Found %d out-of-band ACL rules for ACL %s", len(ruleMap), d.Id())
		for uuid, r := range ruleMap {
			log.Printf("[DEBUG] Adding dummy rule for out-of-band rule: uuid=%s, rule_number=%d", uuid, r.Number)

			// Build the rule from the API response to preserve actual values
			// This ensures the diff shows the real cidr_list and protocol values
			// instead of UUIDs, making it clear what's being deleted
			rule := buildRuleFromAPI(r)

			// Add the dummy rule to the rules set
			rules.Add(rule)
		}
	}

	if rules.Len() > 0 {
		d.Set("rule", rules)
	} else if !managed {
		d.SetId("")
	}

	return nil
}

func resourceCloudStackNetworkACLRulesetUpdate(d *schema.ResourceData, meta interface{}) error {
	// Check if the rule set as a whole has changed
	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		// Build maps of rules by rule_number for efficient lookup
		oldRulesByNumber := make(map[int]map[string]interface{})
		newRulesByNumber := make(map[int]map[string]interface{})

		for _, rule := range oldSet.List() {
			ruleMap := rule.(map[string]interface{})
			ruleNum := ruleMap["rule_number"].(int)
			oldRulesByNumber[ruleNum] = ruleMap
		}

		for _, rule := range newSet.List() {
			ruleMap := rule.(map[string]interface{})
			ruleNum := ruleMap["rule_number"].(int)
			newRulesByNumber[ruleNum] = ruleMap
		}

		// Categorize rules into: update, delete, create
		var rulesToUpdate []*ruleUpdatePair
		var rulesToDelete []map[string]interface{}
		var rulesToCreate []map[string]interface{}

		// Find rules to update or delete
		for ruleNum, oldRule := range oldRulesByNumber {
			if newRule, exists := newRulesByNumber[ruleNum]; exists {
				// Rule exists in both old and new - check if it needs updating
				if aclRuleNeedsUpdate(oldRule, newRule) {
					rulesToUpdate = append(rulesToUpdate, &ruleUpdatePair{
						oldRule: oldRule,
						newRule: newRule,
					})
				}
				// If no update needed, the rule stays as-is (UUID preserved)
			} else {
				// Rule only exists in old state - delete it
				rulesToDelete = append(rulesToDelete, oldRule)
			}
		}

		// Find rules to create
		for ruleNum, newRule := range newRulesByNumber {
			if _, exists := oldRulesByNumber[ruleNum]; !exists {
				// Rule only exists in new state - create it
				rulesToCreate = append(rulesToCreate, newRule)
			}
		}

		// We need to start with a rule set containing all the rules we
		// already have and want to keep. Any rules that are not deleted
		// correctly and any newly created rules, will be added to this
		// set to make sure we end up in a consistent state
		rules := resourceCloudStackNetworkACLRuleset().Schema["rule"].ZeroValue().(*schema.Set)

		// Add all rules that will remain (either unchanged or updated)
		for ruleNum := range newRulesByNumber {
			if oldRule, exists := oldRulesByNumber[ruleNum]; exists {
				// This rule will either be updated or kept as-is
				// Start with the old rule (which has the UUID)
				rules.Add(oldRule)
			}
		}

		// First, delete rules that are no longer needed
		if len(rulesToDelete) > 0 {
			deleteSet := &schema.Set{F: rules.F}
			for _, rule := range rulesToDelete {
				deleteSet.Add(rule)
			}
			err := deleteACLRules(d, meta, rules, deleteSet)

			// We need to update this first to preserve the correct state
			d.Set("rule", rules)

			if err != nil {
				return err
			}
		}

		// Second, update rules that have changed
		if len(rulesToUpdate) > 0 {
			err := updateACLRules(d, meta, rules, rulesToUpdate)

			// We need to update this first to preserve the correct state
			d.Set("rule", rules)

			if err != nil {
				return err
			}
		}

		// Finally, create new rules
		if len(rulesToCreate) > 0 {
			createSet := &schema.Set{F: rules.F}
			for _, rule := range rulesToCreate {
				createSet.Add(rule)
			}
			err := createACLRules(d, meta, rules, createSet)

			// We need to update this first to preserve the correct state
			d.Set("rule", rules)

			if err != nil {
				return err
			}
		}
	}

	return resourceCloudStackNetworkACLRulesetRead(d, meta)
}

type ruleUpdatePair struct {
	oldRule map[string]interface{}
	newRule map[string]interface{}
}

func resourceCloudStackNetworkACLRulesetDelete(d *schema.ResourceData, meta interface{}) error {
	// If managed=false, don't delete any rules - just remove from state
	managed := d.Get("managed").(bool)
	if !managed {
		log.Printf("[DEBUG] Managed=false, not deleting ACL rules for %s", d.Id())
		return nil
	}

	// Create an empty rule set to hold all rules that where
	// not deleted correctly
	rules := resourceCloudStackNetworkACLRuleset().Schema["rule"].ZeroValue().(*schema.Set)

	// Delete all rules
	if ors := d.Get("rule").(*schema.Set); ors.Len() > 0 {
		err := deleteACLRules(d, meta, rules, ors)

		// We need to update this first to preserve the correct state
		d.Set("rule", rules)

		if err != nil {
			return err
		}
	}

	return nil
}

func deleteACLRules(d *schema.ResourceData, meta interface{}, rules *schema.Set, ors *schema.Set) error {
	var errs *multierror.Error
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(ors.Len())

	sem := make(chan struct{}, 10)
	for _, rule := range ors.List() {
		// Put a sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(rule map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Delete a single rule
			err := deleteACLRule(d, meta, rule)

			// If we have a UUID, we need to save the rule
			if uuid, ok := rule["uuid"].(string); ok && uuid != "" {
				mu.Lock()
				rules.Add(rule)
				mu.Unlock()
			}

			if err != nil {
				mu.Lock()
				errs = multierror.Append(errs, err)
				mu.Unlock()
			}

			<-sem
		}(rule.(map[string]interface{}))
	}

	wg.Wait()

	return errs.ErrorOrNil()
}

func deleteACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create the parameter struct
	p := cs.NetworkACL.NewDeleteNetworkACLParams(rule["uuid"].(string))

	// Delete the rule
	if _, err := cs.NetworkACL.DeleteNetworkACL(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if !strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", rule["uuid"].(string))) {
			return err
		}
	}

	// Empty the UUID of this rule
	rule["uuid"] = ""

	return nil
}

func updateACLRules(d *schema.ResourceData, meta interface{}, rules *schema.Set, updatePairs []*ruleUpdatePair) error {
	var errs *multierror.Error
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(len(updatePairs))

	sem := make(chan struct{}, 10)
	for _, pair := range updatePairs {
		// Put a sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(pair *ruleUpdatePair) {
			defer wg.Done()
			sem <- struct{}{}

			// Update a single rule
			err := updateACLRule(d, meta, pair.oldRule, pair.newRule)

			// If we have a UUID, we need to save the updated rule
			if uuid, ok := pair.oldRule["uuid"].(string); ok && uuid != "" {
				mu.Lock()
				// Remove the old rule from the set
				rules.Remove(pair.oldRule)
				// Update the old rule with new values (preserving UUID)
				updateRuleValues(pair.oldRule, pair.newRule)
				// Add the updated rule back to the set
				rules.Add(pair.oldRule)
				mu.Unlock()
			}

			if err != nil {
				mu.Lock()
				errs = multierror.Append(errs, err)
				mu.Unlock()
			}

			<-sem
		}(pair)
	}

	wg.Wait()

	return errs.ErrorOrNil()
}

func updateACLRule(d *schema.ResourceData, meta interface{}, oldRule, newRule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	uuid := oldRule["uuid"].(string)
	if uuid == "" {
		return fmt.Errorf("cannot update rule without UUID")
	}

	log.Printf("[DEBUG] Updating ACL rule with UUID: %s", uuid)

	// If the protocol changed, we need to delete and recreate the rule
	// because the CloudStack API doesn't properly clear protocol-specific fields
	// (e.g., ports when changing from TCP to ICMP)
	if oldRule["protocol"].(string) != newRule["protocol"].(string) {
		log.Printf("[DEBUG] Protocol changed, using delete+create approach for rule %s", uuid)

		// Delete the old rule
		p := cs.NetworkACL.NewDeleteNetworkACLParams(uuid)
		if _, err := cs.NetworkACL.DeleteNetworkACL(p); err != nil {
			// Ignore "does not exist" errors
			if !strings.Contains(err.Error(), fmt.Sprintf(
				"Invalid parameter id value=%s due to incorrect long value format, "+
					"or entity does not exist", uuid)) {
				return fmt.Errorf("failed to delete rule during protocol change: %w", err)
			}
		}

		// Create the new rule with the new protocol
		if err := createACLRule(d, meta, newRule); err != nil {
			return fmt.Errorf("failed to create rule during protocol change: %w", err)
		}

		// The new UUID is now in newRule["uuid"], copy it to oldRule so it gets saved
		oldRule["uuid"] = newRule["uuid"]

		return nil
	}

	// Create the parameter struct
	p := cs.NetworkACL.NewUpdateNetworkACLItemParams(uuid)

	// Set the action
	p.SetAction(newRule["action"].(string))

	// Set the CIDR list
	var cidrList []string
	for _, cidr := range newRule["cidr_list"].(*schema.Set).List() {
		cidrList = append(cidrList, cidr.(string))
	}
	p.SetCidrlist(cidrList)

	// Set the description
	if desc, ok := newRule["description"].(string); ok && desc != "" {
		p.SetReason(desc)
	}

	// Set the protocol
	p.SetProtocol(newRule["protocol"].(string))

	// Set the traffic type
	p.SetTraffictype(newRule["traffic_type"].(string))

	// Set the rule number
	p.SetNumber(newRule["rule_number"].(int))

	protocol := newRule["protocol"].(string)
	switch protocol {
	case "icmp":
		// icmp_type and icmp_code default to -1 (all) in the schema
		icmpType := newRule["icmp_type"].(int)
		icmpCode := newRule["icmp_code"].(int)

		p.SetIcmptype(icmpType)
		p.SetIcmpcode(icmpCode)
		// Don't set ports for ICMP - CloudStack API will handle this
	case "all":
		// Don't set ports or ICMP fields for "all" protocol
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
				}
			}
		}
		// If port is empty, don't set start/end port - CloudStack will handle "all ports"
	}

	// Execute the update
	_, err := cs.NetworkACL.UpdateNetworkACLItem(p)
	if err != nil {
		log.Printf("[ERROR] Failed to update ACL rule %s: %v", uuid, err)
		return err
	}

	log.Printf("[DEBUG] Successfully updated ACL rule %s", uuid)
	return nil
}

func verifyACLRuleParams(d *schema.ResourceData, rule map[string]interface{}) error {
	ruleNumber := rule["rule_number"].(int)
	if ruleNumber < 1 || ruleNumber > 65535 {
		return fmt.Errorf("rule_number must be between 1 and 65535, got: %d", ruleNumber)
	}

	action := rule["action"].(string)
	if action != "allow" && action != "deny" {
		return fmt.Errorf("action must be 'allow' or 'deny', got: %s", action)
	}

	protocol := rule["protocol"].(string)
	switch protocol {
	case "icmp":
		// icmp_type and icmp_code are optional - they default to -1 (all) if not specified
	case "all":
		// No additional validation needed
	case "tcp", "udp":
		// Port is optional
		if portStr, ok := rule["port"].(string); ok && portStr != "" {
			m := splitPorts.FindStringSubmatch(portStr)
			if m == nil {
				return fmt.Errorf("%q is not a valid port value. Valid options are '80' or '80-90'", portStr)
			}
		}
	default:
		// Reject numeric protocols - CloudStack API expects protocol names
		if _, err := strconv.Atoi(protocol); err == nil {
			return fmt.Errorf("numeric protocols are not supported, use protocol names instead (tcp, udp, icmp, all). Got: %s", protocol)
		}
		// If not a number, it's an unsupported protocol name
		return fmt.Errorf("%q is not a valid protocol. Valid options are 'tcp', 'udp', 'icmp', 'all'", protocol)
	}

	return nil
}

func aclRuleNeedsUpdate(oldRule, newRule map[string]interface{}) bool {
	// Check basic attributes
	if oldRule["action"].(string) != newRule["action"].(string) {
		return true
	}

	if oldRule["protocol"].(string) != newRule["protocol"].(string) {
		return true
	}

	if oldRule["traffic_type"].(string) != newRule["traffic_type"].(string) {
		return true
	}

	// Check description
	oldDesc, _ := oldRule["description"].(string)
	newDesc, _ := newRule["description"].(string)
	if oldDesc != newDesc {
		return true
	}

	// Check CIDR list
	oldCidrs := oldRule["cidr_list"].(*schema.Set)
	newCidrs := newRule["cidr_list"].(*schema.Set)
	if !oldCidrs.Equal(newCidrs) {
		return true
	}

	// Check protocol-specific attributes
	protocol := newRule["protocol"].(string)
	switch protocol {
	case "icmp":
		if oldRule["icmp_type"].(int) != newRule["icmp_type"].(int) {
			return true
		}
		if oldRule["icmp_code"].(int) != newRule["icmp_code"].(int) {
			return true
		}
	case "tcp", "udp":
		oldPort, _ := oldRule["port"].(string)
		newPort, _ := newRule["port"].(string)
		if oldPort != newPort {
			return true
		}
	}

	return false
}

func updateRuleValues(oldRule, newRule map[string]interface{}) {
	// Update all values from newRule to oldRule, preserving the UUID
	oldRule["action"] = newRule["action"]
	oldRule["cidr_list"] = newRule["cidr_list"]
	oldRule["protocol"] = newRule["protocol"]
	oldRule["icmp_type"] = newRule["icmp_type"]
	oldRule["icmp_code"] = newRule["icmp_code"]
	oldRule["port"] = newRule["port"]
	oldRule["traffic_type"] = newRule["traffic_type"]
	oldRule["description"] = newRule["description"]
	oldRule["rule_number"] = newRule["rule_number"]
	// Note: UUID is NOT updated - it's preserved from oldRule
}

func resourceCloudStackNetworkACLRulesetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Parse the import ID to extract optional project name
	// Format: acl_id or project/acl_id
	s := strings.SplitN(d.Id(), "/", 2)
	if len(s) == 2 {
		d.Set("project", s[0])
		d.SetId(s[1])
	}

	// Set the acl_id field to match the resource ID
	d.Set("acl_id", d.Id())

	// Don't set managed here - let it use the default value from the schema (false)
	// The Read function will be called after this and will populate the rules

	log.Printf("[DEBUG] Imported ACL ruleset with ID: %s", d.Id())

	return []*schema.ResourceData{d}, nil
}
