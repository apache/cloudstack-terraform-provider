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
				Computed: true,
				ForceNew: true,
			},
			"network_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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

func resourceCloudStackDomainRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Reading Domain %s", d.Id())

	p := cs.Domain.NewListDomainsParams()
	p.SetId(d.Id())

	domains, err := cs.Domain.ListDomains(p)
	if err != nil {
		return fmt.Errorf("Error reading Domain %s: %s", d.Id(), err)
	}

	if domains.Count == 0 {
		log.Printf("[DEBUG] Domain %s does no longer exist", d.Id())
		d.SetId("")
		return nil
	}

	domain := domains.Domains[0]
	log.Printf("[DEBUG] Domain %s found: %s", d.Id(), domain.Name)

	d.Set("name", domain.Name)
	d.Set("domain_id", domain.Id)
	d.Set("network_domain", domain.Networkdomain)
	d.Set("parent_domain_id", domain.Parentdomainid)

	return nil
}

func resourceCloudStackDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)

	if d.HasChange("name") || d.HasChange("network_domain") {
		p := cs.Domain.NewUpdateDomainParams(d.Id())

		if d.HasChange("name") {
			p.SetName(name)
		}

		if d.HasChange("network_domain") {
			p.SetNetworkdomain(d.Get("network_domain").(string))
		}

		log.Printf("[DEBUG] Updating Domain %s", name)
		_, err := cs.Domain.UpdateDomain(p)
		if err != nil {
			return fmt.Errorf("Error updating Domain %s: %s", name, err)
		}
		log.Printf("[DEBUG] Domain %s successfully updated", name)
	}

	return resourceCloudStackDomainRead(d, meta)
}

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
