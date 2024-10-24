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
			Description: "The deployment planner for the service offering",
			Optional:    true,
		},
		"disk_offering_id": schema.StringAttribute{
			Description: "The ID of the disk offering",
			Optional:    true,
		},
		"display_text": schema.StringAttribute{
			Description: "The display text of the service offering",
			Required:    true,
		},
		"domain_ids": schema.SetAttribute{
			Description: "the ID of the containing domain(s), null for public offerings",
			Optional:    true,
			ElementType: types.StringType,
		},
		"dynamic_scaling_enabled": schema.BoolAttribute{
			Description: "Enable dynamic scaling of the service offering",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"host_tags": schema.StringAttribute{
			Description: "The host tag for this service offering",
			Optional:    true,
		},
		"id": schema.StringAttribute{
			Description: "uuid of service offering",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"is_volatile": schema.BoolAttribute{
			Description: "Service offering is volatile",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"limit_cpu_use": schema.BoolAttribute{
			Description: "Restrict the CPU usage to committed service offering",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"name": schema.StringAttribute{
			Description: "The name of the service offering",
			Required:    true,
		},
		"network_rate": schema.Int32Attribute{
			Description: "Data transfer rate in megabits per second",
			Optional:    true,
		},
		"offer_ha": schema.BoolAttribute{
			Description: "The HA for the service offering",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
			Default: booldefault.StaticBool(false),
		},
		"storage_tags": schema.StringAttribute{
			Description: "the tags for the service offering",
			Optional:    true,
		},
		"zone_ids": schema.SetAttribute{
			Description: "The ID of the zone(s)",
			Optional:    true,
			ElementType: types.StringType,
		},
		"disk_hypervisor": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"bytes_read_rate": schema.Int64Attribute{
					Description: "io requests read rate of the disk offering",
					Required:    true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_read_rate_max": schema.Int64Attribute{
					Description: "burst requests read rate of the disk offering",
					Required:    true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_read_rate_max_length": schema.Int64Attribute{
					Description: "length (in seconds) of the burst",
					Required:    true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_write_rate": schema.Int64Attribute{
					Description: "io requests write rate of the disk offering",
					Required:    true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_write_rate_max": schema.Int64Attribute{
					Description: "burst io requests write rate of the disk offering",
					Required:    true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"bytes_write_rate_max_length": schema.Int64Attribute{
					Description: "length (in seconds) of the burst",
					Required:    true,
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
					Description: "the cache mode to use for this disk offering. none, writeback or writethrough",
					Required:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"disk_offering_strictness": schema.BoolAttribute{
					Description: "True/False to indicate the strictness of the disk offering association with the compute offering. When set to true, override of disk offering is not allowed when VM is deployed and change disk offering is not allowed for the ROOT disk after the VM is deployed",
					Required:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
					},
				},
				"provisioning_type": schema.StringAttribute{
					Description: "provisioning type used to create volumes. Valid values are thin, sparse, fat.",
					Required:    true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				"root_disk_size": schema.Int64Attribute{
					Description: "the Root disk size in GB.",
					Required:    true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"storage_type": schema.StringAttribute{
					Description: "the storage type of the service offering. Values are local and shared.",
					Required:    true,
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
					Description: "true if disk offering uses custom iops, false otherwise",
					Optional:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
					},
				},
				"hypervisor_snapshot_reserve": schema.Int32Attribute{
					Description: "Hypervisor snapshot reserve space as a percent of a volume (for managed storage using Xen or VMware)",
					Optional:    true,
					PlanModifiers: []planmodifier.Int32{
						int32planmodifier.RequiresReplace(),
					},
				},
				"max_iops": schema.Int64Attribute{
					Description: "max iops of the compute offering",
					Optional:    true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				"min_iops": schema.Int64Attribute{
					Description: "min iops of the compute offering",
					Optional:    true,
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
