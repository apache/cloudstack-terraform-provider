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
	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackLoadBalancerCreate,
		Read:   resourceCloudStackLoadBalancerRead,
		Delete: resourceCloudStackLoadBalancerDelete,

		Schema: map[string]*schema.Schema{
			"algorithm": {
				Description: "load balancer algorithm (source, roundrobin, leastconn)",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"instanceport": {
				Description: "the TCP port of the virtual machine where the network traffic will be load balanced to",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			"name": {
				Description: "name of the load balancer",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"networkid": {
				Description: "The guest network the load balancer will be created for",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"scheme": {
				Description: "the load balancer scheme. Supported value in this release is Internal",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"sourceipaddressnetworkid": {
				Description: "the network id of the source ip address",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"sourceport": {
				Description: "the source port the network traffic will be load balanced from",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			"description": {
				Description: "the description of the load balancer",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"sourceipaddress": {
				Description: "the source IP address the network traffic will be load balanced from",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"virtualmachineids": {
				Description: "the list of IDs of the virtual machine that are being assigned to the load balancer rule(i.e. virtualMachineIds=1,2,3)",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCloudStackLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.LoadBalancer.NewCreateLoadBalancerParams(
		d.Get("algorithm").(string),
		d.Get("instanceport").(int),
		d.Get("name").(string),
		d.Get("networkid").(string),
		d.Get("scheme").(string),
		d.Get("sourceipaddressnetworkid").(string),
		d.Get("sourceport").(int),
	)
	if v, ok := d.GetOk("description"); ok {
		p.SetDescription(v.(string))
	}
	if v, ok := d.GetOk("sourceipaddress"); ok {
		p.SetSourceipaddress(v.(string))
	}

	r, err := cs.LoadBalancer.CreateLoadBalancer(p)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("virtualmachineids"); ok {
		vmIds := v.(*schema.Set).List()
		for _, vmId := range vmIds {
			p_update := cs.LoadBalancer.NewAssignToLoadBalancerRuleParams(r.Id)
			p_update.SetVirtualmachineids([]string{vmId.(string)})
			_, err := cs.LoadBalancer.AssignToLoadBalancerRule(p_update)
			if err != nil {
				return err
			}
		}
	}

	d.SetId(r.Id)

	return resourceCloudStackLoadBalancerRead(d, meta)
}

func resourceCloudStackLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	r, _, err := cs.LoadBalancer.GetLoadBalancerByID(d.Id())
	if err != nil {
		return err
	}

	d.Set("algorithm", r.Algorithm)
	d.Set("name", r.Name)
	d.Set("network_id", r.Networkid)
	d.Set("sourceipaddressnetworkid", r.Sourceipaddressnetworkid)

	var vmIds []string
	for _, vm := range r.Loadbalancerinstance {
		vmIds = append(vmIds, vm.Id)
	}
	d.Set("virtualmachineids", vmIds)

	return nil
}

func resourceCloudStackLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	_, err := cs.LoadBalancer.DeleteLoadBalancer(cs.LoadBalancer.NewDeleteLoadBalancerParams(d.Id()))
	if err != nil {
		return err
	}

	return nil
}
