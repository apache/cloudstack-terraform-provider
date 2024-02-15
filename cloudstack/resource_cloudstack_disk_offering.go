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
	"strconv"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackDiskOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackDiskOfferingCreate,
		Read:   resourceCloudStackDiskOfferingRead,
		Update: resourceCloudStackDiskOfferingUpdate,
		Delete: resourceCloudStackDiskOfferingDelete,
		Schema: map[string]*schema.Schema{
			"display_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			//
			"cache_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disk_offering_strictness": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"domain_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"iops_read_rate": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops_read_rate_max": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops_read_rate_max_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops_write_rate": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops_write_rate_max": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iops_write_rate_max_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"provisioning_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hypervisor": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bytes_read_rate": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"bytes_read_rate_max": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"bytes_read_rate_max_length": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"bytes_write_rate": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"bytes_write_rate_max": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"bytes_write_rate_max_length": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"storage": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_iops": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"max_iops": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"customized_iops": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"hypervisor_snapshot_reserve": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceCloudStackDiskOfferingCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.DiskOffering.NewCreateDiskOfferingParams(d.Get("name").(string), d.Get("display_text").(string))

	if v, ok := d.GetOk("cache_mode"); ok {
		p.SetCachemode(v.(string))
	}
	if v, ok := d.GetOk("disk_size"); ok {
		p.SetDisksize(int64(v.(int)))
		p.SetCustomized(false)
	} else {
		p.SetCustomized(true)
	}
	if v, ok := d.GetOk("disk_offering_strictness"); ok {
		p.SetDisksizestrictness(v.(bool))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		domain_id := v.([]interface{})
		items := make([]string, len(domain_id))
		for i, raw := range domain_id {
			items[i] = raw.(string)
		}
		p.SetDomainid(items)
	}
	if v, ok := d.GetOk("iops_read_rate"); ok {
		p.SetIopsreadrate(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_read_rate_max"); ok {
		p.SetIopsreadratemax(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_read_rate_max_length"); ok {
		p.SetIopsreadratemaxlength(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate"); ok {
		p.SetIopsreadrate(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max"); ok {
		p.SetIopsreadratemax(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max_length"); ok {
		p.SetIopsreadratemaxlength(int64(v.(int)))
	}
	if v, ok := d.GetOk("provisioning_type"); ok {
		p.SetProvisioningtype(v.(string))
	}
	if v, ok := d.GetOk("storage_type"); ok {
		p.SetStoragetype(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(v.(string))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		zone_id := v.([]interface{})
		items := make([]string, len(zone_id))
		for i, raw := range zone_id {
			items[i] = raw.(string)
		}
		p.SetZoneid(items)
	}

	// storage qos
	if v, ok := d.GetOk("storage"); ok {
		storage := v.(map[string]interface{})

		// related? https://github.com/hashicorp/terraform-plugin-sdk/issues/211
		if v2, ok2 := storage["min_iops"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetMiniops(int64(tmp))
		}
		if v2, ok2 := storage["max_iops"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetMaxiops(int64(tmp))
		}
		if v2, ok2 := storage["customized_iops"]; ok2 {
			tmp, _ := strconv.ParseBool(v2.(string))
			p.SetCustomizediops(tmp)
		}
	}

	// hypervisor qos
	if v, ok := d.GetOk("hypervisor"); ok {
		hypervisor := v.(map[string]interface{})

		if v2, ok2 := hypervisor["bytes_read_rate"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_read_rate_max"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_read_rate_max_length"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_write_rate"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_write_rate_max"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_write_rate_max_length"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
	}

	diskOff, err := cs.DiskOffering.CreateDiskOffering(p)
	if err != nil {
		return err
	}

	d.SetId(diskOff.Id)

	return resourceCloudStackDiskOfferingRead(d, meta)
}

func resourceCloudStackDiskOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.DiskOffering.GetDiskOfferingByID(d.Id())
	if err != nil {
		return err
	}
	d.Set("display_text", r.Displaytext)
	d.Set("name", r.Name)

	//
	d.Set("cache_mode", r.CacheMode)
	d.Set("disk_size", r.Disksize)
	d.Set("disk_offering_strictness", r.Disksize)
	d.Set("domain_id", r.Domainid)
	d.Set("iops_read_rate", r.DiskIopsReadRate)
	d.Set("iops_read_rate_max", r.DiskIopsReadRateMax)
	d.Set("iops_read_rate_max_length", r.DiskIopsReadRateMaxLength)
	d.Set("iops_write_rate", r.DiskIopsWriteRate)
	d.Set("iops_write_rate_max", r.DiskIopsWriteRateMax)
	d.Set("iops_write_rate_max_length", r.DiskIopsWriteRateMaxLength)
	d.Set("provisioning_type", r.Provisioningtype)
	d.Set("storage_type", r.Storagetype)
	d.Set("tags", r.Tags)
	d.Set("zone_id", r.Zoneid)

	//
	hypervisor := make(map[string]interface{})
	hypervisor["bytes_read_rate"] = r.DiskBytesReadRate
	hypervisor["bytes_read_rate_max"] = r.DiskBytesReadRateMax
	hypervisor["bytes_read_rate_max_length"] = r.DiskBytesReadRateMaxLength
	hypervisor["bytes_write_rate"] = r.DiskBytesWriteRate
	hypervisor["bytes_write_rate_max"] = r.DiskBytesWriteRateMax
	hypervisor["bytes_write_rate_max_length"] = r.DiskBytesWriteRateMaxLength
	d.Set("hypervisor", hypervisor)

	//
	storage := make(map[string]interface{})
	storage["min_iops"] = r.Miniops
	storage["max_iops"] = r.Maxiops
	storage["customized_iops"] = r.Iscustomizediops
	d.Set("storage", storage)

	return nil

}
func resourceCloudStackDiskOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.DiskOffering.NewUpdateDiskOfferingParams(d.Id())

	if v, ok := d.GetOk("cache_mode"); ok {
		p.SetCachemode(v.(string))
	}
	if v, ok := d.GetOk("display_text"); ok {
		p.SetDisplaytext(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("iops_read_rate"); ok {
		p.SetIopsreadrate(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_read_rate_max"); ok {
		p.SetIopsreadratemax(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_read_rate_max_length"); ok {
		p.SetIopsreadratemaxlength(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate"); ok {
		p.SetIopsreadrate(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max"); ok {
		p.SetIopsreadratemax(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max_length"); ok {
		p.SetIopsreadratemaxlength(int64(v.(int)))
	}
	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(v.(string))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		p.SetZoneid(fmt.Sprintf("%v", v))
	}

	// hypervisor qos
	if v, ok := d.GetOk("hypervisor"); ok {
		hypervisor := v.(map[string]interface{})

		if v2, ok2 := hypervisor["bytes_read_rate"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_read_rate_max"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_read_rate_max_length"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_write_rate"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_write_rate_max"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
		if v2, ok2 := hypervisor["bytes_write_rate_max_length"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetBytesreadrate(int64(tmp))
		}
	}

	return resourceCloudStackDiskOfferingRead(d, meta)

}

func resourceCloudStackDiskOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.DiskOffering.DeleteDiskOffering(cs.DiskOffering.NewDeleteDiskOfferingParams(d.Id()))
	if err != nil {
		return err
	}

	return nil
}
