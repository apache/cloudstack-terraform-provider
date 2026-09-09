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

func resourceCloudStackDiskOffering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackDiskOfferingCreate,
		Read:   resourceCloudStackDiskOfferingRead,
		Update: resourceCloudStackDiskOfferingUpdate,
		Delete: resourceCloudStackDiskOfferingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the disk offering",
				Type:        schema.TypeString,
				Required:    true,
			},
			"display_text": {
				Description: "The display text of the disk offering",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cache_mode": {
				Description: "The cache mode to use for this disk offering. Values are none, writeback or writethrough",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"disk_size": {
				Description: "The size of the disk offering in GB.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"disk_offering_strictness": {
				Description: "Whether the disk offering size is strictly enforced and cannot be changed at deployment time",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"domain_id": {
				Description: "The IDs of the domains that can use this disk offering",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"iops_read_rate": {
				Description: "The IOPS read rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"iops_read_rate_max": {
				Description: "The maximum IOPS read rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"iops_read_rate_max_length": {
				Description: "The length (in seconds) of the maximum IOPS read rate burst",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"iops_write_rate": {
				Description: "The IOPS write rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"iops_write_rate_max": {
				Description: "The maximum IOPS write rate of the disk offering",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"iops_write_rate_max_length": {
				Description: "The length (in seconds) of the maximum IOPS write rate burst",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"provisioning_type": {
				Description: "Provisioning type used to create volumes. Values are thin, sparse and fat",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "thin" || v == "sparse" || v == "fat" {
						return
					}
					errs = append(errs, fmt.Errorf("provisioning type should be one of thin, sparse or fat, got %s", v))
					return
				},
			},
			"storage_type": {
				Description: "The storage type of the disk offering. Values are local and shared",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "local" || v == "shared" {
						return
					}
					errs = append(errs, fmt.Errorf("storage type should be either local or shared, got %s", v))
					return
				},
			},
			"tags": {
				Description: "The storage tags for the disk offering",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"display_offering": {
				Description: "Whether the disk offering is displayed to the end user",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"zone_id": {
				Description: "The IDs of the zones the disk offering is restricted to. Leave empty to make the offering available in all zones.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"hypervisor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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

	name := d.Get("name").(string)
	displayText := d.Get("display_text").(string)

	// Create a new parameter struct
	p := cs.DiskOffering.NewCreateDiskOfferingParams(displayText, name)

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
		domainIDs := v.([]interface{})
		items := make([]string, len(domainIDs))
		for i, raw := range domainIDs {
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
		p.SetIopswriterate(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max"); ok {
		p.SetIopswriteratemax(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max_length"); ok {
		p.SetIopswriteratemaxlength(int64(v.(int)))
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
	// display_offering defaults to true, so read it directly rather than via GetOk
	p.SetDisplayoffering(d.Get("display_offering").(bool))
	if v, ok := d.GetOk("zone_id"); ok {
		zoneIDs := v.([]interface{})
		items := make([]string, len(zoneIDs))
		for i, raw := range zoneIDs {
			items[i] = raw.(string)
		}
		p.SetZoneid(items)
	}

	// storage qos
	if v, ok := d.GetOk("storage"); ok {
		storageList := v.([]interface{})
		if len(storageList) > 0 && storageList[0] != nil {
			storage := storageList[0].(map[string]interface{})

			if v2, ok2 := storage["min_iops"]; ok2 {
				p.SetMiniops(int64(v2.(int)))
			}
			if v2, ok2 := storage["max_iops"]; ok2 {
				p.SetMaxiops(int64(v2.(int)))
			}
			if v2, ok2 := storage["customized_iops"]; ok2 {
				p.SetCustomizediops(v2.(bool))
			}
			if v2, ok2 := storage["hypervisor_snapshot_reserve"]; ok2 {
				p.SetHypervisorsnapshotreserve(v2.(int))
			}
		}
	}

	// hypervisor qos
	if v, ok := d.GetOk("hypervisor"); ok {
		hypervisorList := v.([]interface{})
		if len(hypervisorList) > 0 && hypervisorList[0] != nil {
			hypervisor := hypervisorList[0].(map[string]interface{})

			if v2, ok2 := hypervisor["bytes_read_rate"]; ok2 {
				p.SetBytesreadrate(int64(v2.(int)))
			}
			if v2, ok2 := hypervisor["bytes_read_rate_max"]; ok2 {
				p.SetBytesreadratemax(int64(v2.(int)))
			}
			if v2, ok2 := hypervisor["bytes_read_rate_max_length"]; ok2 {
				p.SetBytesreadratemaxlength(int64(v2.(int)))
			}
			if v2, ok2 := hypervisor["bytes_write_rate"]; ok2 {
				p.SetByteswriterate(int64(v2.(int)))
			}
			if v2, ok2 := hypervisor["bytes_write_rate_max"]; ok2 {
				p.SetByteswriteratemax(int64(v2.(int)))
			}
			if v2, ok2 := hypervisor["bytes_write_rate_max_length"]; ok2 {
				p.SetByteswriteratemaxlength(int64(v2.(int)))
			}
		}
	}

	log.Printf("[DEBUG] Creating Disk Offering %s", name)
	diskOff, err := cs.DiskOffering.CreateDiskOffering(p)
	if err != nil {
		return fmt.Errorf("Error creating Disk Offering %s: %s", name, err)
	}

	d.SetId(diskOff.Id)

	return resourceCloudStackDiskOfferingRead(d, meta)
}

func resourceCloudStackDiskOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Retrieving Disk Offering %s", d.Id())

	r, count, err := cs.DiskOffering.GetDiskOfferingByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Disk Offering %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", r.Name)
	d.Set("display_text", r.Displaytext)
	d.Set("cache_mode", r.CacheMode)
	d.Set("disk_size", int(r.Disksize))
	d.Set("disk_offering_strictness", r.Disksizestrictness)
	d.Set("iops_read_rate", r.DiskIopsReadRate)
	d.Set("iops_read_rate_max", r.DiskIopsReadRateMax)
	d.Set("iops_read_rate_max_length", r.DiskIopsReadRateMaxLength)
	d.Set("iops_write_rate", r.DiskIopsWriteRate)
	d.Set("iops_write_rate_max", r.DiskIopsWriteRateMax)
	d.Set("iops_write_rate_max_length", r.DiskIopsWriteRateMaxLength)
	d.Set("provisioning_type", r.Provisioningtype)
	d.Set("storage_type", r.Storagetype)
	d.Set("tags", r.Tags)
	d.Set("display_offering", r.Displayoffering)

	// domainid and zoneid are returned as comma-separated strings
	if r.Domainid != "" {
		d.Set("domain_id", strings.Split(r.Domainid, ","))
	} else {
		d.Set("domain_id", []string{})
	}
	if r.Zoneid != "" && r.Zoneid != "all" {
		d.Set("zone_id", strings.Split(r.Zoneid, ","))
	} else {
		d.Set("zone_id", []string{})
	}

	// Only emit the hypervisor block when the API returns non-default QoS values,
	// otherwise leave it null so configs that omit the block don't show perpetual drift.
	if r.DiskBytesReadRate > 0 || r.DiskBytesReadRateMax > 0 || r.DiskBytesReadRateMaxLength > 0 ||
		r.DiskBytesWriteRate > 0 || r.DiskBytesWriteRateMax > 0 || r.DiskBytesWriteRateMaxLength > 0 {
		hypervisor := make(map[string]interface{})
		hypervisor["bytes_read_rate"] = r.DiskBytesReadRate
		hypervisor["bytes_read_rate_max"] = r.DiskBytesReadRateMax
		hypervisor["bytes_read_rate_max_length"] = r.DiskBytesReadRateMaxLength
		hypervisor["bytes_write_rate"] = r.DiskBytesWriteRate
		hypervisor["bytes_write_rate_max"] = r.DiskBytesWriteRateMax
		hypervisor["bytes_write_rate_max_length"] = r.DiskBytesWriteRateMaxLength
		d.Set("hypervisor", []interface{}{hypervisor})
	} else {
		d.Set("hypervisor", []interface{}{})
	}

	// Only emit the storage block when the API returns non-default QoS values.
	if r.Miniops > 0 || r.Maxiops > 0 || r.Iscustomizediops || r.Hypervisorsnapshotreserve > 0 {
		storage := make(map[string]interface{})
		storage["min_iops"] = r.Miniops
		storage["max_iops"] = r.Maxiops
		storage["customized_iops"] = r.Iscustomizediops
		storage["hypervisor_snapshot_reserve"] = r.Hypervisorsnapshotreserve
		d.Set("storage", []interface{}{storage})
	} else {
		d.Set("storage", []interface{}{})
	}

	return nil
}
func resourceCloudStackDiskOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	// Create a new parameter struct
	p := cs.DiskOffering.NewUpdateDiskOfferingParams(d.Id())

	p.SetName(name)
	p.SetDisplaytext(d.Get("display_text").(string))
	p.SetDisplayoffering(d.Get("display_offering").(bool))

	if v, ok := d.GetOk("cache_mode"); ok {
		p.SetCachemode(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		domainIDs := v.([]interface{})
		items := make([]string, len(domainIDs))
		for i, raw := range domainIDs {
			items[i] = raw.(string)
		}
		p.SetDomainid(strings.Join(items, ","))
	}
	if v, ok := d.GetOk("zone_id"); ok {
		zoneIDs := v.([]interface{})
		items := make([]string, len(zoneIDs))
		for i, raw := range zoneIDs {
			items[i] = raw.(string)
		}
		p.SetZoneid(strings.Join(items, ","))
	} else {
		// Empty config means all zones; send "all" so update clears any prior zone restriction.
		p.SetZoneid("all")
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
		p.SetIopswriterate(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max"); ok {
		p.SetIopswriteratemax(int64(v.(int)))
	}
	if v, ok := d.GetOk("iops_write_rate_max_length"); ok {
		p.SetIopswriteratemaxlength(int64(v.(int)))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(v.(string))
	}

	log.Printf("[DEBUG] Updating Disk Offering %s", name)
	_, err := cs.DiskOffering.UpdateDiskOffering(p)
	if err != nil {
		return fmt.Errorf("Error updating Disk Offering %s: %s", name, err)
	}

	return resourceCloudStackDiskOfferingRead(d, meta)
}

func resourceCloudStackDiskOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Deleting Disk Offering %s", d.Get("name").(string))
	_, err := cs.DiskOffering.DeleteDiskOffering(cs.DiskOffering.NewDeleteDiskOfferingParams(d.Id()))
	if err != nil {
		return fmt.Errorf("Error deleting Disk Offering: %s", err)
	}

	return nil
}
