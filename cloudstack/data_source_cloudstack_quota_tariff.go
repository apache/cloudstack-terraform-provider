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

func dataSourceCloudStackQuotaTariff() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudStackQuotaTariffRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the quota tariff",
			},

			"usage_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Usage type for the quota tariff",
			},

			// Computed values
			"tariffs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of quota tariffs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tariff ID",
						},

						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tariff name",
						},

						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tariff description",
						},

						"usage_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Usage type",
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

						"tariff_value": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Tariff value",
						},

						"end_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tariff end date",
						},

						"effective_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tariff effective date",
						},

						"activation_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tariff activation rule",
						},

						"removed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the tariff is removed",
						},

						"currency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Currency for the tariff",
						},

						"position": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Position of the tariff",
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudStackQuotaTariffRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cs := meta.(*cloudstack.CloudStackClient)

	p := cs.Quota.NewQuotaTariffListParams()

	if name, ok := d.GetOk("name"); ok {
		p.SetName(name.(string))
	}

	if usagetype, ok := d.GetOk("usage_type"); ok {
		p.SetUsagetype(usagetype.(int))
	}

	r, err := cs.Quota.QuotaTariffList(p)
	if err != nil {
		return diag.Errorf("Error retrieving quota tariff list: %s", err)
	}

	tariffs := make([]map[string]interface{}, 0, len(r.QuotaTariffList))

	for _, tariff := range r.QuotaTariffList {
		t := map[string]interface{}{
			"id":              tariff.Id,
			"name":            tariff.Name,
			"description":     tariff.Description,
			"usage_type":      tariff.UsageType,
			"usage_name":      tariff.UsageName,
			"usage_unit":      tariff.UsageUnit,
			"tariff_value":    tariff.TariffValue,
			"end_date":        tariff.EndDate,
			"effective_date":  tariff.EffectiveDate,
			"activation_rule": tariff.ActivationRule,
			"removed":         tariff.Removed,
			"currency":        tariff.Currency,
			"position":        tariff.Position,
		}
		tariffs = append(tariffs, t)
	}

	if err := d.Set("tariffs", tariffs); err != nil {
		return diag.Errorf("Error setting tariffs: %s", err)
	}

	// Generate a unique ID for this data source
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("quota-tariff"))
	if name, ok := d.GetOk("name"); ok {
		buf.WriteString(fmt.Sprintf("-name-%s", name.(string)))
	}
	if usagetype, ok := d.GetOk("usage_type"); ok {
		buf.WriteString(fmt.Sprintf("-usagetype-%d", usagetype.(int)))
	}

	sha := sha1.Sum([]byte(buf.String()))
	d.SetId(hex.EncodeToString(sha[:]))

	return nil
}
