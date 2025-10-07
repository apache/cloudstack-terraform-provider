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
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/go-multierror"
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

						"port": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"ports": {
							Type:       schema.TypeSet,
							Optional:   true,
							Elem:       &schema.Schema{Type: schema.TypeString},
							Set:        schema.HashString,
							Deprecated: "Use 'port' instead. 'ports' will be removed in a future version.",
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
	// Make sure all required parameters are there
	if err := verifyNetworkACLParams(d); err != nil {
		return err
	}

	// We need to set this upfront in order to be able to save a partial state
	d.SetId(d.Get("acl_id").(string))

	// Create all rules that are configured
	if nrs := d.Get("rule").([]interface{}); len(nrs) > 0 {
		rules := make([]interface{}, 0, len(nrs))

		err := createNetworkACLRules(d, meta, &rules, nrs)

		// We need to update this first to preserve the correct state
		d.Set("rule", rules)

		if err != nil {
			return err
		}
	}

	return resourceCloudStackNetworkACLRuleRead(d, meta)
}

func createNetworkACLRules(d *schema.ResourceData, meta interface{}, rules *[]interface{}, nrs []interface{}) error {
	var errs *multierror.Error

	var wg sync.WaitGroup
	wg.Add(len(nrs))

	sem := make(chan struct{}, d.Get("parallelism").(int))
	for _, rule := range nrs {
		// Put in a tiny sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(rule map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Create a single rule
			err := createNetworkACLRule(d, meta, rule)

			// If we have at least one UUID, we need to save the rule
			if len(rule["uuids"].(map[string]interface{})) > 0 {
				*rules = append(*rules, rule)
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

func createNetworkACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	uuids := rule["uuids"].(map[string]interface{})

	// Make sure all required parameters are there
	if err := verifyNetworkACLRuleParams(d, rule); err != nil {
		return err
	}

	// Create a new parameter struct
	p := cs.NetworkACL.NewCreateNetworkACLParams(rule["protocol"].(string))

	// If a rule ID is specified, set it
	if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
		p.SetNumber(ruleNum)
	}

	// Set the acl ID
	p.SetAclid(d.Id())

	// Set the action
	p.SetAction(rule["action"].(string))

	// Set the CIDR list
	var cidrList []string
	for _, cidr := range rule["cidr_list"].(*schema.Set).List() {
		cidrList = append(cidrList, cidr.(string))
	}
	p.SetCidrlist(cidrList)

	// Set the traffic type
	p.SetTraffictype(rule["traffic_type"].(string))

	// Set the description
	if desc, ok := rule["description"].(string); ok && desc != "" {
		p.SetReason(desc)
	}

	// If the protocol is ICMP set the needed ICMP parameters
	if rule["protocol"].(string) == "icmp" {
		p.SetIcmptype(rule["icmp_type"].(int))
		p.SetIcmpcode(rule["icmp_code"].(int))

		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			return err
		}

		uuids["icmp"] = r.(*cloudstack.CreateNetworkACLResponse).Id
		rule["uuids"] = uuids
	}

	// If the protocol is ALL set the needed parameters
	if rule["protocol"].(string) == "all" {
		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			return err
		}

		uuids["all"] = r.(*cloudstack.CreateNetworkACLResponse).Id
		rule["uuids"] = uuids
	}

	var portStr string
	if port, ok := rule["port"].(string); ok && port != "" {
		portStr = port
		if ports, ok := rule["ports"].(*schema.Set); ok && ports.Len() > 0 {
			log.Printf("[WARN] Deprecated 'ports' is ignored. Only 'port' is used. Remove 'ports' from your config.")
		}
	} else if ports, ok := rule["ports"].(*schema.Set); ok && ports.Len() > 0 {
		// Deprecated: use first port or join as range if two values
		if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
			return fmt.Errorf("'ports' cannot be used with 'rule_number'. Please migrate to 'port' (string) for numbered rules.")
		}
		list := ports.List()
		if len(list) == 1 {
			portStr = list[0].(string)
		} else if len(list) == 2 {
			start := list[0].(string)
			end := list[1].(string)
			if strings.Contains(start, "-") || strings.Contains(end, "-") {
				return fmt.Errorf("If specifying a port range, use a single string like '1000-2000' in 'port'. Do not mix ranges and single ports.")
			}
			portStr = fmt.Sprintf("%s-%s", start, end)
		} else {
			return fmt.Errorf("'ports' must have one or two values only. Got: %v", list)
		}
		log.Printf("[WARN] 'ports' is deprecated. Use 'port' instead.")
	} else {
		if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
			return fmt.Errorf("Parameter port is required for protocol 'tcp' or 'udp'. Use 'port' (string) for new configs.")
		}
	}
	if portStr != "" {
		m := splitPorts.FindStringSubmatch(portStr)
		if m == nil {
			return fmt.Errorf("%q is not a valid port value. Valid options are '80' or '80-90'", portStr)
		}
		startPort, _ := strconv.Atoi(m[1])
		endPort := startPort
		if m[2] != "" {
			endPort, _ = strconv.Atoi(m[2])
		}
		p.SetStartport(startPort)
		p.SetEndport(endPort)

		r, err := Retry(4, retryableACLCreationFunc(cs, p))
		if err != nil {
			return err
		}

		uuids[portStr] = r.(*cloudstack.CreateNetworkACLResponse).Id
		rule["uuids"] = uuids
	}

	return nil
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
			d.SetId("")
			return nil
		}
		return err
	}

	p := cs.NetworkACL.NewListNetworkACLsParams()
	p.SetAclid(d.Id())
	p.SetListall(true)

	l, err := cs.NetworkACL.ListNetworkACLs(p)
	if err != nil {
		return err
	}

	ruleMap := make(map[string]*cloudstack.NetworkACL, l.Count)
	for _, r := range l.NetworkACLs {
		ruleMap[r.Id] = r
	}

	rules := make([]interface{}, 0, len(ruleMap))

	if rs, ok := d.Get("rule").([]interface{}); ok && len(rs) > 0 {
		for _, rule := range rs {
			rule := rule.(map[string]interface{})
			matched := false
			// 1. Try to match by UUID
			if uuids, ok := rule["uuids"].(map[string]interface{}); ok {
				for uuid := range uuids {
					if r, ok := ruleMap[uuid]; ok {
						cidrs := &schema.Set{F: schema.HashString}
						for _, cidr := range strings.Split(r.Cidrlist, ",") {
							cidrs.Add(cidr)
						}
						rule["action"] = strings.ToLower(r.Action)
						rule["protocol"] = r.Protocol
						rule["traffic_type"] = strings.ToLower(r.Traffictype)
						rule["cidr_list"] = cidrs
						rule["rule_number"] = int(r.Number)
						if desc, ok := rule["description"].(string); ok && desc != "" {
							if desc == r.Reason {
								rule["description"] = r.Reason
							}
						} else if r.Reason != "" {
							rule["description"] = r.Reason
						} else {
							rule["description"] = ""
						}
						if r.Protocol == "tcp" || r.Protocol == "udp" {
							if r.Startport == r.Endport {
								rule["port"] = r.Startport
							} else {
								rule["port"] = r.Startport + "-" + r.Endport
							}
						} else {
							rule["port"] = ""
						}
						rule["uuids"] = map[string]interface{}{r.Id: r.Id}
						matched = true
						delete(ruleMap, r.Id)
						break
					}
				}
			}
			// 2. If not found by UUID, match by all identity fields (rule_number, protocol, port, cidr, traffic_type)
			if !matched {
				for _, r := range ruleMap {
					if int(r.Number) == rule["rule_number"].(int) &&
						strings.ToLower(r.Protocol) == strings.ToLower(rule["protocol"].(string)) &&
						strings.ToLower(r.Traffictype) == strings.ToLower(rule["traffic_type"].(string)) &&
						matchACLRuleByNumberAndFields(r, rule) {
						cidrs := &schema.Set{F: schema.HashString}
						for _, cidr := range strings.Split(r.Cidrlist, ",") {
							cidrs.Add(cidr)
						}
						rule["action"] = strings.ToLower(r.Action)
						rule["protocol"] = r.Protocol
						rule["traffic_type"] = strings.ToLower(r.Traffictype)
						rule["cidr_list"] = cidrs
						rule["rule_number"] = int(r.Number)
						if desc, ok := rule["description"].(string); ok && desc != "" {
							if desc == r.Reason {
								rule["description"] = r.Reason
							}
						} else if r.Reason != "" {
							rule["description"] = r.Reason
						} else {
							rule["description"] = ""
						}
						if r.Protocol == "tcp" || r.Protocol == "udp" {
							if r.Startport == r.Endport {
								rule["port"] = r.Startport
							} else {
								rule["port"] = r.Startport + "-" + r.Endport
							}
						} else {
							rule["port"] = ""
						}
						rule["uuids"] = map[string]interface{}{r.Id: r.Id}
						matched = true
						delete(ruleMap, r.Id)
						break
					}
				}
			}
			// 3. If not found, do NOT update rule_number or other identity fields; just keep config values
			rules = append(rules, rule)
		}
	}

	managed := d.Get("managed").(bool)
	if managed && len(ruleMap) > 0 {
		for uuid := range ruleMap {
			cidrs := &schema.Set{F: schema.HashString}
			cidrs.Add(uuid)
			rule := map[string]interface{}{
				"cidr_list": cidrs,
				"protocol":  uuid,
				"uuids":     map[string]interface{}{uuid: uuid},
			}
			rules = append(rules, rule)
		}
	}

	if len(rules) > 0 {
		d.Set("rule", rules)
	} else if !managed {
		d.SetId("")
	}

	return nil
}

// Matches a CloudStack rule to a Terraform rule by rule_number, protocol, cidr, traffic_type, and port
func matchACLRuleByNumberAndFields(r *cloudstack.NetworkACL, rule map[string]interface{}) bool {
	if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
		if int(r.Number) != ruleNum {
			return false
		}
	}
	if strings.ToLower(r.Protocol) != strings.ToLower(rule["protocol"].(string)) {
		return false
	}
	if strings.ToLower(r.Traffictype) != strings.ToLower(rule["traffic_type"].(string)) {
		return false
	}
	cidrSet := map[string]struct{}{}
	for _, c := range strings.Split(r.Cidrlist, ",") {
		cidrSet[strings.TrimSpace(c)] = struct{}{}
	}
	for _, c := range rule["cidr_list"].(*schema.Set).List() {
		if _, ok := cidrSet[c.(string)]; !ok {
			return false
		}
	}
	portStr := ""
	if p, ok := rule["port"].(string); ok {
		portStr = p
	}
	startPort, _ := strconv.Atoi(r.Startport)
	endPort, _ := strconv.Atoi(r.Endport)
	if portStr != "" {
		if strings.Contains(portStr, "-") {
			parts := strings.SplitN(portStr, "-", 2)
			sp, _ := strconv.Atoi(parts[0])
			ep, _ := strconv.Atoi(parts[1])
			if sp != startPort || ep != endPort {
				return false
			}
		} else {
			sp, _ := strconv.Atoi(portStr)
			if sp != startPort || sp != endPort {
				return false
			}
		}
	}
	return true
}

func resourceCloudStackNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	// Make sure all required parameters are there
	if err := verifyNetworkACLParams(d); err != nil {
		return err
	}

	// Check if the rule set as a whole has changed
	if d.HasChange("rule") {
		o, n := d.GetChange("rule")
		oldRules := o.([]interface{}) // remote (from CloudStack)
		newRules := n.([]interface{}) // config (planned)

		// Build UUID -> rule maps for old (remote) and new (config) rules using CloudStack rule IDs as keys
		oldUUIDMap := map[string]map[string]interface{}{}
		newUUIDMap := map[string]map[string]interface{}{}
		for _, rule := range oldRules {
			r := rule.(map[string]interface{})
			for _, id := range r["uuids"].(map[string]interface{}) {
				oldUUIDMap[id.(string)] = r
			}
		}
		for _, rule := range newRules {
			r := rule.(map[string]interface{})
			for _, id := range r["uuids"].(map[string]interface{}) {
				newUUIDMap[id.(string)] = r
			}
		}

		// For each CloudStack rule ID present in both old and new, update if any relevant field changed (compare config to remote)
		updateFields := []string{"action", "cidr_list", "icmp_code", "icmp_type", "protocol", "description", "port", "traffic_type", "rule_number"} // add rule_number
		updated := map[string]bool{}
		for uuid, oldRule := range oldUUIDMap { // oldRule = remote, newRule = config
			if newRule, ok := newUUIDMap[uuid]; ok && !updated[uuid] {
				for _, field := range updateFields {
					oldVal := normalizeField(oldRule[field])
					newVal := normalizeField(newRule[field])
					if oldVal != newVal {
						log.Printf("[DEBUG] Updating ACL rule for UUID %s: field %s changed (remote=%v, config=%v)", uuid, field, oldRule[field], newRule[field])
						if err := updateNetworkACLRule(d, meta, newRule); err != nil {
							return err
						}
						break
					}
				}
				updated[uuid] = true
			}
		}
		// Remove delete+create for rule_number change
		// Create only truly new rules (no UUID)
		for _, rule := range newRules {
			r := rule.(map[string]interface{})
			if len(r["uuids"].(map[string]interface{})) == 0 {
				if err := createNetworkACLRule(d, meta, r); err != nil {
					return err
				}
			}
		}
		// Delete only truly removed rules
		// Find rules in oldRules not present in newRules by UUID
		toDelete := make([]map[string]interface{}, 0)
		for uuid, oldRule := range oldUUIDMap {
			if _, ok := newUUIDMap[uuid]; !ok {
				toDelete = append(toDelete, oldRule)
			}
		}
		if len(toDelete) > 0 {
			for _, rule := range toDelete {
				if err := deleteNetworkACLRule(d, meta, rule); err != nil {
					return err
				}
			}
		}
		// Update state
		d.Set("rule", newRules)
	}

	return resourceCloudStackNetworkACLRuleRead(d, meta)
}

// updateNetworkACLRule updates a single ACL rule in CloudStack
func updateNetworkACLRule(d *schema.ResourceData, meta interface{}, rule map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	uuids := rule["uuids"].(map[string]interface{})

	for _, id := range uuids {
		p := cs.NetworkACL.NewUpdateNetworkACLItemParams(id.(string))
		p.SetAction(rule["action"].(string))
		p.SetCidrlist(expandStringSet(rule["cidr_list"].(*schema.Set)))
		p.SetProtocol(rule["protocol"].(string))
		p.SetTraffictype(rule["traffic_type"].(string))
		if desc, ok := rule["description"].(string); ok && desc != "" {
			p.SetReason(desc)
		}
		if rule["protocol"].(string) == "icmp" {
			p.SetIcmptype(rule["icmp_type"].(int))
			p.SetIcmpcode(rule["icmp_code"].(int))
		}
		if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
			if port, ok := rule["port"].(string); ok && port != "" {
				m := splitPorts.FindStringSubmatch(port)
				startPort, _ := strconv.Atoi(m[1])
				endPort := startPort
				if m[2] != "" {
					endPort, _ = strconv.Atoi(m[2])
				}
				p.SetStartport(startPort)
				p.SetEndport(endPort)
			}
		}
		if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
			p.SetNumber(ruleNum)
		}
		_, err := cs.NetworkACL.UpdateNetworkACLItem(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceCloudStackNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// Create an empty rule slice to hold all rules that were not deleted correctly
	rules := make([]interface{}, 0)

	// Delete all rules
	if ors, ok := d.Get("rule").([]interface{}); ok && len(ors) > 0 {
		for _, rule := range ors {
			if err := deleteNetworkACLRule(d, meta, rule.(map[string]interface{})); err != nil {
				// If we have at least one UUID, we need to save the rule
				if len(rule.(map[string]interface{})["uuids"].(map[string]interface{})) > 0 {
					rules = append(rules, rule)
				}
				return err
			}
		}
		// We need to update this first to preserve the correct state
		d.Set("rule", rules)
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
				err = multierror.Append(errs, err)
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
	// Disallow 'ports' for anything except deletes (backward compatibility)
	if ports, ok := rule["ports"].(*schema.Set); ok && ports.Len() > 0 {
		// Only allow deletes (rule is being removed)
		if d != nil && d.HasChange("rule") {
			o, n := d.GetChange("rule")
			ors := o.(*schema.Set).Difference(n.(*schema.Set))
			isDelete := false
			for _, r := range ors.List() {
				rm := r.(map[string]interface{})
				if rm["ports"].(*schema.Set).Len() > 0 {
					isDelete = true
					break
				}
			}
			if !isDelete {
				return fmt.Errorf("The 'ports' attribute is deprecated and not allowed for new or updated rules. Please migrate to 'port' (string). Only deletion of existing rules with 'ports' is allowed.")
			}
		} else {
			return fmt.Errorf("The 'ports' attribute is deprecated and not allowed for new or updated rules. Please migrate to 'port' (string). Only deletion of existing rules with 'ports' is allowed.")
		}
	}
	// Disallow 'ports' with 'rule_number'
	if ports, ok := rule["ports"].(*schema.Set); ok && ports.Len() > 0 {
		if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
			return fmt.Errorf("'ports' cannot be used with 'rule_number'. Please migrate to 'port' (schema.TypeString) for numbered rules.")
		}
	}

	action := rule["action"].(string)
	if action != "allow" && action != "deny" {
		return fmt.Errorf("Parameter action only accepts 'allow' or 'deny' as values")
	}

	protocol := rule["protocol"].(string)
	switch protocol {
	case "icmp":
		if _, ok := rule["icmp_type"]; !ok {
			return fmt.Errorf(
				"Parameter icmp_type is a required parameter when using protocol 'icmp'")
		}
		if _, ok := rule["icmp_code"]; !ok {
			return fmt.Errorf(
				"Parameter icmp_code is a required parameter when using protocol 'icmp'")
		}
	case "all":
		// No additional test are needed, so just leave this empty...
	case "tcp", "udp":
		// Error if both ports and rule_number are set (must be first check)
		if ports, ok := rule["ports"].(*schema.Set); ok && ports.Len() > 0 {
			if ruleNum, ok := rule["rule_number"].(int); ok && ruleNum > 0 {
				return fmt.Errorf("'ports' cannot be used with 'rule_number'. Please migrate to 'port' (schema.TypeString) for numbered rules.")
			}
		}
		var portStr string
		if port, ok := rule["port"].(string); ok && port != "" {
			portStr = port
		} else if ports, ok := rule["ports"].(*schema.Set); ok && ports.Len() > 0 {
			list := ports.List()
			if len(list) == 1 {
				portStr = list[0].(string)
			} else if len(list) == 2 {
				start := list[0].(string)
				end := list[1].(string)
				if strings.Contains(start, "-") || strings.Contains(end, "-") {
					return fmt.Errorf("If specifying a port range, use a single string like '1000-2000' in 'port'. Do not mix ranges and single ports.")
				}
				portStr = fmt.Sprintf("%s-%s", start, end)
			} else {
				return fmt.Errorf("'ports' must have one or two values only. Got: %v", list)
			}
			log.Printf("[WARN] 'ports' is deprecated. Use 'port' instead.")
		}
		if portStr != "" {
			m := splitPorts.FindStringSubmatch(portStr)
			if m == nil {
				return fmt.Errorf("%q is not a valid port value. Valid options are '80' or '80-90'", portStr)
			}
		} else {
			return fmt.Errorf("Parameter port is a required parameter when *not* using protocol 'icmp'")
		}
	default:
		_, err := strconv.ParseInt(protocol, 0, 0)
		if err != nil {
			return fmt.Errorf(
				"%q is not a valid protocol. Valid options are 'tcp', 'udp', "+
					"'icmp', 'all' or a valid protocol number", protocol)
		}
	}

	traffic := rule["traffic_type"].(string)
	if traffic != "ingress" && traffic != "egress" {
		return fmt.Errorf(
			"Parameter traffic_type only accepts 'ingress' or 'egress' as values")
	}

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

// expandStringSet converts a *schema.Set to a []string
func expandStringSet(set *schema.Set) []string {
	var out []string
	for _, v := range set.List() {
		out = append(out, v.(string))
	}
	return out
}

// normalizeField returns a comparable value for a field (handles nil, empty, etc)
func normalizeField(v interface{}) interface{} {
	switch val := v.(type) {
	case nil:
		return ""
	case *schema.Set:
		list := val.List()
		strs := make([]string, len(list))
		for i, s := range list {
			strs[i] = fmt.Sprintf("%v", s)
		}
		return strings.Join(strs, ",")
	default:
		return val
	}
}
