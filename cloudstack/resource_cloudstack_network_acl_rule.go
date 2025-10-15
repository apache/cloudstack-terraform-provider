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

				log.Printf("[DEBUG] CustomizeDiff: checking %d old rules -> %d new rules for migration", len(oldRulesList), len(newRulesList))

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

func resourceCloudStackNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Entering resourceCloudStackNetworkACLRuleCreate with acl_id=%s", d.Get("acl_id").(string))

	// Make sure all required parameters are there
	if err := verifyNetworkACLParams(d); err != nil {
		log.Printf("[ERROR] Failed parameter verification: %v", err)
		return err
	}

	// Create all rules that are configured
	if nrs := d.Get("rule").([]interface{}); len(nrs) > 0 {
		// Create an empty rule list to hold all newly created rules
		rules := make([]interface{}, 0)

		log.Printf("[DEBUG] Processing %d rules", len(nrs))
		err := createNetworkACLRules(d, meta, &rules, nrs)
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
				errs = multierror.Append(errs, fmt.Errorf("rule #%d: %v", index+1, err))
			} else if len(rule["uuids"].(map[string]interface{})) > 0 {
				log.Printf("[DEBUG] Successfully created rule #%d, adding to rules list", index+1)
				*rules = append(*rules, rule)
			} else {
				log.Printf("[WARN] Rule #%d created but has no UUIDs", index+1)
			}

			<-sem
		}(rule.(map[string]interface{}), i)
	}

	wg.Wait()

	if err := errs.ErrorOrNil(); err != nil {
		log.Printf("[ERROR] Errors occurred while creating rules: %v", err)
		return err
	}

	log.Printf("[DEBUG] Successfully created all rules")
	return nil
}

func createNetworkACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	uuids := rule["uuids"].(map[string]interface{})
	log.Printf("[DEBUG] Creating network ACL rule with protocol=%s", rule["protocol"].(string))

	// Make sure all required parameters are there
	if err := verifyNetworkACLRuleParams(d, rule); err != nil {
		log.Printf("[ERROR] Failed to verify rule parameters: %v", err)
		return err
	}

	// Create a new parameter struct
	p := cs.NetworkACL.NewCreateNetworkACLParams(rule["protocol"].(string))
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
	p.SetAction(rule["action"].(string))
	log.Printf("[DEBUG] Set action=%s", rule["action"].(string))

	// Set the CIDR list
	var cidrList []string
	for _, cidr := range rule["cidr_list"].([]interface{}) {
		cidrList = append(cidrList, cidr.(string))
	}
	p.SetCidrlist(cidrList)
	log.Printf("[DEBUG] Set cidr_list=%v", cidrList)

	// Set the traffic type
	p.SetTraffictype(rule["traffic_type"].(string))
	log.Printf("[DEBUG] Set traffic_type=%s", rule["traffic_type"].(string))

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
		uuids["icmp"] = r.(*cloudstack.CreateNetworkACLResponse).Id
		rule["uuids"] = uuids
		log.Printf("[DEBUG] Created ICMP rule with ID=%s", r.(*cloudstack.CreateNetworkACLResponse).Id)
	}

	// If the protocol is ALL set the needed parameters
	if rule["protocol"].(string) == "all" {
		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			log.Printf("[ERROR] Failed to create ALL rule: %v", err)
			return err
		}
		uuids["all"] = r.(*cloudstack.CreateNetworkACLResponse).Id
		rule["uuids"] = uuids
		log.Printf("[DEBUG] Created ALL rule with ID=%s", r.(*cloudstack.CreateNetworkACLResponse).Id)
	}

	// If protocol is TCP or UDP, create the rule (with or without port)
	if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
		// Check if deprecated ports field is used and reject it
		if portsSet, hasPortsSet := rule["ports"].(*schema.Set); hasPortsSet && portsSet.Len() > 0 {
			log.Printf("[ERROR] Attempt to create rule with deprecated ports field")
			return fmt.Errorf("The 'ports' field is no longer supported for creating new rules. Please use the 'port' field with separate rules for each port/range.")
		}

		portStr, hasPort := rule["port"].(string)

		if hasPort && portStr != "" {
			// Handle single port
			log.Printf("[DEBUG] Processing single port for TCP/UDP rule: %s", portStr)

			if _, ok := uuids[portStr]; !ok {
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

				uuids[portStr] = r.(*cloudstack.CreateNetworkACLResponse).Id
				rule["uuids"] = uuids
				log.Printf("[DEBUG] Created TCP/UDP rule for port %s with ID=%s", portStr, r.(*cloudstack.CreateNetworkACLResponse).Id)
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
			uuids["all_ports"] = r.(*cloudstack.CreateNetworkACLResponse).Id
			rule["uuids"] = uuids
			log.Printf("[DEBUG] Created TCP/UDP rule for all ports with ID=%s", r.(*cloudstack.CreateNetworkACLResponse).Id)
		}
	}

	log.Printf("[DEBUG] Successfully created rule with uuids=%+v", uuids)
	return nil
}

func processTCPUDPRule(rule map[string]interface{}, ruleMap map[string]*cloudstack.NetworkACL, uuids map[string]interface{}, rules *[]interface{}) {
	// Check for deprecated ports field first (for reading existing state during migration)
	ps, hasPortsSet := rule["ports"].(*schema.Set)
	portStr, hasPort := rule["port"].(string)

	if hasPortsSet && ps.Len() > 0 {
		log.Printf("[DEBUG] Processing deprecated ports field with %d ports during state read", ps.Len())

		// Process each port in the deprecated ports set during state read
		for _, port := range ps.List() {
			portStr := port.(string)

			if processPortForRule(portStr, rule, ruleMap, uuids) {
				log.Printf("[DEBUG] Processed deprecated port %s during state read", portStr)
			}
		}

		// Only add the rule once with all processed ports
		if len(uuids) > 0 {
			*rules = append(*rules, rule)
			log.Printf("[DEBUG] Added TCP/UDP rule with deprecated ports to state during read: %+v", rule)
		}

	} else if hasPort && portStr != "" {
		log.Printf("[DEBUG] Processing single port for TCP/UDP rule: %s", portStr)

		if processPortForRule(portStr, rule, ruleMap, uuids) {
			rule["port"] = portStr
			*rules = append(*rules, rule)
			log.Printf("[DEBUG] Added TCP/UDP rule with single port to state: %+v", rule)
		}

	} else {
		log.Printf("[DEBUG] Processing TCP/UDP rule with no port specified")

		id, ok := uuids["all_ports"]
		if !ok {
			log.Printf("[DEBUG] No UUID for all_ports, skipping rule")
			return
		}

		r, ok := ruleMap[id.(string)]
		if !ok {
			log.Printf("[DEBUG] TCP/UDP rule for all_ports with ID %s not found, removing UUID", id.(string))
			delete(uuids, "all_ports")
			return
		}

		delete(ruleMap, id.(string))

		var cidrs []interface{}
		for _, cidr := range strings.Split(r.Cidrlist, ",") {
			cidrs = append(cidrs, cidr)
		}

		rule["action"] = strings.ToLower(r.Action)
		rule["protocol"] = r.Protocol
		rule["traffic_type"] = strings.ToLower(r.Traffictype)
		rule["cidr_list"] = cidrs
		rule["rule_number"] = r.Number
		*rules = append(*rules, rule)
		log.Printf("[DEBUG] Added TCP/UDP rule with no port to state: %+v", rule)
	}
}

func processPortForRule(portStr string, rule map[string]interface{}, ruleMap map[string]*cloudstack.NetworkACL, uuids map[string]interface{}) bool {
	id, ok := uuids[portStr]
	if !ok {
		log.Printf("[DEBUG] No UUID for port %s, skipping", portStr)
		return false
	}

	r, ok := ruleMap[id.(string)]
	if !ok {
		log.Printf("[DEBUG] TCP/UDP rule for port %s with ID %s not found, removing UUID", portStr, id.(string))
		delete(uuids, portStr)
		return false
	}

	// Delete the known rule so only unknown rules remain in the ruleMap
	delete(ruleMap, id.(string))

	var cidrs []interface{}
	for _, cidr := range strings.Split(r.Cidrlist, ",") {
		cidrs = append(cidrs, cidr)
	}

	rule["action"] = strings.ToLower(r.Action)
	rule["protocol"] = r.Protocol
	rule["traffic_type"] = strings.ToLower(r.Traffictype)
	rule["cidr_list"] = cidrs
	rule["rule_number"] = r.Number

	return true
}

func resourceCloudStackNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	log.Printf("[DEBUG] Entering resourceCloudStackNetworkACLRuleRead with acl_id=%s", d.Id())

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

	// Read all rules that are configured
	if rs := d.Get("rule").([]interface{}); len(rs) > 0 {
		for _, rule := range rs {
			rule := rule.(map[string]interface{})
			uuids := rule["uuids"].(map[string]interface{})
			log.Printf("[DEBUG] Processing rule with protocol=%s, uuids=%+v", rule["protocol"].(string), uuids)

			if rule["protocol"].(string) == "icmp" {
				id, ok := uuids["icmp"]
				if !ok {
					log.Printf("[DEBUG] No ICMP UUID found, skipping rule")
					continue
				}

				// Get the rule
				r, ok := ruleMap[id.(string)]
				if !ok {
					log.Printf("[DEBUG] ICMP rule with ID %s not found, removing UUID", id.(string))
					delete(uuids, "icmp")
					continue
				}

				// Delete the known rule so only unknown rules remain in the ruleMap
				delete(ruleMap, id.(string))

				// Create a list with all CIDR's
				var cidrs []interface{}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs = append(cidrs, cidr)
				}

				// Update the values
				rule["action"] = strings.ToLower(r.Action)
				rule["protocol"] = r.Protocol
				rule["icmp_type"] = r.Icmptype
				rule["icmp_code"] = r.Icmpcode
				rule["traffic_type"] = strings.ToLower(r.Traffictype)
				rule["cidr_list"] = cidrs
				rule["rule_number"] = r.Number
				rules = append(rules, rule)
				log.Printf("[DEBUG] Added ICMP rule to state: %+v", rule)
			}

			if rule["protocol"].(string) == "all" {
				id, ok := uuids["all"]
				if !ok {
					log.Printf("[DEBUG] No ALL UUID found, skipping rule")
					continue
				}

				// Get the rule
				r, ok := ruleMap[id.(string)]
				if !ok {
					log.Printf("[DEBUG] ALL rule with ID %s not found, removing UUID", id.(string))
					delete(uuids, "all")
					continue
				}

				// Delete the known rule so only unknown rules remain in the ruleMap
				delete(ruleMap, id.(string))

				// Create a list with all CIDR's
				var cidrs []interface{}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs = append(cidrs, cidr)
				}

				// Update the values
				rule["action"] = strings.ToLower(r.Action)
				rule["protocol"] = r.Protocol
				rule["traffic_type"] = strings.ToLower(r.Traffictype)
				rule["cidr_list"] = cidrs
				rule["rule_number"] = r.Number
				rules = append(rules, rule)
				log.Printf("[DEBUG] Added ALL rule to state: %+v", rule)
			}

			if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
				uuids := rule["uuids"].(map[string]interface{})
				processTCPUDPRule(rule, ruleMap, uuids, &rules)
			}
		}
	}

	// If this is a managed firewall, add all unknown rules into dummy rules
	managed := d.Get("managed").(bool)
	if managed && len(ruleMap) > 0 {
		for uuid := range ruleMap {
			// We need to create and add a dummy value to a list as the
			// cidr_list is a required field and thus needs a value
			cidrs := []interface{}{uuid}

			// Make a dummy rule to hold the unknown UUID
			rule := map[string]interface{}{
				"cidr_list": cidrs,
				"protocol":  uuid,
				"uuids":     map[string]interface{}{uuid: uuid},
			}

			// Add the dummy rule to the rules list
			rules = append(rules, rule)
			log.Printf("[DEBUG] Added managed dummy rule for UUID %s", uuid)
		}
	}

	if len(rules) > 0 {
		log.Printf("[DEBUG] Setting %d rules in state", len(rules))
		if err := d.Set("rule", rules); err != nil {
			log.Printf("[ERROR] Failed to set rule attribute: %v", err)
			return err
		}
	} else if !managed {
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
		err := updateNetworkACLRules(d, meta, oldRules, newRules)
		if err != nil {
			return err
		}
	}

	return resourceCloudStackNetworkACLRuleRead(d, meta)
}

func resourceCloudStackNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// Delete all rules
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

	return nil
}

func deleteNetworkACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	uuids := rule["uuids"].(map[string]interface{})

	for k, id := range uuids {
		// We don't care about the count here, so just continue
		if k == "%" {
			continue
		}

		// Create the parameter struct
		p := cs.NetworkACL.NewDeleteNetworkACLParams(id.(string))

		// Delete the rule
		if _, err := cs.NetworkACL.DeleteNetworkACL(p); err != nil {

			// This is a very poor way to be told the ID does no longer exist :(
			if strings.Contains(err.Error(), fmt.Sprintf(
				"Invalid parameter id value=%s due to incorrect long value format, "+
					"or entity does not exist", id.(string))) {
				delete(uuids, k)
				rule["uuids"] = uuids
				continue
			}

			return err
		}

		// Delete the UUID of this rule
		delete(uuids, k)
		rule["uuids"] = uuids
	}

	return nil
}

func verifyNetworkACLParams(d *schema.ResourceData) error {
	managed := d.Get("managed").(bool)
	_, rules := d.GetOk("rule")

	if !rules && !managed {
		return fmt.Errorf(
			"You must supply at least one 'rule' when not using the 'managed' firewall feature")
	}

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
		// The deprecated 'ports' field is no longer supported in any scenario
		portsSet, hasPortsSet := rule["ports"].(*schema.Set)
		portStr, hasPort := rule["port"].(string)

		// Block deprecated ports field completely
		if hasPortsSet && portsSet.Len() > 0 {
			log.Printf("[ERROR] Attempt to use deprecated ports field")
			return fmt.Errorf("The 'ports' field is no longer supported. Please use the 'port' field instead.")
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
		} else {
			log.Printf("[DEBUG] No port specified for TCP/UDP, allowing empty port")
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
			if rulesMatch(oldRuleMap, newRuleMap) {
				log.Printf("[DEBUG] Found matching new rule for old rule")

				if oldUUIDs, ok := oldRuleMap["uuids"].(map[string]interface{}); ok {
					newRuleMap["uuids"] = oldUUIDs
				}

				if ruleNeedsUpdate(oldRuleMap, newRuleMap) {
					log.Printf("[DEBUG] Rule needs updating")
					if uuids, ok := oldRuleMap["uuids"].(map[string]interface{}); ok {
						for _, uuid := range uuids {
							if uuid != nil {
								rulesToUpdate[uuid.(string)] = newRuleMap
								break
							}
						}
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

func rulesMatch(oldRule, newRule map[string]interface{}) bool {
	if oldRule["protocol"].(string) != newRule["protocol"].(string) ||
		oldRule["traffic_type"].(string) != newRule["traffic_type"].(string) ||
		oldRule["action"].(string) != newRule["action"].(string) {
		return false
	}

	protocol := newRule["protocol"].(string)

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
	if oldRule["action"].(string) != newRule["action"].(string) {
		log.Printf("[DEBUG] Action changed: %s -> %s", oldRule["action"].(string), newRule["action"].(string))
		return true
	}

	if oldRule["protocol"].(string) != newRule["protocol"].(string) {
		log.Printf("[DEBUG] Protocol changed: %s -> %s", oldRule["protocol"].(string), newRule["protocol"].(string))
		return true
	}

	if oldRule["traffic_type"].(string) != newRule["traffic_type"].(string) {
		log.Printf("[DEBUG] Traffic type changed: %s -> %s", oldRule["traffic_type"].(string), newRule["traffic_type"].(string))
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

	protocol := newRule["protocol"].(string)
	switch protocol {
	case "icmp":
		if oldRule["icmp_type"].(int) != newRule["icmp_type"].(int) {
			log.Printf("[DEBUG] ICMP type changed: %d -> %d", oldRule["icmp_type"].(int), newRule["icmp_type"].(int))
			return true
		}
		if oldRule["icmp_code"].(int) != newRule["icmp_code"].(int) {
			log.Printf("[DEBUG] ICMP code changed: %d -> %d", oldRule["icmp_code"].(int), newRule["icmp_code"].(int))
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

	oldCidrs := oldRule["cidr_list"].([]interface{})
	newCidrs := newRule["cidr_list"].([]interface{})
	if len(oldCidrs) != len(newCidrs) {
		log.Printf("[DEBUG] CIDR list length changed: %d -> %d", len(oldCidrs), len(newCidrs))
		return true
	}

	oldCidrStrs := make([]string, len(oldCidrs))
	newCidrStrs := make([]string, len(newCidrs))
	for i, cidr := range oldCidrs {
		oldCidrStrs[i] = cidr.(string)
	}
	for i, cidr := range newCidrs {
		newCidrStrs[i] = cidr.(string)
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
		for _, cidr := range newRule["cidr_list"].([]interface{}) {
			cidrList = append(cidrList, cidr.(string))
		}
		p.SetCidrlist(cidrList)

		if desc, ok := newRule["description"].(string); ok && desc != "" {
			p.SetReason(desc)
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

		var createdRules []interface{}
		err := createNetworkACLRules(d, meta, &createdRules, rulesToCreate)
		if err != nil {
			return fmt.Errorf("failed to create new rules during migration: %v", err)
		}

		log.Printf("[DEBUG] Successfully created %d new rules during migration", len(createdRules))
	}

	log.Printf("[DEBUG] Ports->port migration completed successfully")
	return nil
}
