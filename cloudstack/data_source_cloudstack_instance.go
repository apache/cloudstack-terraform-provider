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

func dataSourceCloudstackInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackInstanceRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),

			//Computed values
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),

			"nic": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudstackInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Instance Data Source Read Started")

	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.VirtualMachine.NewListVirtualMachinesParams()
	csInstances, err := cs.VirtualMachine.ListVirtualMachines(p)

	if err != nil {
		return fmt.Errorf("Failed to list instances: %s", err)
	}

	filters := d.Get("filter")
	nic := d.Get("nic").([]interface{})
	var instances []*cloudstack.VirtualMachine

	//the if-else block to check whether to filter the data source by an IP address
	// or by any other exported attributes
	if len(nic) != 0 {
		ip_address := nic[0].(map[string]interface{})["ip_address"]
		for _, i := range csInstances.VirtualMachines {
			if ip_address == i.Nic[0].Ipaddress {
				instances = append(instances, i)
			}
		}
	} else {
		for _, i := range csInstances.VirtualMachines {
			match, err := applyInstanceFilters(i, filters.(*schema.Set))
			if err != nil {
				return err
			}

			if match {
				instances = append(instances, i)
			}
		}
	}

	if len(instances) == 0 {
		return fmt.Errorf("No instance is matching with the specified regex")
	}
	//return the latest instance from the list of filtered instances according
	//to its creation date
	instance, err := latestInstance(instances)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Selected instances: %s\n", instance.Displayname)

	return instanceDescriptionAttributes(d, instance)
}

func instanceDescriptionAttributes(d *schema.ResourceData, instance *cloudstack.VirtualMachine) error {
	d.SetId(instance.Id)
	d.Set("instance_id", instance.Id)
	d.Set("account", instance.Account)
	d.Set("created", instance.Created)
	d.Set("display_name", instance.Displayname)
	d.Set("state", instance.State)
	d.Set("host_id", instance.Hostid)
	d.Set("zone_id", instance.Zoneid)
	d.Set("nic", []interface{}{map[string]string{"ip_address": instance.Nic[0].Ipaddress}})

	tags := make(map[string]interface{})
	for _, tag := range instance.Tags {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func latestInstance(instances []*cloudstack.VirtualMachine) (*cloudstack.VirtualMachine, error) {
	var latest time.Time
	var instance *cloudstack.VirtualMachine

	for _, i := range instances {
		created, err := time.Parse("2006-01-02T15:04:05-0700", i.Created)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse creation date of an instance: %s", err)
		}

		if created.After(latest) {
			latest = created
			instance = i
		}
	}

	return instance, nil
}

func applyInstanceFilters(instance *cloudstack.VirtualMachine, filters *schema.Set) (bool, error) {
	var instanceJSON map[string]interface{}
	i, _ := json.Marshal(instance)
	err := json.Unmarshal(i, &instanceJSON)
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
		instanceField := instanceJSON[updatedName].(string)
		if !r.MatchString(instanceField) {
			return false, nil
		}

	}
	return true, nil
}
