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

	csGpuCards, err := cs.GPU.ListGpuCards(p)
	if err != nil {
		return fmt.Errorf("failed to list GPU cards: %s", err)
	}

	filters := d.Get("filter")

	for _, card := range csGpuCards.GpuCards {
		match, err := applyGpuCardFilters(card, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			return gpuCardDescriptionAttributes(d, card)
		}
	}

	return fmt.Errorf("no GPU cards found")
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

func applyGpuCardFilters(card *cloudstack.GpuCard, filters *schema.Set) (bool, error) {
	val := reflect.ValueOf(card).Elem()

	for _, f := range filters.List() {
		filter := f.(map[string]interface{})
		r, err := regexp.Compile(filter["value"].(string))
		if err != nil {
			return false, fmt.Errorf("invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(filter["name"].(string), "_", "")
		cardField := val.FieldByNameFunc(func(fieldName string) bool {
			if strings.EqualFold(fieldName, updatedName) {
				updatedName = fieldName
				return true
			}
			return false
		}).String()

		if !r.MatchString(cardField) {
			return false, nil
		}
	}

	return true, nil
}
