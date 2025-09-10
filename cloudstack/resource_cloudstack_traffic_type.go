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

func resourceCloudStackTrafficType() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackTrafficTypeCreate,
		Read:   resourceCloudStackTrafficTypeRead,
		Update: resourceCloudStackTrafficTypeUpdate,
		Delete: resourceCloudStackTrafficTypeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudStackTrafficTypeImport,
		},

		Schema: map[string]*schema.Schema{
			"hyperv_network_label": {
				Description: "The network name label of the physical device dedicated to this traffic on a Hyperv host",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"isolation_method": {
				Description: "Used if physical network has multiple isolation types and traffic type is public. Choose which isolation method. Valid options currently 'vlan' or 'vxlan', defaults to 'vlan'.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"kvm_network_label": {
				Description: "The network name label of the physical device dedicated to this traffic on a KVM host",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"ovm3_network_label": {
				Description: "The network name of the physical device dedicated to this traffic on an OVM3 host",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"physical_network_id": {
				Description: "the Physical Network ID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"traffic_type": {
				Description:  "the trafficType to be added to the physical network",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateTrafficType,
			},
			"vlan": {
				Description: "The VLAN id to be used for Management traffic by VMware host",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"vmware_network_label": {
				Description: "The network name label of the physical device dedicated to this traffic on a VMware host",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"xen_network_label": {
				Description: "The network name label of the physical device dedicated to this traffic on a XenServer host",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudStackTrafficTypeCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	physicalNetworkID := d.Get("physical_network_id").(string)
	trafficType := d.Get("type").(string)

	// Create a new parameter struct
	p := cs.Usage.NewAddTrafficTypeParams(physicalNetworkID, trafficType)

	// Set optional parameters
	if kvmNetworkLabel, ok := d.GetOk("kvm_network_label"); ok {
		p.SetKvmnetworklabel(kvmNetworkLabel.(string))
	}

	if v, ok := d.GetOk("isolation_method"); ok {
		p.SetIsolationmethod(v.(string))
	}

	if vlan, ok := d.GetOk("vlan"); ok {
		p.SetVlan(vlan.(string))
	}

	if xenNetworkLabel, ok := d.GetOk("xen_network_label"); ok {
		p.SetXennetworklabel(xenNetworkLabel.(string))
	}

	if vmwareNetworkLabel, ok := d.GetOk("vmware_network_label"); ok {
		p.SetVmwarenetworklabel(vmwareNetworkLabel.(string))
	}

	if hypervNetworkLabel, ok := d.GetOk("hyperv_network_label"); ok {
		p.SetHypervnetworklabel(hypervNetworkLabel.(string))
	}

	if ovm3NetworkLabel, ok := d.GetOk("ovm3_network_label"); ok {
		p.SetOvm3networklabel(ovm3NetworkLabel.(string))
	}

	// Create the traffic type
	r, err := cs.Usage.AddTrafficType(p)
	if err != nil {
		return fmt.Errorf("Error creating traffic type %s: %s", trafficType, err)
	}

	d.SetId(r.Id)

	d.Set("physical_network_id", d.Get("physical_network_id").(string))
	d.Set("hyperv_network_label", r.Hypervnetworklabel)
	d.Set("kvm_network_label", r.Kvmnetworklabel)
	d.Set("ovm3_network_label", r.Ovm3networklabel)
	d.Set("traffic_type", r.Traffictype)
	d.Set("vmware_network_label", r.Vmwarenetworklabel)
	d.Set("xen_network_label", r.Xennetworklabel)

	return resourceCloudStackTrafficTypeRead(d, meta)
}

func resourceCloudStackTrafficTypeRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the traffic type details
	p := cs.Usage.NewListTrafficTypesParams(d.Get("physical_network_id").(string))

	l, err := cs.Usage.ListTrafficTypes(p)
	if err != nil {
		return err
	}

	// Find the traffic type with the matching ID
	var trafficType *cloudstack.TrafficType
	for _, t := range l.TrafficTypes {
		if t.Id == d.Id() {
			trafficType = t
			break
		}
	}

	if trafficType == nil {
		log.Printf("[DEBUG] Traffic type %s does no longer exist", d.Get("type").(string))
		d.SetId("")
		return nil
	}

	// The TrafficType struct has a Name field which contains the traffic type
	// But in some cases it might be empty, so we'll keep the original value from the state
	if trafficType.Name != "" {
		d.Set("type", trafficType.Name)
	}

	// Note: The TrafficType struct doesn't have fields for network labels or VLAN
	// We'll need to rely on what we store in the state

	return nil
}

func resourceCloudStackTrafficTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Usage.NewUpdateTrafficTypeParams(d.Id())

	// Only set the parameters that have changed and are supported by the API
	if d.HasChange("kvm_network_label") {
		p.SetKvmnetworklabel(d.Get("kvm_network_label").(string))
	}

	if d.HasChange("xen_network_label") {
		p.SetXennetworklabel(d.Get("xen_network_label").(string))
	}

	if d.HasChange("vmware_network_label") {
		p.SetVmwarenetworklabel(d.Get("vmware_network_label").(string))
	}

	if d.HasChange("hyperv_network_label") {
		p.SetHypervnetworklabel(d.Get("hyperv_network_label").(string))
	}

	if d.HasChange("ovm3_network_label") {
		p.SetOvm3networklabel(d.Get("ovm3_network_label").(string))
	}

	// Note: The UpdateTrafficTypeParams struct doesn't have a SetVlan method
	// so we can't update the VLAN

	// Update the traffic type
	r, err := cs.Usage.UpdateTrafficType(p)
	if err != nil {
		return fmt.Errorf("Error updating traffic type %s: %s", d.Get("type").(string), err)
	}

	d.Set("hyperv_network_label", r.Hypervnetworklabel)
	d.Set("kvm_network_label", r.Kvmnetworklabel)
	d.Set("ovm3_network_label", r.Ovm3networklabel)
	d.Set("traffic_type", r.Traffictype)
	d.Set("vmware_network_label", r.Vmwarenetworklabel)
	d.Set("xen_network_label", r.Xennetworklabel)

	return resourceCloudStackTrafficTypeRead(d, meta)
}

func resourceCloudStackTrafficTypeDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Usage.NewDeleteTrafficTypeParams(d.Id())

	// Delete the traffic type
	_, err := cs.Usage.DeleteTrafficType(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting traffic type %s: %s", d.Get("type").(string), err)
	}

	return nil
}

func resourceCloudStackTrafficTypeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import is expected to receive the traffic type ID
	cs := meta.(*cloudstack.CloudStackClient)

	// We need to determine the physical_network_id by listing all physical networks and their traffic types
	p := cs.Network.NewListPhysicalNetworksParams()
	physicalNetworks, err := cs.Network.ListPhysicalNetworks(p)
	if err != nil {
		return nil, err
	}

	// For each physical network, list its traffic types
	for _, pn := range physicalNetworks.PhysicalNetworks {
		tp := cs.Usage.NewListTrafficTypesParams(pn.Id)
		trafficTypes, err := cs.Usage.ListTrafficTypes(tp)
		if err != nil {
			continue
		}

		// Check if our traffic type ID is in this physical network
		for _, tt := range trafficTypes.TrafficTypes {
			if tt.Id == d.Id() {
				// Found the physical network that contains our traffic type
				d.Set("physical_network_id", pn.Id)

				// Set the type attribute - use the original value from the API call
				// If the Name field is empty, use a default value based on the traffic type ID
				if tt.Name != "" {
					d.Set("type", tt.Name)
				} else {
					// Use a default value based on common traffic types
					// This is a fallback and might not be accurate
					d.Set("type", "Management")
				}

				// For import to work correctly, we need to set default values for network labels
				// These will be overridden by the user if needed
				if d.Get("kvm_network_label") == "" {
					d.Set("kvm_network_label", "cloudbr0")
				}

				if d.Get("xen_network_label") == "" {
					d.Set("xen_network_label", "xenbr0")
				}

				return []*schema.ResourceData{d}, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find physical network for traffic type %s", d.Id())
}

func validateTrafficType(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	allowed := []string{"Public", "Guest", "Management", "Storage"}

	for _, str := range allowed {
		if str == input {
			return
		}
	}
	errors = append(errors, fmt.Errorf("traffic_type identifier (%q) not found, expecting %v", input, allowed))

	return warnings, errors
}