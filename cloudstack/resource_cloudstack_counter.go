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

func resourceCloudStackCounter() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackCounterCreate,
		Read:   resourceCloudStackCounterRead,
		Delete: resourceCloudStackCounterDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the counter",
				ForceNew:    true,
			},

			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source of the counter",
				ForceNew:    true,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value of the counter e.g. oid in case of snmp",
				ForceNew:    true,
			},
			"counter_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Provider of the counter",
				ForceNew:    true,
			},
		},
	}
}

func resourceCloudStackCounterCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	source := d.Get("source").(string)
	value := d.Get("value").(string)
	provider := d.Get("counter_provider").(string)

	p := cs.AutoScale.NewCreateCounterParams(name, provider, source, value)

	log.Printf("[DEBUG] Creating counter: %s", name)
	resp, err := cs.AutoScale.CreateCounter(p)
	if err != nil {
		return fmt.Errorf("Error creating counter: %s", err)
	}

	d.SetId(resp.Id)
	log.Printf("[DEBUG] Counter created with ID: %s", resp.Id)

	return resourceCloudStackCounterRead(d, meta)
}

func resourceCloudStackCounterRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewListCountersParams()
	p.SetId(d.Id())

	resp, err := cs.AutoScale.ListCounters(p)
	if err != nil {
		return fmt.Errorf("Error retrieving counter: %s", err)
	}

	if resp.Count == 0 {
		log.Printf("[DEBUG] Counter %s no longer exists", d.Id())
		d.SetId("")
		return nil
	}

	counter := resp.Counters[0]
	d.Set("name", counter.Name)
	d.Set("source", counter.Source)
	d.Set("value", counter.Value)
	d.Set("counter_provider", counter.Provider)

	return nil
}

func resourceCloudStackCounterDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.AutoScale.NewDeleteCounterParams(d.Id())

	log.Printf("[DEBUG] Deleting counter: %s", d.Id())
	_, err := cs.AutoScale.DeleteCounter(p)
	if err != nil {
		return fmt.Errorf("Error deleting counter: %s", err)
	}

	return nil
}
