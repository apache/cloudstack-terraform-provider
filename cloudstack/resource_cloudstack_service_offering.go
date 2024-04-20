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

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewCloudstackServiceOfferingResource() resource.Resource {
	return &resourceCloudstackServiceOffering{}
}

type resourceCloudstackServiceOffering struct {
	ResourceWithConfigure
}

type resourceCloudStackServiceOfferingModel struct {
	Name                      types.String `tfsdk:"name"`
	DisplayText               types.String `tfsdk:"display_text"`
	Customized                types.Bool   `tfsdk:"customized"`
	CpuNumber                 types.Int64  `tfsdk:"cpu_number"`
	CpuNumberMin              types.Int64  `tfsdk:"cpu_number_min"`
	CpuNumberMax              types.Int64  `tfsdk:"cpu_number_max"`
	CpuSpeed                  types.Int64  `tfsdk:"cpu_speed"`
	Memory                    types.Int64  `tfsdk:"memory"`
	MemoryMin                 types.Int64  `tfsdk:"memory_min"`
	MemoryMax                 types.Int64  `tfsdk:"memory_max"`
	HostTags                  types.String `tfsdk:"host_tags"`
	NetworkRate               types.Int64  `tfsdk:"network_rate"`
	OfferHa                   types.Bool   `tfsdk:"offer_ha"`
	DynamicScaling            types.Bool   `tfsdk:"dynamic_scaling"`
	LimitCpuUse               types.Bool   `tfsdk:"limit_cpu_use"`
	Volatile                  types.Bool   `tfsdk:"volatile"`
	DeploymentPlanner         types.String `tfsdk:"deployment_planner"`
	ZoneId                    types.List   `tfsdk:"zone_id"`
	DiskOfferingId            types.String `tfsdk:"disk_offering_id"`
	StorageType               types.String `tfsdk:"storage_type"`
	ProvisioningType          types.String `tfsdk:"provisioning_type"`
	WriteCacheType            types.String `tfsdk:"write_cache_type"`
	QosType                   types.String `tfsdk:"qos_type"`
	DiskReadRateBps           types.Int64  `tfsdk:"disk_read_rate_bps"`
	DiskWriteRateBps          types.Int64  `tfsdk:"disk_write_rate_bps"`
	DiskReadRateIops          types.Int64  `tfsdk:"disk_read_rate_iops"`
	DiskWriteRateIops         types.Int64  `tfsdk:"disk_write_rate_iops"`
	CustomIops                types.Bool   `tfsdk:"custom_iops"`
	MinIops                   types.Int64  `tfsdk:"min_iops"`
	MaxIops                   types.Int64  `tfsdk:"max_iops"`
	HypervisorSnapshotReserve types.Int64  `tfsdk:"hypervisor_snapshot_reserve"`
	RootDiskSize              types.Int64  `tfsdk:"root_disk_size"`
	StorageTags               types.String `tfsdk:"storage_tags"`
	Encrypt                   types.Bool   `tfsdk:"encrypt"`
	DiskOfferingStrictness    types.Bool   `tfsdk:"disk_offering_strictness"`
	Id                        types.String `tfsdk:"id"`
}

func (r *resourceCloudstackServiceOffering) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_service_offering", req.ProviderTypeName)
}

func (r *resourceCloudstackServiceOffering) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the service offering",
				Required:    true,
			},
			"display_text": schema.StringAttribute{
				Description: "The display text of the service offering",
				Required:    true,
			},
			"customized": schema.BoolAttribute{
				Description: "Is the service offering customized",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("cpu_number"), path.MatchRoot("memory")),
				},
			},
			"cpu_number": schema.Int64Attribute{
				Description: "Number of CPU cores",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("cpu_number_min"), path.MatchRoot("cpu_number_max")),
				},
			},
			"cpu_number_min": schema.Int64Attribute{
				Description: "Minimum number of CPU cores",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("cpu_number_max")),
					int64validator.AtMostSumOf(path.MatchRoot("cpu_number_max")),
					int64validator.ConflictsWith(path.MatchRoot("cpu_number")),
				},
			},
			"cpu_number_max": schema.Int64Attribute{
				Description: "Maximum number of CPU cores",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("cpu_number_min")),
					int64validator.AtLeastSumOf(path.MatchRoot("cpu_number_min")),
					int64validator.ConflictsWith(path.MatchRoot("cpu_number")),
				},
			},
			"cpu_speed": schema.Int64Attribute{
				Description: "Speed of CPU in Mhz",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"memory": schema.Int64Attribute{
				Description: "The total memory of the service offering in MB",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("memory_min"), path.MatchRoot("memory_max")),
				},
			},
			"memory_min": schema.Int64Attribute{
				Description: "Minimum memory of the service offering in MB",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("memory_max")),
					int64validator.AtMostSumOf(path.MatchRoot("memory_max")),
					int64validator.ConflictsWith(path.MatchRoot("memory")),
				},
			},
			"memory_max": schema.Int64Attribute{
				Description: "Maximum memory of the service offering in MB",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.MatchRoot("memory_min")),
					int64validator.AtLeastSumOf(path.MatchRoot("memory_min")),
					int64validator.ConflictsWith(path.MatchRoot("memory")),
				},
			},
			"host_tags": schema.StringAttribute{
				Description: "The host tag for this service offering",
				Optional:    true,
				Computed:    true,
			},
			"network_rate": schema.Int64Attribute{
				Description: "Data transfer rate in megabits per second",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
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
			"dynamic_scaling": schema.BoolAttribute{
				Description: "Enable dynamic scaling of the service offering",
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
			"volatile": schema.BoolAttribute{
				Description: "Service offering is volatile",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
			},
			"deployment_planner": schema.StringAttribute{
				Description: "The deployment planner for the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("FirstFitPlanner", "UserDispersingPlanner", "UserConcentratedPodPlanner", "ImplicitDedicationPlanner", "BareMetalPlanner"),
				},
			},
			"zone_id": schema.ListAttribute{
				Description: "The ID of the zone(s)",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"disk_offering_id": schema.StringAttribute{
				Description: "The ID of the disk offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"storage_type": schema.StringAttribute{
				Description: "The storage type of the service offering. Values are local and shared",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("local", "shared"),
					stringvalidator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
				Default: stringdefault.StaticString("shared"),
			},
			"provisioning_type": schema.StringAttribute{
				Description: "The provisioning type of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("thin", "sparse", "fat"),
					stringvalidator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
				Default: stringdefault.StaticString("thin"),
			},
			"write_cache_type": schema.StringAttribute{
				Description: "The write cache type of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("none", "writeback", "writethrough"),
					stringvalidator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
				Default: stringdefault.StaticString("none"),
			},
			"qos_type": schema.StringAttribute{
				Description: "The QoS type of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("none", "hypervisor", "storage"),
					stringvalidator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
				Default: stringdefault.StaticString("none"),
			},
			"disk_read_rate_bps": schema.Int64Attribute{
				Description: "The read rate of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
			},
			"disk_write_rate_bps": schema.Int64Attribute{
				Description: "The write rate of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
			},
			"disk_read_rate_iops": schema.Int64Attribute{
				Description: "The read rate of the service offering in IOPS",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
			},
			"disk_write_rate_iops": schema.Int64Attribute{
				Description: "The write rate of the service offering in IOPS",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
			},
			"custom_iops": schema.BoolAttribute{
				Description: "Custom IOPS",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("disk_offering_id")),
					boolvalidator.ConflictsWith(path.MatchRoot("disk_read_rate_iops"), path.MatchRoot("disk_write_rate_iops")),
					boolvalidator.ConflictsWith(path.MatchRoot("disk_read_rate_bps"), path.MatchRoot("disk_write_rate_bps")),
				},
			},
			"min_iops": schema.Int64Attribute{
				Description: "The minimum IOPS of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
					int64validator.ConflictsWith(path.MatchRoot("disk_read_rate_iops"), path.MatchRoot("disk_write_rate_iops")),
					int64validator.ConflictsWith(path.MatchRoot("disk_read_rate_bps"), path.MatchRoot("disk_write_rate_bps")),
					int64validator.ConflictsWith(path.MatchRoot("custom_iops")),
				},
			},
			"max_iops": schema.Int64Attribute{
				Description: "The maximum IOPS of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
					int64validator.ConflictsWith(path.MatchRoot("disk_read_rate_iops"), path.MatchRoot("disk_write_rate_iops")),
					int64validator.ConflictsWith(path.MatchRoot("disk_read_rate_bps"), path.MatchRoot("disk_write_rate_bps")),
					int64validator.ConflictsWith(path.MatchRoot("custom_iops")),
				},
			},
			"hypervisor_snapshot_reserve": schema.Int64Attribute{
				Description: "The hypervisor snapshot reserve of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
					int64validator.ConflictsWith(path.MatchRoot("disk_read_rate_iops"), path.MatchRoot("disk_write_rate_iops")),
					int64validator.ConflictsWith(path.MatchRoot("disk_read_rate_bps"), path.MatchRoot("disk_write_rate_bps")),
				},
			},
			"root_disk_size": schema.Int64Attribute{
				Description: "The size of the root disk in GB",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
			},
			"storage_tags": schema.StringAttribute{
				Description: "The storage tags of the service offering",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
			},
			"encrypt": schema.BoolAttribute{
				Description: "Encrypt the service offering storage",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("disk_offering_id")),
				},
			},
			"disk_offering_strictness": schema.BoolAttribute{
				Description: "Disk offering strictness",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Default: booldefault.StaticBool(false),
			},
			"id": schema.StringAttribute{
				Description: "The ID of the service offering",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *resourceCloudstackServiceOffering) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resourceCloudStackServiceOfferingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p := r.client.ServiceOffering.NewCreateServiceOfferingParams(data.DisplayText.ValueString(), data.Name.ValueString())

	if !data.Customized.ValueBool() {
		if !data.CpuNumber.IsNull() && !data.CpuNumber.IsUnknown() {
			p.SetCpunumber(int(data.CpuNumber.ValueInt64()))
		}

		if !data.Memory.IsNull() && !data.Memory.IsUnknown() {
			p.SetMemory(int(data.Memory.ValueInt64()))
		}

		data.CpuNumberMax = types.Int64Null()
		data.CpuNumberMin = types.Int64Null()
		data.MemoryMax = types.Int64Null()
		data.MemoryMin = types.Int64Null()
	} else {
		if !data.CpuNumberMin.IsNull() && !data.CpuNumberMin.IsUnknown() {
			p.SetMincpunumber(int(data.CpuNumberMin.ValueInt64()))
		}

		if !data.CpuNumberMax.IsNull() && !data.CpuNumberMax.IsUnknown() {
			p.SetMaxcpunumber(int(data.CpuNumberMax.ValueInt64()))
		}

		if !data.MemoryMin.IsNull() && !data.MemoryMin.IsUnknown() {
			p.SetMinmemory(int(data.MemoryMin.ValueInt64()))
		}

		if !data.MemoryMax.IsNull() && !data.MemoryMax.IsUnknown() {
			p.SetMaxmemory(int(data.MemoryMax.ValueInt64()))
		}

		data.CpuNumber = types.Int64Null()
		data.Memory = types.Int64Null()
	}

	if !data.CpuSpeed.IsNull() && !data.CpuSpeed.IsUnknown() {
		p.SetCpuspeed(int(data.CpuSpeed.ValueInt64()))
	}

	if !data.HostTags.IsNull() && !data.HostTags.IsUnknown() {
		p.SetHosttags(data.HostTags.ValueString())
	} else {
		data.HostTags = types.StringNull()
	}

	if !data.NetworkRate.IsNull() && !data.NetworkRate.IsUnknown() {
		p.SetNetworkrate(int(data.NetworkRate.ValueInt64()))
	} else {
		data.NetworkRate = types.Int64Null()
	}

	if !data.OfferHa.IsNull() && !data.OfferHa.IsUnknown() {
		p.SetOfferha(data.OfferHa.ValueBool())
	}

	if !data.DynamicScaling.IsNull() && !data.DynamicScaling.IsUnknown() {
		p.SetDynamicscalingenabled(data.DynamicScaling.ValueBool())
	}

	if !data.LimitCpuUse.IsNull() && !data.LimitCpuUse.IsUnknown() {
		p.SetLimitcpuuse(data.LimitCpuUse.ValueBool())
	}

	if !data.Volatile.IsNull() && !data.Volatile.IsUnknown() {
		p.SetIsvolatile(data.Volatile.ValueBool())
	}

	if !data.DeploymentPlanner.IsNull() && !data.DeploymentPlanner.IsUnknown() {
		p.SetDeploymentplanner(data.DeploymentPlanner.ValueString())
	} else {
		data.DeploymentPlanner = types.StringNull()
	}

	if !data.ZoneId.IsNull() && !data.ZoneId.IsUnknown() {
		zoneIds := make([]string, 0, len(data.ZoneId.Elements()))
		diags := data.ZoneId.ElementsAs(ctx, &zoneIds, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		p.SetZoneid(zoneIds)
	} else {
		data.ZoneId = types.ListNull(types.StringType)
	}

	if !data.DiskOfferingId.IsNull() && !data.DiskOfferingId.IsUnknown() {
		p.SetDiskofferingid(data.DiskOfferingId.ValueString())
	} else {
		p.SetDiskofferingid("")
		data.DiskOfferingId = types.StringNull()

		if !data.StorageType.IsNull() && !data.StorageType.IsUnknown() {
			p.SetStoragetype(data.StorageType.ValueString())
		} else {
			data.StorageType = types.StringNull()
		}

		if !data.ProvisioningType.IsNull() && !data.ProvisioningType.IsUnknown() {
			p.SetProvisioningtype(data.ProvisioningType.ValueString())
		}

		if !data.WriteCacheType.IsNull() && !data.WriteCacheType.IsUnknown() {
			p.SetCachemode(data.WriteCacheType.ValueString())
		}

		if data.QosType.ValueString() == "hypervisor" {
			if !data.DiskReadRateBps.IsNull() && !data.DiskReadRateBps.IsUnknown() {
				p.SetBytesreadrate(data.DiskReadRateBps.ValueInt64())
			}

			if !data.DiskWriteRateBps.IsNull() && !data.DiskWriteRateBps.IsUnknown() {
				p.SetByteswriterate(data.DiskWriteRateBps.ValueInt64())
			}

			if !data.DiskReadRateIops.IsNull() && !data.DiskReadRateIops.IsUnknown() {
				p.SetIopsreadrate(data.DiskReadRateIops.ValueInt64())
			}

			if !data.DiskWriteRateIops.IsNull() && !data.DiskWriteRateIops.IsUnknown() {
				p.SetIopswriterate(data.DiskWriteRateIops.ValueInt64())
			}
		} else if data.QosType.ValueString() == "storage" {
			p.SetCustomizediops(data.CustomIops.ValueBool())

			if !data.CustomIops.ValueBool() {
				if !data.MinIops.IsNull() && !data.MinIops.IsUnknown() {
					p.SetMiniops(data.MinIops.ValueInt64())
				}

				if !data.MaxIops.IsNull() && !data.MaxIops.IsUnknown() {
					p.SetMaxiops(data.MaxIops.ValueInt64())
				}
			} else {
				p.SetMiniops(0)
				p.SetMaxiops(0)
			}

			if !data.HypervisorSnapshotReserve.IsNull() && !data.HypervisorSnapshotReserve.IsUnknown() {
				p.SetHypervisorsnapshotreserve(int(data.HypervisorSnapshotReserve.ValueInt64()))
			}
		} else {
			data.DiskReadRateBps = types.Int64Null()
			data.DiskWriteRateBps = types.Int64Null()
			data.DiskReadRateIops = types.Int64Null()
			data.DiskWriteRateIops = types.Int64Null()
			data.MinIops = types.Int64Null()
			data.MaxIops = types.Int64Null()
			data.HypervisorSnapshotReserve = types.Int64Null()
		}
	}

	if !data.RootDiskSize.IsNull() && !data.RootDiskSize.IsUnknown() {
		p.SetRootdisksize(data.RootDiskSize.ValueInt64())
	} else {
		data.RootDiskSize = types.Int64Null()
	}

	if !data.StorageTags.IsNull() && !data.StorageTags.IsUnknown() {
		p.SetTags(data.StorageTags.ValueString())
	} else {
		data.StorageTags = types.StringNull()
	}

	if !data.Encrypt.IsNull() && !data.Encrypt.IsUnknown() {
		p.SetEncryptroot(data.Encrypt.ValueBool())
	}

	if !data.DiskOfferingStrictness.IsNull() && !data.DiskOfferingStrictness.IsUnknown() {
		p.SetDiskofferingstrictness(data.DiskOfferingStrictness.ValueBool())
	}

	s, err := r.client.ServiceOffering.CreateServiceOffering(p)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service offering",
			fmt.Sprintf("Error while trying to create service offering: %s", err),
		)
		return
	}

	data.Id = types.StringValue(s.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceCloudstackServiceOffering) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resourceCloudStackServiceOfferingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	s, count, err := r.client.ServiceOffering.GetServiceOfferingByID(data.Id.ValueString())

	if count == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading service offering",
			fmt.Sprintf("Error while trying to read service offering: %s", err),
		)
		return
	}

	data.Id = types.StringValue(s.Id)
	data.Name = types.StringValue(s.Name)
	data.DisplayText = types.StringValue(s.Displaytext)
	data.CpuNumber = types.Int64Value(int64(s.Cpunumber))
	data.CpuSpeed = types.Int64Value(int64(s.Cpuspeed))
	data.HostTags = types.StringValue(s.Hosttags)
	data.LimitCpuUse = types.BoolValue(s.Limitcpuuse)
	data.Memory = types.Int64Value(int64(s.Memory))
	data.OfferHa = types.BoolValue(s.Offerha)
	data.StorageType = types.StringValue(s.Storagetype)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceCloudstackServiceOffering) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resourceCloudStackServiceOfferingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	p := r.client.ServiceOffering.NewUpdateServiceOfferingParams(data.Id.ValueString())

	p.SetName(data.Name.ValueString())
	p.SetDisplaytext(data.DisplayText.ValueString())
	p.SetHosttags(data.HostTags.ValueString())
	p.SetStoragetags(data.HostTags.ValueString())

	_, err := r.client.ServiceOffering.UpdateServiceOffering(p)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating service offering",
			fmt.Sprintf("Error while trying to update service offering: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *resourceCloudstackServiceOffering) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resourceCloudStackServiceOfferingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	p := r.client.ServiceOffering.NewDeleteServiceOfferingParams(data.Id.ValueString())
	_, err := r.client.ServiceOffering.DeleteServiceOffering(p)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting service offering",
			fmt.Sprintf("Error while trying to delete service offering: %s", err),
		)
		return
	}
}

func (r *resourceCloudstackServiceOffering) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
