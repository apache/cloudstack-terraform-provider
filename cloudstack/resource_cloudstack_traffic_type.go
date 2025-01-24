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

func resourceCloudStackTrafficType() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackTrafficTypeCreate,
		Read:   resourceCloudStackTrafficTypeRead,
		// Update: resourceCloudStackTrafficTypeUpdate,
		Delete: resourceCloudStackTrafficTypeDelete,
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

	p := cs.Usage.NewAddTrafficTypeParams(d.Get("physical_network_id").(string), d.Get("traffic_type").(string))
	if v, ok := d.GetOk("hyperv_network_label"); ok {
		p.SetHypervnetworklabel(v.(string))
	}
	if v, ok := d.GetOk("isolation_method"); ok {
		p.SetIsolationmethod(v.(string))
	}
	if v, ok := d.GetOk("kvm_network_label"); ok {
		p.SetKvmnetworklabel(v.(string))
	}
	if v, ok := d.GetOk("ovm3_network_label"); ok {
		p.SetOvm3networklabel(v.(string))
	}
	if v, ok := d.GetOk("vlan"); ok {
		p.SetVlan(v.(string))
	}
	if v, ok := d.GetOk("vmware_network_label"); ok {
		p.SetVmwarenetworklabel(v.(string))
	}
	if v, ok := d.GetOk("xen_network_label"); ok {
		p.SetXennetworklabel(v.(string))
	}

	r, err := cs.Usage.AddTrafficType(p)
	if err != nil {
		return err
	}
	d.SetId(r.Id)
	d.Set("physical_network_id", d.Get("physical_network_id").(string))

	//
	d.Set("hyperv_network_label", r.Hypervnetworklabel)
	d.Set("kvm_network_label", r.Kvmnetworklabel)
	d.Set("ovm3_network_label", r.Ovm3networklabel)
	d.Set("traffic_type", r.Traffictype)
	d.Set("vmware_network_label", r.Vmwarenetworklabel)
	d.Set("xen_network_label", r.Xennetworklabel)

	return resourceCloudStackTrafficTypeRead(d, meta)
}

func resourceCloudStackTrafficTypeRead(d *schema.ResourceData, meta interface{}) error {
	// Nothing to read.  While these fields are returned by the API
	// they are not documented in the client or ListApi spec.
	// see https://github.com/apache/cloudstack/issues/7837

	return nil
}

func resourceCloudStackTrafficTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	// All fields are ForceNew or Computed w/out Optional, Update is superfluous
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Usage.NewUpdateTrafficTypeParams(d.Id())
	if v, ok := d.GetOk("hyperv_network_label"); ok {
		p.SetHypervnetworklabel(v.(string))
	}
	if v, ok := d.GetOk("kvm_network_label"); ok {
		p.SetKvmnetworklabel(v.(string))
	}
	if v, ok := d.GetOk("ovm3_network_label"); ok {
		p.SetOvm3networklabel(v.(string))
	}
	if v, ok := d.GetOk("vmware_network_label"); ok {
		p.SetVmwarenetworklabel(v.(string))
	}
	if v, ok := d.GetOk("xen_network_label"); ok {
		p.SetXennetworklabel(v.(string))
	}

	r, err := cs.Usage.UpdateTrafficType(p)
	if err != nil {
		return err
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

	_, err := cs.Usage.DeleteTrafficType(cs.Usage.NewDeleteTrafficTypeParams(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting traffic type: %s", err)
	}

	return nil
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
