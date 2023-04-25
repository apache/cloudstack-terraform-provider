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
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudStackDomain() *schema.Resource {
	return &schema.Resource{
		Read:   resourceCloudStackDomainRead,
		Update: resourceCloudStackDomainUpdate,
		Create: resourceCloudStackDomainCreate,
		Delete: resourceCloudStackDomainDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCloudStackDomainCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	name := d.Get("name").(string)
	domain_id := d.Get("domain_id").(string)
	network_domain := d.Get("network_domain").(string)
	parent_domain_id := d.Get("parent_domain_id").(string)

	// Create a new parameter struct
	p := cs.Domain.NewCreateDomainParams(name)

	if domain_id != "" {
		p.SetDomainid(domain_id)
	}

	if network_domain != "" {
		p.SetNetworkdomain(network_domain)
	}

	if parent_domain_id != "" {
		p.SetParentdomainid(parent_domain_id)
	}

	log.Printf("[DEBUG] Creating Domain %s", name)
	domain, err := cs.Domain.CreateDomain(p)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Domain %s successfully created", name)
	d.SetId(domain.Id)

	return resourceCloudStackDomainRead(d, meta)
}

func resourceCloudStackDomainRead(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceCloudStackDomainUpdate(d *schema.ResourceData, meta interface{}) error { return nil }

func resourceCloudStackDomainDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Domain.NewDeleteDomainParams(d.Id())
	_, err := cs.Domain.DeleteDomain(p)

	if err != nil {
		return fmt.Errorf("Error deleting Domain: %s", err)
	}

	return nil
}
