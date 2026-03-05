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

func resourceCloudStackPortForward() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackPortForwardCreate,
		Read:   resourceCloudStackPortForwardRead,
		Update: resourceCloudStackPortForwardUpdate,
		Delete: resourceCloudStackPortForwardDelete,

		Schema: map[string]*schema.Schema{
			"ip_address_id": {
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
				Computed: true,
				ForceNew: true,
			},

			"forward": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},

						"private_port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"private_end_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"public_port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"public_end_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"virtual_machine_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"vm_guest_ip": {
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

func resourceCloudStackPortForwardCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// We need to set this upfront in order to be able to save a partial state
	d.SetId(d.Get("ip_address_id").(string))

	// If no project is explicitly set, try to inherit it from the IP address
	if _, ok := d.GetOk("project"); !ok {
		// Get the IP address to retrieve its project
		// Use projectid=-1 to search across all projects
		ip, count, err := cs.Address.GetPublicIpAddressByID(d.Id(), cloudstack.WithProject("-1"))
		if err == nil && count > 0 && ip.Projectid != "" {
			log.Printf("[DEBUG] Inheriting project %s from IP address %s", ip.Projectid, d.Id())
			// Set the project in the resource data for state management
			d.Set("project", ip.Project)
		}
	}

	// Create all forwards that are configured
	if nrs := d.Get("forward").(*schema.Set); nrs.Len() > 0 {
		// Create an empty schema.Set to hold all forwards
		forwards := resourceCloudStackPortForward().Schema["forward"].ZeroValue().(*schema.Set)

		err := createPortForwards(d, meta, forwards, nrs)
		if err != nil {
			return err
		}

		// We need to update this first to preserve the correct state
		d.Set("forward", forwards)
	}

	return resourceCloudStackPortForwardRead(d, meta)
}

func createPortForwards(d *schema.ResourceData, meta interface{}, forwards *schema.Set, nrs *schema.Set) error {
	var errs *multierror.Error

	var wg sync.WaitGroup
	wg.Add(nrs.Len())

	sem := make(chan struct{}, 10)
	for _, forward := range nrs.List() {
		// Put in a tiny sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(forward map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Create a single forward
			err := createPortForward(d, meta, forward)

			// If we have a UUID, we need to save the forward
			if forward["uuid"].(string) != "" {
				forwards.Add(forward)
			}

			if err != nil {
				errs = multierror.Append(errs, err)
			}

			<-sem
		}(forward.(map[string]interface{}))
	}

	wg.Wait()

	return errs.ErrorOrNil()
}

func createPortForward(d *schema.ResourceData, meta interface{}, forward map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Make sure all required parameters are there
	if err := verifyPortForwardParams(d, forward); err != nil {
		return err
	}

	// Query VM without project filter - it will be found regardless of project
	vm, _, err := cs.VirtualMachine.GetVirtualMachineByID(
		forward["virtual_machine_id"].(string),
	)
	if err != nil {
		return err
	}

	// Create a new parameter struct
	p := cs.Firewall.NewCreatePortForwardingRuleParams(d.Id(), forward["private_port"].(int),
		forward["protocol"].(string), forward["public_port"].(int), vm.Id)
	if val, ok := forward["private_end_port"]; ok && val != nil && val.(int) != 0 {
		p.SetPrivateendport(val.(int))
	}
	if val, ok := forward["public_end_port"]; ok && val != nil && val.(int) != 0 {
		p.SetPublicendport(val.(int))
	}

	if vmGuestIP, ok := forward["vm_guest_ip"]; ok && vmGuestIP.(string) != "" {
		p.SetVmguestip(vmGuestIP.(string))

		// Set the network ID based on the guest IP, needed when the public IP address
		// is not associated with any network yet
	NICS:
		for _, nic := range vm.Nic {
			if vmGuestIP.(string) == nic.Ipaddress {
				p.SetNetworkid(nic.Networkid)
				break NICS
			}
			for _, ip := range nic.Secondaryip {
				if vmGuestIP.(string) == ip.Ipaddress {
					p.SetNetworkid(nic.Networkid)
					break NICS
				}
			}
		}
	} else {
		// If no guest IP is configured, use the primary NIC
		p.SetNetworkid(vm.Nic[0].Networkid)
	}

	// Do not open the firewall automatically in any case
	p.SetOpenfirewall(false)

	r, err := cs.Firewall.CreatePortForwardingRule(p)
	if err != nil {
		return err
	}

	forward["uuid"] = r.Id

	return nil
}

func resourceCloudStackPortForwardRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// First check if the IP address is still associated
	ip, count, err := cs.Address.GetPublicIpAddressByID(
		d.Id(),
		cloudstack.WithProject(d.Get("project").(string)),
	)

	if err != nil {
		if count == 0 {
			log.Printf(
				"[DEBUG] IP address with ID %s is no longer associated", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	// Set the project if the IP address belongs to one
	setValueOrID(d, "project", ip.Project, ip.Projectid)

	// Get all the forwards from the running environment
	p := cs.Firewall.NewListPortForwardingRulesParams()
	p.SetIpaddressid(d.Id())
	p.SetListall(true)

	// Use the project from the IP address if it belongs to one
	// If no project, don't set projectid (use default scope)
	if ip.Projectid != "" {
		p.SetProjectid(ip.Projectid)
	}

	l, err := cs.Firewall.ListPortForwardingRules(p)
	if err != nil {
		return err
	}

	// Make a map of all the forwards so we can easily find a forward
	forwardMap := make(map[string]*cloudstack.PortForwardingRule, l.Count)
	for _, f := range l.PortForwardingRules {
		forwardMap[f.Id] = f
	}

	// Create an empty schema.Set to hold all forwards
	forwards := resourceCloudStackPortForward().Schema["forward"].ZeroValue().(*schema.Set)

	// Read all forwards that are configured
	if rs := d.Get("forward").(*schema.Set); rs.Len() > 0 {
		for _, forward := range rs.List() {
			forward := forward.(map[string]interface{})

			id, ok := forward["uuid"]
			if !ok || id.(string) == "" {
				// Forward doesn't have a UUID yet (shouldn't happen after Create, but handle gracefully)
				log.Printf("[DEBUG] Skipping forward without UUID: %+v", forward)
				continue
			}

			// Get the forward
			f, ok := forwardMap[id.(string)]
			if !ok {
				// Forward not found in API response - the rule was deleted outside of Terraform
				log.Printf("[WARN] Port forwarding rule %s not found in API response, removing from state", id.(string))
				// Don't add this forward to the new set - it will be removed from state
				continue
			}

			// Delete the known rule so only unknown rules remain in the ruleMap
			delete(forwardMap, id.(string))

			privPort, err := strconv.Atoi(f.Privateport)
			if err != nil {
				return err
			}

			pubPort, err := strconv.Atoi(f.Publicport)
			if err != nil {
				return err
			}

			// Update the values
			forward["protocol"] = f.Protocol
			forward["private_port"] = privPort
			forward["public_port"] = pubPort
			// Only set end ports if they differ from start ports (indicating a range)
			if f.Privateendport != "" && f.Privateendport != f.Privateport {
				privEndPort, err := strconv.Atoi(f.Privateendport)
				if err != nil {
					return err
				}
				forward["private_end_port"] = privEndPort
			}
			if f.Publicendport != "" && f.Publicendport != f.Publicport {
				pubEndPort, err := strconv.Atoi(f.Publicendport)
				if err != nil {
					return err
				}
				forward["public_end_port"] = pubEndPort
			}
			forward["virtual_machine_id"] = f.Virtualmachineid

			// This one is a bit tricky. We only want to update this optional value
			// if we've set one ourselves. If not this would become a computed value
			// and that would mess up the calculated hash of the set item.
			if forward["vm_guest_ip"].(string) != "" {
				forward["vm_guest_ip"] = f.Vmguestip
			}

			forwards.Add(forward)
		}
	}

	// If this is a managed resource, add all unknown forwards to dummy forwards
	managed := d.Get("managed").(bool)
	if managed && len(forwardMap) > 0 {
		for uuid := range forwardMap {
			// Make a dummy forward to hold the unknown UUID
			forward := map[string]interface{}{
				"protocol":           uuid,
				"private_port":       0,
				"public_port":        0,
				"virtual_machine_id": uuid,
				"uuid":               uuid,
			}

			// Add the dummy forward to the forwards set
			forwards.Add(forward)
		}
	}

	// Always set the forward attribute to maintain consistent state
	d.Set("forward", forwards)

	// Only remove the resource from state if it's not managed and has no forwards
	if forwards.Len() == 0 && !managed {
		d.SetId("")
	}

	return nil
}

func resourceCloudStackPortForwardUpdate(d *schema.ResourceData, meta interface{}) error {
	// Check if the forward set as a whole has changed
	if d.HasChange("forward") {
		o, n := d.GetChange("forward")
		ors := o.(*schema.Set).Difference(n.(*schema.Set))
		nrs := n.(*schema.Set).Difference(o.(*schema.Set))

		// We need to start with a rule set containing all the rules we
		// already have and want to keep. Any rules that are not deleted
		// correctly and any newly created rules, will be added to this
		// set to make sure we end up in a consistent state
		forwards := o.(*schema.Set).Intersection(n.(*schema.Set))

		// First loop through all the old forwards and delete them
		if ors.Len() > 0 {
			err := deletePortForwards(d, meta, forwards, ors)

			// We need to update this first to preserve the correct state
			d.Set("forward", forwards)

			if err != nil {
				return err
			}
		}

		// Then loop through all the new forwards and create them
		if nrs.Len() > 0 {
			err := createPortForwards(d, meta, forwards, nrs)

			// We need to update this first to preserve the correct state
			d.Set("forward", forwards)

			if err != nil {
				return err
			}
		}
	}

	return resourceCloudStackPortForwardRead(d, meta)
}

func resourceCloudStackPortForwardDelete(d *schema.ResourceData, meta interface{}) error {
	// Create an empty rule set to hold all rules that where
	// not deleted correctly
	forwards := resourceCloudStackPortForward().Schema["forward"].ZeroValue().(*schema.Set)

	// Delete all forwards
	if ors := d.Get("forward").(*schema.Set); ors.Len() > 0 {
		err := deletePortForwards(d, meta, forwards, ors)

		// We need to update this first to preserve the correct state
		d.Set("forward", forwards)

		if err != nil {
			return err
		}
	}

	return nil
}

func deletePortForwards(d *schema.ResourceData, meta interface{}, forwards *schema.Set, ors *schema.Set) error {
	var errs *multierror.Error

	var wg sync.WaitGroup
	wg.Add(ors.Len())

	sem := make(chan struct{}, 10)
	for _, forward := range ors.List() {
		// Put a sleep here to avoid DoS'ing the API
		time.Sleep(500 * time.Millisecond)

		go func(forward map[string]interface{}) {
			defer wg.Done()
			sem <- struct{}{}

			// Delete a single forward
			err := deletePortForward(d, meta, forward)

			// If we have a UUID, we need to save the forward
			if forward["uuid"].(string) != "" {
				forwards.Add(forward)
			}

			if err != nil {
				errs = multierror.Append(errs, err)
			}

			<-sem
		}(forward.(map[string]interface{}))
	}

	wg.Wait()

	return errs.ErrorOrNil()
}

func deletePortForward(d *schema.ResourceData, meta interface{}, forward map[string]interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create the parameter struct
	p := cs.Firewall.NewDeletePortForwardingRuleParams(forward["uuid"].(string))

	// Delete the forward
	if _, err := cs.Firewall.DeletePortForwardingRule(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if !strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", forward["uuid"].(string))) {
			return err
		}
	}

	// Empty the UUID of this rule
	forward["uuid"] = ""

	return nil
}

func verifyPortForwardParams(d *schema.ResourceData, forward map[string]interface{}) error {
	protocol := forward["protocol"].(string)
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf(
			"%s is not a valid protocol. Valid options are 'tcp' and 'udp'", protocol)
	}
	return nil
}
