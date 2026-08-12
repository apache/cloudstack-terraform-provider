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
	"reflect"
	"regexp"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackVgpuProfile() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackVgpuProfileRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "the name of the vGPU profile",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "the description of the vGPU profile",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"device_id": {
				Description: "the device id of the GPU card",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"device_name": {
				Description: "the device name of the GPU card",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"gpu_card_id": {
				Description: "the GPU card id of the vGPU profile",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"gpu_card_name": {
				Description: "the GPU card name of the vGPU profile",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"max_heads": {
				Description: "the maximum displays per vGPU instance",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"max_resolution_x": {
				Description: "the maximum X resolution per display",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"max_resolution_y": {
				Description: "the maximum Y resolution per display",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"max_vgpu_per_physical_gpu": {
				Description: "the maximum number of vGPU instances per physical GPU",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"vendor_id": {
				Description: "the vendor id of the GPU card",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vendor_name": {
				Description: "the vendor name of the GPU card",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"video_ram": {
				Description: "the video RAM size in MB for the vGPU profile",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func datasourceCloudStackVgpuProfileRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.GPU.NewListVgpuProfilesParams()

	csVgpuProfiles, err := cs.GPU.ListVgpuProfiles(p)
	if err != nil {
		return fmt.Errorf("failed to list vGPU profiles: %s", err)
	}

	filters := d.Get("filter")

	for _, profile := range csVgpuProfiles.VgpuProfiles {
		match, err := applyVgpuProfileFilters(profile, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			return vgpuProfileDescriptionAttributes(d, profile)
		}
	}

	return fmt.Errorf("no vGPU profiles found")
}

func vgpuProfileDescriptionAttributes(d *schema.ResourceData, profile *cloudstack.VgpuProfile) error {
	d.SetId(profile.Id)

	fields := map[string]interface{}{
		"id":                        profile.Id,
		"name":                      profile.Name,
		"description":               profile.Description,
		"device_id":                 profile.Deviceid,
		"device_name":               profile.Devicename,
		"gpu_card_id":               profile.Gpucardid,
		"gpu_card_name":             profile.Gpucardname,
		"max_heads":                 profile.Maxheads,
		"max_resolution_x":          profile.Maxresolutionx,
		"max_resolution_y":          profile.Maxresolutiony,
		"max_vgpu_per_physical_gpu": profile.Maxvgpuperphysicalgpu,
		"vendor_id":                 profile.Vendorid,
		"vendor_name":               profile.Vendorname,
		"video_ram":                 profile.Videoram,
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			log.Printf("[WARN] Error setting %s: %s", k, err)
		}
	}

	return nil
}

func applyVgpuProfileFilters(profile *cloudstack.VgpuProfile, filters *schema.Set) (bool, error) {
	val := reflect.ValueOf(profile).Elem()

	for _, f := range filters.List() {
		filter := f.(map[string]interface{})
		r, err := regexp.Compile(filter["value"].(string))
		if err != nil {
			return false, fmt.Errorf("invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(filter["name"].(string), "_", "")
		profileField := val.FieldByNameFunc(func(fieldName string) bool {
			if strings.EqualFold(fieldName, updatedName) {
				updatedName = fieldName
				return true
			}
			return false
		}).String()

		if !r.MatchString(profileField) {
			return false, nil
		}
	}

	return true, nil
}
