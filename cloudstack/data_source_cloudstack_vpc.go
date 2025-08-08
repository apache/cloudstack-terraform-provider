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

func dataSourceCloudstackVPC() *schema.Resource {
	return &schema.Resource{
		Read: datasourceCloudStackVPCRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			"projectid": {
				Type:     schema.TypeString,
				Required: true,
			},

			//Computed values
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_text": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpc_offering_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func datasourceCloudStackVPCRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.VPC.NewListVPCsParams()
	p.SetProjectid(d.Get("projectid").(string))
	csVPCs, err := cs.VPC.ListVPCs(p)

	if err != nil {
		return fmt.Errorf("Failed to list VPCs: %s", err)
	}

	filters := d.Get("filter")
	var vpcs []*cloudstack.VPC

	for _, v := range csVPCs.VPCs {
		match, err := applyVPCFilters(v, filters.(*schema.Set))
		if err != nil {
			return err
		}
		if match {
			vpcs = append(vpcs, v)
		}
	}

	if len(vpcs) == 0 {
		return fmt.Errorf("No VPC is matching with the specified regex")
	}
	//return the latest VPC from the list of filtered VPCs according
	//to its creation date
	vpc, err := latestVPC(vpcs)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected VPCs: %s\n", vpc.Displaytext)

	return vpcDescriptionAttributes(d, vpc)
}

func vpcDescriptionAttributes(d *schema.ResourceData, vpc *cloudstack.VPC) error {
	d.SetId(vpc.Id)
	d.Set("name", vpc.Name)
	d.Set("display_text", vpc.Displaytext)
	d.Set("cidr", vpc.Cidr)
	d.Set("vpc_offering_name", vpc.Vpcofferingname)
	d.Set("network_domain", vpc.Networkdomain)
	d.Set("project", vpc.Project)
	d.Set("zone_name", vpc.Zonename)
	d.Set("tags", tagsToMap(vpc.Tags))

	return nil
}

func latestVPC(vpcs []*cloudstack.VPC) (*cloudstack.VPC, error) {
	var latest time.Time
	var vpc *cloudstack.VPC

	for _, v := range vpcs {
		created, err := time.Parse("2006-01-02T15:04:05-0700", v.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of a VPC: %s", err)
		}

		if created.After(latest) {
			latest = created
			vpc = v
		}
	}

	return vpc, nil
}

func applyVPCFilters(vpc *cloudstack.VPC, filters *schema.Set) (bool, error) {
	var vpcJSON map[string]interface{}
	k, _ := json.Marshal(vpc)
	err := json.Unmarshal(k, &vpcJSON)
	if err != nil {
		return false, err
	}

	for _, f := range filters.List() {
		m := f.(map[string]interface{})
		log.Print(m)
		r, err := regexp.Compile(m["value"].(string))
		if err != nil {
			return false, fmt.Errorf("Invalid regex: %s", err)
		}
		updatedName := strings.ReplaceAll(m["name"].(string), "_", "")
		log.Print(updatedName)
		vpcField := vpcJSON[updatedName].(string)
		if !r.MatchString(vpcField) {
			return false, nil
		}
	}
	return true, nil
}
