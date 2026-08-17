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

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudstackCounter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackCounterRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"counter_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudstackCounterRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	id, idOk := d.GetOk("id")
	name, nameOk := d.GetOk("name")

	if !idOk && !nameOk {
		return fmt.Errorf("either 'id' or 'name' must be specified")
	}

	var counter *cloudstack.Counter

	if idOk {
		// Get counter by ID
		p := cs.AutoScale.NewListCountersParams()
		p.SetId(id.(string))

		resp, err := cs.AutoScale.ListCounters(p)
		if err != nil {
			return fmt.Errorf("failed to list counters: %s", err)
		}

		if resp.Count == 0 {
			return fmt.Errorf("counter with ID %s not found", id.(string))
		}

		counter = resp.Counters[0]
	} else {
		// Get counter by name
		p := cs.AutoScale.NewListCountersParams()

		resp, err := cs.AutoScale.ListCounters(p)
		if err != nil {
			return fmt.Errorf("failed to list counters: %s", err)
		}

		for _, c := range resp.Counters {
			if c.Name == name.(string) {
				counter = c
				break
			}
		}

		if counter == nil {
			return fmt.Errorf("counter with name %s not found", name.(string))
		}
	}

	log.Printf("[DEBUG] Found counter: %s", counter.Name)

	d.SetId(counter.Id)
	d.Set("name", counter.Name)
	d.Set("source", counter.Source)
	d.Set("value", counter.Value)
	d.Set("counter_provider", counter.Provider)

	return nil
}
