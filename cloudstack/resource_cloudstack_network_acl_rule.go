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
				Type:     schema.TypeSet,
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
							Computed: true,
						},

						"icmp_code": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"ports": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
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
	if nrs := d.Get("rule").(*schema.Set); nrs.Len() > 0 {
		// Create an empty rule set to hold all newly created rules
		rules := resourceCloudStackNetworkACLRule().Schema["rule"].ZeroValue().(*schema.Set)

		log.Printf("[DEBUG] Processing %d rules", nrs.Len())
		err := createNetworkACLRules(d, meta, rules, nrs)
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

func createNetworkACLRules(d *schema.ResourceData, meta interface{}, rules *schema.Set, nrs *schema.Set) error {
	log.Printf("[DEBUG] Creating %d network ACL rules", nrs.Len())
	var errs *multierror.Error

	var wg sync.WaitGroup
	wg.Add(nrs.Len())

	sem := make(chan struct{}, d.Get("parallelism").(int))
	for i, rule := range nrs.List() {
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
				log.Printf("[DEBUG] Successfully created rule #%d, adding to rules set", index+1)
				rules.Add(rule)
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
	for _, cidr := range rule["cidr_list"].(*schema.Set).List() {
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

	// If protocol is TCP or UDP, create the rule (with or without ports)
	if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
		ps, ok := rule["ports"].(*schema.Set)
		if !ok || ps == nil {
			log.Printf("[DEBUG] No ports specified for TCP/UDP rule, creating rule for all ports")
			ps = &schema.Set{F: schema.HashString}
		}

		// Create an empty schema.Set to hold all processed ports
		ports := &schema.Set{F: schema.HashString}
		log.Printf("[DEBUG] Processing %d ports for TCP/UDP rule", ps.Len())

		if ps.Len() == 0 {
			// Create a rule for all ports
			r, err := Retry(4, retryableACLCreationFunc(cs, p))
			if err != nil {
				log.Printf("[ERROR] Failed to create TCP/UDP rule for all ports: %v", err)
				return err
			}
			uuids["all_ports"] = r.(*cloudstack.CreateNetworkACLResponse).Id
			rule["uuids"] = uuids
			log.Printf("[DEBUG] Created TCP/UDP rule for all ports with ID=%s", r.(*cloudstack.CreateNetworkACLResponse).Id)
		} else {
			// Process specified ports
			for _, port := range ps.List() {
				if _, ok := uuids[port.(string)]; ok {
					ports.Add(port)
					rule["ports"] = ports
					log.Printf("[DEBUG] Port %s already has UUID, skipping", port.(string))
					continue
				}

				m := splitPorts.FindStringSubmatch(port.(string))
				if m == nil {
					log.Printf("[ERROR] Invalid port format: %s", port.(string))
					return fmt.Errorf("%q is not a valid port value. Valid options are '80' or '80-90'", port.(string))
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
				log.Printf("[DEBUG] Set ports start=%d, end=%d", startPort, endPort)

				r, err := Retry(4, retryableACLCreationFunc(cs, p))
				if err != nil {
					log.Printf("[ERROR] Failed to create TCP/UDP rule for port %s: %v", port.(string), err)
					return err
				}

				ports.Add(port)
				rule["ports"] = ports
				uuids[port.(string)] = r.(*cloudstack.CreateNetworkACLResponse).Id
				rule["uuids"] = uuids
				log.Printf("[DEBUG] Created TCP/UDP rule for port %s with ID=%s", port.(string), r.(*cloudstack.CreateNetworkACLResponse).Id)
			}
		}
	}

	log.Printf("[DEBUG] Successfully created rule with uuids=%+v", uuids)
	return nil
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

	// Create an empty schema.Set to hold all rules
	rules := resourceCloudStackNetworkACLRule().Schema["rule"].ZeroValue().(*schema.Set)

	// Read all rules that are configured
	if rs := d.Get("rule").(*schema.Set); rs.Len() > 0 {
		for _, rule := range rs.List() {
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

				// Create a set with all CIDR's
				cidrs := &schema.Set{F: schema.HashString}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs.Add(cidr)
				}

				// Update the values
				rule["action"] = strings.ToLower(r.Action)
				rule["protocol"] = r.Protocol
				rule["icmp_type"] = r.Icmptype
				rule["icmp_code"] = r.Icmpcode
				rule["traffic_type"] = strings.ToLower(r.Traffictype)
				rule["cidr_list"] = cidrs
				rules.Add(rule)
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

				// Create a set with all CIDR's
				cidrs := &schema.Set{F: schema.HashString}
				for _, cidr := range strings.Split(r.Cidrlist, ",") {
					cidrs.Add(cidr)
				}

				// Update the values
				rule["action"] = strings.ToLower(r.Action)
				rule["protocol"] = r.Protocol
				rule["traffic_type"] = strings.ToLower(r.Traffictype)
				rule["cidr_list"] = cidrs
				rules.Add(rule)
				log.Printf("[DEBUG] Added ALL rule to state: %+v", rule)
			}

			if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
				ps, ok := rule["ports"].(*schema.Set)
				if !ok || ps == nil {
					log.Printf("[DEBUG] No ports specified for TCP/UDP rule, initializing empty set")
					ps = &schema.Set{F: schema.HashString}
				}

				// Create an empty schema.Set to hold all ports
				ports := &schema.Set{F: schema.HashString}
				log.Printf("[DEBUG] Processing %d ports for TCP/UDP rule", ps.Len())

				// Loop through all ports and retrieve their info
				for _, port := range ps.List() {
					id, ok := uuids[port.(string)]
					if !ok {
						log.Printf("[DEBUG] No UUID for port %s, skipping", port.(string))
						continue
					}

					// Get the rule
					r, ok := ruleMap[id.(string)]
					if !ok {
						log.Printf("[DEBUG] TCP/UDP rule for port %s with ID %s not found, removing UUID", port.(string), id.(string))
						delete(uuids, port.(string))
						continue
					}

					// Delete the known rule so only unknown rules remain in the ruleMap
					delete(ruleMap, id.(string))

					// Create a set with all CIDR's
					cidrs := &schema.Set{F: schema.HashString}
					for _, cidr := range strings.Split(r.Cidrlist, ",") {
						cidrs.Add(cidr)
					}

					// Update the values
					rule["action"] = strings.ToLower(r.Action)
					rule["protocol"] = r.Protocol
					rule["traffic_type"] = strings.ToLower(r.Traffictype)
					rule["cidr_list"] = cidrs
					ports.Add(port)
					log.Printf("[DEBUG] Added port %s to TCP/UDP rule", port.(string))
				}

				// If there is at least one port found, add this rule to the rules set
				if ports.Len() > 0 {
					rule["ports"] = ports
					rules.Add(rule)
					log.Printf("[DEBUG] Added TCP/UDP rule to state: %+v", rule)
				} else {
					// Add the rule even if no ports are specified, as ports are optional
					rule["ports"] = ports
					rules.Add(rule)
					log.Printf("[DEBUG] Added TCP/UDP rule with no ports to state: %+v", rule)
				}
			}
		}
	}

	// If this is a managed firewall, add all unknown rules into dummy rules
	managed := d.Get("managed").(bool)
	if managed && len(ruleMap) > 0 {
		for uuid := range ruleMap {
			// We need to create and add a dummy value to a schema.Set as the
			// cidr_list is a required field and thus needs a value
			cidrs := &schema.Set{F: schema.HashString}
			cidrs.Add(uuid)

			// Make a dummy rule to hold the unknown UUID
			rule := map[string]interface{}{
				"cidr_list": cidrs,
				"protocol":  uuid,
				"uuids":     map[string]interface{}{uuid: uuid},
			}

			// Add the dummy rule to the rules set
			rules.Add(rule)
			log.Printf("[DEBUG] Added managed dummy rule for UUID %s", uuid)
		}
	}

	if rules.Len() > 0 {
		log.Printf("[DEBUG] Setting %d rules in state", rules.Len())
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

	// Check if the rule set as a whole has changed
	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// We need to start with a rule set containing all the rules we
		// already have and want to keep. Any rules that are not deleted
		// correctly and any newly created rules, will be added to this
		// set to make sure we end up in a consistent state
		rules := o.(*schema.Set).Intersection(n.(*schema.Set))

		// First loop through all the new rules and create (before destroy) them
		if nrs.Len() > 0 {
			err := createNetworkACLRules(d, meta, rules, nrs)

			// We need to update this first to preserve the correct state
			d.Set("rule", rules)

			if err != nil {
				return err
			}
		}

		// Then loop through all the old rules and delete them
		if ors.Len() > 0 {
			err := deleteNetworkACLRules(d, meta, rules, ors)

			// We need to update this first to preserve the correct state
			d.Set("rule", rules)

			if err != nil {
				return err
			}
		}
	}

	return resourceCloudStackNetworkACLRuleRead(d, meta)
}

func resourceCloudStackNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// Create an empty rule set to hold all rules that where
	// not deleted correctly
	rules := resourceCloudStackNetworkACLRule().Schema["rule"].ZeroValue().(*schema.Set)

	// Delete all rules
	if ors := d.Get("rule").(*schema.Set); ors.Len() > 0 {
		err := deleteNetworkACLRules(d, meta, rules, ors)

		// We need to update this first to preserve the correct state
		d.Set("rule", rules)

		if err != nil {
			return err
		}
	}

	return nil
}

func deleteNetworkACLRules(d *schema.ResourceData, meta interface{}, rules *schema.Set, ors *schema.Set) error {
	var errs *multierror.Error

	var wg sync.WaitGroup
	wg.Add(ors.Len())

	sem := make(chan struct{}, d.Get("parallelism").(int))
	for _, rule := range ors.List() {
		// Put a sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(rule map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Delete a single rule
			err := deleteNetworkACLRule(d, meta, rule)

			// If we have at least one UUID, we need to save the rule
			if len(rule["uuids"].(map[string]interface{})) > 0 {
				rules.Add(rule)
			}

			if err != nil {
				errs = multierror.Append(errs, err)
			}

			<-sem
		}(rule.(map[string]interface{}))
	}

	wg.Wait()

	return errs.ErrorOrNil()
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
		if ports, ok := rule["ports"].(*schema.Set); ok {
			log.Printf("[DEBUG] Found %d ports for TCP/UDP", ports.Len())
			for _, port := range ports.List() {
				m := splitPorts.FindStringSubmatch(port.(string))
				if m == nil {
					log.Printf("[ERROR] Invalid port format: %s", port.(string))
					return fmt.Errorf(
						"%q is not a valid port value. Valid options are '80' or '80-90'", port.(string))
				}
			}
		} else {
			log.Printf("[DEBUG] No ports specified for TCP/UDP, assuming empty set")
			// Allow empty ports for TCP/UDP (your config has no ports)
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
