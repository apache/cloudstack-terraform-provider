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
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackServiceOffering() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackServiceOfferingRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_speed": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_customized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_system": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_volatile": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"limit_cpu_use": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_vm_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_planner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"offer_ha": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hypervisor_snapshot_reserve": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_iops": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_bytes_read_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_bytes_write_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_iops_read_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_iops_write_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"root_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gpu_card_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_card_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gpu_display": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_use": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dynamic_scaling_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encrypt_root": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_annotations": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_customized_iops": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lease_duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lease_expiry_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_heads": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_resolution_x": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_resolution_y": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_bytes_read_rate_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_bytes_write_rate_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_iops_read_rate_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disk_iops_write_rate_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func datasourceCloudStackServiceOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.ServiceOffering.NewListServiceOfferingsParams()
	csServiceOfferings, err := cs.ServiceOffering.ListServiceOfferings(p)

	if err != nil {
		return fmt.Errorf("Failed to list service offerings: %s", err)
	}

	filters := d.Get("filter")
	var serviceOfferings []*cloudstack.ServiceOffering

	for _, s := range csServiceOfferings.ServiceOfferings {
		match, err := applyServiceOfferingFilters(s, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			serviceOfferings = append(serviceOfferings, s)
		}
	}

	if len(serviceOfferings) == 0 {
		return fmt.Errorf("No service offering is matching with the specified regex")
	}
	//return the latest service offering from the list of filtered service according
	//to its creation date
	serviceOffering, err := latestServiceOffering(serviceOfferings)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected service offerings: %s\n", serviceOffering.Displaytext)

	return serviceOfferingDescriptionAttributes(d, serviceOffering)
}

func serviceOfferingDescriptionAttributes(d *schema.ResourceData, serviceOffering *cloudstack.ServiceOffering) error {
	d.SetId(serviceOffering.Id)
	d.Set("name", serviceOffering.Name)
	d.Set("display_text", serviceOffering.Displaytext)
	d.Set("cpu_number", serviceOffering.Cpunumber)
	d.Set("cpu_speed", serviceOffering.Cpuspeed)
	d.Set("memory", serviceOffering.Memory)
	d.Set("created", serviceOffering.Created)
	d.Set("domain_id", serviceOffering.Domainid)
	d.Set("domain", serviceOffering.Domain)
	d.Set("host_tags", serviceOffering.Hosttags)
	d.Set("is_customized", serviceOffering.Iscustomized)
	d.Set("is_system", serviceOffering.Issystem)
	d.Set("is_volatile", serviceOffering.Isvolatile)
	d.Set("limit_cpu_use", serviceOffering.Limitcpuuse)
	d.Set("network_rate", serviceOffering.Networkrate)
	d.Set("storage_type", serviceOffering.Storagetype)
	d.Set("system_vm_type", serviceOffering.Systemvmtype)
	d.Set("deployment_planner", serviceOffering.Deploymentplanner)
	d.Set("offer_ha", serviceOffering.Offerha)
	d.Set("tags", serviceOffering.Storagetags)
	d.Set("provisioning_type", serviceOffering.Provisioningtype)

	// IOPS limits - only set if returned by API (> 0)
	if serviceOffering.Miniops > 0 {
		d.Set("min_iops", int(serviceOffering.Miniops))
	}
	if serviceOffering.Maxiops > 0 {
		d.Set("max_iops", int(serviceOffering.Maxiops))
	}

	d.Set("hypervisor_snapshot_reserve", serviceOffering.Hypervisorsnapshotreserve)
	d.Set("disk_bytes_read_rate", serviceOffering.DiskBytesReadRate)
	d.Set("disk_bytes_write_rate", serviceOffering.DiskBytesWriteRate)
	d.Set("disk_iops_read_rate", serviceOffering.DiskIopsReadRate)
	d.Set("disk_iops_write_rate", serviceOffering.DiskIopsWriteRate)
	d.Set("root_disk_size", serviceOffering.Rootdisksize)
	d.Set("gpu_card_id", serviceOffering.Gpucardid)
	d.Set("gpu_card_name", serviceOffering.Gpucardname)
	d.Set("gpu_count", serviceOffering.Gpucount)
	d.Set("gpu_display", serviceOffering.Gpudisplay)

	// New fields
	d.Set("default_use", serviceOffering.Defaultuse)
	d.Set("dynamic_scaling_enabled", serviceOffering.Dynamicscalingenabled)
	d.Set("encrypt_root", serviceOffering.Encryptroot)
	d.Set("has_annotations", serviceOffering.Hasannotations)
	d.Set("is_customized_iops", serviceOffering.Iscustomizediops)
	d.Set("lease_duration", serviceOffering.Leaseduration)
	d.Set("lease_expiry_action", serviceOffering.Leaseexpiryaction)
	d.Set("max_heads", serviceOffering.Maxheads)
	d.Set("max_resolution_x", serviceOffering.Maxresolutionx)
	d.Set("max_resolution_y", serviceOffering.Maxresolutiony)
	d.Set("disk_bytes_read_rate_max", serviceOffering.DiskBytesReadRateMax)
	d.Set("disk_bytes_write_rate_max", serviceOffering.DiskBytesWriteRateMax)
	d.Set("disk_iops_read_rate_max", serviceOffering.DiskIopsReadRateMax)
	d.Set("disk_iops_write_rate_max", serviceOffering.DiskIopsWriteRateMax)

	return nil
}

func latestServiceOffering(serviceOfferings []*cloudstack.ServiceOffering) (*cloudstack.ServiceOffering, error) {
	var latest time.Time
	var serviceOffering *cloudstack.ServiceOffering

	for _, s := range serviceOfferings {
		created, err := time.Parse("2006-01-02T15:04:05-0700", s.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of an service offering: %s", err)
		}

		if created.After(latest) {
			latest = created
			serviceOffering = s
		}
	}

	return serviceOffering, nil
}

func applyServiceOfferingFilters(serviceOffering *cloudstack.ServiceOffering, filters *schema.Set) (bool, error) {
	var serviceOfferingJSON map[string]interface{}
	k, _ := json.Marshal(serviceOffering)
	err := json.Unmarshal(k, &serviceOfferingJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		serviceOfferingField := serviceOfferingJSON[updatedName].(string)
		if !r.MatchString(serviceOfferingField) {
			return false, nil
		}

	}
	return true, nil
}
