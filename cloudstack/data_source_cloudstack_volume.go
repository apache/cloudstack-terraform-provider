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

func dataSourceCloudstackVolume() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackVolumeRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_offering_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceCloudStackVolumeRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Volume.NewListVolumesParams()
	csVolumes, err := cs.Volume.ListVolumes(p)

	if err != nil {
		return fmt.Errorf("Failed to list volumes: %s", err)
	}

	filters := d.Get("filter")
	var volumes []*cloudstack.Volume

	for _, v := range csVolumes.Volumes {
		match, err := applyVolumeFilters(v, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			volumes = append(volumes, v)
		}
	}

	if len(volumes) == 0 {
		return fmt.Errorf("No volume is matching with the specified regex")
	}
	//return the latest volume from the list of filtered volumes according
	//to its creation date
	volume, err := latestVolume(volumes)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected volume: %s\n", volume.Name)

	return volumeDescriptionAttributes(d, volume)
}

func volumeDescriptionAttributes(d *schema.ResourceData, volume *cloudstack.Volume) error {
	d.SetId(volume.Id)
	d.Set("name", volume.Name)
	d.Set("disk_offering_id", volume.Diskofferingid)
	d.Set("zone_id", volume.Zoneid)

	return nil
}

func latestVolume(volumes []*cloudstack.Volume) (*cloudstack.Volume, error) {
	var latest time.Time
	var volume *cloudstack.Volume

	for _, v := range volumes {
		created, err := time.Parse("2006-01-02T15:04:05-0700", v.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a volume: %s", err)
		}

		if created.After(latest) {
			latest = created
			volume = v
		}
	}

	return volume, nil
}

func applyVolumeFilters(volume *cloudstack.Volume, filters *schema.Set) (bool, error) {
	var volumeJSON map[string]interface{}
	v, _ := json.Marshal(volume)
	err := json.Unmarshal(v, &volumeJSON)
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
		volume := volumeJSON[updatedName].(string)
		if !r.MatchString(volume) {
			return false, nil
		}

	}
	return true, nil
}
