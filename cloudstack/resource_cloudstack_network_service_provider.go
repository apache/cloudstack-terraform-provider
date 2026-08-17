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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackNetworkServiceProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNetworkServiceProviderCreate,
		Read:   resourceCloudStackNetworkServiceProviderRead,
		Update: resourceCloudStackNetworkServiceProviderUpdate,
		Delete: resourceCloudStackNetworkServiceProviderDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudStackNetworkServiceProviderImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"physical_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_physical_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"service_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Enabled" && v != "Disabled" {
						errs = append(errs, fmt.Errorf("%q must be either 'Enabled' or 'Disabled', got: %s", key, v))
					}
					return
				},
			},
		},
	}
}

func resourceCloudStackNetworkServiceProviderCreate(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	physicalNetworkID := d.Get("physical_network_id").(string)

	// Check if the provider already exists
	p := cs.Network.NewListNetworkServiceProvidersParams()
	p.SetPhysicalnetworkid(physicalNetworkID)
	p.SetName(name)

	l, err := cs.Network.ListNetworkServiceProviders(p)
	if err != nil {
		return fmt.Errorf("Error checking for existing network service provider %s: %s", name, err)
	}

	if l.Count > 0 {
		// Provider already exists, use its ID
		d.SetId(l.NetworkServiceProviders[0].Id)

		// Update the provider if needed
		needsUpdate := false
		up := cs.Network.NewUpdateNetworkServiceProviderParams(d.Id())

		// Update service list if provided and not SecurityGroupProvider
		if serviceList, ok := d.GetOk("service_list"); ok && name != "SecurityGroupProvider" {
			services := make([]string, len(serviceList.([]any)))
			for i, v := range serviceList.([]any) {
				services[i] = v.(string)
			}
			up.SetServicelist(services)
			needsUpdate = true
		}

		// Update state if provided
		if state, ok := d.GetOk("state"); ok {
			up.SetState(state.(string))
			needsUpdate = true
		}

		// Perform the update if needed
		if needsUpdate {
			_, err := cs.Network.UpdateNetworkServiceProvider(up)
			if err != nil {
				return fmt.Errorf("Error updating network service provider %s: %s", name, err)
			}
		}
	} else {
		// Provider doesn't exist, create a new one
		cp := cs.Network.NewAddNetworkServiceProviderParams(name, physicalNetworkID)

		// Set optional parameters
		if destinationPhysicalNetworkID, ok := d.GetOk("destination_physical_network_id"); ok {
			cp.SetDestinationphysicalnetworkid(destinationPhysicalNetworkID.(string))
		}

		if serviceList, ok := d.GetOk("service_list"); ok {
			services := make([]string, len(serviceList.([]any)))
			for i, v := range serviceList.([]any) {
				services[i] = v.(string)
			}
			cp.SetServicelist(services)
		}

		// Create the network service provider
		r, err := cs.Network.AddNetworkServiceProvider(cp)
		if err != nil {
			return fmt.Errorf("Error creating network service provider %s: %s", name, err)
		}

		d.SetId(r.Id)
	}

	return resourceCloudStackNetworkServiceProviderRead(d, meta)
}

func resourceCloudStackNetworkServiceProviderRead(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the network service provider details
	p := cs.Network.NewListNetworkServiceProvidersParams()
	p.SetPhysicalnetworkid(d.Get("physical_network_id").(string))

	l, err := cs.Network.ListNetworkServiceProviders(p)
	if err != nil {
		return err
	}

	// Find the network service provider with the matching ID
	var provider *cloudstack.NetworkServiceProvider
	for _, p := range l.NetworkServiceProviders {
		if p.Id == d.Id() {
			provider = p
			break
		}
	}

	if provider == nil {
		log.Printf("[DEBUG] Network service provider %s does no longer exist", d.Get("name").(string))
		d.SetId("")
		return nil
	}

	d.Set("name", provider.Name)
	d.Set("physical_network_id", provider.Physicalnetworkid)
	d.Set("state", provider.State)

	// Special handling for SecurityGroupProvider - don't set service_list to avoid drift
	if provider.Name == "SecurityGroupProvider" {
		// For SecurityGroupProvider, we don't manage the service list
		// as it's predefined and can't be modified
		if _, ok := d.GetOk("service_list"); ok {
			// If service_list was explicitly set in config, keep it for consistency
			// but don't update it from the API response
		} else {
			// If service_list wasn't in config, don't set it to avoid drift
		}
	} else {
		// For other providers, set service list if available
		if len(provider.Servicelist) > 0 {
			d.Set("service_list", provider.Servicelist)
		}
	}

	return nil
}

func resourceCloudStackNetworkServiceProviderUpdate(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Check if we need to update the provider
	if d.HasChange("service_list") || d.HasChange("state") {
		// Create a new parameter struct
		p := cs.Network.NewUpdateNetworkServiceProviderParams(d.Id())

		// Update service list if changed and not SecurityGroupProvider
		if d.HasChange("service_list") && d.Get("name").(string) != "SecurityGroupProvider" {
			if serviceList, ok := d.GetOk("service_list"); ok {
				services := make([]string, len(serviceList.([]any)))
				for i, v := range serviceList.([]any) {
					services[i] = v.(string)
				}
				p.SetServicelist(services)
			}
		}

		// Update state if changed
		if d.HasChange("state") {
			state := d.Get("state").(string)
			p.SetState(state)
		}

		// Update the network service provider
		_, err := cs.Network.UpdateNetworkServiceProvider(p)
		if err != nil {
			return fmt.Errorf("Error updating network service provider %s: %s", d.Get("name").(string), err)
		}
	}

	return resourceCloudStackNetworkServiceProviderRead(d, meta)
}

func resourceCloudStackNetworkServiceProviderDelete(d *schema.ResourceData, meta any) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Network.NewDeleteNetworkServiceProviderParams(d.Id())

	// Delete the network service provider
	_, err := cs.Network.DeleteNetworkServiceProvider(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting network service provider %s: %s", d.Get("name").(string), err)
	}

	return nil
}

func resourceCloudStackNetworkServiceProviderImport(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	// Import is expected to receive the network service provider ID
	cs := meta.(*cloudstack.CloudStackClient)

	// We need to determine the physical_network_id by listing all physical networks and their service providers
	p := cs.Network.NewListPhysicalNetworksParams()
	physicalNetworks, err := cs.Network.ListPhysicalNetworks(p)
	if err != nil {
		return nil, err
	}

	// For each physical network, list its service providers
	for _, pn := range physicalNetworks.PhysicalNetworks {
		sp := cs.Network.NewListNetworkServiceProvidersParams()
		sp.SetPhysicalnetworkid(pn.Id)
		serviceProviders, err := cs.Network.ListNetworkServiceProviders(sp)
		if err != nil {
			continue
		}

		// Check if our service provider ID is in this physical network
		for _, provider := range serviceProviders.NetworkServiceProviders {
			if provider.Id == d.Id() {
				// Found the physical network that contains our service provider
				d.Set("physical_network_id", pn.Id)
				d.Set("name", provider.Name)
				d.Set("state", provider.State)

				// Set service list if available
				if len(provider.Servicelist) > 0 {
					d.Set("service_list", provider.Servicelist)
				}

				return []*schema.ResourceData{d}, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find physical network for network service provider %s", d.Id())
}
