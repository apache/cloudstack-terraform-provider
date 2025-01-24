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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackStoragePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackStoragePoolCreate,
		Read:   resourceCloudStackStoragePoolRead,
		Update: resourceCloudStackStoragePoolUpdate,
		Delete: resourceCloudStackStoragePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Description: "the cluster ID for the storage pool",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"hypervisor": {
				Description: "hypervisor type of the hosts in zone that will be attached to this storage pool. KVM, VMware supported as of now.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "the name for the storage pool",
				Type:        schema.TypeString,
				Required:    true,
			},
			"pod_id": {
				Description: "the Pod ID for the storage pool",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"storage_provider": {
				Description: "Storage provider for this pool",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"scope": {
				Description: "the scope of the storage: cluster or zone",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"state": {
				Description: "the state of the storage pool",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				// "Maintenance","Disabled","Up",
			},
			"tags": {
				Description: "the tags for the storage pool",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"url": {
				Description: "the URL of the storage pool",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"zone_id": {
				Description: "the Zone ID for the storage pool",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudStackStoragePoolCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Pool.NewCreateStoragePoolParams(d.Get("name").(string), d.Get("url").(string), d.Get("zone_id").(string))
	if v, ok := d.GetOk("cluster_id"); ok {
		p.SetClusterid(v.(string))
	}
	if v, ok := d.GetOk("hypervisor"); ok {
		p.SetHypervisor(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}
	if v, ok := d.GetOk("pod_id"); ok {
		p.SetPodid(v.(string))
	}
	if v, ok := d.GetOk("storage_provider"); ok {
		p.SetProvider(v.(string))
	}
	if v, ok := d.GetOk("scope"); ok {
		p.SetScope(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(v.(string))
	}
	if v, ok := d.GetOk("url"); ok {
		p.SetUrl(v.(string))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		p.SetZoneid(v.(string))
	}

	r, err := cs.Pool.CreateStoragePool(p)
	if err != nil {
		return err
	}

	d.SetId(r.Id)

	return resourceCloudStackStoragePoolRead(d, meta)
}

func resourceCloudStackStoragePoolRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.Pool.GetStoragePoolByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("cluster_id", r.Clusterid)
	d.Set("hypervisor", r.Hypervisor)
	d.Set("name", r.Name)
	d.Set("pod_id", r.Podid)
	d.Set("storage_provider", r.Provider)
	d.Set("scope", r.Scope)
	d.Set("state", r.State)
	d.Set("tags", r.Tags)
	d.Set("zone_id", r.Zoneid)

	return nil
}

func resourceCloudStackStoragePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Pool.NewUpdateStoragePoolParams(d.Id())
	if v, ok := d.GetOk("capacity_bytes"); ok {
		p.SetCapacitybytes(v.(int64))
	}
	if v, ok := d.GetOk("capacity_iops"); ok {
		p.SetCapacityiops(v.(int64))
	}
	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(strings.Split(v.(string), ","))
	}

	if v, ok := d.GetOk("state"); ok {
		if v == "Up" {
			p.SetEnabled(true)
		} else if v == "Disabled" {
			p.SetEnabled(false)
		} else if v == "Maintenance" {
			_, err := cs.StoragePool.EnableStorageMaintenance(cs.StoragePool.NewEnableStorageMaintenanceParams(d.Id()))
			if err != nil {
				return err
			}
		}
	}

	_, err := cs.Pool.UpdateStoragePool(p)
	if err != nil {
		return err
	}

	return resourceCloudStackStoragePoolRead(d, meta)
}

func resourceCloudStackStoragePoolDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.Pool.DeleteStoragePool(cs.Pool.NewDeleteStoragePoolParams(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting storage pool: %s", err)
	}

	return nil
}
