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
				rules = append(rules, rule)
				log.Printf("[DEBUG] Added ALL rule to state: %+v", rule)
			}

			if rule["protocol"].(string) == "tcp" || rule["protocol"].(string) == "udp" {
				// Check for deprecated ports field first (for backward compatibility)
				ps, hasPortsSet := rule["ports"].(*schema.Set)
				portStr, hasPort := rule["port"].(string)

				if hasPortsSet && ps.Len() > 0 {
					// Handle deprecated ports field (multiple ports)
					log.Printf("[DEBUG] Processing %d ports for TCP/UDP rule (deprecated field)", ps.Len())

					// Create an empty list to hold all ports
					var ports []interface{}

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
						ports = append(ports, port)
						log.Printf("[DEBUG] Added port %s to TCP/UDP rule", port.(string))
					}

					// Add this rule to the rules list with ports
					rule["ports"] = schema.NewSet(schema.HashString, ports)
					rules = append(rules, rule)
					log.Printf("[DEBUG] Added TCP/UDP rule with deprecated ports to state: %+v", rule)

				} else if hasPort && portStr != "" {
					log.Printf("[DEBUG] Processing single port for TCP/UDP rule: %s", portStr)

					id, ok := uuids[portStr]
					if !ok {
						log.Printf("[DEBUG] No UUID for port %s, skipping rule", portStr)
						continue
					}

					r, ok := ruleMap[id.(string)]
					if !ok {
						log.Printf("[DEBUG] TCP/UDP rule for port %s with ID %s not found, removing UUID", portStr, id.(string))
						delete(uuids, portStr)
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
					rule["port"] = portStr
					rules = append(rules, rule)
					log.Printf("[DEBUG] Added TCP/UDP rule with single port to state: %+v", rule)

				} else {
					log.Printf("[DEBUG] Processing TCP/UDP rule with no port specified")

					id, ok := uuids["all_ports"]
					if !ok {
						log.Printf("[DEBUG] No UUID for all_ports, skipping rule")
						continue
					}

					r, ok := ruleMap[id.(string)]
					if !ok {
						log.Printf("[DEBUG] TCP/UDP rule for all_ports with ID %s not found, removing UUID", id.(string))
						delete(uuids, "all_ports")
						continue
					}

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
					rules = append(rules, rule)
					log.Printf("[DEBUG] Added TCP/UDP rule with no port to state: %+v", rule)
				}
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
		// Check if deprecated ports field is used (not allowed for new configurations)
		portsSet, hasPortsSet := rule["ports"].(*schema.Set)
		portStr, hasPort := rule["port"].(string)

		if hasPortsSet && portsSet.Len() > 0 {
			log.Printf("[ERROR] Deprecated ports field used in new configuration")
			return fmt.Errorf("The 'ports' field is deprecated. Use 'port' instead for new configurations.")
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

	oldRuleMap := make(map[string]map[string]interface{})
	newRuleMap := make(map[string]map[string]interface{})

	for _, rule := range oldRules {
		ruleMap := rule.(map[string]interface{})
		key := createRuleKey(ruleMap)
		oldRuleMap[key] = ruleMap
		log.Printf("[DEBUG] Old rule key: %s", key)
	}

	for _, rule := range newRules {
		ruleMap := rule.(map[string]interface{})
		key := createRuleKey(ruleMap)
		newRuleMap[key] = ruleMap
		log.Printf("[DEBUG] New rule key: %s", key)
	}

	for key, oldRule := range oldRuleMap {
		if _, exists := newRuleMap[key]; !exists {
			log.Printf("[DEBUG] Deleting rule: %s", key)
			err := deleteNetworkACLRule(d, meta, oldRule)
			if err != nil {
				return fmt.Errorf("failed to delete rule %s: %v", key, err)
			}
		}
	}

	var rulesToCreate []interface{}
	for key, newRule := range newRuleMap {
		if _, exists := oldRuleMap[key]; !exists {
			log.Printf("[DEBUG] Creating new rule: %s", key)
			rulesToCreate = append(rulesToCreate, newRule)
		}
	}

	if len(rulesToCreate) > 0 {
		var createdRules []interface{}
		err := createNetworkACLRules(d, meta, &createdRules, rulesToCreate)
		if err != nil {
			return fmt.Errorf("failed to create new rules: %v", err)
		}
	}

	for key, newRule := range newRuleMap {
		if oldRule, exists := oldRuleMap[key]; exists {
			if ruleNeedsUpdate(oldRule, newRule) {
				log.Printf("[DEBUG] Updating rule: %s", key)
				err := updateNetworkACLRule(cs, oldRule, newRule)
				if err != nil {
					return fmt.Errorf("failed to update rule %s: %v", key, err)
				}
			}
		}
	}

	return nil
}

func createRuleKey(rule map[string]interface{}) string {
	protocol := rule["protocol"].(string)
	trafficType := rule["traffic_type"].(string)

	if protocol == "icmp" {
		icmpType := rule["icmp_type"].(int)
		icmpCode := rule["icmp_code"].(int)
		return fmt.Sprintf("%s-%s-icmp-%d-%d", protocol, trafficType, icmpType, icmpCode)
	}

	if protocol == "all" {
		return fmt.Sprintf("%s-%s-all", protocol, trafficType)
	}

	if protocol == "tcp" || protocol == "udp" {
		portStr, hasPort := rule["port"].(string)
		if hasPort && portStr != "" {
			return fmt.Sprintf("%s-%s-port-%s", protocol, trafficType, portStr)
		} else {
			return fmt.Sprintf("%s-%s-noport", protocol, trafficType)
		}
	}

	// For numeric protocols
	return fmt.Sprintf("%s-%s", protocol, trafficType)
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
