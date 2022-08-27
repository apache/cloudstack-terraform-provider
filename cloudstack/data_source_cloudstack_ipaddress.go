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
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudstackIPAddress() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackIPAddressRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"is_portable": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_source_nat": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func datasourceCloudStackIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Address.NewListPublicIpAddressesParams()
	csPublicIPAddresses, err := cs.Address.ListPublicIpAddresses(p)

	if err != nil {
		return fmt.Errorf("Failed to list ip addresses: %s", err)
	}

	filters := d.Get("filter")
	var publicIpAddresses []*cloudstack.PublicIpAddress

	for _, ip := range csPublicIPAddresses.PublicIpAddresses {
		match, err := applyIPAddressFilters(ip, filters.(*schema.Set))

		if err != nil {
			return err
		}
		if match {
			publicIpAddresses = append(publicIpAddresses, ip)
		}
	}

	if len(publicIpAddresses) == 0 {
		return fmt.Errorf("No ip address is matching with the specified regex")
	}
	//return the latest ip address from the list of filtered ip addresses according
	//to its creation date
	publicIpAddress, err := latestIPAddress(publicIpAddresses)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected ip addresses: %s\n", publicIpAddress.Ipaddress)

	return ipAddressDescriptionAttributes(d, publicIpAddress)
}

func ipAddressDescriptionAttributes(d *schema.ResourceData, publicIpAddress *cloudstack.PublicIpAddress) error {
	d.SetId(publicIpAddress.Id)
	d.Set("is_portable", publicIpAddress.Isportable)
	d.Set("network_id", publicIpAddress.Networkid)
	d.Set("vpc_id", publicIpAddress.Vpcid)
	d.Set("zone_name", publicIpAddress.Zonename)
	d.Set("project", publicIpAddress.Project)
	d.Set("ip_address", publicIpAddress.Ipaddress)
	d.Set("is_source_nat", publicIpAddress.Issourcenat)
	d.Set("tags", publicIpAddress.Tags)

	return nil
}

func latestIPAddress(publicIpAddresses []*cloudstack.PublicIpAddress) (*cloudstack.PublicIpAddress, error) {
	var latest time.Time
	var publicIpAddress *cloudstack.PublicIpAddress

	for _, ip := range publicIpAddresses {
		created, err := time.Parse("2006-01-02T15:04:05-0700", ip.Allocated)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse allocation date of the ip address: %s", err)
		}

		if created.After(latest) {
			latest = created
			publicIpAddress = ip
		}
	}

	return publicIpAddress, nil
}

func applyIPAddressFilters(publicIpAddress *cloudstack.PublicIpAddress, filters *schema.Set) (bool, error) {
	var publicIPAdressJSON map[string]interface{}
	k, _ := json.Marshal(publicIpAddress)
	err := json.Unmarshal(k, &publicIPAdressJSON)

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
		publicIPAdressField := fmt.Sprintf("%v", publicIPAdressJSON[updatedName])
		if !r.MatchString(publicIPAdressField) {
			return false, nil
		}
	}

	return true, nil
}
