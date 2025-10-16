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
	"fmt"
	"log"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackQuotaTariff() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudStackQuotaTariffCreate,
		ReadContext:   resourceCloudStackQuotaTariffRead,
		UpdateContext: resourceCloudStackQuotaTariffUpdate,
		DeleteContext: resourceCloudStackQuotaTariffDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // CloudStack recreates tariffs internally anyway
				Description: "Name of the quota tariff",
			},

			"usage_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Usage type for the quota tariff",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 25 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 25, got: %d", key, v))
					}
					return
				},
			},

			"value": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "Value of the quota tariff",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(float64)
					if v < 0 {
						errs = append(errs, fmt.Errorf("%q cannot be negative, got: %f", key, v))
					}
					return
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the quota tariff",
			},

			"start_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start date for the quota tariff (format: yyyy-MM-dd)",
			},

			"end_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "End date for the quota tariff (format: yyyy-MM-dd)",
			},

			"activation_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Activation rule for the quota tariff",
			},

			// Computed values
			"currency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Currency for the tariff",
			},

			"effective_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Effective date of the tariff",
			},

			"usage_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Usage name",
			},

			"usage_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Usage unit",
			},

			"position": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Position of the tariff",
			},

			"removed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the tariff is removed",
			},
		},
	}
}

func resourceCloudStackQuotaTariffCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	usagetype := d.Get("usage_type").(int)
	value := d.Get("value").(float64)

	p := cs.Quota.NewQuotaTariffCreateParams(name, usagetype, value)

	if description, ok := d.GetOk("description"); ok {
		p.SetDescription(description.(string))
	}

	if startdate, ok := d.GetOk("start_date"); ok {
		p.SetStartdate(startdate.(string))
	}

	if enddate, ok := d.GetOk("end_date"); ok {
		p.SetEnddate(enddate.(string))
	}

	if activationrule, ok := d.GetOk("activation_rule"); ok {
		p.SetActivationrule(activationrule.(string))
	}

	r, err := cs.Quota.QuotaTariffCreate(p)
	if err != nil {
		return diag.Errorf("Error creating quota tariff %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCloudStackQuotaTariffRead(ctx, d, meta)
}

func resourceCloudStackQuotaTariffRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cs := meta.(*cloudstack.CloudStackClient)

	log.Printf("[DEBUG] Reading quota tariff with ID: %s", d.Id())

	// Always try to find by name and usage type first (more reliable)
	name := d.Get("name").(string)
	usageType := d.Get("usage_type").(int)

	if name != "" && usageType > 0 {
		return findTariffByNameAndType(ctx, d, meta, cs)
	}

	// Fallback to ID search if name/usage_type not available
	if d.Id() != "" {
		p := cs.Quota.NewQuotaTariffListParams()
		p.SetId(d.Id())

		r, err := cs.Quota.QuotaTariffList(p)
		if err != nil {
			log.Printf("[DEBUG] Error searching by ID, trying by name: %s", err)
			return findTariffByNameAndType(ctx, d, meta, cs)
		}

		if len(r.QuotaTariffList) == 0 {
			log.Printf("[DEBUG] Quota tariff %s no longer exists", d.Id())
			d.SetId("")
			return nil
		}

		tariff := r.QuotaTariffList[0]
		return setTariffData(d, tariff)
	}

	return diag.Errorf("Cannot read tariff: no ID, name, or usage_type available")
}

func findTariffByNameAndType(ctx context.Context, d *schema.ResourceData, meta interface{}, cs *cloudstack.CloudStackClient) diag.Diagnostics {
	// List all tariffs and find by name and usage type
	p := cs.Quota.NewQuotaTariffListParams()

	r, err := cs.Quota.QuotaTariffList(p)
	if err != nil {
		return diag.Errorf("Error listing quota tariffs: %s", err)
	}

	name := d.Get("name").(string)
	usageType := d.Get("usage_type").(int)

	for _, tariff := range r.QuotaTariffList {
		if tariff.Name == name && tariff.UsageType == usageType {
			d.SetId(tariff.Id)
			return setTariffData(d, tariff)
		}
	}

	log.Printf("[DEBUG] Quota tariff %s (usage type %d) not found", name, usageType)
	d.SetId("")
	return nil
}

func setTariffData(d *schema.ResourceData, tariff *cloudstack.QuotaTariffList) diag.Diagnostics {
	d.Set("name", tariff.Name)
	d.Set("usage_type", tariff.UsageType)
	d.Set("value", tariff.TariffValue)
	d.Set("description", tariff.Description)
	d.Set("end_date", tariff.EndDate)
	d.Set("activation_rule", tariff.ActivationRule)
	d.Set("currency", tariff.Currency)
	d.Set("effective_date", tariff.EffectiveDate)
	d.Set("usage_name", tariff.UsageName)
	d.Set("usage_unit", tariff.UsageUnit)
	d.Set("position", tariff.Position)
	d.Set("removed", tariff.Removed)

	return nil
}

func resourceCloudStackQuotaTariffUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cs := meta.(*cloudstack.CloudStackClient)

	// Use current name in state for update (before any changes)
	oldName, newName := d.GetChange("name")

	// Determine which name to use for the update API call
	updateName := oldName.(string)
	if updateName == "" {
		// If no old name (shouldn't happen), use current name
		updateName = newName.(string)
	}

	log.Printf("[DEBUG] Updating quota tariff '%s' (ID: %s)", updateName, d.Id())

	p := cs.Quota.NewQuotaTariffUpdateParams(updateName)

	if d.HasChange("name") {
		p.SetName(d.Get("name").(string))
	}

	if d.HasChange("description") {
		p.SetDescription(d.Get("description").(string))
	}

	if d.HasChange("value") {
		p.SetValue(d.Get("value").(float64))
	}

	if d.HasChange("start_date") {
		p.SetStartdate(d.Get("start_date").(string))
	}

	if d.HasChange("end_date") {
		p.SetEnddate(d.Get("end_date").(string))
	}

	if d.HasChange("activation_rule") {
		p.SetActivationrule(d.Get("activation_rule").(string))
	}

	_, err := cs.Quota.QuotaTariffUpdate(p)
	if err != nil {
		return diag.Errorf("Error updating quota tariff %s: %s", d.Id(), err)
	}

	return resourceCloudStackQuotaTariffRead(ctx, d, meta)
}

func resourceCloudStackQuotaTariffDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Quota.NewQuotaTariffDeleteParams(d.Id())

	_, err := cs.Quota.QuotaTariffDelete(p)
	if err != nil {
		return diag.Errorf("Error deleting quota tariff %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
