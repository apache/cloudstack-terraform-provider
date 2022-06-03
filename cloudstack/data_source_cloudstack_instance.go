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
	"log"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCloudstackInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudstackInstanceRead,
		Schema: map[string]*schema.Schema{

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed values
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
			// "ip_address": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
		},
	}
}

func dataSourceCloudstackInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Instance Data Source Read Started")

	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.VirtualMachine.NewListVirtualMachinesParams()
	instance_id := d.Get("instance_id").(string)

	log.Printf("Instance ID =====================> %v", instance_id)

	csInstances, _ := cs.VirtualMachine.ListVirtualMachines(p)

	for _, instance := range csInstances.VirtualMachines {
		log.Printf("Instance ===============> %v", instance)
		if instance.Id == instance_id {
			d.SetId(instance.Id)
			d.Set("instance_id", instance.Id)
			d.Set("account", instance.Account)
			d.Set("created", instance.Created)
			d.Set("display_name", instance.Displayname)
			d.Set("state", instance.State)
			d.Set("host_id", instance.Hostid)
			d.Set("zone_id", instance.Zoneid)

			log.Printf("=========================================================================")

			log.Printf("Instance-ID being Set To: ===============> %v", d.Get("instance_id").(string))
			log.Printf("Instance-Created being Set To: ===============> %v", d.Get("created").(string))
			log.Printf("Instance-Display-Name being Set To: ===============> %v", d.Get("display_name").(string))
			log.Printf("Instance-State being Set To: ===============> %v", d.Get("state").(string))
			log.Printf("Instance-Host-ID being Set To: ===============> %v", d.Get("host_id").(string))
			log.Printf("Instance-Zone-ID being Set To: ===============> %v", d.Get("zone_id").(string))

			log.Printf("=========================================================================")
		}
	}
	log.Printf("Instance Data Source Read Ended")

	return nil
}
