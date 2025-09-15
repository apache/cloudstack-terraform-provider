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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloudStackCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackClusterCreate,
		Read:   resourceCloudStackClusterRead,
		Update: resourceCloudStackClusterUpdate,
		Delete: resourceCloudStackClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CloudManaged",
					"ExternalManaged",
				}, false),
			},
			"hypervisor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"XenServer",
					"KVM",
					"VMware",
					"Hyperv",
					"BareMetal",
					"Simulator",
					"Ovm3",
				}, false),
			},
			"pod_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"allocation_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"arch": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"x86_64",
					"aarch64",
				}, false),
			},
			"guest_vswitch_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"guest_vswitch_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vmwaresvs",
					"vmwaredvs",
				}, false),
			},
			"ovm3cluster": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ovm3pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ovm3vip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"public_vswitch_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_vswitch_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vmwaresvs",
					"vmwaredvs",
				}, false),
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vsm_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vsm_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"vsm_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pod_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"managed_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_overcommit_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_overcommit_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackClusterCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	clusterType := d.Get("cluster_type").(string)
	hypervisor := d.Get("hypervisor").(string)
	podID := d.Get("pod_id").(string)
	zoneID := d.Get("zone_id").(string)

	// Create a new parameter struct
	p := cs.Cluster.NewAddClusterParams(name, clusterType, hypervisor, podID, zoneID)

	// Set optional parameters
	if allocationState, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(allocationState.(string))
	}

	if arch, ok := d.GetOk("arch"); ok {
		p.SetArch(arch.(string))
	}

	if guestVSwitchName, ok := d.GetOk("guest_vswitch_name"); ok {
		p.SetGuestvswitchname(guestVSwitchName.(string))
	}

	if guestVSwitchType, ok := d.GetOk("guest_vswitch_type"); ok {
		p.SetGuestvswitchtype(guestVSwitchType.(string))
	}

	if ovm3cluster, ok := d.GetOk("ovm3cluster"); ok {
		p.SetOvm3cluster(ovm3cluster.(string))
	}

	if ovm3pool, ok := d.GetOk("ovm3pool"); ok {
		p.SetOvm3pool(ovm3pool.(string))
	}

	if ovm3vip, ok := d.GetOk("ovm3vip"); ok {
		p.SetOvm3vip(ovm3vip.(string))
	}

	if password, ok := d.GetOk("password"); ok {
		p.SetPassword(password.(string))
	}

	if publicVSwitchName, ok := d.GetOk("public_vswitch_name"); ok {
		p.SetPublicvswitchname(publicVSwitchName.(string))
	}

	if publicVSwitchType, ok := d.GetOk("public_vswitch_type"); ok {
		p.SetPublicvswitchtype(publicVSwitchType.(string))
	}

	if url, ok := d.GetOk("url"); ok {
		p.SetUrl(url.(string))
	}

	if username, ok := d.GetOk("username"); ok {
		p.SetUsername(username.(string))
	}

	if vsmIPAddress, ok := d.GetOk("vsm_ip_address"); ok {
		p.SetVsmipaddress(vsmIPAddress.(string))
	}

	if vsmPassword, ok := d.GetOk("vsm_password"); ok {
		p.SetVsmpassword(vsmPassword.(string))
	}

	if vsmUsername, ok := d.GetOk("vsm_username"); ok {
		p.SetVsmusername(vsmUsername.(string))
	}

	log.Printf("[DEBUG] Creating Cluster %s", name)
	r, err := cs.Cluster.AddCluster(p)
	if err != nil {
		return fmt.Errorf("Error creating Cluster %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCloudStackClusterRead(d, meta)
}

func resourceCloudStackClusterRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the Cluster details
	c, count, err := cs.Cluster.GetClusterByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Cluster %s does no longer exist", d.Get("name").(string))
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", c.Name)
	d.Set("cluster_type", c.Clustertype)
	d.Set("hypervisor", c.Hypervisortype)
	d.Set("pod_id", c.Podid)
	d.Set("pod_name", c.Podname)
	d.Set("zone_id", c.Zoneid)
	d.Set("zone_name", c.Zonename)
	d.Set("allocation_state", c.Allocationstate)
	d.Set("managed_state", c.Managedstate)
	d.Set("cpu_overcommit_ratio", c.Cpuovercommitratio)
	d.Set("memory_overcommit_ratio", c.Memoryovercommitratio)
	d.Set("arch", c.Arch)
	d.Set("ovm3vip", c.Ovm3vip)

	return nil
}

func resourceCloudStackClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Cluster.NewUpdateClusterParams(d.Id())

	if d.HasChange("name") {
		p.SetClustername(d.Get("name").(string))
	}

	if d.HasChange("allocation_state") {
		p.SetAllocationstate(d.Get("allocation_state").(string))
	}

	// Note: managed_state is a computed field and cannot be set directly

	_, err := cs.Cluster.UpdateCluster(p)
	if err != nil {
		return fmt.Errorf("Error updating Cluster %s: %s", d.Get("name").(string), err)
	}

	return resourceCloudStackClusterRead(d, meta)
}

func resourceCloudStackClusterDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Cluster.NewDeleteClusterParams(d.Id())

	log.Printf("[DEBUG] Deleting Cluster %s", d.Get("name").(string))
	_, err := cs.Cluster.DeleteCluster(p)

	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting Cluster %s: %s", d.Get("name").(string), err)
	}

	return nil
}
