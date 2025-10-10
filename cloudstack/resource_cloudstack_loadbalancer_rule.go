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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackLoadBalancerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackLoadBalancerRuleCreate,
		Read:   resourceCloudStackLoadBalancerRuleRead,
		Update: resourceCloudStackLoadBalancerRuleUpdate,
		Delete: resourceCloudStackLoadBalancerRuleDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_address_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"algorithm": {
				Type:     schema.TypeString,
				Required: true,
			},

			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"private_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"public_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"member_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"cidrlist": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackLoadBalancerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Make sure all required parameters are there
	if err := verifyLoadBalancerRule(d); err != nil {
		return err
	}

	// Create a new parameter struct
	p := cs.LoadBalancer.NewCreateLoadBalancerRuleParams(
		d.Get("algorithm").(string),
		d.Get("name").(string),
		d.Get("private_port").(int),
		d.Get("public_port").(int),
	)

	// Don't autocreate a firewall rule, use a resource if needed
	p.SetOpenfirewall(false)

	// Set the description
	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	} else {
		p.SetDescription(d.Get("name").(string))
	}

	if networkid, ok := d.GetOk("network_id"); ok {
		// Set the network id
		p.SetNetworkid(networkid.(string))
	}

	// Set the protocol
	if protocol, ok := d.GetOk("protocol"); ok {
		p.SetProtocol(protocol.(string))
	}

	// Set CIDR list
	if cidr, ok := d.GetOk("cidrlist"); ok {
		var cidrList []string
		for _, id := range cidr.(*schema.Set).List() {
			cidrList = append(cidrList, id.(string))
		}

		p.SetCidrlist(cidrList)
	}

	// Set the ipaddress id
	p.SetPublicipid(d.Get("ip_address_id").(string))

	// Create the load balancer rule
	r, err := cs.LoadBalancer.CreateLoadBalancerRule(p)
	if err != nil {
		return err
	}

	// Set the load balancer rule ID and set partials
	d.SetId(r.Id)

	if certificateID, ok := d.GetOk("certificate_id"); ok {
		// Create a new parameter struct
		cp := cs.LoadBalancer.NewAssignCertToLoadBalancerParams(certificateID.(string), r.Id)
		if _, err := cs.LoadBalancer.AssignCertToLoadBalancer(cp); err != nil {
			return err
		}
	}

	// Create a new parameter struct
	mp := cs.LoadBalancer.NewAssignToLoadBalancerRuleParams(r.Id)

	var mbs []string
	for _, id := range d.Get("member_ids").(*schema.Set).List() {
		mbs = append(mbs, id.(string))
	}

	mp.SetVirtualmachineids(mbs)

	_, err = cs.LoadBalancer.AssignToLoadBalancerRule(mp)
	if err != nil {
		return err
	}

	return resourceCloudStackLoadBalancerRuleRead(d, meta)
}

func resourceCloudStackLoadBalancerRuleRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the load balancer details
	lb, count, err := cs.LoadBalancer.GetLoadBalancerRuleByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project").(string)),
	)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Load balancer rule %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}

		return err
	}

	public_port, err := strconv.Atoi(lb.Publicport)
	if err != nil {
		return err
	}

	private_port, err := strconv.Atoi(lb.Privateport)
	if err != nil {
		return err
	}

	d.Set("name", lb.Name)
	d.Set("ip_address_id", lb.Publicipid)
	d.Set("algorithm", lb.Algorithm)
	d.Set("public_port", public_port)
	d.Set("private_port", private_port)
	d.Set("protocol", lb.Protocol)

	// Only set cidr if user specified it to avoid spurious diffs
	delimiters := regexp.MustCompile(`\s*,\s*|\s+`)
	if _, ok := d.GetOk("cidrlist"); ok {
		d.Set("cidrlist", delimiters.Split(lb.Cidrlist, -1))
	}

	// Only set network if user specified it to avoid spurious diffs
	if _, ok := d.GetOk("network_id"); ok {
		d.Set("network_id", lb.Networkid)
	}

	setValueOrID(d, "project", lb.Project, lb.Projectid)

	p := cs.LoadBalancer.NewListLoadBalancerRuleInstancesParams(d.Id())
	l, err := cs.LoadBalancer.ListLoadBalancerRuleInstances(p)
	if err != nil {
		return err
	}

	var mbs []string
	for _, i := range l.LoadBalancerRuleInstances {
		mbs = append(mbs, i.Id)
	}

	asgCheckParams := cs.AutoScale.NewListAutoScaleVmGroupsParams()
	asgCheckParams.SetLbruleid(d.Id())

	asgGroups, err := cs.AutoScale.ListAutoScaleVmGroups(asgCheckParams)
	if err != nil {
		log.Printf("[WARN] Failed to check for autoscale VM groups during read: %s", err)
	}

	if len(asgGroups.AutoScaleVmGroups) > 0 {
		log.Printf("[DEBUG] Load balancer rule %s is managed by %d autoscale VM group(s), current members: %v",
			d.Id(), len(asgGroups.AutoScaleVmGroups), mbs)

		if currentMemberIds, ok := d.GetOk("member_ids"); ok {
			currentSet := currentMemberIds.(*schema.Set)
			if currentSet.Len() == 0 && len(mbs) > 0 {
				d.Set("member_ids", []string{})
				return nil
			}
		}
	}

	d.Set("member_ids", mbs)

	return nil
}

func waitForASGsToBeDisabled(cs *cloudstack.CloudStackClient, lbRuleID string) error {
	log.Printf("[DEBUG] Waiting for autoscale VM groups using load balancer rule %s to be disabled", lbRuleID)

	maxRetries := 60 // 60 * 2 seconds = 120 seconds max wait (longer for Terraform-driven changes)
	for i := 0; i < maxRetries; i++ {
		listParams := cs.AutoScale.NewListAutoScaleVmGroupsParams()
		listParams.SetLbruleid(lbRuleID)

		groups, err := cs.AutoScale.ListAutoScaleVmGroups(listParams)
		if err != nil {
			log.Printf("[WARN] Failed to list autoscale VM groups: %s", err)
			time.Sleep(2 * time.Second)
			continue
		}

		allDisabled := true
		var enabledGroups []string

		for _, group := range groups.AutoScaleVmGroups {
			if group.State != "disabled" && group.State != "disable" {
				allDisabled = false
				enabledGroups = append(enabledGroups, fmt.Sprintf("%s(%s:%s)", group.Name, group.Id, group.State))
			}
		}

		if allDisabled {
			log.Printf("[INFO] All autoscale VM groups using load balancer rule %s are now disabled", lbRuleID)
			return nil
		}

		if i < maxRetries-1 {
			log.Printf("[DEBUG] Waiting for autoscale VM groups to be disabled (attempt %d/%d). Groups still enabled: %v",
				i+1, maxRetries, enabledGroups)
			time.Sleep(2 * time.Second)
		}
	}

	return fmt.Errorf("Timeout waiting for autoscale VM groups to be disabled after %d seconds", maxRetries*2)
}

func resourceCloudStackLoadBalancerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Make sure all required parameters are there
	if err := verifyLoadBalancerRule(d); err != nil {
		return err
	}

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("algorithm") {
		name := d.Get("name").(string)

		// Create new parameter struct
		p := cs.LoadBalancer.NewUpdateLoadBalancerRuleParams(d.Id())

		if d.HasChange("name") {
			log.Printf("[DEBUG] Name has changed for load balancer rule %s, starting update", name)

			p.SetName(name)
		}

		if d.HasChange("description") {
			log.Printf(
				"[DEBUG] Description has changed for load balancer rule %s, starting update", name)

			p.SetDescription(d.Get("description").(string))
		}

		if d.HasChange("algorithm") {
			algorithm := d.Get("algorithm").(string)

			log.Printf(
				"[DEBUG] Algorithm has changed to %s for load balancer rule %s, starting update",
				algorithm,
				name,
			)

			// Set the new Algorithm
			p.SetAlgorithm(algorithm)
		}

		_, err := cs.LoadBalancer.UpdateLoadBalancerRule(p)
		if err != nil {
			return fmt.Errorf(
				"Error updating load balancer rule %s", name)
		}
	}

	if d.HasChange("certificate_id") {
		p := cs.LoadBalancer.NewRemoveCertFromLoadBalancerParams(d.Id())
		if _, err := cs.LoadBalancer.RemoveCertFromLoadBalancer(p); err != nil {
			return err
		}

		_, certificateID := d.GetChange("certificate_id")
		cp := cs.LoadBalancer.NewAssignCertToLoadBalancerParams(certificateID.(string), d.Id())
		if _, err := cs.LoadBalancer.AssignCertToLoadBalancer(cp); err != nil {
			return err
		}
	}

	if d.HasChange("member_ids") {
		log.Printf("[DEBUG] Load balancer rule %s member_ids change detected", d.Id())

		asgCheckParams := cs.AutoScale.NewListAutoScaleVmGroupsParams()
		asgCheckParams.SetLbruleid(d.Id())

		asgGroups, err := cs.AutoScale.ListAutoScaleVmGroups(asgCheckParams)
		if err != nil {
			log.Printf("[WARN] Failed to check for autoscale VM groups: %s", err)
		}

		if len(asgGroups.AutoScaleVmGroups) > 0 {
			log.Printf("[INFO] Load balancer rule %s is managed by %d autoscale VM group(s), handling member updates carefully",
				d.Id(), len(asgGroups.AutoScaleVmGroups))

			o, n := d.GetChange("member_ids")
			ombs, nmbs := o.(*schema.Set), n.(*schema.Set)

			setToStringList := func(s *schema.Set) []string {
				l := make([]string, s.Len())
				for i, v := range s.List() {
					l[i] = v.(string)
				}
				return l
			}

			oldMembers := setToStringList(ombs)
			newMembers := setToStringList(nmbs)

			log.Printf("[DEBUG] Terraform state - old members: %v, new members: %v", oldMembers, newMembers)

			p := cs.LoadBalancer.NewListLoadBalancerRuleInstancesParams(d.Id())
			currentInstances, err := cs.LoadBalancer.ListLoadBalancerRuleInstances(p)
			if err != nil {
				return fmt.Errorf("Error listing current load balancer members: %s", err)
			}

			var currentMembers []string
			for _, i := range currentInstances.LoadBalancerRuleInstances {
				currentMembers = append(currentMembers, i.Id)
			}

			log.Printf("[DEBUG] CloudStack actual members: %v", currentMembers)

			// If Terraform state is empty but CloudStack has members, it means autoscale is managing them
			if len(oldMembers) == 0 && len(currentMembers) > 0 {
				log.Printf("[INFO] Detected autoscale-managed members in load balancer. Skipping member updates to avoid conflicts.")
				log.Printf("[INFO] Autoscale VM groups will manage the member lifecycle automatically.")

				d.Set("member_ids", currentMembers)
				return resourceCloudStackLoadBalancerRuleRead(d, meta)
			}

			if len(newMembers) > 0 {
				log.Printf("[WARN] Explicit member_ids specified for autoscale-managed load balancer. This may conflict with autoscale operations.")

				if err := waitForASGsToBeDisabled(cs, d.Id()); err != nil {
					return fmt.Errorf("Autoscale VM groups must be disabled before modifying load balancer members: %s", err)
				}

				membersToAdd := setToStringList(nmbs.Difference(ombs))
				membersToRemove := setToStringList(ombs.Difference(nmbs))

				log.Printf("[DEBUG] Explicit member changes - to add: %v, to remove: %v", membersToAdd, membersToRemove)

				if len(membersToRemove) > 0 {
					log.Printf("[DEBUG] Removing %d explicit members from load balancer rule %s", len(membersToRemove), d.Id())
					removeParams := cs.LoadBalancer.NewRemoveFromLoadBalancerRuleParams(d.Id())
					removeParams.SetVirtualmachineids(membersToRemove)
					if _, err := cs.LoadBalancer.RemoveFromLoadBalancerRule(removeParams); err != nil {
						return fmt.Errorf("Error removing explicit members from load balancer rule %s: %s. Members: %v", d.Id(), err, membersToRemove)
					}
				}

				if len(membersToAdd) > 0 {
					log.Printf("[DEBUG] Adding %d explicit members to load balancer rule %s", len(membersToAdd), d.Id())
					addParams := cs.LoadBalancer.NewAssignToLoadBalancerRuleParams(d.Id())
					addParams.SetVirtualmachineids(membersToAdd)
					if _, err := cs.LoadBalancer.AssignToLoadBalancerRule(addParams); err != nil {
						return fmt.Errorf("Error adding explicit members to load balancer rule %s: %s. Members: %v", d.Id(), err, membersToAdd)
					}
				}
			}
		} else {
			// No autoscale groups, proceed with normal member management
			log.Printf("[DEBUG] No autoscale groups found, proceeding with normal member management")

			o, n := d.GetChange("member_ids")
			ombs, nmbs := o.(*schema.Set), n.(*schema.Set)

			setToStringList := func(s *schema.Set) []string {
				l := make([]string, s.Len())
				for i, v := range s.List() {
					l[i] = v.(string)
				}
				return l
			}

			membersToAdd := setToStringList(nmbs.Difference(ombs))
			membersToRemove := setToStringList(ombs.Difference(nmbs))

			log.Printf("[DEBUG] Members to add: %v, remove: %v", membersToAdd, membersToRemove)

			if len(membersToRemove) > 0 {
				log.Printf("[DEBUG] Removing %d members from load balancer rule %s", len(membersToRemove), d.Id())
				p := cs.LoadBalancer.NewRemoveFromLoadBalancerRuleParams(d.Id())
				p.SetVirtualmachineids(membersToRemove)
				if _, err := cs.LoadBalancer.RemoveFromLoadBalancerRule(p); err != nil {
					return fmt.Errorf("Error removing members from load balancer rule %s: %s. Members to remove: %v", d.Id(), err, membersToRemove)
				}
				log.Printf("[DEBUG] Successfully removed members from load balancer rule")
			}

			if len(membersToAdd) > 0 {
				log.Printf("[DEBUG] Adding %d members to load balancer rule %s", len(membersToAdd), d.Id())
				p := cs.LoadBalancer.NewAssignToLoadBalancerRuleParams(d.Id())
				p.SetVirtualmachineids(membersToAdd)
				if _, err := cs.LoadBalancer.AssignToLoadBalancerRule(p); err != nil {
					return fmt.Errorf("Error adding members to load balancer rule %s: %s. Members to add: %v", d.Id(), err, membersToAdd)
				}
				log.Printf("[DEBUG] Successfully added members to load balancer rule")
			}
		}
	}

	return resourceCloudStackLoadBalancerRuleRead(d, meta)
}

func resourceCloudStackLoadBalancerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.LoadBalancer.NewDeleteLoadBalancerRuleParams(d.Id())

	log.Printf("[INFO] Deleting load balancer rule: %s", d.Get("name").(string))
	if _, err := cs.LoadBalancer.DeleteLoadBalancerRule(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if !strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return err
		}
	}

	return nil
}

func verifyLoadBalancerRule(d *schema.ResourceData) error {
	if protocol, ok := d.GetOk("protocol"); ok {
		protocol := protocol.(string)

		switch protocol {
		case "tcp", "udp", "tcp-proxy":
			// These are supported
		default:
			return fmt.Errorf(
				"%q is not a valid protocol. Valid options are 'tcp', 'udp' of 'tcp-proxy'", protocol)
		}
	}

	return nil
}
