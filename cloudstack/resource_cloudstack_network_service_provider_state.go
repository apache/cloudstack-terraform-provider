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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackNetworkServiceProviderState() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackNetworkServiceProviderStateUpdate,
		Read:   resourceCloudStackNetworkServiceProviderStateRead,
		Update: resourceCloudStackNetworkServiceProviderStateUpdate,
		Delete: resourceCloudStackNetworkServiceProviderStateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"physical_network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceCloudStackNetworkServiceProviderStateRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	nsp, err := getNetworkServiceProvider(cs, d)
	if err != nil {
		return err
	}

	d.SetId(nsp.Id)
	if nsp.State == "Enabled" {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	// d.Set("enabled", nsp.State)

	return nil
}

func resourceCloudStackNetworkServiceProviderStateUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Network service provider
	nsp, err := getNetworkServiceProvider(cs, d)
	if err != nil {
		return err
	}

	switch d.Get("name") {
	case "VirtualRouter", "VpcVirtualRouter":
		err := virtualRouterElementState(cs, nsp, d.Get("enabled").(bool))
		if err != nil {
			return err
		}
	case "InternalLbVm":
		err := internalLbVmElementState(cs, nsp, d.Get("enabled").(bool))
		if err != nil {
			return err
		}
	case "ConfigDrive":
		// No elements to configure
	default:
		return fmt.Errorf("Service provider (%s) name not supported.", d.Get("name"))
	}

	// Service provider state
	pUNSPP := cs.Network.NewUpdateNetworkServiceProviderParams(nsp.Id)
	if d.Get("enabled").(bool) {
		pUNSPP.SetState("Enabled")
	} else {
		pUNSPP.SetState("Disabled")
	}
	_, err = cs.Network.UpdateNetworkServiceProvider(pUNSPP)
	if err != nil {
		return err
	}

	return resourceCloudStackNetworkServiceProviderStateRead(d, meta)
}

func resourceCloudStackNetworkServiceProviderStateDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Network.NewUpdateNetworkServiceProviderParams(d.Id())
	p.SetState("Disabled")

	_, err := cs.Network.UpdateNetworkServiceProvider(p)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func getNetworkServiceProvider(cs *cloudstack.CloudStackClient, d *schema.ResourceData) (*cloudstack.NetworkServiceProvider, error) {
	p := cs.Network.NewListNetworkServiceProvidersParams()
	if _, ok := d.GetOk("name"); ok {
		p.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("physical_network_id"); ok {
		p.SetPhysicalnetworkid(d.Get("physical_network_id").(string))
	}

	// list network service providers
	r, err := cs.Network.ListNetworkServiceProviders(p)
	if err != nil {
		return nil, err
	}

	for _, nsp := range r.NetworkServiceProviders {
		if _, ok := d.GetOk("name"); ok {
			if nsp.Name == d.Get("name").(string) {
				return nsp, nil
			}
		} else if nsp.Id == d.Id() {
			return nsp, nil
		}
	}

	return nil, fmt.Errorf("Service provider element id not found.")
}

func virtualRouterElementState(cs *cloudstack.CloudStackClient, nsp *cloudstack.NetworkServiceProvider, state bool) error {
	// VirtualRouterElement state
	p := cs.Router.NewListVirtualRouterElementsParams()
	p.SetNspid(nsp.Id)

	vre, err := cs.Router.ListVirtualRouterElements(p)
	if err != nil {
		return err
	}

	var vreID string
	for _, e := range vre.VirtualRouterElements {
		if nsp.Id == e.Nspid {
			vreID = e.Id
			break
		}
		return fmt.Errorf("Service provider element id (nspod) not found: %s.", nsp.Id)
	}

	_, err = cs.Router.ConfigureVirtualRouterElement(cs.Router.NewConfigureVirtualRouterElementParams(state, vreID))
	if err != nil {
		return err
	}

	return nil

}

func internalLbVmElementState(cs *cloudstack.CloudStackClient, nsp *cloudstack.NetworkServiceProvider, state bool) error {
	// InternalLoadBalancerElement state
	p := cs.InternalLB.NewListInternalLoadBalancerElementsParams()
	p.SetNspid(nsp.Id)

	ilbe, err := cs.InternalLB.ListInternalLoadBalancerElements(p)
	if err != nil {
		return err
	}

	var ilbeID string
	for _, e := range ilbe.InternalLoadBalancerElements {
		if nsp.Id == e.Nspid {
			ilbeID = e.Id
			break
		}
		return fmt.Errorf("Service provider element id (nspod) not found: %s.", nsp.Id)
	}

	ilvm, err := cs.InternalLB.ConfigureInternalLoadBalancerElement(cs.InternalLB.NewConfigureInternalLoadBalancerElementParams(state, ilbeID))
	if err != nil {
		return err
	}

	_, err = cs.InternalLB.ConfigureInternalLoadBalancerElement(cs.InternalLB.NewConfigureInternalLoadBalancerElementParams(state, ilvm.Id))
	if err != nil {
		return err
	}

	return nil
}
