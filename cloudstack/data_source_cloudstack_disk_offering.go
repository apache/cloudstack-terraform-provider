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

func dataSourceCloudstackDiskOffering() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackDiskOfferingRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			// Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_offering": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func datasourceCloudStackDiskOfferingRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.DiskOffering.NewListDiskOfferingsParams()
	csDiskOfferings, err := cs.DiskOffering.ListDiskOfferings(p)
	if err != nil {
		return fmt.Errorf("Failed to list disk offerings: %s", err)
	}

	filters := d.Get("filter")
	var diskOfferings []*cloudstack.DiskOffering

	for _, o := range csDiskOfferings.DiskOfferings {
		match, err := applyDiskOfferingFilters(o, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			diskOfferings = append(diskOfferings, o)
		}
	}

	if len(diskOfferings) == 0 {
		return fmt.Errorf("No disk offering is matching with the specified regex")
	}

	diskOffering, err := latestDiskOffering(diskOfferings)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected disk offering: %s\n", diskOffering.Displaytext)

	return diskOfferingDescriptionAttributes(d, diskOffering)
}

func diskOfferingDescriptionAttributes(d *schema.ResourceData, diskOffering *cloudstack.DiskOffering) error {
	d.SetId(diskOffering.Id)
	d.Set("name", diskOffering.Name)
	d.Set("display_text", diskOffering.Displaytext)
	d.Set("disk_size", int(diskOffering.Disksize))
	d.Set("customized", diskOffering.Iscustomized)
	d.Set("storage_type", diskOffering.Storagetype)
	d.Set("provisioning_type", diskOffering.Provisioningtype)
	d.Set("tags", diskOffering.Tags)
	d.Set("display_offering", diskOffering.Displayoffering)

	return nil
}

func latestDiskOffering(diskOfferings []*cloudstack.DiskOffering) (*cloudstack.DiskOffering, error) {
	var latest time.Time
	var diskOffering *cloudstack.DiskOffering

	for _, o := range diskOfferings {
		created, err := time.Parse("2006-01-02T15:04:05-0700", o.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a disk offering: %s", err)
		}

		if created.After(latest) {
			latest = created
			diskOffering = o
		}
	}

	return diskOffering, nil
}

func applyDiskOfferingFilters(diskOffering *cloudstack.DiskOffering, filters *schema.Set) (bool, error) {
	var diskOfferingJSON map[string]interface{}
	k, _ := json.Marshal(diskOffering)
	err := json.Unmarshal(k, &diskOfferingJSON)
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
		raw, ok := diskOfferingJSON[updatedName]
		if !ok {
			return false, fmt.Errorf("Unknown filter field %s", m["name"].(string))
		}
		diskOfferingField := fmt.Sprint(raw)
		if !r.MatchString(diskOfferingField) {
			return false, nil
		}
	}
	return true, nil
}
