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
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceCloudStackLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudStackLimitsRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"instance", "ip", "volume", "snapshot", "template", "project", "network", "vpc",
					"cpu", "memory", "primarystorage", "secondarystorage",
				}, false), // false disables case-insensitive matching
				Description: "The type of resource to list the limits. Available types are: " +
					"instance, ip, volume, snapshot, template, project, network, vpc, cpu, memory, " +
					"primarystorage, secondarystorage",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List resources by account. Must be used with the domainid parameter.",
			},
			"domainid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List only resources belonging to the domain specified.",
			},
			"projectid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "List resource limits by project.",
			},
			"limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resourcetype": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resourcetypename": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domainid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"projectid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudStackLimitsRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Limit.NewListResourceLimitsParams()

	// Set optional parameters
	if v, ok := d.GetOk("type"); ok {
		typeStr := v.(string)
		if resourcetype, ok := resourceTypeMap[typeStr]; ok {
			p.SetResourcetype(resourcetype)
		} else {
			return fmt.Errorf("invalid type value: %s", typeStr)
		}
	}

	if v, ok := d.GetOk("account"); ok {
		p.SetAccount(v.(string))
	}

	if v, ok := d.GetOk("domainid"); ok {
		p.SetDomainid(v.(string))
	}

	if v, ok := d.GetOk("projectid"); ok {
		p.SetProjectid(v.(string))
	}

	// Retrieve the resource limits
	l, err := cs.Limit.ListResourceLimits(p)
	if err != nil {
		return fmt.Errorf("Error retrieving resource limits: %s", err)
	}

	// Generate a unique ID for this data source
	id := generateDataSourceID(d)
	d.SetId(id)

	limits := make([]map[string]interface{}, 0, len(l.ResourceLimits))

	// Set the resource data
	for _, limit := range l.ResourceLimits {
		limitMap := map[string]interface{}{
			"resourcetype":     limit.Resourcetype,
			"resourcetypename": limit.Resourcetypename,
			"max":              limit.Max,
		}

		if limit.Account != "" {
			limitMap["account"] = limit.Account
		}

		if limit.Domain != "" {
			limitMap["domain"] = limit.Domain
		}

		if limit.Domainid != "" {
			limitMap["domainid"] = limit.Domainid
		}

		if limit.Project != "" {
			limitMap["project"] = limit.Project
		}

		if limit.Projectid != "" {
			limitMap["projectid"] = limit.Projectid
		}

		limits = append(limits, limitMap)
	}

	if err := d.Set("limits", limits); err != nil {
		return fmt.Errorf("Error setting limits: %s", err)
	}

	return nil
}

// generateDataSourceID generates a unique ID for the data source based on its parameters
func generateDataSourceID(d *schema.ResourceData) string {
	var buf bytes.Buffer

	if v, ok := d.GetOk("type"); ok {
		typeStr := v.(string)
		if resourcetype, ok := resourceTypeMap[typeStr]; ok {
			buf.WriteString(fmt.Sprintf("%d-", resourcetype))
		}
	}

	if v, ok := d.GetOk("account"); ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := d.GetOk("domainid"); ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := d.GetOk("projectid"); ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	// Generate a SHA-256 hash of the buffer content
	hash := sha256.Sum256(buf.Bytes())
	return fmt.Sprintf("limits-%s", hex.EncodeToString(hash[:])[:8])
}
