package cloudstack

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func serviceOfferingMergeCommonSchema(s1 map[string]schema.Attribute) map[string]schema.Attribute {
	common := map[string]schema.Attribute{
		"deployment_planner": schema.StringAttribute{
			Optional: true,
		},
		"disk_offering_id": schema.StringAttribute{
			Optional: true,
		},
		"display_text": schema.StringAttribute{
			Required: true,
		},
		"domain_ids": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
		},
		"dynamic_scaling_enabled": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"host_tags": schema.StringAttribute{
			Optional: true,
		},
		"id": schema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"is_volatile": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"limit_cpu_use": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"network_rate": schema.Int32Attribute{
			Optional: true,
		},
		"offer_ha": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"storage_tags": schema.StringAttribute{
			Optional: true,
		},
		"zone_ids": schema.SetAttribute{
			Optional:    true,
			ElementType: types.StringType,
		},
		"disk_hypervisor": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"bytes_read_rate": schema.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_read_rate_max": schema.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_read_rate_max_length": schema.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_write_rate": schema.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_write_rate_max": schema.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_write_rate_max_length": schema.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
			},
		},
		"disk_offering": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"cache_mode": schema.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"disk_offering_strictness": schema.BoolAttribute{
					Required: true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
					},
				},
				"provisioning_type": schema.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"root_disk_size": schema.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"storage_type": schema.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		},
		"disk_storage": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"customized_iops": schema.BoolAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
					},
				},
				"hypervisor_snapshot_reserve": schema.Int32Attribute{
					Optional: true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
				"max_iops": schema.Int64Attribute{
					Optional: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"min_iops": schema.Int64Attribute{
					Optional: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
			},
		},
	}

	for key, value := range s1 {
		common[key] = value
	}

	return common

}
