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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackStaticRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackStaticRouteCreate,
		Read:   resourceCloudStackStaticRouteRead,
		Delete: resourceCloudStackStaticRouteDelete,

		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudStackStaticRouteCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.VPC.NewCreateStaticRouteParams(
		d.Get("cidr").(string),
		d.Get("gateway_id").(string),
	)

	// Create the new private gateway
	r, err := cs.VPC.CreateStaticRoute(p)
	if err != nil {
		return fmt.Errorf("Error creating static route for %s: %s", d.Get("cidr").(string), err)
	}

	d.SetId(r.Id)

	return resourceCloudStackStaticRouteRead(d, meta)
}

func resourceCloudStackStaticRouteRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the virtual machine details
	r, count, err := cs.VPC.GetStaticRouteByID(d.Id())
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Static route %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("cidr", r.Cidr)

	return nil
}

func resourceCloudStackStaticRouteDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.VPC.NewDeleteStaticRouteParams(d.Id())

	// Delete the private gateway
	_, err := cs.VPC.DeleteStaticRoute(p)
	if err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error deleting static route for %s: %s", d.Get("cidr").(string), err)
	}

	return nil
}
