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
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudStackQuota() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudStackQuotaRead,

		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Lists quota for the specified account",
			},

			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Lists quota for the specified domain ID",
			},

			// Computed values
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of quota summaries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account ID",
						},

						"account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account name",
						},

						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain ID",
						},

						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name",
						},

						"quota_value": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Current quota value",
						},

						"quota_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether quota is enabled for this account",
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudStackQuotaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Quota.NewQuotaSummaryParams()

	if account, ok := d.GetOk("account"); ok {
		p.SetAccount(account.(string))
	}

	if domainid, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(domainid.(string))
	}

	r, err := cs.Quota.QuotaSummary(p)
	if err != nil {
		return diag.Errorf("Error retrieving quota summary: %s", err)
	}

	quotas := make([]map[string]interface{}, 0, len(r.QuotaSummary))

	for _, summary := range r.QuotaSummary {
		q := map[string]interface{}{
			"account_id":    summary.Accountid,
			"account":       summary.Account,
			"domain_id":     summary.Domainid,
			"domain":        summary.Domain,
			"quota_value":   summary.Quota,
			"quota_enabled": summary.Quotaenabled,
		}
		quotas = append(quotas, q)
	}

	if err := d.Set("quotas", quotas); err != nil {
		return diag.Errorf("Error setting quotas: %s", err)
	}

	// Generate a unique ID for this data source
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("quota-summary"))
	if account, ok := d.GetOk("account"); ok {
		buf.WriteString(fmt.Sprintf("-account-%s", account.(string)))
	}
	if domainid, ok := d.GetOk("domain_id"); ok {
		buf.WriteString(fmt.Sprintf("-domain-%s", domainid.(string)))
	}

	sha := sha1.Sum([]byte(buf.String()))
	d.SetId(hex.EncodeToString(sha[:]))

	return nil
}
