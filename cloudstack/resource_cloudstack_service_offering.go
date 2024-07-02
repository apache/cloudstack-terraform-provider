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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	s, _, err := cs.ServiceOffering.GetServiceOfferingByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", s.Name)
	d.Set("display_text", s.Displaytext)
	d.Set("bytes_read_rate", s.DiskBytesReadRate)
	d.Set("bytes_read_rate_max", s.DiskBytesReadRateMax)
	d.Set("bytes_read_rate_max_length", s.DiskBytesReadRateMaxLength)
	d.Set("bytes_write_rate", s.DiskBytesWriteRate)
	d.Set("bytes_write_rate_max", s.DiskBytesWriteRateMax)
	d.Set("bytes_write_rate_max_length", s.DiskBytesWriteRateMaxLength)
	d.Set("cache_mode", s.CacheMode)
	d.Set("cpu_number", s.Cpunumber)
	d.Set("cpu_speed", s.Cpuspeed)
	d.Set("customized", s.Iscustomized)
	d.Set("customized_iops", s.Iscustomizediops)
	d.Set("deployment_planner", s.Deploymentplanner)
	d.Set("disk_offering_id", s.Diskofferingid)
	d.Set("disk_offering_strictness", s.Diskofferingstrictness)
	d.Set("domain_id", s.Domainid)
	d.Set("dynamic_scaling_enabled", s.Dynamicscalingenabled)
	// Not available in cloudstack client
	// d.Set("encrypt_root", s.EncryptRoot)
	d.Set("host_tags", s.Hosttags)
	d.Set("hypervisor_snapshot_reserve", s.Hypervisorsnapshotreserve)
	d.Set("iops_read_rate", s.DiskIopsReadRate)
	d.Set("iops_read_rate_max", s.DiskIopsReadRateMax)
	d.Set("iops_read_rate_max_length", s.DiskIopsReadRateMaxLength)
	d.Set("iops_write_rate", s.DiskIopsWriteRate)
	d.Set("iops_write_rate_max", s.DiskIopsWriteRateMax)
	d.Set("iops_write_rate_max_length", s.DiskIopsWriteRateMaxLength)
	d.Set("is_system", s.Issystem)
	d.Set("is_volatile", s.Isvolatile)
	d.Set("limit_cpu_use", s.Limitcpuuse)
	d.Set("max_cpu_number", s.Serviceofferingdetails["maxcpunumber"])
	d.Set("max_iops", s.Maxiops)
	d.Set("max_memory", s.Serviceofferingdetails["maxmemory"])
	d.Set("memory", s.Memory)
	d.Set("min_cpu_number", s.Serviceofferingdetails["mincpunumber"])
	d.Set("min_iops", s.Miniops)
	d.Set("min_memory", s.Serviceofferingdetails["minmemory"])
	d.Set("network_rate", s.Networkrate)
	d.Set("offer_ha", s.Offerha)
	d.Set("provisioning_type", s.Provisioningtype)
	d.Set("root_disk_size", s.Rootdisksize)
	d.Set("storage_policy", s.Vspherestoragepolicy)
	d.Set("storage_type", s.Storagetype)
	d.Set("system_vm_type", s.Systemvmtype)
	d.Set("tags", s.Storagetags)
	if len(s.Zoneid) > 0 {
		d.Set("zone_id", strings.Split(s.Zoneid, ","))
	}

	return nil
}

func resourceCloudStackServiceOfferingUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.ServiceOffering.NewUpdateServiceOfferingParams(d.Id())
	if v, ok := d.GetOk("display_text"); ok {
		p.SetDisplaytext(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("host_tags"); ok {
		p.SetHosttags(v.(string))
	}
	if v, ok := d.GetOk("name"); ok {
		p.SetName(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetStoragetags(v.(string))
	}

	if v, ok := d.GetOk("zone_id"); ok {
		zone_id := v.(*schema.Set).List()
		items := make([]string, len(zone_id))
		for i, raw := range zone_id {
			items[i] = raw.(string)
		}
		p.SetZoneid(strings.Join(items, ","))
	} else {
		// Special parameter not documented in spec.
		p.SetZoneid("all")
	}

	_, err := cs.ServiceOffering.UpdateServiceOffering(p)
	if err != nil {
		return err
	}

	return resourceCloudStackServiceOfferingRead(d, meta)
}

func resourceCloudStackServiceOfferingDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.ServiceOffering.NewDeleteServiceOfferingParams(d.Id())
	_, err := cs.ServiceOffering.DeleteServiceOffering(p)

	if err != nil {
		return fmt.Errorf("Error deleting Service Offering: %s", err)
	}

	return nil
}

func serviceOfferingMergeCommonSchema(s1 map[string]*schema.Schema) map[string]*schema.Schema {
	common := map[string]*schema.Schema{
		// required
		"display_text": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		// optional
		"deployment_planner": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"disk_offering_id": {
			Type:     schema.TypeString,
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
		"dynamic_scaling_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"host_tags": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"is_volatile": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"limit_cpu_use": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"network_rate": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
		},
		"offer_ha": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"tags": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"zone_id": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"disk_hypervisor": serviceOfferingDiskQosHypervisor(),
		"disk_offering":   serviceOfferingDisk(),
		"disk_storage":    serviceOfferingDiskQosStorage(),
	}

	for k, v := range s1 {
		common[k] = v
	}

	return common

}

func serviceOfferingCreateParams(p *cloudstack.CreateServiceOfferingParams, d *schema.ResourceData) *cloudstack.CreateServiceOfferingParams {
	// other
	if v, ok := d.GetOk("host_tags"); ok {
		p.SetHosttags(v.(string))
	}
	if v, ok := d.GetOk("network_rate"); ok {
		p.SetNetworkrate(v.(int))
	}
	if v, ok := d.GetOk("deployment_planner"); ok {
		p.SetDeploymentplanner(v.(string))
	}
	if v, ok := d.GetOk("disk_offering_id"); ok {
		p.SetDiskofferingid(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		p.SetTags(v.(string))
	}

	// Features flags
	p.SetDynamicscalingenabled(d.Get("dynamic_scaling_enabled").(bool))
	p.SetIsvolatile(d.Get("is_volatile").(bool))
	p.SetLimitcpuuse(d.Get("limit_cpu_use").(bool))
	p.SetOfferha(d.Get("offer_ha").(bool))

	// access
	if v, ok := d.GetOk("domain_id"); ok {
		domain_id := v.([]interface{})
		items := make([]string, len(domain_id))
		for i, raw := range domain_id {
			items[i] = raw.(string)
		}
		p.SetDomainid(items)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		zone_id := v.(*schema.Set).List()
		items := make([]string, len(zone_id))
		for i, raw := range zone_id {
			items[i] = raw.(string)
		}
		p.SetZoneid(items)
	}

	// disk offering
	if v, ok := d.GetOk("disk_offering"); ok {
		offering := v.(map[string]interface{})

		if v2, ok2 := offering["storage_type"]; ok2 {
			p.SetStoragetype(v2.(string))
		}
		if v2, ok2 := offering["provisioning_type"]; ok2 {
			p.SetProvisioningtype(v2.(string))
		}
		if v2, ok2 := offering["cache_mode"]; ok2 {
			p.SetCachemode(v2.(string))
		}
		if v2, ok2 := offering["root_disk_size"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetRootdisksize(int64(tmp))
		}
		if v2, ok2 := offering["disk_offering_strictness"]; ok2 {
			tmp, _ := strconv.ParseBool(v2.(string))
			p.SetDiskofferingstrictness(tmp)
		}
	}

	// hypervisor qos
	if v, ok := d.GetOk("disk_hypervisor"); ok {
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

	// storage qos
	if v, ok := d.GetOk("disk_storage"); ok {
		storage := v.(map[string]interface{})

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
		if v2, ok2 := storage["hypervisor_snapshot_reserve"]; ok2 {
			tmp, _ := strconv.Atoi(v2.(string))
			p.SetHypervisorsnapshotreserve(tmp)
		}
	}

	return p
}

func serviceOfferingDiskQosHypervisor() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
			},
		},
	}
}

func serviceOfferingDisk() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cache_mode": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"disk_offering_strictness": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
				"provisioning_type": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"root_disk_size": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
				"storage_type": {
					Type:     schema.TypeInt,
					Required: true,
					ForceNew: true,
					Default:  "shared",
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(string)

						if v == "local" || v == "shared" {
							return
						}

						errs = append(errs, fmt.Errorf("storage type should be either local or shared, got %s", v))

						return
					},
				},
			},
		},
	}
}

func serviceOfferingDiskQosStorage() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"storage": {
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
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
							"max_iops": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
								ForceNew: true,
							},
							"min_iops": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
								ForceNew: true,
							},
						},
					},
				},
			},
		},
	}
}
