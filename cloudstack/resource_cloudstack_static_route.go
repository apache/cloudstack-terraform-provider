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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"nexthop", "vpc_id"},
			},

			"nexthop": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"gateway_id"},
				RequiredWith:  []string{"vpc_id"},
			},

			"vpc_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"gateway_id"},
				RequiredWith:  []string{"nexthop"},
			},
		},
	}
}

func resourceCloudStackStaticRouteCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Verify that required parameters are set
	if err := verifyStaticRouteParams(d); err != nil {
		return err
	}

	// Create a new parameter struct
	p := cs.VPC.NewCreateStaticRouteParams(
		d.Get("cidr").(string),
	)

	// Set either gateway_id or nexthop+vpc_id (they are mutually exclusive)
	if v, ok := d.GetOk("gateway_id"); ok {
		p.SetGatewayid(v.(string))
	}

	if v, ok := d.GetOk("nexthop"); ok {
		p.SetNexthop(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		p.SetVpcid(v.(string))
	}

	// Create the new static route
	r, err := cs.VPC.CreateStaticRoute(p)
	if err != nil {
		return fmt.Errorf("Error creating static route for %s: %s", d.Get("cidr").(string), err)
	}

	d.SetId(r.Id)

	return resourceCloudStackStaticRouteRead(d, meta)
}

func resourceCloudStackStaticRouteRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Get the static route details
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

	// Set gateway_id if it's not empty (indicates this route uses a gateway)
	if r.Vpcgatewayid != "" {
		d.Set("gateway_id", r.Vpcgatewayid)
	}

	// Set nexthop and vpc_id if nexthop is not empty (indicates this route uses nexthop)
	if r.Nexthop != "" {
		d.Set("nexthop", r.Nexthop)
		if r.Vpcid != "" {
			d.Set("vpc_id", r.Vpcid)
		}
	}

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

func verifyStaticRouteParams(d *schema.ResourceData) error {
	_, hasGatewayID := d.GetOk("gateway_id")
	_, hasNexthop := d.GetOk("nexthop")
	_, hasVpcID := d.GetOk("vpc_id")

	// Check that either gateway_id or (nexthop + vpc_id) is provided
	if !hasGatewayID && !hasNexthop {
		return fmt.Errorf(
			"You must supply either 'gateway_id' or 'nexthop' (with 'vpc_id')")
	}

	// Check that nexthop and vpc_id are used together
	if hasNexthop && !hasVpcID {
		return fmt.Errorf(
			"You must supply 'vpc_id' when using 'nexthop'")
	}

	if hasVpcID && !hasNexthop {
		return fmt.Errorf(
			"You must supply 'nexthop' when using 'vpc_id'")
	}

	return nil
}
