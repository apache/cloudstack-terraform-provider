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
	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"allocation_state": {
				Description: "Allocation state of this cluster for allocation of new resources",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"cluster_name": {
				Description: "the cluster name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cluster_type": {
				Description: "ype of the cluster: CloudManaged, ExternalManaged",
				Type:        schema.TypeString,
				Required:    true,
			},
			"guest_vswitch_name": {
				Description: "Name of virtual switch used for guest traffic in the cluster. This would override zone wide traffic label setting.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"guest_vswitch_type": {
				Description: "Type of virtual switch used for guest traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch)",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"hypervisor": {
				Description: "hypervisor type of the cluster: XenServer,KVM,VMware,Hyperv,BareMetal,Simulator,Ovm3",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ovm3_cluster": {
				Description: "Ovm3 native OCFS2 clustering enabled for cluster",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"ovm3_pool": {
				Description: "Ovm3 native pooling enabled for cluster",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"ovm3_vip": {
				Description: "Ovm3 vip to use for pool (and cluster)",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"password": {
				Description: "the password for the host",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"public_vswitch_name": {
				Description: "Name of virtual switch used for public traffic in the cluster. This would override zone wide traffic label setting.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"public_vswitch_type": {
				Description: "Type of virtual switch used for public traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch)",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"pod_id": {
				Description: "Type of virtual switch used for public traffic in the cluster. Allowed values are, vmwaresvs (for VMware standard vSwitch) and vmwaredvs (for VMware distributed vSwitch)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "the URL",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"username": {
				Description: "the username for the cluster",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"vsm_ip_address": {
				Description: "the ipaddress of the VSM associated with this cluster",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"vsm_password": {
				Description: "the password for the VSM associated with this cluster",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"vsm_username": {
				Description: "the username for the VSM associated with this cluster",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"zone_id": {
				Description: "the Zone ID for the cluster",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCloudStackClusterCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Cluster.NewAddClusterParams(d.Get("cluster_name").(string), d.Get("cluster_type").(string), d.Get("hypervisor").(string), d.Get("pod_id").(string), d.Get("zone_id").(string))
	if v, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(v.(string))
	}
	if v, ok := d.GetOk("guest_vswitch_name"); ok {
		p.SetGuestvswitchname(v.(string))
	}
	if v, ok := d.GetOk("guest_vswitch_type"); ok {
		p.SetGuestvswitchtype(v.(string))
	}
	if v, ok := d.GetOk("hypervisor"); ok {
		p.SetHypervisor(v.(string))
	}
	if v, ok := d.GetOk("ovm3_cluster"); ok {
		p.SetOvm3cluster(v.(string))
	}
	if v, ok := d.GetOk("ovm3_pool"); ok {
		p.SetOvm3pool(v.(string))
	}
	if v, ok := d.GetOk("ovm3_vip"); ok {
		p.SetOvm3vip(v.(string))
	}
	if v, ok := d.GetOk("password"); ok {
		p.SetPassword(v.(string))
	}
	if v, ok := d.GetOk("public_vswitch_name"); ok {
		p.SetPublicvswitchname(v.(string))
	}
	if v, ok := d.GetOk("public_vswitch_type"); ok {
		p.SetPublicvswitchtype(v.(string))
	}
	if v, ok := d.GetOk("url"); ok {
		p.SetUrl(v.(string))
	}
	if v, ok := d.GetOk("username"); ok {
		p.SetUsername(v.(string))
	}
	if v, ok := d.GetOk("vsm_ip_address"); ok {
		p.SetVsmipaddress(v.(string))
	}
	if v, ok := d.GetOk("vsm_password"); ok {
		p.SetVsmpassword(v.(string))
	}
	if v, ok := d.GetOk("vsm_username"); ok {
		p.SetVsmusername(v.(string))
	}

	r, err := cs.Cluster.AddCluster(p)
	if err != nil {
		return err
	}
	d.SetId(r.Id)

	return resourceCloudStackClusterRead(d, meta)
}

func resourceCloudStackClusterRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.Cluster.GetClusterByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("allocation_state", r.Allocationstate)
	d.Set("cluster_type", r.Clustertype)
	d.Set("hypervisor", r.Hypervisortype)
	d.Set("cluster_name", r.Name)
	d.Set("ovm3_vip", r.Ovm3vip)
	d.Set("pod_id", r.Podid)
	d.Set("zone_id", r.Zoneid)

	return nil
}

func resourceCloudStackClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Cluster.NewUpdateClusterParams(d.Id())
	if v, ok := d.GetOk("allocation_state"); ok {
		p.SetAllocationstate(v.(string))
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		p.SetClustername(v.(string))
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		p.SetClustertype(v.(string))
	}
	if v, ok := d.GetOk("hypervisor"); ok {
		p.SetHypervisor(v.(string))
	}

	_, err := cs.Cluster.UpdateCluster(p)
	if err != nil {
		return err
	}

	return resourceCloudStackClusterRead(d, meta)
}

func resourceCloudStackClusterDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.Cluster.DeleteCluster(cs.Cluster.NewDeleteClusterParams(d.Id()))
	if err != nil {
		return err
	}

	return nil
}
