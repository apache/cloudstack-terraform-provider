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
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceCloudStackResourceLimit manages a single CloudStack resource limit at project, domain, or account scope.
func resourceCloudStackResourceLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackResourceLimitCreate,
		Read:   resourceCloudStackResourceLimitRead,
		Update: resourceCloudStackResourceLimitUpdate,
		Delete: resourceCloudStackResourceLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Scope selectors (exactly one of project_id or domain/account is required)
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Project ID to apply the limit to",
				ConflictsWith: []string{"domain_id", "account"},
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain ID to apply the limit to; required if 'account' is set",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Account name to apply the limit to (requires 'domain_id')",
			},

			// Resource type
			"resourcetype": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Numeric CloudStack resource type (0..11). See API docs for mapping.",
				ConflictsWith: []string{"resourcetype_name"},
			},
			"resourcetype_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Resource type name: user_vm, public_ip, volume, snapshot, template, project, network, vpc, cpu, memory, primary_storage, secondary_storage",
				ConflictsWith: []string{"resourcetype"},
			},

			// Limit value
			"max": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum resource limit; use -1 for unlimited",
			},

			// Optional tag discriminator
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional tag for the resource type",
			},

			// Computed fields from response
			"resourcetypename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resolved resource type name reported by CloudStack",
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudStackResourceLimitCreate(d *schema.ResourceData, meta interface{}) error {
	// ID is synthetic; set after update succeeds
	if err := resourceCloudStackResourceLimitUpdate(d, meta); err != nil {
		return err
	}
	d.SetId(buildResourceLimitID(d))
	return resourceCloudStackResourceLimitRead(d, meta)
}

func resourceCloudStackResourceLimitRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)
	p := cs.Limit.NewListResourceLimitsParams()

	// Scope
	if v, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("account"); ok {
		p.SetAccount(v.(string))
	}
	if v, ok := d.GetOk("tag"); ok {
		p.SetTag(v.(string))
	}

	// Type filter
	if v, ok := d.GetOk("resourcetype"); ok {
		p.SetResourcetype(v.(int))
	} else if v, ok := d.GetOk("resourcetype_name"); ok {
		p.SetResourcetypename(v.(string))
	}

	resp, err := cs.Limit.ListResourceLimits(p)
	if err != nil {
		return err
	}

	if resp == nil || len(resp.ResourceLimits) == 0 {
		// If nothing found, unset ID
		d.SetId("")
		return nil
	}

	// Expect single match with the given filters
	rl := resp.ResourceLimits[0]

	_ = d.Set("max", rl.Max)
	_ = d.Set("resourcetypename", rl.Resourcetypename)
	_ = d.Set("domain", rl.Domain)
	_ = d.Set("domain_id", rl.Domainid)
	_ = d.Set("domain_path", rl.Domainpath)
	_ = d.Set("project", rl.Project)
	_ = d.Set("project_id", rl.Projectid)
	_ = d.Set("tag", rl.Tag)

	// Keep ID stable
	d.SetId(buildResourceLimitID(d))
	return nil
}

func resourceCloudStackResourceLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	rtype, err := resolveResourceType(d)
	if err != nil {
		return err
	}

	// Validate scope
	projectID, projectSet := d.GetOk("project_id")
	domainID, domainSet := d.GetOk("domain_id")
	account, accountSet := d.GetOk("account")

	if !projectSet && !domainSet {
		return fmt.Errorf("either 'project_id' or 'domain_id' must be set")
	}
	if accountSet && !domainSet {
		return fmt.Errorf("'account' requires 'domain_id'")
	}

	p := cs.Limit.NewUpdateResourceLimitParams(rtype)

	if v, ok := d.GetOk("max"); ok {
		p.SetMax(int64(v.(int)))
	}
	if projectSet {
		p.SetProjectid(projectID.(string))
	}
	if domainSet {
		p.SetDomainid(domainID.(string))
	}
	if accountSet {
		p.SetAccount(account.(string))
	}
	if v, ok := d.GetOk("tag"); ok {
		p.SetTag(v.(string))
	}

	_, err = cs.Limit.UpdateResourceLimit(p)
	if err != nil {
		return err
	}

	return resourceCloudStackResourceLimitRead(d, meta)
}

func resourceCloudStackResourceLimitDelete(d *schema.ResourceData, meta interface{}) error {
	// Best-effort revert to unlimited (-1) on delete to avoid leaving unwanted caps
	cs := meta.(*cloudstack.CloudStackClient)

	rtype, err := resolveResourceType(d)
	if err != nil {
		return err
	}

	p := cs.Limit.NewUpdateResourceLimitParams(rtype)
	p.SetMax(-1)

	if v, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("account"); ok {
		p.SetAccount(v.(string))
	}
	if v, ok := d.GetOk("tag"); ok {
		p.SetTag(v.(string))
	}

	if _, err := cs.Limit.UpdateResourceLimit(p); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resolveResourceType(d *schema.ResourceData) (int, error) {
	if v, ok := d.GetOk("resourcetype"); ok {
		return v.(int), nil
	}
	if v, ok := d.GetOk("resourcetype_name"); ok {
		name := strings.ToLower(v.(string))
		if t, ok := resourceTypeNameToID[name]; ok {
			return t, nil
		}
		return 0, fmt.Errorf("unknown resourcetype_name: %s", v.(string))
	}
	return 0, fmt.Errorf("either 'resourcetype' or 'resourcetype_name' must be set")
}

func buildResourceLimitID(d *schema.ResourceData) string {
	parts := []string{}
	if v, ok := d.GetOk("project_id"); ok {
		parts = append(parts, fmt.Sprintf("project=%s", v.(string)))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		parts = append(parts, fmt.Sprintf("domain=%s", v.(string)))
	}
	if v, ok := d.GetOk("account"); ok {
		parts = append(parts, fmt.Sprintf("account=%s", v.(string)))
	}
	if v, ok := d.GetOk("resourcetype"); ok {
		parts = append(parts, fmt.Sprintf("type=%d", v.(int)))
	} else if v, ok := d.GetOk("resourcetype_name"); ok {
		parts = append(parts, fmt.Sprintf("type=%s", v.(string)))
	}
	if v, ok := d.GetOk("tag"); ok {
		parts = append(parts, fmt.Sprintf("tag=%s", v.(string)))
	}
	return strings.Join(parts, "|")
}

var resourceTypeNameToID = map[string]int{
	// Official resourcetypename values
	"user_vm":           0, // instance count
	"public_ip":         1,
	"volume":            2,
	"snapshot":          3,
	"template":          4,
	"project":           5,
	"network":           6,
	"vpc":               7,
	"cpu":               8,
	"memory":            9,
	"primary_storage":   10,
	"secondary_storage": 11,
}
