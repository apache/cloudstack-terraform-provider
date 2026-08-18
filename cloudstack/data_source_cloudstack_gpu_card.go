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
	"strconv"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackGpuCard() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackGpuCardRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "the name of the GPU card",
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
		},
	}
}

func datasourceCloudStackGpuCardRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.GPU.NewListGpuCardsParams()

	if err := applyGpuCardFilters(p, d.Get("filter").(*schema.Set)); err != nil {
		return err
	}

	csGpuCards, err := cs.GPU.ListGpuCards(p)
	if err != nil {
		return fmt.Errorf("failed to list GPU cards: %s", err)
	}

	switch len(csGpuCards.GpuCards) {
	case 0:
		return fmt.Errorf("no GPU cards found")
	case 1:
		return gpuCardDescriptionAttributes(d, csGpuCards.GpuCards[0])
	default:
		return fmt.Errorf("%d GPU cards matched the given filters; "+
			"refine the filters (e.g. add device_id or vendor_id) to match exactly one card", len(csGpuCards.GpuCards))
	}
}

func gpuCardDescriptionAttributes(d *schema.ResourceData, card *cloudstack.GpuCard) error {
	d.SetId(card.Id)

	fields := map[string]interface{}{
		"id":          card.Id,
		"name":        card.Name,
		"device_id":   card.Deviceid,
		"device_name": card.Devicename,
		"vendor_id":   card.Vendorid,
		"vendor_name": card.Vendorname,
	}

	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			log.Printf("[WARN] Error setting %s: %s", k, err)
		}
	}

	return nil
}

func applyGpuCardFilters(p *cloudstack.ListGpuCardsParams, filters *schema.Set) error {
	seen := make(map[string]bool)
	for _, f := range filters.List() {
		filter := f.(map[string]interface{})
		name := filter["name"].(string)
		value := filter["value"].(string)

		if seen[name] {
			return fmt.Errorf("duplicate filter %q; each filter name may only be specified once", name)
		}
		seen[name] = true

		switch name {
		case "id":
			p.SetId(value)
		case "device_id":
			p.SetDeviceid(value)
		case "device_name":
			p.SetDevicename(value)
		case "vendor_id":
			p.SetVendorid(value)
		case "vendor_name":
			p.SetVendorname(value)
		case "keyword":
			p.SetKeyword(value)
		case "active_only":
			b, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid boolean value %q for filter %q: %s", value, name, err)
			}
			p.SetActiveonly(b)
		default:
			return fmt.Errorf("unsupported filter %q; supported filters: id, device_id, device_name, vendor_id, vendor_name, keyword, active_only", name)
		}
	}

	return nil
}
